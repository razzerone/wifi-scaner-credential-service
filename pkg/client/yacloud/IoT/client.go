package IoT

import (
	"fmt"
	"io"
	"net/http"
	"time"
	"wifi-scaner-credentials/pkg/client/yacloud/autorisation"
)

type Client struct {
	Client     *http.Client
	RegisterID string
	key        *autorisation.Key
}

func (c *Client) MakeRequest(req *http.Request) (*http.Response, error) {
	iam := autorisation.GetIAMToken(autorisation.SignedToken(c.key))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", iam))

	res, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func NewClient(key *autorisation.Key, registerID string) *Client {
	return &Client{
		&http.Client{
			Timeout: 15 * time.Second,
		},
		registerID,
		key,
	}
}

func ParseResponse(res *http.Response) ([]byte, error) {
	resBody, err := io.ReadAll(res.Body)
	debug := string(resBody)

	_ = debug

	if err != nil {
		return nil, err
	}

	return resBody, nil
}
