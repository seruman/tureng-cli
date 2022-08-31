package tureng

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	apiAddr     = "http://api.tureng.com/v1/dictionary"
	contentType = "application/json"

	TurkishEnglish = "entr"
	GermanEnglish  = "ende"
	SpanishEnglish = "enes"
	FrenchEnglish  = "enfr"
)

var defaultClient = NewClient()

type Client struct {
	httpClient *http.Client
	dictionary string
}

func NewClient(opts ...clientOpts) *Client {
	c := &Client{
		httpClient: http.DefaultClient,
		dictionary: TurkishEnglish,
	}
	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Translate(word string) ([]TermResult, error) {
	apiResponse, err := c.doRequest(apiAddr, word)
	if err != nil {
		return nil, err
	}

	if !apiResponse.IsFound {
		if len(apiResponse.Suggestions) == 0 {
			return nil, fmt.Errorf("%q not found", word)
		}

		return nil, fmt.Errorf("%q not found, suggestions;\n %s", word, strings.Join(apiResponse.Suggestions, "\n "))
	}

	return apiResponse.AFullTextResults, nil
}

func (c *Client) doRequest(baseAddr string, word string) (*apiResponse, error) {
	addr, err := url.JoinPath(baseAddr, c.dictionary, word)
	if err != nil {
		return nil, fmt.Errorf("url-join : %w", err)
	}
	req, err := http.NewRequest(http.MethodGet, addr, nil)
	if err != nil {
		return nil, fmt.Errorf("new-request : %w", err)
	}

	// NOTE: Tureng's API responds with Cloudflare challenge if Go's default
	// HTTP user agent is used. Mimic iOS app's request headers.
	req.Header.Set("User-Agent", "Tureng/2012061663 CFNetwork/1335.0.3 Darwin/21.6.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-GB,en;q=0.9")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http: %w", err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}
	var responsePayload apiResponse
	err = json.Unmarshal(b, &responsePayload)
	if err != nil {
		return nil, fmt.Errorf("json: %w, body:\n%s", err, b)
	}

	return &responsePayload, nil
}

func Translate(word string) ([]TermResult, error) {
	return defaultClient.Translate(word)
}

type clientOpts func(*Client)

func WithHttpClient(httpClient *http.Client) clientOpts {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithDictionary(dictionary string) clientOpts {
	return func(c *Client) {
		c.dictionary = dictionary
	}
}

type apiResponse struct {
	SearchedTerm                    string       `json:"SearchedTerm"`
	IsFound                         bool         `json:"IsFound"`
	AResults                        []TermResult `json:"AResults"`
	BResults                        []TermResult `json:"BResults"`
	AFullTextResults                []TermResult `json:"AFullTextResults"`
	BFullTextResults                []TermResult `json:"BFullTextResults"`
	AccentInsensitive               bool
	AvailabilityOnOtherDictionaries map[string]bool
	PrimeATerm                      string
	Suggestions                     []string
}

type TermResult struct {
	TermA         string `json:"TermA"`
	TermB         string `json:"TermB"`
	CategoryTextA string `json:"CategoryTextA"`
	CategoryTextB string `json:"CategoryTextB"`
	TermTypeTextA string `json:"TermTypeTextA"`
	TermTypeTextB string `json:"TermTypeTextB"`
	IsSlang       bool
}
