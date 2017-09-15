#!/bin/bash

sudo service elasticsearch start
sudo service iglu_server_0.2.0 start
sudo service snowplow_stream_collector start
sudo service snowplow_stream_enrich start
sudo service snowplow_elasticsearch_loader_good start
sudo service snowplow_elasticsearch_loader_bad start
sudo service kibana start
sudo service nginx start
sleep 15

# Send good and bad events
COUNTER=0
while [  $COUNTER -lt 10 ]; do
  curl http://localhost:8080/i?e=pv
  curl http://localhost:8080/i
  let COUNTER=COUNTER+1
done
sleep 60

# Assertions
good_count="$(curl --silent -XGET 'http://localhost:9200/good/good/_count' | python -c 'import json,sys;obj=json.load(sys.stdin);print obj["count"]')"
bad_count="$(curl --silent -XGET 'http://localhost:9200/bad/bad/_count' | python -c 'import json,sys;obj=json.load(sys.stdin);print obj["count"]')"

echo "Event Counts:"
echo " - Good: ${good_count}"
echo " - Bad: ${bad_count}"

org_control_plane_dir="provisioning/resources/control-plane"
test_control_plane_dir="/home/ubuntu/snowplow/control-plane"
$org_control_plane_dir/test/run_tests.sh $org_control_plane_dir $test_control_plane_dir
control_plane_test_res=$?

if [[ "${good_count}" -eq "10" ]] && [[ "${bad_count}" -eq "10" ]] &&
   [[ "${control_plane_test_res}" -eq "0" ]]; then
  exit 0
else
  exit 1
fi
