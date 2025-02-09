# golflog [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci]

`golflog` is a logging utility package built around `go-logr`. It's main use
is a higher level api for building context based logging iteratively and
stored in `context.Context`.

## Installation

	go get -u github.com/floatme-corp/golflog

## Use

As close to your application entrypoint or initial context creation, create
a `Configurator` and `logr.Logger`, then set that logger in the
`context.Context`:
```golang
import (
    "context"
    "fmt"

    "github.com/floatme-corp/golflog
)

func main() {
    configurator := golflog.NewZapProductionConfigurator()
    verbosity := 0

    logger, err := golflog.NewLogger(configurator, "RootLoggerName", verbosity)
    if err != nil {
        panic(fmt.Sprintf("runtime error: failed to setup logging: %s", err))
    }

    ctx := golflog.NewContext(context.Background(), logger)

    ...
}
```
Later if your application has another section, such as a `Queue` handler, you
can set that in the logger name:
```golang
func HandleQueue(ctx context.Context) {
    ctx = golflog.ContextWithName(ctx, "Queue")
}
```
All messages logged from that point forware will have the name added to the
existing name: `RootLoggerName.Queue`.

Additional values can be setup as well for future logging:
```golang
func HandleUser(ctx context.Context, userID string) {
    ctx = golflog.ContextWithValues(ctx, "user_id", userID)
}
```
Funcitons are guaranteed to be able to get a logger from any context:
```golang
func randoFunc(ctx context.Context) {
    log := golflog.AlwaysFromContext(ctx)
    log.Info("my log message")
}
```
If the context does not have a logger associated with it `golflog` will
create a fallback logger with the default configuration. If that fails
it will fallback to logging via `fmt.Fprintln` to `os.Stdout`

### Env setup

An alternative to calling `golflog.NewLogger` with the parameters, is to call
`NewLoggerFromEnv` and give it only the root name:
```golang
import (
    "context"
    "fmt"

    "github.com/floatme-corp/golflog
)

func main() {
    logger, err := golflog.NewLoggerFromEnv("RootLoggerName")
    if err != nil {
        panic(fmt.Sprintf("runtime error: failed to setup logging: %s", err))
    }

    ctx := golflog.NewContext(context.Background(), logger)

    ...
}
```
`NewLoggerFromEnv` uses the environment variables `LOG_PRODUCTION`,
`LOG_IMPLEMENTATION`, and `LOG_VERBOSITY` to configure the logger. If they
do not exist, it will default to configuring a production logger, with
a `zap` `Configurator` at `0` verbosity (normal / info level).

-------------------------------------------------------------------------------

Released under the [Apache 2.0 License].

[Apache 2.0 License]: LICENSE
[doc-img]: https://pkg.go.dev/badge/github.com/floatme-corp/golflog
[doc]: https://pkg.go.dev/github.com/floatme-corp/golflog
[ci-img]: https://github.com/floatme-corp/golflog/actions/workflows/test.yaml/badge.svg
[ci]: https://github.com/floatme-corp/golflog/actions/workflows/test.yaml
