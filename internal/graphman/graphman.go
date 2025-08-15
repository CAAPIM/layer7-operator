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
package graphman

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// Implode - convert an exploded Graphman directory into a single JSON file.
func Implode(path string) ([]byte, error) {
	bundle, err := implodeBundle(path)
	if err != nil {
		return nil, err
	}

	properties, err := parseBundleProperties(path)
	if err != nil {
		return nil, err
	}

	bundle.Properties = &properties

	bundleBytes, err := json.Marshal(bundle)

	if err != nil {
		return nil, err
	}

	// if len(bundleBytes) <= 40 {
	// 	return nil, errors.New("repository not synced yet")
	// }

	return bundleBytes, nil
}

func RemoveL7PortalApi(username string, password string, target string, apiName string, policyFragmentName string, secretNames []string) ([]byte, error) {
	resp, err := deleteL7PortalApi(context.Background(), gqlClient(username, password, target, ""), []string{apiName}, []string{policyFragmentName}, secretNames)
	if err != nil {
		return nil, err
	}
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	return respBytes, nil
}

func ApplyDynamicBundle(username string, password string, target string, encpass string, bundleBytes []byte) (interface{}, error) {
	bundle := Bundle{}

	err := json.Unmarshal(bundleBytes, &bundle)
	if err != nil {
		return nil, err
	}
	resp, applyErr := installGenericBundle(context.Background(), gqlClient(username, password, target, encpass), &bundle)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	if applyErr != nil {
		bundleApplyErrors, err := CheckDetailedStatus(respBytes)
		if err != nil {
			return nil, err
		}
		if bundleApplyErrors == nil {
			return nil, applyErr
		}
		bundleApplyErrorBytes, err := json.Marshal(bundleApplyErrors)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", string(bundleApplyErrorBytes))
	}

	return resp, nil
}

func DeleteDynamicBundle(username string, password string, target string, encpass string, bundleBytes []byte) (interface{}, error) {
	bundle := Bundle{}

	err := json.Unmarshal(bundleBytes, &bundle)
	if err != nil {
		return nil, err
	}
	resp, applyErr := deleteGenericBundle(context.Background(), gqlClient(username, password, target, encpass), &bundle)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, err
	}

	if applyErr != nil {
		bundleApplyErrors, err := CheckDetailedStatus(respBytes)
		if err != nil {
			return nil, err
		}
		if bundleApplyErrors == nil {
			return nil, applyErr
		}
		bundleApplyErrorBytes, err := json.Marshal(bundleApplyErrors)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("%s", string(bundleApplyErrorBytes))
	}

	return resp, nil
}

func CheckDetailedStatus(respBytes []byte) (*[]BundleApplyError, error) {
	mutationDetailedStatuses := BundleResponseDetailedStatus{}
	bundleApplyErrors := []BundleApplyError{}

	err := json.Unmarshal(respBytes, &mutationDetailedStatuses)
	if err != nil {
		return nil, err
	}

	v := reflect.ValueOf(mutationDetailedStatuses)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		mutationDetailedStatus := MutationDetailedStatus{}
		bundleApplyError := BundleApplyError{
			Entity: typeOfS.Field(i).Name,
		}
		detailedStatusBytes, err := json.Marshal(v.Field(i).Interface())
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(detailedStatusBytes, &mutationDetailedStatus)
		if err != nil {
			return nil, err
		}
		if !reflect.DeepEqual(mutationDetailedStatus, (MutationDetailedStatus{})) {
			for _, ds := range mutationDetailedStatus.DetailedStatus {

				if ds.Status == MutationStatusError {
					bundleApplyError.Error = ds
					bundleApplyErrors = append(bundleApplyErrors, bundleApplyError)
				}
			}
		}
	}

	if len(bundleApplyErrors) > 0 {
		return &bundleApplyErrors, nil
	}

	return nil, nil
}
