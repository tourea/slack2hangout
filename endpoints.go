package main

import (
    "context"

    "github.com/go-kit/kit/endpoint"
)

// Endpoints REST endpoints
type Endpoints struct {
    ForwardEndpoint endpoint.Endpoint
}

type forwardRequest struct {
    space             string
    key               string
    token             string
    slackNotification SlackNotification
}

// GetServerEndpoints returns all endpoints
func GetServerEndpoints(s Service) Endpoints {
    return Endpoints{
        ForwardEndpoint: makeForwardEndpoint(s),
    }
}

func makeForwardEndpoint(s Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        r := request.(forwardRequest)
        err := s.Forward(ctx, r.space, r.key, r.token, r.slackNotification)
        if err != nil {
            return nil, err
        }
        return nil, nil  // TODO: return gchat response?
    }
}
