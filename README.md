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
    log.SetOutput(mujlog.Mujlog{Output: os.Stdout})

    log.Println("Hello\nWorld!")
}
```

```json
{
    "short_message":"Hello",
    "full_message":"Hello\nWorld!"
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

    log.Println("Hello\nGELF!")
}
```

```json
{
    "short_message":"Hello",
    "full_message":"Hello\nGELF!"
}
```
