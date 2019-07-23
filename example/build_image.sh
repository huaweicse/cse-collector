#!/bin/bash


export ROOT=$PWD
cd order
set -e 
set -x 
#GO111MODULE=on go mod vendor
build_image(){
    cd ${ROOT}/$1
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main
    sudo docker build --build-arg VERSION=$2 -t ${REPO}/$1:$2 .
    sudo docker push  ${REPO}/$1:$2
}

build_image order 1.0.0
build_image restaurant 1.0.0
build_image restaurant 1.0.1
