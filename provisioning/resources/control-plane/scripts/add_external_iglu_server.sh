#!/bin/bash
iglu_server_uri=$1
iglu_server_api_key=$2
config_dir=$3
control_api_scripts_dir=$4

iglu_resolver_config_dir="$config_dir/iglu-resolver.json"
template_iglu_server="$control_api_scripts_dir/template_iglu_server"

##first change uri and apikey in the template_iglu_server file
sudo sed -i 's/\(.*"uri":\)\(.*\)/\1 "'$iglu_server_uri'",/' $template_iglu_server
sudo sed -i 's/\(.*"apikey":\)\(.*\)/\1 "'$iglu_server_api_key'"/' $template_iglu_server
##secondly write content in the template_iglu_server to iglu_resolver.json
sudo sed -i -E '/.*"repositories":.*/r '$template_iglu_server'' $iglu_resolver_config_dir
