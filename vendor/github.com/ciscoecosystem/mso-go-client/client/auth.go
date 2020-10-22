package client

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Auth struct {
	Token  string
	Expiry time.Time
}

func (au *Auth) IsValid() bool {
	if au.Token != "" && au.Expiry.Unix() > au.estimateExpireTime() {
		return true
	}
	return false
}

func (t *Auth) CalculateExpiry(willExpire int64) {
	t.Expiry = time.Unix((time.Now().Unix() + willExpire), 0)
}

func (t *Auth) estimateExpireTime() int64 {
	return time.Now().Unix() + 3
}

func (client *Client) InjectAuthenticationHeader(req *http.Request, path string) (*http.Request, error) {
	log.Printf("[DEBUG] Begin Injection")
	if client.AuthToken == nil || !client.AuthToken.IsValid() {

		err := client.Authenticate()

		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.AuthToken.Token))

		return req, nil
	}

	return req, nil
}
