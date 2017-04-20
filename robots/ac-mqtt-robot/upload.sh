#!/bin/bash

cp ./robot.so ~/ablecloud.cn/src/deploys/ac-bench/roles/upload-ac-mqtt-robot/files/
cp ./schedulers.json ~/ablecloud.cn/src/deploys/ac-bench/roles/upload-ac-mqtt-robot/files/
cp ./robot.json.j2 ~/ablecloud.cn/src/deploys/ac-bench/roles/upload-ac-mqtt-robot/templates/robot.json

cd ~/ablecloud.cn/src/deploys/ac-bench/; ./upload-ac-mqtt-robot.sh
