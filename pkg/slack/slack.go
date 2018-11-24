package slack

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURL string = "https://slack.com/api"

// Client type
type Client struct {
	Token      string
	httpClient *http.Client // default http.DefaultClient
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

// NewClient returns a new slack client instance.
func NewClient(token string) (*Client, error) {
	if token == "" {
		return nil, errors.New("missing token")
	}

	return &Client{
		Token:      token,
		httpClient: http.DefaultClient,
	}, nil
}

// doRequest makes a request to the API
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	c.httpClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if 200 != resp.StatusCode {
		return nil, fmt.Errorf("%s", body)
	}
	return body, nil
}

// PostMessage API `chat.postMessage`: Sends a message to a channel.
func (c *Client) PostMessage(channelID string, text string) (*PostMessageResponse, error) {
	textEncode := &url.URL{Path: text}
	url := baseURL + "/chat.postMessage?token=" + c.Token + "&channel=" + channelID + "&text=" + textEncode.String()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data PostMessageResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// OpenIMChannel API `im.open`: Opens a direct message channel.
func (c *Client) OpenIMChannel(userID string) (*OpenIMChannelResponse, error) {
	url := baseURL + "/im.open?token=" + c.Token + "&user=" + userID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data OpenIMChannelResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

// GetUsersList API `users.list`: Lists all users in a Slack team.
func (c *Client) GetUsersList() (*GetUsersListResponse, error) {
	url := baseURL + "/users.list?token=" + c.Token
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data GetUsersListResponse
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
