package slack

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var DefaultBaseURL = "https://slack.com/api"

type ClientInterface interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	Token      string
	BaseURL    string
	HttpClient ClientInterface
}

// Profile type for Member
type Profile struct {
	Title                 string `json:"title"`
	Phone                 string `json:"phone"`
	Skype                 string `json:"skype"`
	RealName              string `json:"real_name"`
	RealNameNormalized    string `json:"real_name_normalized"`
	DisplayName           string `json:"display_name"`
	DisplayNameNormalized string `json:"display_name_normalized"`
	StatusText            string `json:"status_text"`
	StatusEmoji           string `json:"status_emoji"`
	StatusExpiration      int    `json:"status_expiration"`
	AvatarHash            string `json:"avatar_hash"`
	Email                 string `json:"email"`
	Image24               string `json:"image_24"`
	Image32               string `json:"image_32"`
	Image48               string `json:"image_48"`
	Image72               string `json:"image_72"`
	Image192              string `json:"image_192"`
	Image512              string `json:"image_512"`
	StatusTextCanonical   string `json:"status_text_canonical"`
	Team                  string `json:"team"`
}

// Member type for GetUsersListResponse
type Member struct {
	ID                string  `json:"id"`
	TeamID            string  `json:"team_id"`
	Name              string  `json:"name"`
	Deleted           bool    `json:"deleted"`
	Color             string  `json:"color"`
	RealName          string  `json:"real_name"`
	Tz                string  `json:"tz"`
	TzLabel           string  `json:"tz_label"`
	TzOffset          int     `json:"tz_offset"`
	Profile           Profile `json:"profile"`
	IsAdmin           bool    `json:"is_admin"`
	IsOwner           bool    `json:"is_owner"`
	IsPrimaryOwner    bool    `json:"is_primary_owner"`
	IsRestricted      bool    `json:"is_restricted"`
	IsUltraRestricted bool    `json:"is_ultra_restricted"`
	IsBot             bool    `json:"is_bot"`
	IsAppUser         bool    `json:"is_app_user"`
	Updated           int     `json:"updated"`
	Has2Fa            bool    `json:"has_2fa,omitempty"`
}

// GetUsersListResponse type for `users.list` api
type GetUsersListResponse struct {
	Ok      bool     `json:"ok"`
	Members []Member `json:"members"`
	CacheTs int      `json:"cache_ts"`
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	httpClient := c.HttpClient
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	if httpClient == nil {
		httpClient = &http.Client{Transport: tr}
	}

	return httpClient.Do(req)
}

func (c *Client) url(path string) string {
	if c.BaseURL == "" {
		c.BaseURL = DefaultBaseURL
	}
	return fmt.Sprintf("%s%s", c.BaseURL, path)
}

// GetUserList API `users.list`: Lists all users in a Slack team.
func (c *Client) GetUserList() (*GetUsersListResponse, error) {
	endpoint := c.url("/users.list")
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("token", c.Token)
	req.URL.RawQuery = q.Encode()
	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var usersListResponse GetUsersListResponse
	err = json.Unmarshal(body, &usersListResponse)
	if err != nil {
		return nil, err
	}
	return &usersListResponse, nil
}

// Channel type for OpenIMChannelResponse
type Channel struct {
	ID string `json:"id"`
}

// OpenIMChannelResponse type for `im.open` api
type OpenIMChannelResponse struct {
	Ok          bool    `json:"ok"`
	NoOp        bool    `json:"no_op"`
	AlreadyOpen bool    `json:"already_open"`
	Channel     Channel `json:"channel"`
}

// OpenIMChannel API `im.open`: Opens a direct message channel.
func (c *Client) OpenIMChannel(userID string) (*OpenIMChannelResponse, error) {
	endpoint := c.url("/im.open")
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("token", c.Token)
	q.Add("user", userID)
	req.URL.RawQuery = q.Encode()
	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var openIMChannelResponse OpenIMChannelResponse
	err = json.Unmarshal(body, &openIMChannelResponse)
	if err != nil {
		return nil, err
	}
	return &openIMChannelResponse, nil
}

// Message type for PostMessageResponse
type Message struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	BotID    string `json:"bot_id"`
	Type     string `json:"type"`
	Subtype  string `json:"subtype"`
	Ts       string `json:"ts"`
}

// PostMessageResponse type for `chat.postMessage` api
type PostMessageResponse struct {
	Ok      bool    `json:"ok"`
	Channel string  `json:"channel"`
	Ts      string  `json:"ts"`
	Message Message `json:"message"`
}

// PostMessage API `chat.postMessage`: Sends a message to a channel.
func (c *Client) PostMessage(channelID string, text string) (*PostMessageResponse, error) {
	endpoint := c.url("/chat.postMessage")
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	q := req.URL.Query()
	q.Add("token", c.Token)
	q.Add("channel", channelID)
	q.Add("text", text)
	req.URL.RawQuery = q.Encode()
	res, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var postMessageResponse PostMessageResponse
	err = json.Unmarshal(body, &postMessageResponse)
	if err != nil {
		return nil, err
	}
	return &postMessageResponse, nil
}
