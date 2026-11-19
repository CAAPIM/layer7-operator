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

// Repository-scoped /tmp layout must stay in sync between the repository
// controller and gateway controller (cluster-wide operator, namespaced CRs).
//
// Every per-repository directory nests namespace and name as separate path
// segments ({base}/{namespace}/{name}) rather than joining them with '-'.
// '-' is legal inside both Kubernetes names and namespaces, so a hyphen-joined
// key is not injective ("apis"/"team-a" and "apis-team"/"a" would both map to
// "apis-team-a") and lets one tenant collide with another tenant's cache.
// '/' is invalid in Kubernetes names and namespaces, so nesting is injective.

// These are vars (not consts) only so tests can redirect them; treat them as
// constants in production code.
var (
	repoCacheBase          = "/tmp/repo-cache"
	stateStoreCacheBase    = "/tmp/statestore"
	gatewayBundleWorkBase  = "/tmp/bundles"
	repositoryCheckoutBase = "/tmp/repositories"
)

// emptyRefSegment is the path segment used for a repository with no branch or
// tag. "@" on its own is not a valid git ref name (git-check-ref-format), so it
// cannot collide with a real ref, and it keeps the checkout a leaf directory
// rather than the parent of other refs' checkouts.
const emptyRefSegment = "@"

// RepoCacheDir returns the directory for per-commit bundle map JSON (repo-cache).
func RepoCacheDir(name, namespace string) string {
	return filepath.Join(repoCacheBase, namespace, name)
}

// StateStoreCacheDir returns the directory for statestore-backed repository temp files.
func StateStoreCacheDir(name, namespace string) string {
	return filepath.Join(stateStoreCacheBase, namespace, name)
}

// GatewayBundleWorkDir returns the per-gateway-repository working directory under /tmp/bundles.
func GatewayBundleWorkDir(name, namespace string) string {
	return filepath.Join(gatewayBundleWorkBase, namespace, name)
}

// RepoArtifactPath returns the path of a downloaded HTTP repository artifact, or the
// folder it is extracted to, alongside the git checkouts for the same Repository.
func RepoArtifactPath(name, namespace, fileName string) string {
	return filepath.Join(repositoryCheckoutBase, namespace, name, fileName)
}

// RepoCheckoutDir returns the on-disk git working directory for a repository
// checkout of ref (branch or tag): {base}/{namespace}/{name}/{SafeRef(ref)}.
// ref must already have passed ValidateRef. Every site that derives the
// checkout path (CloneRepository, the repository finalizer, storage-secret and
// statestore sync) must call this so they all resolve to the same directory.
func RepoCheckoutDir(name, namespace, ref string) string {
	refSegment := SafeRef(ref)
	// SafeRef flattens separators, so "." and ".." are the only values that
	// could resolve outside a single child directory; ValidateRef already
	// rejects "..", this keeps the helper contained even if called without it.
	if refSegment == "" || refSegment == "." || refSegment == ".." {
		refSegment = emptyRefSegment
	}
	return filepath.Join(repositoryCheckoutBase, namespace, name, refSegment)
}
