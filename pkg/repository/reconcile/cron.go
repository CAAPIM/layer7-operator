package reconcile

import (
	"context"
	"time"

	v1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/pkg/util"
	"github.com/go-co-op/gocron"
)

var s = gocron.NewScheduler(time.Local)

var syncCache = util.NewSyncCache(3 * time.Second)

func ScheduledJobs(ctx context.Context, params Params) error {

	registerJobs(ctx, params)

	for _, j := range s.Jobs() {
		for _, t := range j.Tags() {
			if t == params.Instance.Name+"-sync-repository" {
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
	repoSyncInterval := 10

	if params.Instance.Spec.RepositorySyncConfig.IntervalSeconds != 0 {
		repoSyncInterval = params.Instance.Spec.RepositorySyncConfig.IntervalSeconds
	}

	if params.Instance.Spec.Type != v1.RepositoryTypeLocal {
		_, err := s.Every(repoSyncInterval).Seconds().Tag(params.Instance.Name+"-"+params.Instance.Namespace+"-sync-repository").Do(syncRepository, ctx, params)
		if err != nil {
			params.Log.V(2).Info("repository sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		}
	}
}
