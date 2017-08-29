#!/bin/bash

orgControlPlaneDir=$1
testControlPlaneDir=$2

$orgControlPlaneDir/test/test_rest_api.sh $orgControlPlaneDir $testControlPlaneDir
testRestApiRes=$?

if [[ "${testRestApiRes}" -eq 0 ]];then
   echo "Rest API test is successful"
else
   echo "Rest API test is not successful"
fi

$orgControlPlaneDir/test/test_scripts.sh $orgControlPlaneDir $testControlPlaneDir
testScriptsRes=$?

if [[ "${testScriptsRes}" -eq 0 ]];then
   echo "Control Plane API scripts test is successful"
else
   echo "Control Plane API scripts test is not successful"
fi

if [[ "${testScriptsRes}" -eq 0 ]] && [[ "${testRestApiRes}" -eq 0 ]];then
   exit 0
else
   exit 1
fi