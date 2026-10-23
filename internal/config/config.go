package config

import (
	"fmt"
	"strconv"

	"github.com/ciscoecosystem/mso-go-client/client"
)

// Config is the provider runtime configuration shared by the SDKv2 and
// Framework provider implementations.
type Config struct {
	Username   string
	Password   string
	IsInsecure bool
	ProxyUrl   string
	URL        string
	Domain     string
	Platform   string
	MaxRetries int
}

func BuildConfig(username, password, url string, isInsecure bool, proxyURL, domain, platform, retries string) (Config, error) {
	maxRetries := 2
	if retries != "" {
		var err error
		maxRetries, err = strconv.Atoi(retries)
		if err != nil {
			return Config{}, fmt.Errorf("invalid value for retries")
		}
	}

	if platform != "" && platform != "mso" && platform != "nd" {
		return Config{}, fmt.Errorf("invalid value for platform: must be one of mso or nd")
	}

	config := Config{
		Username:   username,
		Password:   password,
		URL:        url,
		IsInsecure: isInsecure,
		ProxyUrl:   proxyURL,
		Domain:     domain,
		Platform:   platform,
		MaxRetries: maxRetries,
	}

	if err := config.Valid(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) Valid() error {
	if c.Username == "" {
		return fmt.Errorf("Username must be provided for the MSO provider")
	}

	if c.Password == "" {
		return fmt.Errorf("Password must be provided for the MSO provider")
	}

	if c.URL == "" {
		return fmt.Errorf("URL must be provided for MSO provider")
	}

	return nil
}

func (c Config) GetClient() any {
	if c.Password == "" {
		return nil
	}

	return client.GetClient(c.URL, c.Username, client.Password(c.Password), client.Insecure(c.IsInsecure), client.ProxyUrl(c.ProxyUrl), client.Domain(c.Domain), client.Platform(c.Platform), client.MaxRetries(c.MaxRetries))
}
