package reconcile

import (
	"context"
	"time"

	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-co-op/gocron"
)

var s = gocron.NewScheduler(time.Local)

var syncCache = util.NewSyncCache(3 * time.Second)

func Jobs(ctx context.Context, params Params) error {

	registerJobs(ctx, params)

	for _, j := range s.Jobs() {
		for _, t := range j.Tags() {
			params.Log.V(2).Info("starting job", "job", t, "namespace", params.Instance.Namespace)
			err := s.RunByTag(t)
			if err != nil {
				params.Log.V(2).Info("no job with given tag", "job", t, "namespace", params.Instance.Namespace)
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
	portalSyncInterval := 10

	if params.Instance.Spec.SyncIntervalSeconds != 0 {
		portalSyncInterval = params.Instance.Spec.SyncIntervalSeconds
	}

	_, err := s.Every(portalSyncInterval).Seconds().Tag(params.Instance.Name+"-sync-portal").Do(syncPortal, ctx, params)

	if err != nil {
		params.Log.V(2).Info("portal sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

	_, err = s.Every(portalSyncInterval).Seconds().Tag(params.Instance.Name+"-sync-portal-apis").Do(syncPortalApis, ctx, params)

	if err != nil {
		params.Log.V(2).Info("portal api sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

}

func removeJob(tag string) error {
	err := s.RemoveByTag(tag)
	if err != nil {
		return err
	}
	return nil
}
