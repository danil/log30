# Logastic

[![Build Status](https://cloud.drone.io/api/badges/danil/logastic/status.svg)](https://cloud.drone.io/danil/logastic)
[![Go Reference](https://pkg.go.dev/badge/github.com/danil/logastic.svg)](https://pkg.go.dev/github.com/danil/logastic)

JSON logging for Go.

* [About](#about)
* [Install](#install)
* [Usage](#usage)
* [Use as GELF formater](#use-as-gelf-formater)
* [Caveat: numeric types appears in the message as a string](#caveat-numeric-types-appears-in-the-message-as-a-string)
* [Benchmark](#benchmark)
* [License](#license)

## About

The software is considered to be at a alpha level of readiness -
its extremely slow and allocates a lots of memory)

## Install

    go get github.com/danil/logastic@v0.84.0

## Usage

Set Logastic as global logger

```go
package main

import (
    "os"
    "log"

    "github.com/danil/logastic"
)

func main() {
    l := logastic.Log{
        Output: os.Stdout,
        Trunc: 12,
        Keys: [4]json.Marshaler{logastic.String("message"), logastic.String("excerpt")},
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

## Use as GELF formater

```go
package main

import (
    "log"
    "os"

    "github.com/danil/logastic"
)

func main() {
    l := logastic.GELF()
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

## Caveat: numeric types appears in the message as a string

```go
package main

import (
    "log"
    "os"

    "github.com/danil/logastic"
)

func main() {
    l := logastic.Log{
        Output: os.Stdout,
        Keys: [4]json.Marshaler{logastic.String("message")},
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

## Benchmark

```
go test -bench=. ./...
goos: linux
goarch: amd64
pkg: github.com/danil/logastic
BenchmarkLogastic/io.Writer_36-8         	  323197	      3678 ns/op
BenchmarkLogastic/fmt.Fprint_io.Writer_1009-8         	  121657	     10417 ns/op
```

## License

Copyright (C) 2021 [Danil Kutkevich](https://github.com/danil)  
See the [LICENSE](./LICENSE) file for license rights and limitations (MIT)
