# mbo

A command-line interface to [MINDBODY Online](https://clients.mindbodyonline.com), a website people use to reserve spots in gym classes.

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
