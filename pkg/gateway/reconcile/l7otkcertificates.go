package reconcile

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"strconv"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
	"github.com/caapim/layer7-operator/pkg/gateway"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
)

func syncOtkCertificates(ctx context.Context, params Params) {
	gateway := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gateway)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = removeJob(params.Instance.Name + "-" + params.Instance.Namespace + "-sync-otk-certificates")
		return
	}

	if !gateway.Spec.App.Otk.Enabled {
		_ = removeJob(params.Instance.Name + "-" + params.Instance.Namespace + "-sync-otk-certificates")
		return
	}

	err = applyOtkCertificates(ctx, params, gateway)
	if err != nil {
		params.Log.Info("failed to reconcile otk certificates", "name", gateway.Name, "namespace", gateway.Namespace, "error", err.Error())
	}

}

func applyOtkCertificates(ctx context.Context, params Params, gateway *securityv1.Gateway) error {

	bundle := graphman.Bundle{}
	annotation := ""
	sha1Sum := ""

	switch gateway.Spec.App.Otk.Type {
	case securityv1.OtkTypeDMZ:
		internalSecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.InternalOtkGatewayReference+"-otk-internal-certificates")
		sha1Sum = internalSecret.ObjectMeta.Annotations["checksum/data"]

		if err != nil {
			return err
		}
		annotation = "security.brcmlabs.com/" + gateway.Name + "-" + string(gateway.Spec.App.Otk.Type) + "-certificates"
		for k, v := range internalSecret.Data {
			bundle.TrustedCerts = append(bundle.TrustedCerts, &graphman.TrustedCertInput{
				Name:                      k,
				CertBase64:                base64.StdEncoding.EncodeToString(v),
				TrustAnchor:               true,
				VerifyHostname:            false,
				RevocationCheckPolicyType: "USE_DEFAULT",
				TrustedFor: []graphman.TrustedForType{
					"SSL",
					"SIGNING_SERVER_CERTS",
				},
			})
		}

	case securityv1.OtkTypeInternal:
		dmzSecret, err := getGatewaySecret(ctx, params, gateway.Spec.App.Otk.DmzOtkGatewayReference+"-otk-dmz-certificates")
		sha1Sum = dmzSecret.ObjectMeta.Annotations["checksum/data"]
		if err != nil {
			return err
		}
		annotation = "security.brcmlabs.com/" + gateway.Name + "-" + string(gateway.Spec.App.Otk.Type) + "-fips-users"
		for k, v := range dmzSecret.Data {
			bundle.FipUsers = append(bundle.FipUsers, &graphman.FipUserInput{
				Name:         k,
				ProviderName: "otk-fips-provider",
				SubjectDn:    "cn=" + k,
				CertBase64:   base64.RawStdEncoding.EncodeToString(v),
			})
		}
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return err
	}

	name := gateway.Name
	if gateway.Spec.App.Management.SecretName != "" {
		name = gateway.Spec.App.Management.SecretName
	}
	gwSecret, err := getGatewaySecret(ctx, params, name)

	if err != nil {
		return err
	}

	if !gateway.Spec.App.Management.Database.Enabled {
		podList, err := getGatewayPods(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileEphemeralGateway(ctx, params, "otk certificates", *podList, gateway, gwSecret, "", annotation, sha1Sum, true, "otk certificates", bundleBytes)
		if err != nil {
			return err
		}
	} else {
		gatewayDeployment, err := getGatewayDeployment(ctx, params)
		if err != nil {
			return err
		}
		err = ReconcileDBGateway(ctx, params, "otk certificates", gatewayDeployment, gateway, gwSecret, "", annotation, sha1Sum, false, bundleBytes)
		if err != nil {
			return err
		}
	}

	return nil
}

