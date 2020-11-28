# Logastic

[![build status](https://travis-ci.org/danil/logastic.svg?branch=master)](https://travis-ci.org/danil/logastic)
[![Build Status](https://cloud.drone.io/api/badges/danil/logastic/status.svg)](https://cloud.drone.io/danil/logastic)

JSON logging for Go.
Development of the Logastic is in pre alpha version and
is extremely slow and allocates a lots of memory)

## Install

    go get github.com/danil/logastic@v0.30.1

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

## Caveats

Numeric types appears in the message as a string

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
