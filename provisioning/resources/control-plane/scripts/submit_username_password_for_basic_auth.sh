#!/bin/bash
username=$1
password=$2
config_dir=$3

caddyfile_directory="$config_dir/Caddyfile"

#add username and password to Caddyfile for basic auth
sudo sed -i 's/\(.*basicauth\)\(.*\)/\1 "'$username'" "'$password'" {/' $caddyfile_directory

sudo service caddy_init restart
