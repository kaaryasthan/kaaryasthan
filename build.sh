#!/bin/bash

WD=`pwd`

cd ${WD}/web
npm install
go generate

cd ${WD}
go generate

cd ${WD}
go build
