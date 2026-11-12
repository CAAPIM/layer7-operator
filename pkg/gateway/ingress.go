/*
* Copyright (c) 2026 Broadcom. All rights reserved.
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
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */
package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewIngress(gw *securityv1.Gateway) *networkingv1.Ingress {

	if !gw.Spec.App.Ingress.Enabled {
		return nil
	}

	tls := gw.Spec.App.Ingress.TLS
	rules := []networkingv1.IngressRule{}
	ingressClassName := gw.Spec.App.Ingress.IngressClassName
	annotations := gw.Spec.App.Ingress.Annotations

	const portName = "https"
	pathTypePrefix := networkingv1.PathTypePrefix

	for _, r := range gw.Spec.App.Ingress.Rules {
		rule := networkingv1.IngressRule{
			Host: r.Host,
		}
		paths := []networkingv1.HTTPIngressPath{}

		if r.HTTP != nil {
			if r.HTTP.Paths != nil {
				for _, p := range r.HTTP.Paths {

					if p.Backend.Service != nil {

						if p.Backend.Service.Name == "" {
							p.Backend.Service.Name = gw.Name
						}
						if p.Backend.Service.Name == "management" {
							p.Backend.Service.Name = gw.Name + "-management-service"
						}

						if p.Backend.Service.Port.Name == "" {
							p.Backend.Service.Port.Name = portName
						}

						if p.Backend.Service.Port.Number != 0 {
							p.Backend.Service.Port.Name = ""
						}
					}

					if p.PathType == nil {
						p.PathType = &pathTypePrefix
					}

					if p.Path == "" {
						p.Path = "/"
					}
					paths = append(paths, p)

				}
			}
		} else {
			path := networkingv1.HTTPIngressPath{
				Path:     "/",
				PathType: &pathTypePrefix,
				Backend: networkingv1.IngressBackend{
					Service: &networkingv1.IngressServiceBackend{
						Name: gw.Name,
						Port: networkingv1.ServiceBackendPort{
							Name: portName,
						},
					},
				},
			}
			paths = append(paths, path)
		}

		rule.HTTP = &networkingv1.HTTPIngressRuleValue{
			Paths: paths,
		}

		rules = append(rules, rule)
	}

	ls := util.DefaultLabels(gw.Name, gw.Spec.App.Labels)
	spec := networkingv1.IngressSpec{
		TLS:   tls,
		Rules: rules,
	}
	if ingressClassName != "" {
		name := ingressClassName
		spec.IngressClassName = &name
	}
	ingress := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        gw.Name,
			Namespace:   gw.Namespace,
			Annotations: annotations,
			Labels:      ls,
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: "networking.k8s.io/v1",
			Kind:       "Ingress",
		},
		Spec: spec,
	}
	return ingress
}
