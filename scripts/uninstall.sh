#! /bin/bash

echo "Uninstalling..."

if [ -d "/usr/local/bin/cli-go" ] ; then
    rm /usr/local/bin/cli-go
fi

if [ -d "/usr/local/cli/go" ] ; then
    rm -r /usr/local/cli/go
fi

echo "Uninstalled"