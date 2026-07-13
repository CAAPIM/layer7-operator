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
package reconcile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	k8sClient   client.Client
	testEnv     *envtest.Environment
	testScheme  *runtime.Scheme = scheme.Scheme
	ctx         context.Context
	cancel      context.CancelFunc
	instanceUID = uuid.NewUUID()
	logger      = logf.Log.WithName("unit-tests")
)

func TestMain(m *testing.M) {
	ctx, cancel = context.WithCancel(context.TODO())
	defer cancel()
	flag := false
	testEnv = &envtest.Environment{
		CRDDirectoryPaths:     []string{filepath.Join("..", "..", "..", "config", "crd", "bases")},
		ErrorIfCRDPathMissing: true,
		UseExistingCluster:    &flag,
	}
	cfg, err := testEnv.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: envtest unavailable (%v); using fake client for reconcile tests\n", err)
		if err = securityv1.AddToScheme(testScheme); err != nil {
			fmt.Printf("failed to register scheme: %v", err)
			os.Exit(1)
		}
		k8sClient = fake.NewClientBuilder().
			WithScheme(testScheme).
			WithStatusSubresource(&securityv1.Repository{}).
			Build()
		os.Exit(m.Run())
	}
	if err = securityv1.AddToScheme(testScheme); err != nil {
		fmt.Printf("failed to register scheme: %v", err)
		os.Exit(1)
	}
	k8sClient, err = client.New(cfg, client.Options{Scheme: testScheme})
	if err != nil {
		fmt.Printf("failed to setup a Kubernetes client: %v", err)
		os.Exit(1)
	}
	code := m.Run()

	err = testEnv.Stop()
	if err != nil {
		fmt.Printf("failed to stop testEnv: %v", err)
		os.Exit(1)
	}

	os.Exit(code)
}
func newParams() Params {
	params := Params{
		Client: k8sClient,
		Instance: &securityv1.Repository{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Repository",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "default",
				UID:       instanceUID,
			},
			Spec: securityv1.RepositorySpec{
				Enabled: true,
				// Name:     "test",
				Branch:   "testbranch",
				Endpoint: "testing.com",
				Auth: securityv1.RepositoryAuth{
					Username: "testUser",
					Password: "testPassword",
					Token:    "testToken",
					Type:     securityv1.RepositoryAuthTypeBasic,
				},
			},
		},
		Scheme: testScheme,
		Log:    logger,
	}
	params.Instance.Name = "test"
	params.Instance.Namespace = "default"
	params.Client = k8sClient
	return params
}
