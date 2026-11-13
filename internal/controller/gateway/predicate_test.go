package gateway

import (
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

func TestGatewayAPIUpdate(t *testing.T) {
	meta := func(generation int64) metav1.ObjectMeta {
		return metav1.ObjectMeta{Name: "test", Generation: generation}
	}
	tests := []struct {
		name            string
		oldObj, newObj  client.Object
		wantGatewayAPI  bool
		wantSpecChanged bool
	}{
		{
			name:           "Gateway status only",
			oldObj:         &gatewayv1.Gateway{ObjectMeta: meta(1)},
			newObj:         &gatewayv1.Gateway{ObjectMeta: meta(1), Status: gatewayv1.GatewayStatus{Conditions: []metav1.Condition{{Type: "Programmed"}}}},
			wantGatewayAPI: true,
		},
		{
			name:            "Gateway spec",
			oldObj:          &gatewayv1.Gateway{ObjectMeta: meta(1)},
			newObj:          &gatewayv1.Gateway{ObjectMeta: meta(2)},
			wantGatewayAPI:  true,
			wantSpecChanged: true,
		},
		{
			name:           "HTTPRoute status only",
			oldObj:         &gatewayv1.HTTPRoute{ObjectMeta: meta(3)},
			newObj:         &gatewayv1.HTTPRoute{ObjectMeta: meta(3)},
			wantGatewayAPI: true,
		},
		{
			name:            "TLSRoute spec",
			oldObj:          &gatewayv1.TLSRoute{ObjectMeta: meta(1)},
			newObj:          &gatewayv1.TLSRoute{ObjectMeta: meta(2)},
			wantGatewayAPI:  true,
			wantSpecChanged: true,
		},
		{
			name:           "BackendTLSPolicy status only",
			oldObj:         &gatewayv1.BackendTLSPolicy{ObjectMeta: meta(1)},
			newObj:         &gatewayv1.BackendTLSPolicy{ObjectMeta: meta(1)},
			wantGatewayAPI: true,
		},
		{
			name:   "Layer7 Gateway is not a Gateway API resource",
			oldObj: &securityv1.Gateway{ObjectMeta: meta(1)},
			newObj: &securityv1.Gateway{ObjectMeta: meta(2)},
		},
		{
			name:   "Deployment is not a Gateway API resource",
			oldObj: &appsv1.Deployment{ObjectMeta: meta(1)},
			newObj: &appsv1.Deployment{ObjectMeta: meta(2)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isGatewayAPI, specChanged := gatewayAPIUpdate(event.UpdateEvent{ObjectOld: tt.oldObj, ObjectNew: tt.newObj})
			if isGatewayAPI != tt.wantGatewayAPI || specChanged != tt.wantSpecChanged {
				t.Fatalf("gatewayAPIUpdate() = (%v, %v), want (%v, %v)", isGatewayAPI, specChanged, tt.wantGatewayAPI, tt.wantSpecChanged)
			}
		})
	}
}
