# templates

To install:

```
go get github.com/ncw/gotemplate/...
go get github.com/govlas/templates/...
```
For more information about templates [see here](https://github.com/ncw/gotemplate).


## amap
template for async, channel based, thread-safe map
To use it using a special comment in your code. For example:
```
//go:generate gotemplate "github.com/govlas/tempates/amap" MyMap(int,string)
```
and run `go generate` in your code directory.
