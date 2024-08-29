#!/bin/sh
# need to wait for a while until the connection is ready
sleep 5;
# /scripts is defined at docker-compose.yml
mysql -h mysql --port 3306 -u root < /scripts/initialize.sql;