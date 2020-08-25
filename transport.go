package main


import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/transport"
    httptransport "github.com/go-kit/kit/transport/http"
    "github.com/gorilla/mux"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// MakeHTTPHandler creates http handler for REST app
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
    r := mux.NewRouter()
    endpoint := GetServerEndpoints(s)

    options := []httptransport.ServerOption{
        httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
        httptransport.ServerErrorEncoder(encodeError),
    }

    // POST    /{space}        forwards slack message to google hangout space
    // GET     /metrics        retrieves monitoring metrics

    r.Methods("POST").Path("/v1/spaces/{space}/messages").Handler(httptransport.NewServer(
        endpoint.ForwardEndpoint,
        decodeForwardRequest,
        encodeResponse,
        options...,
    ))

    r.Methods("GET").Path("/metrics").Handler(promhttp.Handler())

    return r
}

func decodeForwardRequest(_ context.Context, r *http.Request) (interface{}, error) {
    vars := mux.Vars(r)
    s := vars["space"]
    k := r.FormValue("key")
    t := r.FormValue("token")
    var n SlackNotification
    err := json.NewDecoder(r.Body).Decode(&n)
    if err != nil {
        return nil, err
    }

    return forwardRequest{space: s, key: k, token: t, slackNotification: n}, nil
}

// errorer is implemented by all concrete response types that may contain errors. It allows us to change the
// HTTP response code without needing to trigger an endpoint (transport-level) error.
type errorer interface {
    error() error
}

// encodeResponse is the common method to encode all response types to the client.
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
    if e, ok := response.(errorer); ok && e.error() != nil {
        encodeError(ctx, e.error(), w)
        return nil
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
    if err == nil {
        panic("encodeError with nil error")
    }
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    // TODO: what is alertmanager expecting (if anything)
    w.WriteHeader(http.StatusInternalServerError)
    _ = json.NewEncoder(w).Encode(map[string]interface{}{
        "error":        err.Error(),
        "responseTime": time.Now().Unix(),
    })
}
