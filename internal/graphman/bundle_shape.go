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
* AI assistance has been used to generate some or all contents of this file. That includes, but is not limited to, new code, modifying existing code, stylistic edits.
 */

package graphman

import (
	"encoding/json"
)

// fullGraphmanBundleJSONKeys lists JSON object keys that may appear at the top level of a
// Graphman Bundle document (see Bundle struct json tags). Used to detect loose files that are
// intended as full bundles vs unrelated JSON.
var fullGraphmanBundleJSONKeys = map[string]struct{}{
	"webApiServices":                      {},
	"internalWebApiServices":              {},
	"soapServices":                        {},
	"internalSoapServices":                {},
	"policyFragments":                     {},
	"encassConfigs":                       {},
	"clusterProperties":                   {},
	"jdbcConnections":                     {},
	"trustedCerts":                        {},
	"schemas":                             {},
	"dtds":                                {},
	"fips":                                {},
	"ldaps":                               {},
	"internalGroups":                      {},
	"fipGroups":                           {},
	"internalUsers":                       {},
	"fipUsers":                            {},
	"secrets":                             {},
	"keys":                                {},
	"cassandraConnections":                {},
	"jmsDestinations":                     {},
	"kerberosConfigs":                     {},
	"globalPolicies":                      {},
	"backgroundTaskPolicies":              {},
	"scheduledTasks":                      {},
	"serverModuleFiles":                   {},
	"smConfigs":                           {},
	"activeConnectors":                    {},
	"emailListeners":                      {},
	"listenPorts":                         {},
	"administrativeUserAccountProperties": {},
	"passwordPolicies":                    {},
	"revocationCheckPolicies":             {},
	"logSinks":                            {},
	"httpConfigurations":                  {},
	"customKeyValues":                     {},
	"serviceResolutionConfigs":            {},
	"folders":                             {},
	"federatedIdps":                       {},
	"federatedGroups":                     {},
	"federatedUsers":                      {},
	"internalIdps":                        {},
	"ldapIdps":                            {},
	"simpleLdapIdps":                      {},
	"policyBackedIdps":                    {},
	"policies":                            {},
	"services":                            {},
	"roles":                               {},
	"genericEntities":                     {},
	"auditConfigurations":                 {},
	"policyBackedServices":                {},
	"policyAliases":                       {},
	"serviceAliases":                      {},
	"sampleMessages":                      {},
	"properties":                          {},
}

// IsFullGraphmanBundleJSON reports whether data is a JSON object whose top-level keys suggest a
// full Graphman Bundle document (as opposed to unrelated repo JSON or a standalone alias file).
func IsFullGraphmanBundleJSON(data []byte) bool {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(data, &m); err != nil {
		return false
	}
	if len(m) == 0 {
		return false
	}
	for k := range m {
		if _, ok := fullGraphmanBundleJSONKeys[k]; ok {
			return true
		}
	}
	return false
}
