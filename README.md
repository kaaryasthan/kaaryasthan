# Kaaryasthan - Task Management for Small Teams

[![Go Report Card](https://goreportcard.com/badge/github.com/kaaryasthan/kaaryasthan)](https://goreportcard.com/report/github.com/kaaryasthan/kaaryasthan)
[![Build Status](https://travis-ci.org/kaaryasthan/kaaryasthan.svg?branch=master)](https://travis-ci.org/kaaryasthan/kaaryasthan)
[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

**Disclaimer:** This is [my pet project] for learning web application
development using Go & Angular.  See the [FAQ](#why-this-project) for
more details.

Kaaryasthan helps you to manage private projects.  You can add tasks &
issues called items to your project.  An item will have title,
description, discussions, creator, assignees and labels.  Multiple
milestones can be created with due date and priority-ordered items.

This repository contains the source code of Kaaryasthan.  It has
source code for both front-end and back-end.  The front-end of
Kaaryasthan is written in [Angular], back-end in [Go], and
[PostgreSQL] is used as the persistent data store.

> Kaaryasthan (കാര്യസ്ഥൻ) is a [Malayalam] word with meaning "manager".

## Key Features

- All URLs will require require authentication except login &
  registration.  (Note: This makes it unsuitable for open source
  projects with public issues.  Only private projects can be hosted
  using Kaaryasthan)
- Threaded discussions.  Discussions can be added under items.
  Comments can be added under each discussion.
- Easy deployment.  Entire application including front-end is
  available as a single binary. You also need to install PostgreSQL
  and NGINX.

## Development

This project requires [Go] version 1.8 or above.  This project also
requires [Node] version 6.11 or above, preferably an LTS release.

Once Go and Node is installed, you can install these utilities:

- <https://github.com/pilu/fresh>
- <https://github.com/jteeuwen/go-bindata>
- <https://github.com/elazarl/go-bindata-assetfs>
- <https://github.com/alecthomas/gometalinter>
- <https://glide.sh>
- <https://cli.angular.io>

To install the above utilities:

    curl https://glide.sh/get | sh
    go get github.com/pilu/fresh
    go get github.com/jteeuwen/go-bindata/...
    go get github.com/elazarl/go-bindata-assetfs/...
    go get github.com/alecthomas/gometalinter
    gometalinter --install
    npm install -g @angular/cli

You can clone [Kaaryasthan] repository inside `$GOPATH` using these
commands (Note: `$GOPATH` should be pointing to a single directory):

    mkdir -p $GOPATH/src/github.com/kaaryasthan
    cd $GOPATH/src/github.com/kaaryasthan
    git clone https://github.com/kaaryasthan/kaaryasthan.git

Now you can run `./build.sh` command.

    cd $GOPATH/src/github.com/kaaryasthan/kaaryasthan
    ./build.sh

Install [Docker] and [Docker Compose] and then run:

    docker-compose up -d

Migrate database schema by running:

    ./kaaryasthan -migrate

To run test cases:

    ./runtests.sh

### Running development servers

To run the back-end development server (from the top-level directory):

    fresh

To run the front-end development server (from the `./web` directory):

    ng serve

## License

    Kaaryasthan - Task Management for Small Teams
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

## FAQ

### Why this project?

This is [my pet project] for learning web application development
using Go & Angular.  Here I am trying to explore various ideas about
web application development.  I started this [project in 2014] with
another name.  Few things I am learning through this project:

- Designing RESTful APIs
- Database design
- Front-end development
- Web application security
- Load testing

As part of my professional work, I can work on some of these things.
But this side project gives me more freedom and flexibility.

### Are you seeking external contributions?

As of now, I am not really looking for any external contributions.  I
am not really sure how much time I can spend on pull request reviews.
Though if you have any feedback, I would be happy to listen.  You can
send your feedback to: baiju.m.mail AT gmail.com

### Why did you choose AGPLv3+ as the license?

I believe that's the best license for a web application.  From the
[GNU website]:

> The GNU Affero General Public License is a modified version of the
> ordinary GNU GPL version 3.  It has one added requirement: if you
> run a modified program on a server and let other users communicate
> with it there, your server must also allow them to download the
> source code corresponding to the modified version running there.

### Can I use Kaaryasthan to manage open source project tasks & issues?

I think that may not be possible.  Because all URLs will require
authentication except login & registration.  Since open source
projects requires public trackers, this won't be a desirable solution.
Kaaryasthan is designed exclusively for private projects.  And I am
not planning to change it. These are few reasons for not supporting:

1. There are many good trackers available for open source projects
2. Kaaryasthan is not exclusively designed for software projects
3. Don't want to make the system complex to handle Slashdot effects
4. No plan for search engine optimization (SEO) for Google, Bing etc.
5. Reducing scope will make the system simple

The architecture doesn't support Kaaryasthan being used as a public
tracker.  Despite this limitation, if you really want to use, a work
around could be using a reverse proxy.

---

> [kaaryasthan.org](https://kaaryasthan.org) &nbsp;&middot;&nbsp;
> Demo [demo.kaaryasthan.org](https://demo.kaaryasthan.org) &nbsp;&middot;&nbsp;
> IRC [#kaaryasthan@freenode](https://riot.im/app/#/room/#freenode_#kaaryasthan:matrix.org) &nbsp;&middot;&nbsp;
> [Mailing List](https://groups.google.com/forum/#!forum/kaaryasthan) &nbsp;&middot;&nbsp;
> Twitter [@kaaryasthan](https://twitter.com/kaaryasthan)

[my pet project]: https://team-coder.com/pet-project
[Node]: https://nodejs.org/en
[Angular]: https://angular.io
[Go]: https://golang.org
[PostgreSQL]: https://www.postgresql.org
[Malayalam]: https://en.wikipedia.org/wiki/Malayalam
[Docker]: https://docs.docker.com
[Docker Compose]: https://docs.docker.com/compose
[Kaaryasthan]: https://github.com/kaaryasthan/kaaryasthan
[project in 2014]: https://github.com/baijum/pitracker
[GNU website]: https://www.gnu.org/licenses/why-affero-gpl.en.html
