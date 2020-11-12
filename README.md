# cli-go
a cli for go

## Roadmap

- implementation source install
- switch to use go modules
- implement with git library

## Build

```bash
make deps
make all
```

## Installation
Follow build instructions first

```bash
sudo ./install.sh
```

## Usage
### List current go version
```bash
cli-go version
```
### List available go versions
```bash
cli-go list   
```   
### Install specific version of go from source or binary
```bash
cli-go install --{source|binary} --version {version}
```
### Init a new go project in current directory
```bash
cli-go init
```