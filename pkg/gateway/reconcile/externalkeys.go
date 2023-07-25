package reconcile

// func reconcileExternalKeys(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway) error {
// 	keySecretMap := []util.GraphmanKey{}
// 	bundleBytes := []byte{}

// 	podList, err := getGatewayPods(r, ctx, gw)
// 	if err != nil {
// 		return err
// 	}

// 	for _, externalKey := range gw.Spec.App.ExternalKeys {
// 		if externalKey.Enabled {

// 			secret, err := getSecret(r, ctx, gw, externalKey.Name)
// 			if err != nil {
// 				if k8serrors.IsNotFound(err) {
// 					r.Log.Info("Secret not found", "Name", gw.Name, "Namespace", gw.Namespace, "External Key Ref", externalKey.Name)
// 				} else {
// 					r.Log.Info("Can't retrieve secret", "Name", gw.Name, "Namespace", gw.Namespace, "Error", err.Error())
// 				}
// 			}

// 			if secret.Type == corev1.SecretTypeTLS {
// 				keySecretMap = append(keySecretMap, util.GraphmanKey{
// 					Name: secret.Name,
// 					Crt:  string(secret.Data["tls.crt"]),
// 					Key:  string(secret.Data["tls.key"]),
// 				})
// 			}

// 		}
// 	}

// 	if len(keySecretMap) > 0 {
// 		bundleBytes, err = util.ConvertX509ToGraphmanBundle(keySecretMap)
// 		if err != nil {
// 			r.Log.Info("Can't convert secrets to Graphman bundle", "Name", gw.Name, "Namespace", gw.Namespace, "Error", err.Error())
// 		}
// 	} else {
// 		return nil
// 	}

// 	sort.Slice(keySecretMap, func(i, j int) bool {
// 		return keySecretMap[i].Name < keySecretMap[j].Name
// 	})

// 	keySecretMapBytes, err := json.Marshal(keySecretMap)

// 	if err != nil {
// 		return err
// 	}
// 	h := sha1.New()
// 	h.Write(keySecretMapBytes)
// 	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

// 	patch := fmt.Sprintf("{\"metadata\": {\"labels\": {\"%s\": \"%s\"}}}", "security.brcmlabs.com/external-keys", sha1Sum)

// 	name := gw.Name
// 	if gw.Spec.App.Management.SecretName != "" {
// 		name = gw.Spec.App.Management.SecretName
// 	}
// 	gwSecret, err := getSecret(r, ctx, gw, name)

// 	if err != nil {
// 		return err
// 	}

// 	for i, pod := range podList.Items {
// 		ready := false

// 		for _, containerStatus := range pod.Status.ContainerStatuses {
// 			if containerStatus.Name == "gateway" {
// 				ready = containerStatus.Ready
// 			}
// 		}

// 		if ready && pod.Labels["security.brcmlabs.com/external-keys"] != sha1Sum {
// 			endpoint := pod.Status.PodIP + ":9443/graphman"

// 			r.Log.Info("Applying Latest Secret Bundle", "Secret SHA", sha1Sum, "Pod", pod.Name, "Name", gw.Name, "Namespace", gw.Namespace)

// 			err = util.ApplyGraphmanBundle(string(gwSecret.Data["SSG_ADMIN_USERNAME"]), string(gwSecret.Data["SSG_ADMIN_PASSWORD"]), endpoint, "7layer", bundleBytes)
// 			if err != nil {
// 				return err
// 			}

// 			if err := r.Client.Patch(context.Background(), &podList.Items[i],
// 				client.RawPatch(types.StrategicMergePatchType, []byte(patch))); err != nil {
// 				r.Log.Error(err, "Failed to update pod label", "Namespace", gw.Namespace, "Name", gw.Name)
// 				return err
// 			}
// 		}
// 	}

// 	return nil
// }