func manageCertificateSecrets(ctx context.Context, params Params) {
	gw := &securityv1.Gateway{}
	err := params.Client.Get(ctx, types.NamespacedName{Name: params.Instance.Name, Namespace: params.Instance.Namespace}, gw)
	if err != nil && k8serrors.IsNotFound(err) {
		params.Log.Error(err, "gateway not found", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		_ = removeJob(params.Instance.Name + "-sync-otk-certificate-secret")
		return
	}

	if !gw.Spec.App.Otk.Enabled {
		_ = removeJob(params.Instance.Name + "-sync-otk-certificate-secret")
		return
	}
	params.Instance = gw
	internalGatewayPort := 9443
	defaultOtkPort := 8443
	rawInternalCertList := map[string][]byte{}
	rawDMZCertList := map[string][]byte{}
	desiredSecrets := []*corev1.Secret{}
	if gw.Spec.App.Management.Graphman.DynamicSyncPort != 0 {
		internalGatewayPort = gw.Spec.App.Management.Graphman.DynamicSyncPort

	}

	if gw.Spec.App.Otk.InternalGatewayPort != 0 {
		internalGatewayPort = gw.Spec.App.Otk.InternalGatewayPort
	}

	if gw.Spec.App.Otk.OTKPort != 0 {
		defaultOtkPort = gw.Spec.App.Otk.OTKPort
	}
	podList, err := getGatewayPods(ctx, params)
	if err != nil {
		params.Log.Error(err, "failed to retrieve gateway pods", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		return
	}

	for _, pod := range podList.Items {
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.Name == "gateway" {
				if !containerStatus.Ready {
					params.Log.V(2).Info("pod not ready", "pod", pod.Name, "name", params.Instance.Name, "namespace", params.Instance.Namespace)
					return
				}
			}
		}

		switch gw.Spec.App.Otk.Type {
		case securityv1.OtkTypeDMZ:
			rawCert, err := retrieveCertificate(pod.Status.PodIP, strconv.Itoa(defaultOtkPort))
			if err != nil {
				params.Log.Error(err, "failed to retrieve certificate", "pod", pod.Name, "name", params.Instance.Name, "namespace", params.Instance.Namespace)
				return
			}
			if len(rawDMZCertList) > 0 {
				for _, cert := range rawDMZCertList {
					if string(rawCert) != string(cert) {
						rawDMZCertList[pod.Name] = rawCert
					}
				}
			} else {
				rawDMZCertList[pod.Name] = rawCert
			}
		case securityv1.OtkTypeInternal:
			rawCert, err := retrieveCertificate(pod.Status.PodIP, strconv.Itoa(internalGatewayPort))
			if err != nil {
				params.Log.Error(err, "failed to retrieve certificate", "pod", pod.Name, "name", params.Instance.Name, "namespace", params.Instance.Namespace)
				return
			}
			if len(rawInternalCertList) > 0 {
				for _, cert := range rawInternalCertList {
					if string(rawCert) != string(cert) {
						rawInternalCertList[pod.Name] = rawCert
					}
				}
			} else {
				rawInternalCertList[pod.Name] = rawCert
			}

		}
	}

	if gw.Spec.App.Otk.Type == securityv1.OtkTypeDMZ && len(rawDMZCertList) > 0 {
		desiredSecrets = append(desiredSecrets, gateway.NewOtkCertificateSecret(gw, gw.Name+"-otk-dmz-certificates", rawDMZCertList))
	}

	if gw.Spec.App.Otk.Type == securityv1.OtkTypeInternal && len(rawInternalCertList) > 0 {
		desiredSecrets = append(desiredSecrets, gateway.NewOtkCertificateSecret(gw, gw.Name+"-otk-internal-certificates", rawInternalCertList))
	}

	err = reconcileSecrets(ctx, params, desiredSecrets)
	if err != nil {
		params.Log.Error(err, "failed to reconcile otk certificates", "Name", gw.Name, "namespace", gw.Namespace)
		return
	}

}

func retrieveCertificate(host string, port string) ([]byte, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.Dial("tcp", host+":"+port, conf)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	cert := conn.ConnectionState().PeerCertificates[0].Raw
	return cert, nil
}

// func getSha1Thumbprint(rawCert []byte) (string, error) {
// 	fingerprint := sha1.Sum(rawCert)
// 	var buf bytes.Buffer
// 	for _, f := range fingerprint {
// 		fmt.Fprintf(&buf, "%02X", f)
// 	}
// 	hexDump, err := hex.DecodeString(buf.String())
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.StdEncoding.EncodeToString(hexDump), nil
// }
