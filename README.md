# RackspaceFiles
Simple tool for files managment in rackspace cloud. Now it supports `list`, `delete`, `upload`, `download` commands 

## Setup / Install

Get and compile RackspaceFiles:

    apt-get install golang
    export GOPATH=~/golang
    mkdir -p ${GOPATH}
    go get https://github.com/asp24/rackspace-files
    cd $GOPATH/src/https://github.com/asp24/rackspace-files
    go build -v

### How-to use

    ./rackspace-files --help
