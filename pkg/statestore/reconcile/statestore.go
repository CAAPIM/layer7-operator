/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package reconcile

import (
	"context"
	"reflect"

	securityv1alpha1 "github.com/caapim/layer7-operator/api/v1alpha1"
	"github.com/caapim/layer7-operator/pkg/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func RedisStateStore(ctx context.Context, params Params) error {
	statestore := params.Instance
	// Retrieve existing secret for Redis
	// this will need to be updated for multi-state store provider support
	if statestore.Spec.Redis.ExistingSecret != "" {
		stateStoreSecret, err := getStateStoreSecret(ctx, statestore.Spec.Redis.ExistingSecret, *statestore, params)
		if err != nil {
			return err
		}
		statestore.Spec.Redis.Username = string(stateStoreSecret.Data["username"])
		statestore.Spec.Redis.MasterPassword = string(stateStoreSecret.Data["masterPassword"])
	}

	c := util.RedisClient(&statestore.Spec.Redis)
	ping := c.Ping(ctx)
	status := statestore.Status

	pong, err := ping.Result()
	if err != nil {
		params.Recorder.Eventf(statestore, "Warning", "ConnectionFailed", "%s in namespace %s", statestore.Name, statestore.Namespace)
		params.Log.V(2).Info("failed to connect to state store", "name", statestore.Name, "namespace", statestore.Namespace, "message", err.Error())
	}

	status.Ready = false

	if pong == "PONG" {
		params.Recorder.Eventf(statestore, "Normal", "ConnectionSuccess", "%s status in namespace %s", statestore.Name, statestore.Namespace)
		status.Ready = true
	}

	if !reflect.DeepEqual(statestore, status) {
		statestore.Status = status
		err = params.Client.Status().Update(ctx, statestore)
		if err != nil {
			params.Log.V(2).Info("failed to update state store status", "name", statestore.Name, "namespace", statestore.Namespace, "message", err.Error())
			return err
		}
		params.Log.V(2).Info("updated state store status", "name", statestore.Name, "namespace", statestore.Namespace)
	}
	return nil
}

func getStateStoreSecret(ctx context.Context, name string, statestore securityv1alpha1.L7StateStore, params Params) (*corev1.Secret, error) {
	statestoreSecret := &corev1.Secret{}

	err := params.Client.Get(ctx, types.NamespacedName{Name: name, Namespace: statestore.Namespace}, statestoreSecret)
	if err != nil {
		return statestoreSecret, err
	}
	return statestoreSecret, nil
}
