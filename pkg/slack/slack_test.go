package slack_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/slack"
)

func TestClient_GetUserList(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
			"ok": true,
			"members": [
				{
					"id": "UE88NSJRW",
					"team_id": "TE88NSHUG",
					"name": "tong",
					"deleted": false,
					"color": "9f69e7",
					"real_name": "GuutonG",
					"tz": "Asia/Bangkok",
					"tz_label": "Indochina Time",
					"tz_offset": 25200,
					"profile": {
						"title": "",
						"phone": "",
						"skype": "",
						"real_name": "GuutonG",
						"real_name_normalized": "GuutonG",
						"display_name": "",
						"display_name_normalized": "",
						"status_text": "",
						"status_emoji": "",
						"status_expiration": 0,
						"avatar_hash": "g527f563bed3",
						"email": "tong@odds.team",
						"first_name": "GuutonG",
						"last_name": "",
						"image_24": "24_0000-24.png",
						"image_32": "32_0000-32.png",
						"image_48": "48_0000-48.png",
						"image_72": "72_0000-72.png",
						"image_192": "192_0000-192.png",
						"image_512": "512_0000-512.png",
						"status_text_canonical": "",
						"team": "TE88NSHUG"
					},
					"is_admin": true,
					"is_owner": true,
					"is_primary_owner": true,
					"is_restricted": false,
					"is_ultra_restricted": false,
					"is_bot": false,
					"is_app_user": false,
					"updated": 1542855855,
					"has_2fa": false
				}
			],
			"cache_ts": 1542875779
		}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	c := slack.Client{
		Token:   "xx-token-test",
		BaseURL: server.URL,
	}
	users, err := c.GetUserList()
	if err != nil {
		t.Fatalf("err = %v; want nil", err)
	}
	if users.Ok != true {
		t.Fatalf("err = %v; want nil", users.Ok)
	}
}
func TestClient_GetUserListInvalidBaseUrl(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
			"ok": true,
			"members": [
				{
					"id": "UE88NSJRW",
					"team_id": "TE88NSHUG",
					"name": "tong",
					"deleted": false,
					"color": "9f69e7",
					"real_name": "GuutonG",
					"tz": "Asia/Bangkok",
					"tz_label": "Indochina Time",
					"tz_offset": 25200,
					"profile": {
						"title": "",
						"phone": "",
						"skype": "",
						"real_name": "GuutonG",
						"real_name_normalized": "GuutonG",
						"display_name": "",
						"display_name_normalized": "",
						"status_text": "",
						"status_emoji": "",
						"status_expiration": 0,
						"avatar_hash": "g527f563bed3",
						"email": "tong@odds.team",
						"first_name": "GuutonG",
						"last_name": "",
						"image_24": "24_0000-24.png",
						"image_32": "32_0000-32.png",
						"image_48": "48_0000-48.png",
						"image_72": "72_0000-72.png",
						"image_192": "192_0000-192.png",
						"image_512": "512_0000-512.png",
						"status_text_canonical": "",
						"team": "TE88NSHUG"
					},
					"is_admin": true,
					"is_owner": true,
					"is_primary_owner": true,
					"is_restricted": false,
					"is_ultra_restricted": false,
					"is_bot": false,
					"is_app_user": false,
					"updated": 1542855855,
					"has_2fa": false
				}
			],
			"cache_ts": 1542875779
		}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	c := slack.Client{
		Token:   "xx-token-test",
		BaseURL: "",
	}
	_, err := c.GetUserList()
	if err != nil {
		t.Fatalf("err = %v; want nil", err)
	}
}
func TestClient_GetUserListInvalidBody(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, nil)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	c := slack.Client{
		Token:   "xx-token-test",
		BaseURL: server.URL,
	}
	_, err := c.GetUserList()
	if err == nil {
		t.Fatalf("err = %v; want nil", err)
	}
}
func TestClient_GetUserListInvalidAuth(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
			"ok": false,
			"error": "invalid_auth"
			}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	c := slack.Client{
		Token:   "xx-token-test",
		BaseURL: server.URL,
	}
	users, err := c.GetUserList()
	if err != nil {
		t.Fatalf("err = %v; want nil", err)
	}
	if users.Ok != false {
		t.Fatalf("err = %v; want fase", users.Ok)
	}
}
func TestClient_OpenIMChannel(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
			"ok": true,
			"no_op": true,
			"already_open": true,
			"channel": {
			"id": "DEACQJGLW"
			}
		}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	c := slack.Client{
		Token:   "xx-token-test",
		BaseURL: server.URL,
	}
	users, err := c.OpenIMChannel("UE88NSJRW")
	if err != nil {
		t.Fatalf("err = %v; want nil", err)
	}
	if users.Ok != true {
		t.Fatalf("err = %v; want nil", users.Ok)
	}
}
func TestClient_PostMessage(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		fmt.Fprint(w, `{
			"ok": true,
			"channel": "CE8CK8KUH",
			"ts": "1543112823.000100",
			"message": {
				"text": "Hi @here",
				"username": "Slack API Tester",
				"bot_id": "BE8TMD34K",
				"type": "message",
				"subtype": "bot_message",
				"ts": "1543112823.000100"
			}
		}`)
	})

	server := httptest.NewServer(mux)
	defer server.Close()
	c := slack.Client{
		Token:   "xx-token-test",
		BaseURL: server.URL,
	}
	users, err := c.PostMessage("CE8CK8KUH", "Test message")
	if err != nil {
		t.Fatalf("err = %v; want nil", err)
	}
	if users.Ok != true {
		t.Fatalf("err = %v; want nil", users.Ok)
	}
	if users.Channel != "CE8CK8KUH" {
		t.Fatalf("err = %v; want CE8CK8KUH", users.Channel)
	}
}
