#!/bin/bash

cp ../robot.so ~/ablecloud.cn/src/deploys/ac-robots/roles/upload-mqtt-robot/files/
cp ../schedulers.json ~/ablecloud.cn/src/deploys/ac-robots/roles/upload-mqtt-robot/files/
cp ./run-mqtt-robot.sh ~/ablecloud.cn/src/deploys/ac-robots/roles/upload-mqtt-robot/files/
cp ./robot.json ~/ablecloud.cn/src/deploys/ac-robots/roles/upload-mqtt-robot/templates/robot.json

cd ~/ablecloud.cn/src/deploys/ac-robots/; ./upload-mqtt-robot.sh
