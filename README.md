# mbo

A command-line interface to [MINDBODY Online](https://clients.mindbodyonline.com), a website I use to reserve spots in gym classes.

## Installation

Install Go. Install mercurial. Then run:

```
mkdir /tmp/mbotmp
export GOPATH=/tmp/mbotmp
go get github.com/rgarcia/mbo
sudo cp $GOPATH/bin/mbo /usr/local/bin/
```

Make sure `/usr/local/bin` is in your PATH.

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

Pass `-h` to any subcommand for more information.

## Automate it

The script `cron.sh` contains an example script that will automatically check for a class at a certain time and register for it.

## Example

![](https://raw.github.com/rgarcia/mbo/master/mbo.gif)
