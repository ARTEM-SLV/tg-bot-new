package telgram

import (
	"TG_bot_new/lib/e"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Client struct {
	host     string
	basePath string
	client   http.Client
}

const (
	getUpdatesMethod  = "getUpdate"
	sendMessageMethod = "sendMessage"
)

func New(host string, token string) Client {
	return Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

func newBasePath(token string) string {
	return "bot" + token
}

func (c *Client) Updates(offset int, limit int) ([]Update, error) {
	q := url.Values{}
	q.Set("offset", strconv.Itoa(offset))
	q.Set("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) SendMessage(chatID int, message string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", message)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	const errMsg = "can't do request"

	u := url.URL{
		Scheme: c.host,
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap(errMsg, err)
	}

	return body, nil
}
