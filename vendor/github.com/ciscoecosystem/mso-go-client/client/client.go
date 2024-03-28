package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/ciscoecosystem/mso-go-client/container"
	"github.com/ciscoecosystem/mso-go-client/models"
	"github.com/hashicorp/go-version"
)

const msoAuthPayload = `{
	"username": "%s",
	"password": "%s"
}`

const ndAuthPayload = `{
	"userName": "%s",
	"userPasswd": "%s"
}`

// Client is the main entry point
type Client struct {
	BaseURL            *url.URL
	httpClient         *http.Client
	AuthToken          *Auth
	Mutex              sync.Mutex
	username           string
	password           string
	insecure           bool
	proxyUrl           string
	domain             string
	platform           string
	version            string
	skipLoggingPayload bool
}

// singleton implementation of a client
var clientImpl *Client

type Option func(*Client)

func Insecure(insecure bool) Option {
	return func(client *Client) {
		client.insecure = insecure
	}
}

func Password(password string) Option {
	return func(client *Client) {
		client.password = password
	}
}

func ProxyUrl(pUrl string) Option {
	return func(client *Client) {
		client.proxyUrl = pUrl
	}
}

func Domain(domain string) Option {
	return func(client *Client) {
		client.domain = domain
	}
}

func Platform(platform string) Option {
	return func(client *Client) {
		client.platform = platform
	}
}

func Version(version string) Option {
	return func(client *Client) {
		client.version = version
	}
}

func SkipLoggingPayload(skipLoggingPayload bool) Option {
	return func(client *Client) {
		client.skipLoggingPayload = skipLoggingPayload
	}
}

func initClient(clientUrl, username string, options ...Option) *Client {
	var transport *http.Transport
	bUrl, err := url.Parse(clientUrl)
	if err != nil {
		// cannot move forward if url is undefined
		log.Fatal(err)
	}
	client := &Client{
		BaseURL:    bUrl,
		username:   username,
		httpClient: http.DefaultClient,
	}

	for _, option := range options {
		option(client)
	}

	transport = client.useInsecureHTTPClient(client.insecure)
	if client.proxyUrl != "" {
		transport = client.configProxy(transport)
	}

	client.httpClient = &http.Client{
		Transport: transport,
	}

	return client
}

// GetClient returns a singleton
func GetClient(clientUrl, username string, options ...Option) *Client {
	if clientImpl == nil {
		clientImpl = initClient(clientUrl, username, options...)
	}
	return clientImpl
}

func (c *Client) configProxy(transport *http.Transport) *http.Transport {
	pUrl, err := url.Parse(c.proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	transport.Proxy = http.ProxyURL(pUrl)
	return transport
}

func (c *Client) useInsecureHTTPClient(insecure bool) *http.Transport {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       insecure,
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS13,
		},
	}

	return transport
}

func (c *Client) MakeRestRequest(method string, path string, body *container.Container, authenticated bool) (*http.Request, error) {
	if c.platform == "nd" && path != "/login" {
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}
		path = fmt.Sprintf("mso/%v", path)
	}
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	if method == "PATCH" {
		validateString := url.Query()
		validateString.Set("validate", "false")
		url.RawQuery = validateString.Encode()
	}
	fURL := c.BaseURL.ResolveReference(url)
	var req *http.Request
	if method == "GET" || method == "DELETE" {
		req, err = http.NewRequest(method, fURL.String(), nil)
	} else {
		req, err = http.NewRequest(method, fURL.String(), bytes.NewBuffer((body.Bytes())))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	log.Printf("[DEBUG] HTTP request %s %s", method, path)

	if authenticated {

		req, err = c.InjectAuthenticationHeader(req, path)
		if err != nil {
			return req, err
		}
	}
	log.Printf("[DEBUG] HTTP request after injection %s %s", method, path)

	return req, nil
}

