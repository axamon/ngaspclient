# Tizzy

[![Build Status](https://secure.travis-ci.org/panamafrancis/tizzy.png?branch=master)](http://travis-ci.org/panamafrancis/tizzy)
[![GoDoc](https://godoc.org/github.com/panamafrancis/tizzy?status.svg)](https://godoc.org/github.com/panamafrancis/tizzy)
[![License](https://img.shields.io/github/license/panamafrancis/tizzy.svg)](https://github.com/panamafrancis/tizzy/blob/master/LICENSE)

An in-memory copy of time.LoadLocation(), nothing original simply a workaround
for when you find your code calling time.LoadLocation() at 1khz. This avoids
opening a zip file on disk to unmarshal files containing bindata every call.

The timezone locations are loaded from $GOROOT/lib/time/zoneinfo.zip, the
default for unix and windows.

Any copyrights belong to the Go authors.


## Usage

```go
    import "github.com/panamafrancis/tizzy"

    ...

    loc, err := tizzy.LoadLocation("Europe/Berlin")
    
    //or faster...
    
    locv, err := tizzy.LoadLocationValue("Europe/Berlin")
```

## Benchmarks

On an early 2015 macbook pro, 2.9 GHz Intel Core i5, 16GB DDR3

```
go test -v -bench . -benchmem

...

goos: darwin
goarch: amd64
pkg: github.com/panamafrancis/tizzy
BenchmarkLoadLocationValue-4   	20000000	        82.9 ns/op	       0 B/op	       0 allocs/op
BenchmarkLoadLocation-4        	10000000	       186 ns/op	      96 B/op	       1 allocs/op
BenchmarkTimeLoadLocation-4    	  100000	     14176 ns/op	    1922 B/op	      10 allocs/op
```

## TODO

 - Support all operating systems, not just assume $GOROOT/lib/time/zoneinfo.zip exists.
