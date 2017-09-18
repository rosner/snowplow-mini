#!/bin/bash
tls_cond=$1
domain_name=$2
config_dir=$3

caddyfile_directory="$config_dir/Caddyfile"

#add domain name to Caddyfile
inserted_line=""
sudo sed -i '1,2d' $caddyfile_directory #delete first two line of the default Caddyfile 
if [[ "${tls_cond}" == "on" ]]; then
  inserted_line="$domain_name *:80 {\n  tls example@example.com\n"
else
  inserted_line="*:80 {\n  tls off\n"
fi
sudo sed -i "1s/^/${inserted_line}/" $caddyfile_directory 

sudo service caddy_init restart
