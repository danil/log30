# Mujlog

Mujlog (Multiline JSON Log) is a formatter and writer.

Mujlog in pre alpha version is extremely slow and allocates a lots of memory.

## Usage

Set Mujlog as global logger

```go
package main

import (
    "log"

    "gitlab.rocketbank.sexy/backend/random-values/jaunt/mujlog"
)

func main() {
    l := mujlog.Log{
        Output: os.Stdout,
        Short: "shortMessage",
        Full: "fullMessage",
        File: "file",
        TruncateMax: 1024,
        TruncateMin: 120,
    }
    log.SetOutput(l)

    log.Println("Hello,\nWorld!")
}
```

Output:

```json
{
    "shortMessage":"Hello, World!",
    "fullMessage":"Hello,\nWorld!"
}
```

## Use Mujlog as GLEF formater

```go
package main

import (
    "log"

    "gitlab.rocketbank.sexy/backend/random-values/jaunt/mujlog"
)

func main() {
    glf := mujlog.GELF()
    glf.Output = os.Stdout
    glf.Fields["host"] = "example.com"

    log.SetOutput(glf)

    log.Println("Hello,\nGELF!")
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

Numeric types appears in the short/full messages as a string. For example:

```go
package main

import (
    "log"

    "gitlab.rocketbank.sexy/backend/random-values/jaunt/mujlog"
)

func main() {
    l := mujlog.Log{
        Output: os.Stdout,
        Short: "shortMessage",
        Full: "fullMessage",
        File: "file",
        TruncateMax: 1024,
        TruncateMin: 120,
    }
    log.SetOutput(l)

    log.Println(123)
    log.Println(3.21)
}
```

Output:

```json
{
    "shortMessage":"123"
}
```

Output second:

```json
{
    "shortMessage":"3.21"
}
```