// Authenticate is used to
func (c *Client) Authenticate() error {
	method := "POST"
	path := "/api/v1/auth/login"
	var authPayload string

	if c.platform == "nd" {
		authPayload = ndAuthPayload
		if c.domain == "" {
			c.domain = "DefaultAuth"
		}
		path = "/login"
	} else {
		authPayload = msoAuthPayload
	}
	body, err := container.ParseJSON([]byte(fmt.Sprintf(authPayload, c.username, c.password)))
	if err != nil {
		return err
	}

	if c.domain != "" {
		if c.platform == "nd" {
			body.Set(c.domain, "domain")
		} else {
			domainId, err := c.GetDomainId(c.domain)
			if err != nil {
				return err
			}
			body.Set(domainId, "domainId")
		}
	}

	c.skipLoggingPayload = true

	req, err := c.MakeRestRequest(method, path, body, false)
	if err != nil {
		return err
	}

	obj, _, err := c.Do(req)
	c.skipLoggingPayload = false
	if err != nil {
		return err
	}

	if obj == nil {
		return errors.New("Empty response")
	}
	req.Header.Set("Content-Type", "application/json")

	token := models.StripQuotes(obj.S("token").String())

	if token == "" || token == "{}" {
		return errors.New("Invalid Username or Password")
	}

	if c.AuthToken == nil {
		c.AuthToken = &Auth{}
	}
	c.AuthToken.Token = stripQuotes(token)
	c.AuthToken.CalculateExpiry(1200) //refreshTime=1200 Sec

	return nil
}

func (c *Client) GetDomainId(domain string) (string, error) {
	req, err := c.MakeRestRequest("GET", "/api/v1/auth/login-domains", nil, false)
	if err != nil {
		return "", err
	}

	obj, _, err := c.Do(req)

	if err != nil {
		return "", err
	}
	err = CheckForErrors(obj, "GET")
	if err != nil {
		return "", err
	}
	count, err := obj.ArrayCount("domains")
	if err != nil {
		return "", err
	}

	for i := 0; i < count; i++ {
		domainCont, err := obj.ArrayElement(i, "domains")
		if err != nil {
			return "", err
		}
		domainName := models.StripQuotes(domainCont.S("name").String())

		if domainName == domain {
			return models.StripQuotes(domainCont.S("id").String()), nil
		}
	}
	return "", fmt.Errorf("Unable to find domain id for domain %s", domain)
}

func (c *Client) GetVersion() (string, error) {
	req, err := c.MakeRestRequest("GET", "/api/v1/platform/version", nil, true)
	if err != nil {
		return "unknown", err
	}

	obj, _, err := c.Do(req)
	if err != nil {
		return "unknown", err
	}

	err = CheckForErrors(obj, "GET")
	if err != nil {
		return "unknown", err
	}

	version := stripQuotes(obj.Search("version").String())
	if version == "" {
		return "unknown", fmt.Errorf("Unable to identify version")
	}
	c.version = version
	return version, nil
}

// Compares the version to the retrieved version.
// This returns -1, 0, or 1 if this version is smaller, equal, or larger than the retrieved version, respectively.
func (c *Client) CompareVersion(v string) (int, error) {
	if c.version == "" {
		c.GetVersion()
	}
	if c.version == "unknown" {
		return 0, fmt.Errorf("Could not retrieve version")
	}

	v1, err := version.NewVersion(c.version)
	if err != nil {
		return 0, fmt.Errorf("Could not parse retrieved version")
	}
	v2, err := version.NewVersion(v)
	if err != nil {
		return 0, fmt.Errorf("Could not parse version")
	}

	return v2.Compare(v1), nil
}

func StrtoInt(s string, startIndex int, bitSize int) (int64, error) {
	return strconv.ParseInt(s, startIndex, bitSize)
}

func (c *Client) Do(req *http.Request) (*container.Container, *http.Response, error) {
	log.Printf("[DEBUG] Begining DO method %s", req.URL.String())
	log.Printf("[TRACE] HTTP Request Method and URL: %s %s", req.Method, req.URL.String())
	if !c.skipLoggingPayload {
		log.Printf("[TRACE] HTTP Request Body: %v", req.Body)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[DEBUG] HTTP Request: %s %s", req.Method, req.URL.String())
	log.Printf("[DEBUG] HTTP Response: %d %s %v", resp.StatusCode, resp.Status, resp)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	resp.Body.Close()
	log.Printf("[DEBUG] HTTP response unique string %s %s %s", req.Method, req.URL.String(), bodyStr)
	if req.Method != "DELETE" && resp.StatusCode != 204 {
		obj, err := container.ParseJSON(bodyBytes)

		if err != nil {
			log.Printf("Error occured while json parsing %+v", err)
			return nil, resp, err
		}
		log.Printf("[DEBUG] Exit from do method")
		return obj, resp, err
	} else if resp.StatusCode == 204 {
		return nil, nil, nil
	} else {
		return nil, resp, err
	}
}

func stripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}
