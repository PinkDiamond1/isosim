package crypto

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/logrus"
	"github.com/go-kit/kit/transport"
	httptransport "github.com/go-kit/kit/transport/http"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

const URLCryptoPinGen = "/iso/v1/crypto/pin_gen"

func pinGenReqDecoder(ctx context.Context, req *http.Request) (response interface{}, err error) {

	fmt.Println("Decoding...")

	reqData, err := ioutil.ReadAll(req.Body)
	fmt.Println(string(reqData))
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	pgr := &PinGenRequest{}
	if err := json.Unmarshal(reqData, pgr); err != nil {
		return nil, err
	}

	return *pgr, nil

}

// decode the response into JSON - generic decoder
func respEncoder(ctx context.Context, rw http.ResponseWriter, response interface{}) error {

	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		errorEncoder(ctx, f.Failed(), rw)
		return nil
	}
	rw.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
	rw.Header().Add("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(rw).Encode(response)
}

type errorWrapper struct {
	Error string `json:"error"`
}

func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	//TODO:: construct specific error types based on err
	w.Header().Add("Access-Control-Allow-Origin", "http://localhost:3000")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(errorWrapper{Error: err.Error()})
}

func RegisterHTTPTransport() {

	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logrus.NewLogrusLogger(log.New()))),
	}

	service := &serviceImpl{}

	http.Handle(URLCryptoPinGen, httptransport.NewServer(pinGenEndpoint(service), pinGenReqDecoder, respEncoder, options...))

}
