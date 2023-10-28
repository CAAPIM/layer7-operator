package reconcile

import (
	"context"
	"time"

	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-co-op/gocron"
)

var s = gocron.NewScheduler(time.Local)

var syncCache = util.NewSyncCache(3 * time.Second)

func ScheduledJobs(ctx context.Context, params Params) error {

	registerJobs(ctx, params)

	////// TODO - check run count before trying to start a job
	for _, j := range s.Jobs() {
		for _, t := range j.Tags() {
			if !j.IsRunning() && j.RunCount() == 0 {
				params.Log.V(2).Info("starting job", "job", t, "name", params.Instance.Name, "namespace", params.Instance.Namespace, "runCount", j.RunCount())
				err := s.RunByTag(t)
				if err != nil {
					params.Log.V(2).Info("no job with given tag", "job", t, "name", params.Instance.Name, "namespace", params.Instance.Namespace)
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
	repoSyncInterval := 5

	if params.Instance.Spec.RepositorySyncConfig.IntervalSeconds != 0 {
		repoSyncInterval = params.Instance.Spec.RepositorySyncConfig.IntervalSeconds
	}

	_, err := s.Every(repoSyncInterval).Seconds().Tag(params.Instance.Name+"-sync-repository").Do(syncRepository, ctx, params)

	if err != nil {
		params.Log.V(2).Info("repository sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	}

}
