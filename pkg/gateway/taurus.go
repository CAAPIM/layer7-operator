package gateway

// func CreateTaurusTest(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, commitId string) error {
// 	found := &batchv1.Job{}

// 	err := r.Get(ctx, types.NamespacedName{Name: "taurus-" + commitId, Namespace: m.Namespace}, found)
// 	if err != nil && k8serrors.IsNotFound(err) {
// 		r.Log.Info("Creating Job", "Name", "taurus-"+commitId, "Namespace", m.Namespace)
// 		job := taurusJob(m, commitId)
// 		ctrl.SetControllerReference(m, job, r.Scheme)
// 		err = r.Create(ctx, job)
// 		if err != nil {
// 			r.Log.Error(err, "Failed creating Job", "Name", "taurus-"+commitId, "Namespace", m.Namespace)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func CreateTaurusSecret(r *GatewayReconciler, ctx context.Context, gw *securityv1.Gateway, commitId string, testData []byte) error {
// 	found := &corev1.Secret{}
// 	err := r.Get(ctx, types.NamespacedName{Name: "taurus-" + commitId, Namespace: m.Namespace}, found)
// 	if err != nil && k8serrors.IsNotFound(err) {
// 		r.Log.Info("Creating Secret", "Name", "taurus-"+commitId, "Namespace", m.Namespace)
// 		secret := taurusSecret(m, commitId, testData)
// 		ctrl.SetControllerReference(m, secret, r.Scheme)
// 		err = r.Create(ctx, secret)
// 		if err != nil {
// 			r.Log.Error(err, "Failed creating Secret", "Name", "taurus-"+commitId, "Namespace", m.Namespace)
// 			return err
// 		}
// 	}
// 	return nil
// }

// func taurusJob(gw *securityv1.Gateway, commitId string) *batchv1.Job {
// 	var image string = "harbor.sutraone.com/library/taurus:1.15.4"

// 	ls := util.LabelsForGateway("taurus")
// 	var terminationTtl = int32(3600)
// 	job := &batchv1.Job{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "taurus-" + commitId,
// 			Namespace: m.Namespace,
// 		},
// 		TypeMeta: metav1.TypeMeta{
// 			APIVersion: "v1",
// 			Kind:       "Job",
// 		},

// 		Spec: batchv1.JobSpec{
// 			TTLSecondsAfterFinished: &terminationTtl,
// 			Template: corev1.PodTemplateSpec{
// 				ObjectMeta: metav1.ObjectMeta{
// 					Labels: ls,
// 				},

// 				Spec: corev1.PodSpec{
// 					RestartPolicy: corev1.RestartPolicyNever,

// 					Containers: []corev1.Container{{

// 						Image:      image,
// 						Name:       "taurus",
// 						Args:       []string{"test.yaml"},
// 						WorkingDir: "/bzt-configs",
// 						VolumeMounts: []corev1.VolumeMount{{
// 							Name:      "taurus-" + commitId,
// 							MountPath: "/bzt-configs/test.yaml",
// 							SubPath:   "test.yaml",
// 						}},
// 					}},
// 					Volumes: []corev1.Volume{{
// 						Name: "taurus-" + commitId,
// 						VolumeSource: corev1.VolumeSource{
// 							Secret: &corev1.SecretVolumeSource{
// 								SecretName: "taurus-" + commitId,
// 								Items: []corev1.KeyToPath{{
// 									Path: "test.yaml",
// 									Key:  "test.yaml"},
// 								},
// 							},
// 						},
// 					}},
// 				},
// 			},
// 		},
// 	}

// 	return job
// }

// // GatewaySecret returns a Gateway Secret Definition
// func taurusSecret(gw *securityv1.Gateway, commitId string, testData []byte) *corev1.Secret {
// 	secret := &corev1.Secret{
// 		ObjectMeta: metav1.ObjectMeta{
// 			Name:      "taurus-" + commitId,
// 			Namespace: m.Namespace,
// 		},
// 		TypeMeta: metav1.TypeMeta{
// 			APIVersion: "v1",
// 			Kind:       "Secret",
// 		},
// 		Type: corev1.SecretTypeOpaque,
// 		Data: map[string][]byte{
// 			"test.yaml": testData,
// 		},
// 	}

// 	return secret
// }
