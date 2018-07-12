#!/bin/sh
set -e

go get github.com/rcrowley/go-metrics
go get github.com/stretchr/testify/assert
go get gopkg.in/yaml.v2
go get github.com/Masterminds/glide

mkdir -p $HOME/gopath/src/github.com/ServiceComb/
cd $HOME/gopath/src/github.com/ServiceComb
git clone http://github.com/ServiceComb/go-chassis
cd $HOME/gopath/src/github.com/ServiceComb/go-chassis
glide install
mkdir -p $HOME/gopath/src/github.com/ServiceComb/go-chassis/vendor/github.com/ServiceComb/cse-collector/
rsync -az ${TRAVIS_BUILD_DIR}/ $HOME/gopath/src/github.com/ServiceComb/go-chassis/vendor/github.com/ServiceComb/cse-collector/
export TRAVIS_BUILD_DIR=$HOME/gopath/src/github.com/ServiceComb/go-chassis/vendor/github.com/ServiceComb/cse-collector/
cd $HOME/gopath/src/github.com/ServiceComb/go-chassis/vendor/github.com/ServiceComb/cse-collector/


cd $GOPATH/src/github.com/ServiceComb/go-chassis/vendor/github.com/ServiceComb/cse-collector
#Start unit test
for d in $(go list ./...); do
    echo $d
    echo $GOPATH
    cd $GOPATH/src/$d
    if [ $(ls | grep _test.go | wc -l) -gt 0 ]; then
        go test
    fi
done
