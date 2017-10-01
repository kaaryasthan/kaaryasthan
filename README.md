# Kaaryasthan - Collaborative Task Management for Small Teams

[![Go Report Card](https://goreportcard.com/badge/github.com/kaaryasthan/kaaryasthan)](https://goreportcard.com/report/github.com/kaaryasthan/kaaryasthan)
[![Build Status](https://travis-ci.org/kaaryasthan/kaaryasthan.svg?branch=master)](https://travis-ci.org/kaaryasthan/kaaryasthan)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

**Disclaimer:** This is [my pet project] for learning web application
development using Go & Angular.

This repository contains the source code of Kaaryasthan.  It has
source code for both fron-end and back-end.  The front-end of
Kaaryasthan is written in [Angular], back-end in [Go], and
[PostgreSQL] is used as the persistent data store.

> Kaaryasthan (കാര്യസ്ഥൻ) is a [Malayalam] word with meaning "manager".

## Development

This project requires [Go] version 1.8 or above.  This project also
requires [Node] version 6.11 or above, preferably an LTS release.

Once Go and Node is installed, you can install these utilities:

- <https://github.com/pilu/fresh>
- <https://github.com/jteeuwen/go-bindata>
- <https://github.com/elazarl/go-bindata-assetfs>
- <https://glide.sh>
- <https://cli.angular.io>

To install the above utilities:

    curl https://glide.sh/get | sh
    go get -u github.com/pilu/fresh
    go get -u github.com/jteeuwen/go-bindata/...
    go get -u github.com/elazarl/go-bindata-assetfs/...
    npm install -g @angular/cli

You can clone [Kaaryasthan] repository insdide `$GOPATH` using these
commands (Note: `$GOPATH` should be pointing to a single directory):

    mkdir -p $GOPATH/src/github.com/kaaryasthan
    cd $GOPATH/src/github.com/kaaryasthan
    git clone https://github.com/kaaryasthan/kaaryasthan.git

Now you can run `./build.sh` command.

    cd $GOPATH/src/github.com/kaaryasthan/kaaryasthan
    ./build.sh

Install [Docker] and [Docker Compose] and then run:

    docker-compose up -d

Finally, migrate database schema by running:

    ./kaaryasthan -migrate

## Running

To run the back-end development server (from the top-level directory):

    fresh

To run the front-end development server (from the `./web` directory):

    ng serve

## License

    Kaaryasthan - Collaborative Task Management for Small Teams
    Copyright (C) 2017 Baiju Muthukadan

    This program is free software: you can redistribute it and/or
    modify it under the terms of the GNU Affero General Public License
    as published by the Free Software Foundation, either version 3 of
    the License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
    Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public
    License along with this program.  If not, see
    <http://www.gnu.org/licenses/>.

[my pet project]: https://team-coder.com/pet-project
[Node]: https://nodejs.org/en
[Angular]: https://angular.io
[Go]: https://golang.org
[PostgreSQL]: https://www.postgresql.org
[Malayalam]: https://en.wikipedia.org/wiki/Malayalam
[Docker]: https://docs.docker.com
[Docker Compose]: https://docs.docker.com/compose
[Kaaryasthan]: https://github.com/kaaryasthan/kaaryasthan
