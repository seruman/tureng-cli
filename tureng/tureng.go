package tureng

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	url         = "http://ws.tureng.com/TurengSearchServiceV4.svc/Search"
	contentType = "application/json"
)

type Result struct {
	Category string `json:"CategoryEN"`
	Term     string `json:"Term"`
	TypeEN   string `json:"TypeEN"`
}

func Translate(word string) ([]Result, error) {
	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(requestPayload{Term: word})
	if err != nil {
		return nil, err
	}

	apiResponse, err := doRequest(&payload)
	if err != nil {
		return nil, err
	}

	if !apiResponse.IsSuccessful {
		if apiResponse.Exception != "" {
			return nil, fmt.Errorf("api-exception: %s", apiResponse.Exception)
		}

		return nil, fmt.Errorf("api-response: is not successful")
	}

	if apiResponse.MobileResult.IsFound != 1 {
		return nil, fmt.Errorf("api-response: no results")
	}

	return apiResponse.MobileResult.Results, nil
}

func doRequest(body io.Reader) (*apiResponse, error) {
	resp, err := http.Post(url, contentType, body)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}

	defer resp.Body.Close()

	var responsePayload apiResponse
	err = json.NewDecoder(resp.Body).Decode(&responsePayload)
	if err != nil {
		return nil, fmt.Errorf("json: %w", err)
	}

	return &responsePayload, nil
}

type apiResponse struct {
	Exception    string `json:"ExceptionMessage"`
	IsSuccessful bool   `json:"IsSuccessful"`
	MobileResult struct {
		IsFound  int      `json:"IsFound"`
		IsTRToEN int      `json:"IsTRToEN"`
		Results  []Result `json:"Results"`
	} `json:"MobileResult"`
}

type requestPayload struct {
	Term string `json:"Term"`
}
