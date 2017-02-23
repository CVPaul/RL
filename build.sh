#!/bin/bash
CUR_DIR=`pwd`
GOPATH=${CUR_DIR}
THIRDPARTY=${GOPATH}/src/thirdparty
GOPATH=${THIRDPARTY}:$GOPATH
export GOPATH
export LIBRARY_PATH=${THIRDPARTY}/lib
#export LD_LIBRARY_PATH=${THIRDPARTY}/lib   # For Linux
export DYLD_LIBRARY_PATH=${THIRDPARTY}/lib # Export TansorFlow path for OS X
#echo "GOPATH="$GOPATH
#echo "DYLD_LIBRARY_PATH="$DYLD_LIBRARY_PATH
#echo "Starting go test[tensorflow]..."
#go get github.com/tensorflow/tensorflow/tensorflow/go
#go test github.com/tensorflow/tensorflow/tensorflow/go
#echo "End of go test[tensorflow]"

#download the packages
#go get -u -v azul3d.org/examples/...
#go get -u -v azul3d.org/engine/gfx
#go get code.google.com/p/biogo/
# format
gofmt -l -w -s src/
#==================build[main]===================#
cd ${CUR_DIR} && go build -o main ./src/main.go
#build[main] result
ret=$?
if [ $ret -ne 0 ];then
    echo "===== build[main] failure ====="
    exit $ret
else
    echo "===== build[main] successfully ====="
fi
#===================build[train]=================#
cd ${CUR_DIR} && go build -o training ./src/training.go
#build[main] result
ret=$?
if [ $ret -ne 0 ];then
   echo "===== build[training] failure ====="
   exit $ret
else
   echo "===== build[training] successfully ====="
fi
#===================build[test]==================#
cd ${CUR_DIR} && go build -o testing ./src/testing.go
#build[main] result
ret=$?
if [ $ret -ne 0 ];then
   echo "===== build[testing] failure ====="
   exit $ret
else
   echo "===== build[testing] successfully ====="
fi
exit
