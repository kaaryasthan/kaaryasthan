# Kaaryasthan

**WARNING:** This is [my pet](https://team-coder.com/pet-project) project for learning web application development using Go.

Collaborative task management system for small teams

## Goals

- RESTful API using JSONAPI
- Internal authentication source

## Installation

- Install Go 1.8 or above version
- Install Ember 2.10.x and Node 6.5.x
- Install https://github.com/pilu/fresh & https://glide.sh/

First time you need to install packages required for Ember.js
To perform this, run `ember install` from `web` directory:

    cd web
    ember install

Then run the `build.sh` from top directory

Next migrate database schema by running `./kaaryasthan -migrate`

## Running

    fresh

## Running tests

    ./test.sh
