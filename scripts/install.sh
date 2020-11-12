#! /bin/bash

echo "Installing..."

mkdir -p /usr/local/cli/go

cd ../src

cat makefile

make build

cp cli-go /usr/local/cli/go
cp releases.json /usr/local/cli/go
cp -r tpl /usr/local/cli/go

cd ..

if [ -d "/usr/local/bin/cli-go" ] ; then
    rm /usr/local/bin/cli-go
fi

ln -s /usr/local/cli/go/cli-go /usr/local/bin/cli-go

echo "Installed"