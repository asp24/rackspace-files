# RackspaceFiles
Simple tool for files managment in rackspace cloud. Now it supports `list`, `delete`, `upload`, `download` commands 

## Setup / Install

Install go if needed

    apt-get install golang
    export GOPATH=~/golang

Then get and compile RackspaceFiles:

    go get github.com/asp24/rackspace-files
    cd $GOPATH/src/github.com/asp24/rackspace-files
    go build -v

### How-to use

    ./rackspace-files --help
