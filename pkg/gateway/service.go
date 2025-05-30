/*
* Copyright (c) 2025 Broadcom. All rights reserved.
* The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
* All trademarks, trade names, service marks, and logos referenced
* herein belong to their respective companies.
*
* This software and all information contained therein is confidential
* and proprietary and shall not be duplicated, used, disclosed or
* disseminated in any way except as authorized by the applicable
* license agreement, without the express written permission of Broadcom.
* All authorized reproductions must be marked with this language.
*
* EXCEPT AS SET FORTH IN THE APPLICABLE LICENSE AGREEMENT, TO THE
* EXTENT PERMITTED BY APPLICABLE LAW OR AS AGREED BY BROADCOM IN ITS
* APPLICABLE LICENSE AGREEMENT, BROADCOM PROVIDES THIS DOCUMENTATION
* "AS IS" WITHOUT WARRANTY OF ANY KIND, INCLUDING WITHOUT LIMITATION,
* ANY IMPLIED WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR
* PURPOSE, OR. NONINFRINGEMENT. IN NO EVENT WILL BROADCOM BE LIABLE TO
* THE END USER OR ANY THIRD PARTY FOR ANY LOSS OR DAMAGE, DIRECT OR
* INDIRECT, FROM THE USE OF THIS DOCUMENTATION, INCLUDING WITHOUT LIMITATION,
* LOST PROFITS, LOST INVESTMENT, BUSINESS INTERRUPTION, GOODWILL, OR
* LOST DATA, EVEN IF BROADCOM IS EXPRESSLY ADVISED IN ADVANCE OF THE
* POSSIBILITY OF SUCH LOSS OR DAMAGE.
*
 */
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

		port := corev1.ServicePort{
			Name:       gw.Spec.App.Service.Ports[p].Name,
			Port:       gw.Spec.App.Service.Ports[p].Port,
			TargetPort: intstr.FromString(gw.Spec.App.Service.Ports[p].Name),
			Protocol:   corev1.Protocol(gw.Spec.App.Service.Ports[p].Protocol),
		}

		if gw.Spec.App.Service.Type == corev1.ServiceTypeNodePort && gw.Spec.App.Service.Ports[p].NodePort != 0 {
			port.NodePort = gw.Spec.App.Service.Ports[p].NodePort
		}

		ports = append(ports, port)
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
			Selector: util.DefaultLabels(gw.Name, nil),
			Ports:    ports,
			Type:     gw.Spec.App.Service.Type,
		},
	}

	if gw.Spec.App.Service.ClusterIP != "" {
		service.Spec.ClusterIP = gw.Spec.App.Service.ClusterIP
	}

	if gw.Spec.App.Service.ClusterIPs != nil {
		service.Spec.ClusterIPs = gw.Spec.App.Service.ClusterIPs
	}

	if gw.Spec.App.Service.ExternalIPs != nil {
		service.Spec.ExternalIPs = gw.Spec.App.Service.ExternalIPs
	}

	if gw.Spec.App.Service.SessionAffinity != "" {
		service.Spec.SessionAffinity = gw.Spec.App.Service.SessionAffinity
	}

	if gw.Spec.App.Service.LoadBalancerIP != "" {
		service.Spec.LoadBalancerIP = gw.Spec.App.Service.LoadBalancerIP
	}

	if gw.Spec.App.Service.LoadBalancerSourceRanges != nil {
		service.Spec.LoadBalancerSourceRanges = gw.Spec.App.Service.LoadBalancerSourceRanges
	}

	if gw.Spec.App.Service.LoadBalancerClass != "" {
		service.Spec.LoadBalancerClass = &gw.Spec.App.Service.LoadBalancerClass
	}

	if gw.Spec.App.Service.ExternalName != "" {
		service.Spec.ExternalName = gw.Spec.App.Service.ExternalName
	}

	if gw.Spec.App.Service.ExternalTrafficPolicy != "" {
		service.Spec.ExternalTrafficPolicy = gw.Spec.App.Service.ExternalTrafficPolicy
	}

	if gw.Spec.App.Service.HealthCheckNodePort != 0 {
		service.Spec.HealthCheckNodePort = gw.Spec.App.Service.HealthCheckNodePort
	}

	if gw.Spec.App.Service.SessionAffinityConfig != (corev1.SessionAffinityConfig{}) {
		service.Spec.SessionAffinityConfig = &gw.Spec.App.Service.SessionAffinityConfig
	}

	if gw.Spec.App.Service.IPFamilies != nil {
		service.Spec.IPFamilies = gw.Spec.App.Service.IPFamilies
	}

	if gw.Spec.App.Service.IPFamilyPolicy != "" {
		service.Spec.IPFamilyPolicy = &gw.Spec.App.Service.IPFamilyPolicy
	}

	if gw.Spec.App.Service.AllocateLoadBalancerNodePorts != nil {
		service.Spec.AllocateLoadBalancerNodePorts = gw.Spec.App.Service.AllocateLoadBalancerNodePorts
	}

	if gw.Spec.App.Service.InternalTrafficPolicy != "" {
		service.Spec.InternalTrafficPolicy = &gw.Spec.App.Service.InternalTrafficPolicy
	}

	return service
}

