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

	if gateway.Spec.App.Management.Database.Enabled || !gateway.Spec.App.Otk.Enabled || gateway.Spec.App.Otk.Type != v1.OtkTypeSingle {
		return nil
	}

	bundle := graphman.Bundle{}

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

	bundleBytes, err := json.Marshal(bundle)
	if err != nil {
		return err
	}

	h := sha1.New()
	h.Write(bundleBytes)
	sha1Sum := fmt.Sprintf("%x", h.Sum(nil))

	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithBundleType(BundleTypeOTKDatabaseMaintenance),
		WithDelete(false),
		WithBundle(bundleBytes),
		WithChecksum(sha1Sum),
		WithBundleName("otk-db-maintenance-tasks"),
		WithPatchAnnotation("security.brcmlabs.com/otk-db-maintenance-tasks"),
		WithCacheEntry(gateway.Name+"-"+string(BundleTypeOTKDatabaseMaintenance)+"-otk-db-maintenance-tasks-"+sha1Sum),
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
