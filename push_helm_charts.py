# Package the layer7-operator helm chart from this workspace and push it to
# Artifactory as an OCI artifact.
#
# Pre-requisite:
# - Artifactory credentials set in environment: ARTIFACTORY_CREDS_USR, ARTIFACTORY_CREDS_PSW
#
# command examples:
# python3 push_helm_charts.py
# python3 push_helm_charts.py --release
# python3 push_helm_charts.py --chart-dir charts/layer7-operator

import argparse
import os
import subprocess
import sys

parser = argparse.ArgumentParser(description='Package and push the layer7-operator helm chart to artifactory')
parser.add_argument('--chart-dir', default='charts/layer7-operator', help='path to the chart to package')
parser.add_argument('--release', action='store_true', help='flag to push to release repo instead of dev')

args = parser.parse_args()
username = os.getenv('ARTIFACTORY_CREDS_USR')
password = os.getenv('ARTIFACTORY_CREDS_PSW')
if not username or not password:
    sys.exit("please set env for ARTIFACTORY_CREDS_USR and ARTIFACTORY_CREDS_PSW")
helm_stage = "release" if args.release else "dev"
helm_repo = f"apim-docker-{helm_stage}-local.usw1.packages.broadcom.com"
subprocess.run(['docker', 'login', helm_repo, '-u', username, '-p', password], check=True, text=True)

def package_chart(chart_dir):
    result = subprocess.run(['helm', 'package', chart_dir], check=True, text=True, capture_output=True)
    print(result.stdout)
    # `helm package` prints "Successfully packaged chart and saved it to: <path>"
    return result.stdout.strip().rsplit(": ", 1)[-1]

def main():
    print(f"working on operator chart: {args.chart_dir}")
    operator_chart = package_chart(args.chart_dir)
    subprocess.run(['helm', 'push', operator_chart, f"oci://{helm_repo}"], check=True, text=True)

    subprocess.run(['docker', 'logout', helm_repo], check=True, text=True)

if __name__ == "__main__":
    main()
