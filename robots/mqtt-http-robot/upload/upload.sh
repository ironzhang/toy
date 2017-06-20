#!/bin/bash

cp ../robot.so ~/ablecloud.cn/src/deploys/ac-mqtt-bench/roles/upload-mqtt-http-robot/files/
cp ../schedulers.json ~/ablecloud.cn/src/deploys/ac-mqtt-bench/roles/upload-mqtt-http-robot/files/
cp ./run-mqtt-http-robot.sh ~/ablecloud.cn/src/deploys/ac-mqtt-bench/roles/upload-mqtt-http-robot/files/
cp ./robot.json ~/ablecloud.cn/src/deploys/ac-mqtt-bench/roles/upload-mqtt-http-robot/templates/robot.json

cd ~/ablecloud.cn/src/deploys/ac-mqtt-bench/; ./upload-mqtt-http-robot.sh
