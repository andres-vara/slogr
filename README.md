# slogr

A structured logging package built on Go's [slog](https://pkg.go.dev/log/slog) with enhanced features and simpler configuration.

## Features

- **Multiple log levels**: Debug, Info, Warn, Error, Fatal
- **Handler types**: JSON and Text formats
- **Customizable output**: Any `io.Writer` support
- **Level prefixes**: Optional level tags in log messages
- **Runtime configuration**: Change level/output/handler dynamically
- **Context support**: Pass context through logs
- **Global logger**: Default instance for quick setup

## Installation

```bash
go get github.com/andres-vara/logr
```


## Quick Start


```go
package main

import (
    "context"
    "github.com/andres-vara/logr"
)

func main() {
    ctx := context.Background()
    logr.Info(ctx, "Application started")
    logr.Infof(ctx, "User %s logged in", "john@example.com")
}
```


## Basic Usage

You can set the default logger's threshold like so:
```go
logger := logr.New(os.Stdout, nil)
logger.SetLevel(slog.LevelWarn)
```

### Create a Logger

```go
// Custom logger with JSON formatting
logger := logr.New(os.Stdout, &logr.Options{
    Level: slog.LevelDebug,
    AddLevelPrefix: true,
    HandlerType: logr.HandlerTypeJSON,
})
```

### Log messages


### Custom Options

```go
package main

import (
    "context"
    "log/slog"
    "os"
    "github.com/andres-vara/logr"
)

func main() {
    // Create a logger with custom options
    logger := logr.New(os.Stdout, &logr.Options{
        Level: slog.LevelDebug,
        AddLevelPrefix: true,
        HandlerType: logr.HandlerTypeJSON,
    })
    ctx := context.Background()
    logger.Info(ctx, "This will be output as JSON")
}
```

### Runtime Configuration

```go
package main

import (
    "context"
    "github.com/andres-vara/logr"
)

func main() {
    ctx := context.Background()
    logr.Info(ctx, "Application started")
}
```

### Context Support

```go
package main

import (
    "context"
    "github.com/andres-vara/logr"
)

func main() {
    ctx := context.Background()
    logr.Info(ctx, "Application started")
}
```


