#!/bin/bash
iglu_server_super_uid=$1
config_dir=$2

iglu_resolver_config_dir="$config_dir/iglu-resolver.json"

##there will be only one local iglu server, therefore apikey will be added for this iglu server 
n=$(awk '/localhost/{print NR}' $iglu_resolver_config_dir) ## find number of the line which contains "localhost"
n=$((n+1)) ##add one to line number which contains 'localhost' keyword, this line will contain apikey for local iglu server
sudo sed -i ''$n's/\(.*"apikey":\)\(.*\)/\1 "'$iglu_server_super_uid'"/' $iglu_resolver_config_dir ##write apikey for local iglu server

#write super apikey to db
export PGPASSWORD=snowplow
delete_all_apikeys="DELETE FROM apikeys"
iglu_server_setup="INSERT INTO apikeys (uid, vendor_prefix, permission, createdat) VALUES ('${iglu_server_super_uid}','*','super',current_timestamp);"
psql --host=localhost --port=5432 --username=snowplow --dbname=iglu -c "${delete_all_apikeys}"
psql --host=localhost --port=5432 --username=snowplow --dbname=iglu -c "${iglu_server_setup}"
