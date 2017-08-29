#!/bin/bash

#### In this script, endpoints of the Control Plane API is checked
#### whether they are working as expected or not.

orgControlPlaneDir=$1
testControlPlaneDir=$2
testDir="$testControlPlaneDir/test"
testEnv="$testDir/testEnv"
testInit="$testDir/testInit"
testConfigDir="$testEnv/testConfig"

#copy original control plane directory to testing directory
sudo cp -r $orgControlPlaneDir $testControlPlaneDir

sudo cp $testInit/snowplow_mini_control_plane_api_test_init /etc/init.d/snowplow_mini_control_plane_api
sudo /etc/init.d/snowplow_mini_control_plane_api restart $testDir
sleep 2

## restart SP services test
stream_enrich_pid_file=/var/run/snowplow_stream_enrich.pid
stream_collector_pid_file=/var/run/snowplow_stream_collector.pid
sink_bad_pid_file=/var/run/snowplow_elasticsearch_loader_bad.pid
sink_good_pid_file=/var/run/snowplow_elasticsearch_loader_good.pid

stream_enrich_pid_old="$(cat "${stream_enrich_pid_file}")"
stream_collector_pid_old="$(cat "${stream_collector_pid_file}")"
sink_bad_pid_old="$(cat "${sink_bad_pid_file}")"
sink_good_pid_old="$(cat "${sink_good_pid_file}")"

req_result=$(curl -s -o /dev/null -w "%{http_code}" -XPUT http://localhost:10000/restart-services)

stream_enrich_pid_new="$(cat "${stream_enrich_pid_file}")"
stream_collector_pid_new="$(cat "${stream_collector_pid_file}")"
sink_bad_pid_new="$(cat "${sink_bad_pid_file}")"
sink_good_pid_new="$(cat "${sink_good_pid_file}")"

if [[ "${req_result}" -eq 200 ]] &&
   [[ "${stream_enrich_pid_old}" -ne "${stream_enrich_pid_new}" ]] &&
   [[ "${stream_collector_pid_old}" -ne "${stream_collector_pid_new}" ]] &&
   [[ "${sink_bad_pid_old}" -ne "${sink_bad_pid_new}" ]] &&
   [[ "${sink_good_pid_old}" -ne "${sink_good_pid_new}" ]]; then

  echo "Restarting SP services is working correctly."
else
  echo "Restarting SP services is not working correctly."
  exit 1
fi

## upload enrichment test: success
upload_enrichments_result=$(curl -s -o /dev/null -w "%{http_code}"  --header "Content-Type: multipart/form-data" -F "enrichmentjson=@$testEnv/orgEnrichments/enrich.json" localhost:10000/upload-enrichments)
enrichment_diff=$(diff $testEnv/testEnrichments/enrich.json $testEnv/orgEnrichments/enrich.json)
sleep 2

if [[ "${upload_enrichments_result}" -eq 200 ]] && [[ "${enrichment_diff}" == "" ]];then
  echo "Uploading enrichment success test returned as expected."
else
  echo "Uploading enrichment success test did not return as expected ."
  exit 1
fi

## upload enrichment test: invalid JSON fail
upload_enrichments_result=$(curl -s -o /dev/null -w "%{http_code}"  --header "Content-Type: multipart/form-data" -F "enrichmentjson=@$testEnv/orgEnrichments/invalid_enrichment.json" localhost:10000/upload-enrichments)
sleep 2

if [[ "${upload_enrichments_result}" -eq 400 ]] ;then
  echo "Uploading enrichment fail test returned as expected."
else
  echo "Uploading enrichment fail test did return as expected."
  exit 1
fi

sudo cp $testInit/snowplow_mini_control_plane_api_original_init /etc/init.d/snowplow_mini_control_plane_api
sudo /etc/init.d/snowplow_mini_control_plane_api restart

#remove test control plane directory after testing is done
sudo rm -rf $testControlPlaneDir

exit 0
