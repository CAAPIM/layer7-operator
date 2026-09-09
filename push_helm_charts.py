# Package the layer7-operator helm chart from this workspace and push it to
# Artifactory as an OCI artifact - but only if it's actually a new version.
#
# Pre-requisite:
# - Artifactory credentials set in environment: ARTIFACTORY_CREDS_USR, ARTIFACTORY_CREDS_PSW
# - `ct` (chart-testing) on PATH if --check-changed is used
#
# command examples:
# python3 push_helm_charts.py
# python3 push_helm_charts.py --release
# python3 push_helm_charts.py --check-changed main
# python3 push_helm_charts.py --chart-dir charts/layer7-operator

import argparse
import os
import subprocess
import sys
import tempfile

parser = argparse.ArgumentParser(description='Package and push the layer7-operator helm chart to artifactory, if it is a new version')
parser.add_argument('--chart-dir', default='charts/layer7-operator', help='path to the chart to package')
parser.add_argument('--release', action='store_true', help='flag to push to release repo instead of dev')
parser.add_argument('--check-changed', metavar='TARGET_BRANCH', default=None,
                     help='skip entirely if `ct list-changed --target-branch TARGET_BRANCH` does not report '
                          '--chart-dir as changed (cheap short-circuit for PR builds; meaningless when already '
                          'on the target branch itself, so omit it there)')

args = parser.parse_args()
username = os.getenv('ARTIFACTORY_CREDS_USR')
password = os.getenv('ARTIFACTORY_CREDS_PSW')
if not username or not password:
    sys.exit("please set env for ARTIFACTORY_CREDS_USR and ARTIFACTORY_CREDS_PSW")
helm_stage = "release" if args.release else "dev"
helm_repo = f"apim-docker-{helm_stage}-local.usw1.packages.broadcom.com"

def chart_changed(chart_dir, target_branch):
    result = subprocess.run(['ct', 'list-changed', '--target-branch', target_branch],
                             check=True, text=True, capture_output=True)
    changed = [line.strip() for line in result.stdout.splitlines()]
    return chart_dir.rstrip('/') in changed

def chart_metadata(chart_dir):
    name = version = None
    with open(os.path.join(chart_dir, "Chart.yaml")) as f:
        for line in f:
            if line.startswith("name:"):
                name = line.split(":", 1)[1].strip()
            elif line.startswith("version:"):
                version = line.split(":", 1)[1].strip()
    if not name or not version:
        sys.exit(f"could not read name/version from {chart_dir}/Chart.yaml")
    return name, version

def version_exists(chart_name, version):
    # Mirrors chart-releaser's own "don't recreate an existing release" check,
    # just against the Artifactory OCI repo instead of GitHub Releases.
    with tempfile.TemporaryDirectory() as tmpdir:
        result = subprocess.run(
            ['helm', 'pull', f"oci://{helm_repo}/{chart_name}", '--version', version, '-d', tmpdir],
            text=True, capture_output=True)
        return result.returncode == 0

def package_chart(chart_dir):
    result = subprocess.run(['helm', 'package', chart_dir], check=True, text=True, capture_output=True)
    print(result.stdout)
    # `helm package` prints "Successfully packaged chart and saved it to: <path>"
    return result.stdout.strip().rsplit(": ", 1)[-1]

def main():
    if args.check_changed and not chart_changed(args.chart_dir, args.check_changed):
        print(f"no changes under {args.chart_dir} vs {args.check_changed} - skipping publish")
        return

    chart_name, version = chart_metadata(args.chart_dir)
    subprocess.run(['docker', 'login', helm_repo, '-u', username, '-p', password], check=True, text=True)
    try:
        if version_exists(chart_name, version):
            print(f"{chart_name} {version} already published to {helm_repo} - skipping")
            return
        print(f"working on {chart_name} {version} from {args.chart_dir}")
        operator_chart = package_chart(args.chart_dir)
        subprocess.run(['helm', 'push', operator_chart, f"oci://{helm_repo}"], check=True, text=True)
    finally:
        subprocess.run(['docker', 'logout', helm_repo], check=True, text=True)

if __name__ == "__main__":
    main()
