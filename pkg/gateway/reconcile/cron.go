package reconcile

import (
	"context"
	"time"

	securityv1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-co-op/gocron"
)

var s = gocron.NewScheduler(time.Local)
var syncCache = util.NewSyncCache(3 * time.Second)

func ScheduledJobs(ctx context.Context, params Params) error {

	registerJobs(ctx, params)

	for _, j := range s.Jobs() {
		for _, t := range j.Tags() {
			if !j.IsRunning() {
				params.Log.V(2).Info("starting job", "job", t, "namespace", params.Instance.Namespace)
				err := s.RunByTag(t)
				if err != nil {
					params.Log.V(2).Info("no job with given tag", "job", t, "namespace", params.Instance.Namespace)
				}
			}
		}
	}
	if !s.IsRunning() {
		s.StartAsync()
	}
	return nil
}

func registerJobs(ctx context.Context, params Params) {
	s.TagsUnique()
	otkSyncInterval := 10

	if params.Instance.Spec.App.Otk.RuntimeSyncIntervalSeconds != 0 {
		otkSyncInterval = params.Instance.Spec.App.Otk.RuntimeSyncIntervalSeconds
	}

	if params.Instance.Spec.App.Otk.Enabled && (params.Instance.Spec.App.Otk.Type == securityv1.OtkTypeDMZ || params.Instance.Spec.App.Otk.Type == securityv1.OtkTypeInternal) {
		_, err := s.Every(otkSyncInterval).Seconds().Tag(params.Instance.Name+"-sync-otk-policies").Do(syncOtkPolicies, ctx, params)

		if err != nil {
			params.Log.V(2).Info("otk policy sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		}
		if params.Instance.Spec.App.Otk.Type == securityv1.OtkTypeDMZ || params.Instance.Spec.App.Otk.Type == securityv1.OtkTypeInternal {
			_, err = s.Every(otkSyncInterval).Seconds().Tag(params.Instance.Name+"-"+params.Instance.Namespace+"-sync-otk-certificates").Do(syncOtkCertificates, ctx, params)
			if err != nil {
				params.Log.V(2).Info("otk certificate sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
			}
			_, err = s.Every(otkSyncInterval).Seconds().Tag(params.Instance.Name+"-"+params.Instance.Namespace+"-sync-otk-certificate-secret").Do(manageCertificateSecrets, ctx, params)
			if err != nil {
				params.Log.V(2).Info("otk certificate secret sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
			}
		}
	}
}

func removeJob(tag string) error {
	err := s.RemoveByTag(tag)
	if err != nil {
		return err
	}
	return nil
}
