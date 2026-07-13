/*
* Copyright (c) 2026 Broadcom. All rights reserved.
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

package util

import "path/filepath"

// Repository-scoped /tmp layout: {name}-{namespace} must stay in sync between
// the repository controller and gateway controller (cluster-wide operator, namespaced CRs).

// RepoCacheDir returns the directory for per-commit bundle map JSON (repo-cache).
func RepoCacheDir(name, namespace string) string {
	return filepath.Join("/tmp/repo-cache", name+"-"+namespace)
}

// StateStoreCacheDir returns the directory for statestore-backed repository temp files.
func StateStoreCacheDir(name, namespace string) string {
	return filepath.Join("/tmp/statestore", name+"-"+namespace)
}

// GatewayBundleWorkDir returns the per-gateway-repository working directory under /tmp/bundles.
func GatewayBundleWorkDir(name, namespace string) string {
	return filepath.Join("/tmp/bundles", name+"-"+namespace)
}
