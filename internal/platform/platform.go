package platform

import (
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
)

type Platform string

const (
	PlatformTypeKubernetes Platform = "kubernetes"
	PlatformTypeOpenshift  Platform = "openshift"
)

func Detect(restConfig rest.Config) (Platform, error) {
	dClient, err := discovery.NewDiscoveryClientForConfig(&restConfig)

	if err != nil {
		return "", err
	}

	serverGroups, err := dClient.ServerGroups()
	if err != nil {
		return "", err
	}

	for _, sg := range serverGroups.Groups {
		if sg.Name == "route.openshift.io" {
			return PlatformTypeOpenshift, nil
		}
	}

	return PlatformTypeKubernetes, nil

}
