package main

import (
    "fmt"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/log/level"
)

const addr = ":3001"

func main() {
    var logger log.Logger
    {
        logger = log.NewLogfmtLogger(os.Stdout)
        logger = log.With(logger, "ts", log.DefaultTimestampUTC)
        logger = log.With(logger, "caller", log.DefaultCaller)
    }

    var svc Service
    {
        svc = NewService(logger)
        // LoggingMiddleware(logger) returns a function which expects a svc as an argument
        svc = LoggingMiddleware(logger)(svc)
    }

    var handler http.Handler
    {
        handler = MakeHTTPHandler(svc, log.With(logger, "component", "HTTP"))
    }

    errs := make(chan error)

    go func() {
        interrupt := make(chan os.Signal, 1)
        signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)
        errs <- fmt.Errorf("%s", <-interrupt)
    }()

    go func() {
        _ = logger.Log("transport", "HTTP", "addr", addr)
        errs <- http.ListenAndServe(addr, handler)
    }()

    _ = level.Debug(logger).Log("exit", <-errs)

}
