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
	"crypto/sha1"
	"encoding/json"
	"fmt"

	v1 "github.com/caapim/layer7-operator/api/v1"
	"github.com/caapim/layer7-operator/internal/graphman"
)

var scheduledTasks = map[string]string{
	"OTK Database Maintenance - Client":                "0 */31 * * * ?",
	"OTK Database Maintenance - id_token":              "0 */29 * * * ?",
	"OTK Database Maintenance - sessions":              "0 */7 * * * ?",
	"OTK Database Maintenance - token":                 "0 */5 * * * ?",
	"OTK Database Maintenance - token revocation list": "0 0 */1 * * ?",
	"OTK Database Maintenance - miscellaneous tokens":  "0 */5 * * * ?",
}

func OTKDatabaseMaintenanceTasks(ctx context.Context, params Params) error {
	gateway := params.Instance

	if gateway.Spec.App.Management.Database.Enabled || !gateway.Spec.App.Otk.Enabled || gateway.Spec.App.Otk.Type != v1.OtkTypeSingle || gateway.Spec.App.Otk.Database.Type == v1.OtkDatabaseTypeCassandra {
		return nil
	}

	bundle := graphman.Bundle{}
	if gateway.Spec.App.Otk.Database.Type != v1.OtkDatabaseTypeCassandra {
		for name, schedule := range scheduledTasks {
			bundle.ScheduledTasks = append(bundle.ScheduledTasks, &graphman.ScheduledTaskInput{
				Name:                name,
				PolicyName:          name,
				JobType:             graphman.JobTypeRecurring,
				CronExpression:      schedule,
				ExecuteOnSingleNode: true,
				ExecuteOnCreation:   false,
				Status:              graphman.JobStatusScheduled,
			})
		}
	}

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return err
	}

	h := sha1.New()
	h.Write([]byte(gateway.Spec.App.Otk.InitContainerImage))
	checksum := fmt.Sprintf("%x", h.Sum(nil))

	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithBundleType(BundleTypeOTKDatabaseMaintenance),
		WithDelete(false),
		WithBundle(bundleBytes),
		WithChecksum(checksum),
		WithBundleName("otk-db-maintenance-tasks"),
		WithPatchAnnotation("security.brcmlabs.com/otk-db-maintenance-tasks"),
		WithCacheEntry(gateway.Name+"-"+string(BundleTypeOTKDatabaseMaintenance)+"-otk-db-maintenance-tasks-"+checksum),
	)

	if err != nil {
		return err
	}

	err = SyncGateway(ctx, params, *gwUpdReq)
	if err != nil {
		return err
	}

	return nil
}
