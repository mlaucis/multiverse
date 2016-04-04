#!/usr/bin/env sh

set -ex

# add prometheus source

curl -vvv \
  -X POST \
  -u admin:admin \
  -H 'Content-Type: application/json;charset=UTF-8' \
  --data-binary '{"name":"prometheus", "type":"prometheus","url":"http://localhost:9090","access":"proxy","isDefault":true}' \
  'http://127.0.0.1:3000/api/datasources'

# change admin user login
PASSWORD=$(date +%s | sha256sum | base64 | head -c 32 ; echo)

curl -vvvv \
  -X PUT \
  -u admin:admin \
  -H 'Content-Type: application/json;charset=UTF-8' \
  --data-binary "{\"oldPassword\": \"admin\", \"newPassword\": \"$PASSWORD\", \"confirmNew\": \"$PASSWORD\"}" \
  'http://127.0.0.1:3000/api/user/password'

# /etc/grafana/grafana.ini

echo "
[auth.basic]
enabled = false

[auth.google]
enabled = true
client_id = $GOOGLE_CLIENT_ID
client_secret = $GOOGLE_CLIENT_SECRET
scopes = https://www.googleapis.com/auth/userinfo.profile https://www.googleapis.com/auth/userinfo.email
auth_url = https://accounts.google.com/o/oauth2/auth
token_url = https://accounts.google.com/o/oauth2/token
allowed_domains = tapglue.com
allow_sign_up = true

[dashboards.json]
enabled = true
path = /var/lib/grafana/dashboards

[security]
admin_user = admin
admin_password = $PASSWORD

[server]
root_url = https://monitoring-$ENV-$AWS_REGION.tapglue.com

[users]
auto_assign_org = true
auto_assign_org_role = Editor
" | sudo tee /etc/grafana/grafana.ini > /dev/null

# /var/lib/grafana/dashboards/ops.json

sudo chown grafana:grafana /etc/grafana/grafana.ini

sudo mkdir -p /var/lib/grafana/dashboards
sudo mv /tmp/dashboard-ops.json /var/lib/grafana/dashboards/
