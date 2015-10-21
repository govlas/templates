# amap

[![Build Status](https://travis-ci.org/govlas/amap.svg?branch=master)](https://travis-ci.org/govlas/amap)

template for async, channel based, thread-safe map

To install:

```
go get github.com/ncw/gotemplate/...
go get github.com/govlas/amap
```

To use it using a special comment in your code. For example:
```
//go:generate gotemplate "github.com/govlas/amap" MyMap(int,string)
```
and run `go generate` in your code directory.

For more information about templates [see here](https://github.com/ncw/gotemplate).
