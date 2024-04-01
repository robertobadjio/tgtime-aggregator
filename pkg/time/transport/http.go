package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"tgtime-aggregator/internal/util"
	"tgtime-aggregator/pkg/time/endpoints"
)

type errorer interface {
	Error() error
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}

func NewHTTPHandler(ep endpoints.Set) http.Handler {
	router := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Methods(http.MethodGet).Path("/service/status").Handler(
		httptransport.NewServer(
			ep.ServiceStatusEndpoint,
			decodeHTTPServiceStatusRequest,
			encodeResponse,
			opts...,
		),
	)

	var api = router.PathPrefix("/api").Subrouter()

	var api1 = api.
		PathPrefix("/v1").
		Subrouter()

	api1.Methods(http.MethodPost).
		Path("/time").
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

	var tiaErr *util.TimeInvalidArgument

	switch {
	case errors.As(err, &tiaErr), errors.Is(err, util.ErrInvalidMacAddress):
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
