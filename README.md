# mbo

A command-line interface to the client website of [MINDBODY Online](https://clients.mindbodyonline.com)
Allows you to browse and register for classes from the command line.

## Installation

Install Go. Then run:

```
mkdir /tmp/mbotmp
export GOPATH=/tmp/mbotmp
go get github.com/rgarcia/mbo
sudo cp $GOPATH/bin/mbo /usr/local/bin/
```

## Usage

```
$ mbo -h
Usage of mbo:
Commands:
   login     Start a session with MBO
   logout    End session with MBO
   ls        List classes
   register  Register for a class
   schedule  Show your schedule
   cancel    Cancel a visit
```