#!/bin/bash

#### In this script, Bash scripts of the Control Plane API is checked
#### whether they are working as expected or not.

orgControlPlaneDir=$1
testControlPlaneDir=$2
testDir=$testControlPlaneDir/test
testEnv="$testDir/testEnv"
testInit="$testDir/testInit"
testConfigDir="$testEnv/testConfig"
scripts="$testControlPlaneDir/scripts"

#copy original control plane directory to testing directory
sudo cp -r $orgControlPlaneDir $testControlPlaneDir

restartServicesScript="restart_SP_services.sh"
addExternalIgluServerScript="add_external_iglu_server.sh"
addIgluSuperUuidScript="add_iglu_server_super_uuid.sh"
changeUsernameAndPasswordScript="submit_username_password_for_basic_auth.sh"
addDomainNameScript="write_domain_name_to_caddyfile.sh"


## restart SP services test
stream_enrich_pid_file=/var/run/snowplow_stream_enrich.pid
stream_collector_pid_file=/var/run/snowplow_stream_collector.pid
sink_bad_pid_file=/var/run/snowplow_elasticsearch_loader_bad.pid
sink_good_pid_file=/var/run/snowplow_elasticsearch_loader_good.pid

stream_enrich_pid_old="$(cat "${stream_enrich_pid_file}")"
stream_collector_pid_old="$(cat "${stream_collector_pid_file}")"
sink_bad_pid_old="$(cat "${sink_bad_pid_file}")"
sink_good_pid_old="$(cat "${sink_good_pid_file}")"

sudo $scripts/$restartServicesScript >> /dev/null
res=$?

stream_enrich_pid_new="$(cat "${stream_enrich_pid_file}")"
stream_collector_pid_new="$(cat "${stream_collector_pid_file}")"
sink_bad_pid_new="$(cat "${sink_bad_pid_file}")"
sink_good_pid_new="$(cat "${sink_good_pid_file}")"

if [[ "${res}" -eq 0 ]] &&
   [[ "${stream_enrich_pid_old}" -ne "${stream_enrich_pid_new}" ]] &&
   [[ "${stream_collector_pid_old}" -ne "${stream_collector_pid_new}" ]] &&
   [[ "${sink_bad_pid_old}" -ne "${sink_bad_pid_new}" ]] &&
   [[ "${sink_good_pid_old}" -ne "${sink_good_pid_new}" ]]; then

  echo "Restarting SP services script is working correctly."
else
  echo "Restarting SP services script is not working correctly."
  exit 1
fi

## add external iglu server test
sudo cp $testEnv/orgConfig/iglu-resolver.json $testConfigDir/.
external_test_uuid=$(uuidgen)
iglu_server_uri="testigluserveruri.com"

sudo $scripts/$addExternalIgluServerScript $iglu_server_uri $external_test_uuid $testConfigDir $scripts >> /dev/null
res=$?

written_apikey=$(diff $testConfigDir/iglu-resolver.json $testEnv/expectedConfig/iglu-resolver-external-iglu.json | grep $external_test_uuid)

if [[ "${res}" -eq 0 ]] && [[ "${written_apikey}" != "" ]];then
  echo "Adding external Iglu Server script is working correctly."
else
  echo "Adding external Iglu Server script is not working correctly."
  exit 1
fi

#remove test control plane directory after testing is done
sudo rm -rf $testControlPlaneDir

exit 0
