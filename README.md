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

# mbo login
  -p="": Password. Will prompt if not passed.
  -studio="": Studio ID. Will prompt if not passed.
  -u="": Username. Will prompt if not passed.


# mbo logout


# mbo ls
  -date="": list classes as of this date. Format is MM/DD/YYYY. Default is today.


# mbo register
  -date="": Class date
  -id="": Class ID
```