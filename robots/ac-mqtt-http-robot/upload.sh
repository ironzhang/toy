#!/bin/bash

cp ./robot.so ~/ablecloud.cn/src/deploys/ac-bench/roles/upload-ac-mqtt-http-robot/files/
cp ./schedulers.json ~/ablecloud.cn/src/deploys/ac-bench/roles/upload-ac-mqtt-http-robot/files/
cp ./robot.json ~/ablecloud.cn/src/deploys/ac-bench/roles/upload-ac-mqtt-http-robot/templates/robot.json

cd ~/ablecloud.cn/src/deploys/ac-bench/; ./upload-ac-mqtt-http-robot.sh