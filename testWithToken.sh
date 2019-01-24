#!/bin/bash

if [[ $# -lt 2 ]]; then
  cat - <<EOM
Use: bash test.sh ssh-url commons-hostname commons-login
EOM
  exit 0
fi
SSH_URL="$1"
shift
COMMONS_HOST="$1"
shift
COMMONS_LOGIN="$1"
shift
export AUTHPROXY_FENCE_URL="https://${COMMONS_HOST}/user"
export AUTHPROXY_TEST_TOKEN="$(ssh "$SSH_URL" 'set -i; source ~/.bashrc; g3kubectl exec $(gen3 pod fence) -- fence-create token-create --scopes openid,user,fence,data,credentials,google_service_account --type access_token --exp 3600 --username '${COMMONS_LOGIN})"

go test -v