func NewManagementService(gw *securityv1.Gateway) *corev1.Service {
	ports := []corev1.ServicePort{}

	for p := range gw.Spec.App.Management.Service.Ports {

		port := corev1.ServicePort{
			Name:       gw.Spec.App.Management.Service.Ports[p].Name,
			Port:       gw.Spec.App.Management.Service.Ports[p].Port,
			TargetPort: intstr.FromString(gw.Spec.App.Management.Service.Ports[p].Name),
			Protocol:   corev1.Protocol(gw.Spec.App.Management.Service.Ports[p].Protocol),
		}

		if gw.Spec.App.Management.Service.Type == corev1.ServiceTypeNodePort && gw.Spec.App.Management.Service.Ports[p].NodePort != 0 {
			port.NodePort = gw.Spec.App.Management.Service.Ports[p].NodePort
		}

		ports = append(ports, port)
	}

	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	mls := map[string]string{"management-access": "leader"}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gw.Name + "-management-service",
			Namespace:   gw.Namespace,
			Annotations: gw.Spec.App.Management.Service.Annotations,
			Labels:      ls,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Service",
		},
		Spec: corev1.ServiceSpec{
			Selector: util.DefaultLabels(gw.Name, mls),
			Ports:    ports,
			Type:     gw.Spec.App.Management.Service.Type,
		},
	}

	if gw.Spec.App.Management.Service.ClusterIP != "" {
		service.Spec.ClusterIP = gw.Spec.App.Management.Service.ClusterIP
	}

	if gw.Spec.App.Management.Service.ClusterIPs != nil {
		service.Spec.ClusterIPs = gw.Spec.App.Management.Service.ClusterIPs
	}

	if gw.Spec.App.Management.Service.ExternalIPs != nil {
		service.Spec.ExternalIPs = gw.Spec.App.Management.Service.ExternalIPs
	}

	if gw.Spec.App.Management.Service.SessionAffinity != "" {
		service.Spec.SessionAffinity = gw.Spec.App.Management.Service.SessionAffinity
	}

	if gw.Spec.App.Management.Service.LoadBalancerIP != "" {
		service.Spec.LoadBalancerIP = gw.Spec.App.Management.Service.LoadBalancerIP
	}

	if gw.Spec.App.Management.Service.LoadBalancerSourceRanges != nil {
		service.Spec.LoadBalancerSourceRanges = gw.Spec.App.Management.Service.LoadBalancerSourceRanges
	}

	if gw.Spec.App.Management.Service.LoadBalancerClass != "" {
		service.Spec.LoadBalancerClass = &gw.Spec.App.Management.Service.LoadBalancerClass
	}

	if gw.Spec.App.Management.Service.ExternalName != "" {
		service.Spec.ExternalName = gw.Spec.App.Management.Service.ExternalName
	}

	if gw.Spec.App.Management.Service.ExternalTrafficPolicy != "" {
		service.Spec.ExternalTrafficPolicy = gw.Spec.App.Management.Service.ExternalTrafficPolicy
	}

	if gw.Spec.App.Management.Service.HealthCheckNodePort != 0 {
		service.Spec.HealthCheckNodePort = gw.Spec.App.Management.Service.HealthCheckNodePort
	}

	if gw.Spec.App.Management.Service.SessionAffinityConfig != (corev1.SessionAffinityConfig{}) {
		service.Spec.SessionAffinityConfig = &gw.Spec.App.Management.Service.SessionAffinityConfig
	}

	if gw.Spec.App.Management.Service.IPFamilies != nil {
		service.Spec.IPFamilies = gw.Spec.App.Management.Service.IPFamilies
	}

	if gw.Spec.App.Management.Service.IPFamilyPolicy != "" {
		service.Spec.IPFamilyPolicy = &gw.Spec.App.Management.Service.IPFamilyPolicy
	}

	if gw.Spec.App.Management.Service.AllocateLoadBalancerNodePorts != nil {
		service.Spec.AllocateLoadBalancerNodePorts = gw.Spec.App.Management.Service.AllocateLoadBalancerNodePorts
	}

	if gw.Spec.App.Management.Service.InternalTrafficPolicy != "" {
		service.Spec.InternalTrafficPolicy = &gw.Spec.App.Management.Service.InternalTrafficPolicy
	}

	return service
}
