package tonicpow

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is the TonicPow client/configuration
type Client struct {
	httpClient *resty.Client
	options    *clientOptions // Options are all the default settings / configuration
}

// clientOptions holds all the configuration for client requests and default resources
type clientOptions struct {
	apiKey         string              // API key
	apiURL         string              // API endpoint (URL by environment)
	environment    string              // Name of the current environment
	customHeaders  map[string][]string // Custom headers on outgoing requests
	httpTimeout    time.Duration       // Default timeout in seconds for GET requests
	requestTracing bool                // If enabled, it will trace the request timing
	retryCount     int                 // Default retry count for HTTP requests
	userAgent      string              // User agent for all outgoing requests
}

// StandardResponse is the standard fields returned on all responses
type StandardResponse struct {
	Body       []byte          `json:"-"` // Body of the response request
	Error      *Error          `json:"-"` // API error response
	StatusCode int             `json:"-"` // Status code returned on the request
	Tracing    resty.TraceInfo `json:"-"` // Trace information if enabled on the request
}

// ClientOps allow functional options to be supplied
// that overwrite default client options.
type ClientOps func(c *clientOptions)

// WithEnvironment will change the environment
func WithEnvironment(e environment) ClientOps {
	return func(c *clientOptions) {
		c.apiURL = e.apiURL
		c.environment = e.name
	}
}

// WithEnvironmentString will change the environment
func WithEnvironmentString(e string) ClientOps {
	e = strings.ToLower(strings.TrimSpace(e))
	if e == environmentStagingName || e == environmentStagingAlias {
		return WithEnvironment(EnvironmentStaging)
	} else if e == environmentDevelopmentName || e == environmentDevelopmentAlias {
		return WithEnvironment(EnvironmentDevelopment)
	}
	return WithEnvironment(EnvironmentLive)
}

// WithHTTPTimeout can be supplied to adjust the default http client timeouts.
// The http client is used when creating requests
// Default timeout is 10 seconds.
func WithHTTPTimeout(timeout time.Duration) ClientOps {
	return func(c *clientOptions) {
		c.httpTimeout = timeout
	}
}

// WithRequestTracing will enable tracing.
// Tracing is disabled by default.
func WithRequestTracing() ClientOps {
	return func(c *clientOptions) {
		c.requestTracing = true
	}
}

// WithRetryCount will overwrite the default retry count for http requests.
// Default retries is 2.
func WithRetryCount(retries int) ClientOps {
	return func(c *clientOptions) {
		c.retryCount = retries
	}
}

// WithUserAgent will overwrite the default useragent.
// Default is package name + version.
func WithUserAgent(userAgent string) ClientOps {
	return func(c *clientOptions) {
		c.userAgent = userAgent
	}
}

// WithAPIKey provides the API key
func WithAPIKey(appAPIKey string) ClientOps {
	return func(c *clientOptions) {
		c.apiKey = appAPIKey
	}
}

// WithCustomHeaders will add custom headers to outgoing requests
// Custom headers is empty by default
func WithCustomHeaders(headers map[string][]string) ClientOps {
	return func(c *clientOptions) {
		c.customHeaders = headers
	}
}

// WithCustomHTTPClient will overwrite the default client with a custom client.
func (c *Client) WithCustomHTTPClient(client *resty.Client) *Client {
	c.httpClient = client
	return c
}

// GetUserAgent will return the user agent string of the client
func (c *Client) GetUserAgent() string {
	return c.options.userAgent
}

// GetEnvironment will return the environment of the client
func (c *Client) GetEnvironment() string {
	return c.options.environment
}

// defaultClientOptions will return an Options struct with the default settings
//
// Useful for starting with the default and then modifying as needed
func defaultClientOptions() (opts *clientOptions) {
	// Set the default options
	opts = &clientOptions{
		apiURL:         EnvironmentLive.apiURL,
		environment:    EnvironmentLive.name,
		httpTimeout:    defaultHTTPTimeout,
		requestTracing: false,
		retryCount:     defaultRetryCount,
		userAgent:      defaultUserAgent,
	}
	return
}

// NewClient creates a new client for all TonicPow requests
//
// If no options are given, it will use the DefaultClientOptions()
// If no client is supplied it will use a default Resty HTTP client
func NewClient(opts ...ClientOps) (*Client, error) {
	defaults := defaultClientOptions()

	// Create a new client
	client := &Client{
		options: defaults,
	}
	// overwrite defaults with any set by user
	for _, opt := range opts {
		opt(client.options)
	}
	// Check for API key
	if client.options.apiKey == "" {
		return nil, errors.New("missing an API Key")
	}
	// Set the Resty HTTP client
	if client.httpClient == nil {
		client.httpClient = resty.New()
		// Set defaults (for GET requests)
		client.httpClient.SetTimeout(client.options.httpTimeout)
		client.httpClient.SetRetryCount(client.options.retryCount)
	}
	return client, nil
}

// Request is a standard GET / POST / PUT / DELETE request for all outgoing HTTP requests
// Omit the data attribute if using a GET request
func (c *Client) Request(httpMethod string, requestEndpoint string,
	data interface{}, expectedCode int) (response StandardResponse, err error) {

	// Set the user agent
	req := c.httpClient.R().SetHeader("User-Agent", c.options.userAgent)

	// Set the body if (PUT || POST || PATCH)
	if httpMethod != http.MethodGet && httpMethod != http.MethodDelete {
		var j []byte
		if j, err = json.Marshal(data); err != nil {
			return
		}
		req = req.SetBody(string(j))
		req.Header.Add("Content-Length", strconv.Itoa(len(j)))
		req.Header.Set("Content-Type", "application/json")
	}

	// Enable tracing
	if c.options.requestTracing {
		req.EnableTrace()
	}

	// Set the authorization and content type
	req.Header.Set(fieldAPIKey, c.options.apiKey)

	// Custom headers?
	for key, headers := range c.options.customHeaders {
		for _, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// Fire the request
	var resp *resty.Response
	switch httpMethod {
	case http.MethodPost:
		if resp, err = req.Post(c.options.apiURL + requestEndpoint); err != nil {
			return
		}
	case http.MethodPut:
		if resp, err = req.Put(c.options.apiURL + requestEndpoint); err != nil {
			return
		}
	case http.MethodDelete:
		if resp, err = req.Delete(c.options.apiURL + requestEndpoint); err != nil {
			return
		}
	case http.MethodGet:
		if resp, err = req.Get(c.options.apiURL + requestEndpoint); err != nil {
			return
		}
	}

	// Tracing enabled?
	if c.options.requestTracing {
		response.Tracing = resp.Request.TraceInfo()
	}

	// Set the status code & body
	response.StatusCode = resp.StatusCode()
	response.Body = resp.Body()

	// Check expected code if set
	if expectedCode > 0 && response.StatusCode != expectedCode {
		response.Error = new(Error)
		if err = json.Unmarshal(response.Body, &response.Error); err != nil {
			return
		}
		err = fmt.Errorf("%s", response.Error.Message)
	}

	return
}
