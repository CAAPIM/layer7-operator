package reconcile

import (
	"github.com/go-logr/logr"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Params struct {
	Client   client.Client
	Recorder record.EventRecorder
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Instance *securityv1.Gateway
}
