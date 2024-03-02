package main

import (
	"log"
	"net/http"
	"github.com/tiaguinho/gosoap"
)

type sumResponse struct {
	Result string `xml:"result"`
}

var (
	r sumResponse
)

func main() {
	httpClient := &http.Client{	}
	soap, err := gosoap.SoapClient("http://localhost:8001/WSDL", httpClient)
	if err != nil {
		log.Fatalf("SoapClient error: %s", err)
	}
	params := gosoap.Params{
		"term1": "10",
		"term2": "20",
	}

	res, err := soap.Call("sum", params)
	if err != nil {
		log.Fatalf("Call error: %s", err)
	}

	res.Unmarshal(&r)
	log.Println("Result: ", r.Result)
}
