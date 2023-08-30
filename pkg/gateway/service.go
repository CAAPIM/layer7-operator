package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func NewService(gw *securityv1.Gateway) *corev1.Service {

	ports := []corev1.ServicePort{}

	for p := range gw.Spec.App.Service.Ports {
		ports = append(ports, corev1.ServicePort{
			Name:       gw.Spec.App.Service.Ports[p].Name,
			Port:       gw.Spec.App.Service.Ports[p].Port,
			TargetPort: intstr.FromString(gw.Spec.App.Service.Ports[p].Name),
			Protocol:   corev1.Protocol(gw.Spec.App.Service.Ports[p].Protocol),
		})
	}

	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gw.Name,
			Namespace:   gw.Namespace,
			Annotations: gw.Spec.App.Service.Annotations,
			Labels:      ls,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports:    ports,
			Type:     gw.Spec.App.Service.Type,
		},
	}
	return service
}

func NewManagementService(gw *securityv1.Gateway) *corev1.Service {
	ports := []corev1.ServicePort{}

	for p := range gw.Spec.App.Management.Service.Ports {
		ports = append(ports, corev1.ServicePort{
			Name:       gw.Spec.App.Management.Service.Ports[p].Name,
			Port:       gw.Spec.App.Management.Service.Ports[p].Port,
			TargetPort: intstr.FromString(gw.Spec.App.Management.Service.Ports[p].Name),
			Protocol:   corev1.Protocol(gw.Spec.App.Management.Service.Ports[p].Protocol),
		})
	}

	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	mls := map[string]string{"management-access": "leader"}

	for k, v := range mls {
		ls[k] = v
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gw.Name + "-management-service",
			Namespace:   gw.Namespace,
			Annotations: gw.Spec.App.Management.Service.Annotations,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports:    ports,
			Type:     gw.Spec.App.Management.Service.Type,
		},
	}
	return service
}
