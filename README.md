# Logastic

[![Go Reference](https://pkg.go.dev/badge/github.com/danil/logastic.svg)](https://pkg.go.dev/github.com/danil/logastic)
[![Build Status](https://cloud.drone.io/api/badges/danil/logastic/status.svg)](https://cloud.drone.io/danil/logastic)

JSON logging for Go.

## About

The software is considered to be at a alpha level of readiness -
its extremely slow and allocates a lots of memory)

## Install

    go get github.com/danil/logastic@v0.61.0

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
        Keys: [4]string{"message", "excerpt"},
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
    l.KV["version"] = logastic.String("1.1")
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
        Keys: [4]string{"message"},
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

## License

Copyright (C) 2020 Danil Kutkevich <danil@kutkevich.org>  
See the [LICENSE](./LICENSE) file for license rights and limitations (MIT)
