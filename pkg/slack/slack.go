package slack

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURL string = "https://slack.com/api"

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

func (s *Client) PostMessage(channelID string, text string) (*PostMessageResponse, error) {
	textEncode := &url.URL{Path: text}
	url := baseURL + "/chat.postMessage?token=" + s.Token + "&channel=" + channelID + "&text=" + textEncode.String()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
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

func (s *Client) OpenIMChannel(userID string) (*OpenIMChannelResponse, error) {
	url := baseURL + "/im.open?token=" + s.Token + "&user=" + userID
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
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

func (s *Client) GetUsersList() (*GetUsersList, error) {
	url := baseURL + "/users.list?token=" + s.Token
	fmt.Printf(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := s.doRequest(req)
	if err != nil {
		return nil, err
	}
	var data GetUsersList
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *Client) doRequest(req *http.Request) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
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

type GetUsersList struct {
	Ok      bool `json:"ok"`
	Members []struct {
		ID       string `json:"id"`
		TeamID   string `json:"team_id"`
		Name     string `json:"name"`
		Deleted  bool   `json:"deleted"`
		Color    string `json:"color"`
		RealName string `json:"real_name"`
		Tz       string `json:"tz"`
		TzLabel  string `json:"tz_label"`
		TzOffset int    `json:"tz_offset"`
		Profile  struct {
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
		} `json:"profile"`
		IsAdmin           bool `json:"is_admin"`
		IsOwner           bool `json:"is_owner"`
		IsPrimaryOwner    bool `json:"is_primary_owner"`
		IsRestricted      bool `json:"is_restricted"`
		IsUltraRestricted bool `json:"is_ultra_restricted"`
		IsBot             bool `json:"is_bot"`
		IsAppUser         bool `json:"is_app_user"`
		Updated           int  `json:"updated"`
		Has2Fa            bool `json:"has_2fa,omitempty"`
	} `json:"members"`
	CacheTs int `json:"cache_ts"`
}

type OpenIMChannelResponse struct {
	Ok          bool `json:"ok"`
	NoOp        bool `json:"no_op"`
	AlreadyOpen bool `json:"already_open"`
	Channel     struct {
		ID string `json:"id"`
	} `json:"channel"`
}

type PostMessageResponse struct {
	Ok      bool   `json:"ok"`
	Channel string `json:"channel"`
	Ts      string `json:"ts"`
	Message struct {
		Text     string `json:"text"`
		Username string `json:"username"`
		BotID    string `json:"bot_id"`
		Type     string `json:"type"`
		Subtype  string `json:"subtype"`
		Ts       string `json:"ts"`
	} `json:"message"`
}
