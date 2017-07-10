#!/bin/bash

directory="./robots"

function build_robots() {
	for file in `ls $directory`
	do
		if [ -d $directory/$file ]
		then
			cd $directory/$file; ./build.sh; cd ../..
		fi
	done
}

function main() {
	build_robots
	go build
}

main
