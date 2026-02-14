/*
* Copyright (c) 2025 Broadcom. All rights reserved.
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
 */
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

		// Register certificate sync job for DMZ and Internal gateways
		// Use SyncIntervalSeconds if specified, otherwise fall back to RuntimeSyncIntervalSeconds or default
		certSyncInterval := otkSyncInterval
		if params.Instance.Spec.App.Otk.SyncIntervalSeconds != 0 {
			certSyncInterval = params.Instance.Spec.App.Otk.SyncIntervalSeconds
		}

		_, err = s.Every(certSyncInterval).Seconds().Tag(params.Instance.Name+"-sync-otk-certificates").Do(syncOtkCertificates, ctx, params)

		if err != nil {
			params.Log.V(2).Info("otk certificate sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		}

		// Register external keys sync job for certificate publishing between DMZ and Internal
		_, err = s.Every(certSyncInterval).Seconds().Tag(params.Instance.Name+"-sync-otk-external-keys").Do(syncOtkExternalKeys, ctx, params)

		if err != nil {
			params.Log.V(2).Info("otk external keys sync job already registered", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
		}
	}
}

func syncOtkExternalKeys(ctx context.Context, params Params) {
	// Sync certificates between DMZ and Internal gateways via ExternalKeys
	// This handles OTK certificate publishing (publishDmzCertToInternal, publishInternalCertToDmz)
	err := ExternalKeys(ctx, params)
	if err != nil {
		params.Log.Error(err, "failed to sync OTK external keys certificates", "name", params.Instance.Name, "namespace", params.Instance.Namespace)
	} else {
		params.Log.V(2).Info("OTK external keys certificates synced", "name", params.Instance.Name, "namespace", params.Instance.Namespace, "interval", params.Instance.Spec.App.Otk.SyncIntervalSeconds)
	}
}

func removeJob(tag string) error {
	err := s.RemoveByTag(tag)
	if err != nil {
		return err
	}
	return nil
}
