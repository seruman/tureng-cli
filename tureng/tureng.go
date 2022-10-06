package tureng

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
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
	debug      struct {
		enabled bool
		out     io.Writer
	}
}

func NewClient(opts ...clientOpts) *Client {
	c := &Client{
		httpClient: http.DefaultClient,
		dictionary: TurkishEnglish,
		debug: struct {
			enabled bool
			out     io.Writer
		}{
			enabled: false,
			out:     os.Stderr,
		},
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

	c.dumpRequest(req)

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

	c.dumpResponse(resp)

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

func (c *Client) dumpRequest(req *http.Request) {
	if c.debug.enabled {
		b, err := httputil.DumpRequestOut(req, true)
		if err != nil {
			fmt.Fprintf(c.debug.out, "dump request: %v", err)
			return
		}

		fmt.Fprintf(c.debug.out, "%s", b)
	}
}

func (c *Client) dumpResponse(resp *http.Response) {
	if c.debug.enabled {
		b, err := httputil.DumpResponse(resp, true)
		if err != nil {
			fmt.Fprintf(c.debug.out, "dump response: %v", err)
			return
		}

		fmt.Fprintf(c.debug.out, "%s", b)
	}
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

func WithDebug(enabled bool) clientOpts {
	return func(c *Client) {
		c.debug.enabled = enabled
	}
}

func WithDebugOutput(w io.Writer) clientOpts {
	return func(c *Client) {
		c.debug.out = w
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
