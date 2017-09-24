#!/bin/bash

WD=`pwd`

cd ${WD}
glide install

cd ${WD}/web
npm install
go generate

cd ${WD}
go generate

cd ${WD}
go build
