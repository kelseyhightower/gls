gls - Distributed ls

# Build

Create the `/src/github.com/kelseyhightower` directory and clone the gls repo under it:

```
$ mkdir $GOPATH:/src/github.com/kelseyhightower
$ cd $GOPATH:/src/github.com/kelseyhightower
$ git clone https://github.com/kelseyhightower/gls.git
```

Run the build script:

```
$ ./build
```

# Usage

## Server

```
$ sudo glsd
```

## Client

```
$ gls --hostsfile=hosts.txt "/"
```
