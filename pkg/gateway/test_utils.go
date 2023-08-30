package gateway

import (
	securityv1 "github.com/caapim/layer7-operator/api/v1"
)

func getGatewayWitApp() securityv1.Gateway {
	gateway := securityv1.Gateway{}
	gateway.Spec = securityv1.GatewaySpec{}
	gateway.Spec.App = securityv1.App{}
	return gateway
}
