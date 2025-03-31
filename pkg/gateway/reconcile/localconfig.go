package reconcile

import (
	"context"
	"fmt"
)

func ClusterProperties(ctx context.Context, params Params) error {
	deleteDbBacked := false
	gateway := params.Instance
	if !gateway.Spec.App.ClusterProperties.Enabled {
		if len(gateway.Status.LastAppliedClusterProperties) == 0 {
			return nil
		}
		if !gateway.Spec.App.Management.Database.Enabled {
			gateway.Status.LastAppliedClusterProperties = []string{}
			if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
				return fmt.Errorf("failed to remove cluster properties status: %w", err)
			}
			return nil
		}
		deleteDbBacked = true
	}

	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithDelete(deleteDbBacked),
		WithBundleType(BundleTypeClusterProp),
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

func ListenPorts(ctx context.Context, params Params) error {
	deleteDbBacked := false
	gateway := params.Instance
	if !gateway.Spec.App.ListenPorts.Custom.Enabled {
		if len(gateway.Status.LastAppliedListenPorts) == 0 {
			return nil
		}
		if !gateway.Spec.App.Management.Database.Enabled {
			gateway.Status.LastAppliedListenPorts = []string{}
			if err := params.Client.Status().Update(ctx, params.Instance); err != nil {
				return fmt.Errorf("failed to remove listen ports status: %w", err)
			}
			return nil
		}
		deleteDbBacked = true
	}

	gwUpdReq, err := NewGwUpdateRequest(
		ctx,
		gateway,
		params,
		WithDelete(deleteDbBacked),
		WithBundleType(BundleTypeListenPort),
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
