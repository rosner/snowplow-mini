#!/bin/bash

# NOTE: Use `uuidgen` to create new `uid`
iglu_server_super_uid="980ae3ab-3aba-4ffe-a3c2-3b2e24e2ffce"

domain_name=example.com
tls_status="off"

username=USERNAME_PLACEHOLDER
password=PASSWORD_PLACEHOLDER

# DO NOT ALTER BELOW #
sudo service snowplow_mini_control_plane_api start
sleep 2

#add apiKey to iglu-resolver.json for auth in the iglu server
curl -XPOST -d "iglu_server_super_uuid=$iglu_server_super_uid" localhost:10000/add-iglu-server-super-uuid

#add domain name to Caddyfile
curl -XPOST -d "tls_status=$tls_status&domain_name=$domain_name" localhost:10000/add-domain-name

#add username and password to Caddyfile for basic auth
curl -XPOST -d "new_username=$username&new_password=$password" localhost:10000/change-username-and-password
