#!/bin/bash

WD=`pwd`

cd ${WD}/web
go generate

cd ${WD}
go generate

cd ${WD}
go build
