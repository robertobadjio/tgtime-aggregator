package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/robertobadjio/tgtime-aggregator/internal/error_helper"
	"github.com/robertobadjio/tgtime-aggregator/pkg/time/endpoints"
)

const (
	basePostfix          = "/api"
	versionAPIPostfix    = "/v1"
	serviceStatusPostfix = "/service/status"
	timePostfix          = "/time"
)

type errorer interface {
	Error() error
}

// NewHTTPHandler ???
func NewHTTPHandler(ep endpoints.Set) http.Handler {
	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Methods(http.MethodGet).Path(serviceStatusPostfix).Handler(
		httptransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeHTTPServiceStatusRequest,
			encodeResponse,
			opts...,
		),
	)

	var api = router.PathPrefix(basePostfix).Subrouter()

	var api1 = api.
		PathPrefix(versionAPIPostfix).
		Subrouter()

	api1.Methods(http.MethodPost).
		Path(timePostfix).
		Handler(
			httptransport.NewServer(
				ep.CreateTimeEndpoint,
				decodeHTTPCreateTimeRequest,
				encodeResponse,
				opts...,
			),
		)

	return router
}

func decodeHTTPServiceStatusRequest(_ context.Context, _ *http.Request) (interface{}, error) {
	var req endpoints.ServiceStatusRequest
	return req, nil
}

func decodeHTTPCreateTimeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	// TODO: Validate body
	/*_, err := strconv.Atoi(vars["seconds"])
	if err != nil {
		return nil, &util.TimeInvalidArgument{Message: "seconds"}
	}

	_, err = strconv.Atoi(vars["router_id"])
	if err != nil {
		return nil, &util.TimeInvalidArgument{Message: "router_id"}
	}

	if _, err = net.ParseMAC(vars["mac_address"]); err != nil {
		return nil, &util.TimeInvalidArgument{Message: "mac_address"}
	}*/

	var req endpoints.CreateTimeRequest
	err := json.NewDecoder(r.Body).Decode(&req.Time)
	if err != nil {
		return nil, fmt.Errorf("error decode json body")
	}

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if e, ok := response.(errorer); ok && e.Error() != nil {
		encodeError(ctx, e.Error(), w)
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	var tiaErr *error_helper.TimeInvalidArgument

	switch {
	case errors.As(err, &tiaErr), errors.Is(err, error_helper.ErrInvalidMacAddress):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	_ = json.NewEncoder(w).Encode(map[string]interface{}{ // TODO: Handle error
		"error": err.Error(),
	})
}
