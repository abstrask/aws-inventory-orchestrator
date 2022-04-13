#!/bin/bash
set -eu
varsfile="./k8s/vars.env"

cat << EOF > $varsfile
BILLING_ACCOUNT_ID=$BILLING_ACCOUNT_ID
INVENTORY_ORCHESTRATOR_ROLE=$INVENTORY_ORCHESTRATOR_ROLE
RUNNER_ACCOUNT_ID=$RUNNER_ACCOUNT_ID
INVENTORY_RUNNER_ROLE=$INVENTORY_RUNNER_ROLE
BUCKET_NAME=$BUCKET_NAME
CRON_SCHEDULE=$CRON_SCHEDULE
INCLUDE_ACCOUNTS=$INCLUDE_ACCOUNTS
INVENTORY_ROLE=$INVENTORY_ROLE
OIDC_PROVIDER=$OIDC_PROVIDER
EOF