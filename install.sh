#!/bin/sh
#
# This basic install.sh was written in the hopes that will be useful
# and reusable to those developing Go programs, but WITHOUT ANY WARRANTY;
# without even the implied warranty of MERCHANTABILITY or FITNESS FOR A
# PARTICULAR PURPOSE.  See the GNU General Public License for more details.
#
# This program is free software; you can redistribute it and/or
# modify it under the terms of the GNU General Public
# License v2 as published by the Free Software Foundation.
#
# Copyright (C) 2016  Nicholas D Steeves <nsteeves@gmail.com>
#

export GOPATH="$HOME/go"
export PACKAGE="github.com/tshprecher/gospell"

command -v go >/dev/null || { echo "command 'go' not found."; exit 1; }

echo "Creating required dirs"
mkdir -p "$GOPATH"/src/`dirname "$PACKAGE"`
ln -s "$PWD" "$GOPATH"/src/"$PACKAGE"
mkdir "$GOPATH"/bin
mkdir "$GOPATH"/pkg

echo "Installing `basename $PACKAGE`"
go install $PACKAGE

echo "Your PATH is currently: $PATH"
echo "You may need to add $GOPATH/bin to it, like so:"
echo "PATH=\$PATH:$GOPATH/bin"
echo
echo "See $GOPATH/src/$PACKAGE/README for more info"
