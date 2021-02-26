log0
=====

[![Build Status](https://cloud.drone.io/api/badges/danil/log0/status.svg)](https://cloud.drone.io/danil/log0)
[![Go Reference](https://pkg.go.dev/badge/github.com/danil/log0.svg)](https://pkg.go.dev/github.com/danil/log0)

JSON logging for Go.  
Source files are distributed under the BSD-style license
found in the [LICENSE](./LICENSE) file.

<!-- markdown-toc start - Don't edit this section. Run M-x markdown-toc-refresh-toc -->

* [About](#about)
* [Install](#install)
* [Usage](#usage)
* [Use as GELF formater](#use-as-gelf-formater)
* [Caveat: numeric types appears in the message as a string](#caveat-numeric-types-appears-in-the-message-as-a-string)
* [Benchmark](#benchmark)

<!-- markdown-toc end -->

About
-----

The software is considered to be at a alpha level of readiness -
its extremely slow and allocates a lots of memory)

Install
-------

    go get github.com/danil/log0@v0.126.0

Usage
-----

Set log0 as global logger

```go
package main

import (
    "os"
    "log"

    "github.com/danil/log0"
)

func main() {
    l := log0.Log{
        Output: os.Stdout,
        Trunc: 12,
        Keys: [4]json.Marshaler{log0.String("message"), log0.String("excerpt")},
        Marks: [3][]byte{[]byte("…")},
        Replace: [][]byte{[]byte("\n"), []byte(" ")},
    }
    log.SetFlags(0)
    log.SetOutput(l)

    log.Print("Hello,\nWorld!")
}
```

Output:

```json
{
    "message":"Hello,\nWorld!",
    "excerpt":"Hello, World…"
}
```

Use as GELF formater
--------------------

```go
package main

import (
    "log"
    "os"

    "github.com/danil/log0"
)

func main() {
    l := log0.GELF()
    l.Output = os.Stdout
    log.SetFlags(0)
    log.SetOutput(l)
    log.Print("Hello,\nGELF!")
}
```

Output:

```json
{
    "version":"1.1",
    "short_message":"Hello, GELF!",
    "full_message":"Hello,\nGELF!",
    "timestamp":1602785340
}
```

Caveat: numeric types appears in the message as a string
--------------------------------------------------------

```go
package main

import (
    "log"
    "os"

    "github.com/danil/log0"
)

func main() {
    l := log0.Log{
        Output: os.Stdout,
        Keys: [4]json.Marshaler{log0.String("message")},
    }
    log.SetFlags(0)
    log.SetOutput(l)

    log.Print(123)
    log.Print(3.21)
}
```

Output 1:

```json
{
    "message":"123"
}
```

Output 2:

```json
{
    "message":"3.21"
}
```

Benchmark
---------

```
go test -bench=. ./...
goos: linux
goarch: amd64
pkg: github.com/danil/log0
cpu: Intel(R) Core(TM) i7-8565U CPU @ 1.80GHz
BenchmarkLog0/io.Writer_76-8         	  289021	      3951 ns/op
BenchmarkLog0/fmt.Fprint_io.Writer_1126-8         	  106959	      9824 ns/op
```
