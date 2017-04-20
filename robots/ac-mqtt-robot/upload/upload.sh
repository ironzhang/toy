#!/bin/bash

cp ../robot.so ~/ablecloud.cn/src/deploys/ac-mqtt-bench/roles/upload-ac-mqtt-robot/files/
cp ../schedulers.json ~/ablecloud.cn/src/deploys/ac-mqtt-bench/roles/upload-ac-mqtt-robot/files/
cp ./robot.json ~/ablecloud.cn/src/deploys/ac-mqtt-bench/roles/upload-ac-mqtt-robot/templates/

cd ~/ablecloud.cn/src/deploys/ac-mqtt-bench/; ./upload-ac-mqtt-robot.sh
