package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func mockClient(server *httptest.Server) (*Client, error) {
	client, err := NewClient("testtoken")
	if err != nil {
		return nil, err
	}
	return client, nil
}

func TestNewClient(t *testing.T) {
	token := "testtoken"
	client, err := NewClient(token)
	if err != nil {
		t.Fatal(err)
	}
	if client.Token != token {
		t.Errorf("Token %s; want %s", client.Token, token)
	}
	if client.httpClient != http.DefaultClient {
		t.Errorf("httpClient %p; want %p", client.httpClient, http.DefaultClient)
	}
}

func TestNewClientMissingToken(t *testing.T) {
	token := ""
	_, err := NewClient(token)
	if err == nil {
		t.Errorf("Token %s; want %s", err, "missing token")
	}
}

// func TestGetProfile(t *testing.T) {
// 	type want struct {
// 		URLPath     string
// 		RequestBody []byte
// 		Response    *GetUsersListResponse
// 		Error       error
// 	}
// 	var testCases = []struct {
// 		ResponseCode int
// 		Response     []byte
// 		Want         want
// 	}{
// 		{
// 			ResponseCode: 200,
// 			Response: []byte(`{
// 				"ok": true,
// 				"members": [
// 					{
// 						"id": "UE88NSJRW",
// 						"team_id": "TE88NSHUG",
// 						"name": "tong",
// 						"deleted": false,
// 						"color": "9f69e7",
// 						"real_name": "GuutonG",
// 						"tz": "Asia/Bangkok",
// 						"tz_label": "Indochina Time",
// 						"tz_offset": 25200,
// 						"profile": {
// 							"title": "",
// 							"phone": "",
// 							"skype": "",
// 							"real_name": "GuutonG",
// 							"real_name_normalized": "GuutonG",
// 							"display_name": "",
// 							"display_name_normalized": "",
// 							"status_text": "",
// 							"status_emoji": "",
// 							"status_expiration": 0,
// 							"avatar_hash": "g527f563bed3",
// 							"email": "tong@odds.team",
// 							"image_24": "0000-24.png",
// 							"image_32": "0000-32.png",
// 							"image_48": "0000-48.png",
// 							"image_72": "0000-72.png",
// 							"image_192": "0000-192.png",
// 							"image_512": "0000-512.png",
// 							"status_text_canonical": "",
// 							"team": "TE88NSHUG"
// 						},
// 						"is_admin": true,
// 						"is_owner": true,
// 						"is_primary_owner": true,
// 						"is_restricted": false,
// 						"is_ultra_restricted": false,
// 						"is_bot": false,
// 						"is_app_user": false,
// 						"updated": 1542855855,
// 						"has_2fa": false
// 					}
// 				],
// 				"cache_ts": 1542875779
// 			}`),
// 			Want: want{
// 				URLPath:     fmt.Sprintf("https://slack.com/api/users.list?token=testtoken"),
// 				RequestBody: []byte(""),
// 				Response: &GetUsersListResponse{
// 					Ok: true,
// 					Members: []Member{
// 						Member{
// 							ID:       "UE88NSJRW",
// 							TeamID:   "TE88NSHUG",
// 							Name:     "tong",
// 							Deleted:  false,
// 							Color:    "9f69e7",
// 							RealName: "GuutonG",
// 							Tz:       "Asia/Bangkok",
// 							TzLabel:  "Indochina Time",
// 							TzOffset: 25200,
// 							Profile: Profile{
// 								Title:                 "",
// 								Phone:                 "",
// 								Skype:                 "",
// 								RealName:              "GuutonG",
// 								RealNameNormalized:    "GuutonG",
// 								DisplayName:           "",
// 								DisplayNameNormalized: "",
// 								StatusText:            "",
// 								StatusEmoji:           "",
// 								StatusExpiration:      0,
// 								AvatarHash:            "g527f563bed3",
// 								Email:                 "tong@odds.team",
// 								Image24:               "0000-24.png",
// 								Image32:               "0000-32.png",
// 								Image48:               "0000-48.png",
// 								Image72:               "0000-72.png",
// 								Image192:              "0000-192.png",
// 								Image512:              "0000-512.png",
// 								StatusTextCanonical:   "",
// 								Team:                  "TE88NSHUG",
// 							},
// 							IsAdmin:           true,
// 							IsOwner:           true,
// 							IsPrimaryOwner:    true,
// 							IsRestricted:      false,
// 							IsUltraRestricted: false,
// 							IsBot:             false,
// 							IsAppUser:         false,
// 							Updated:           1542855855,
// 							Has2Fa:            false,
// 						},
// 					},
// 					CacheTs: 1542875779,
// 				},
// 			},
// 		},
// 	}

// 	var currentTestIdx int
// 	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		defer r.Body.Close()
// 		tc := testCases[currentTestIdx]
// 		if r.Method != http.MethodGet {
// 			t.Errorf("Method %s; want %s", r.Method, http.MethodGet)
// 		}
// 		if r.URL.Path != tc.Want.URLPath {
// 			t.Errorf("URLPath %s; want %s", r.URL.Path, tc.Want.URLPath)
// 		}
// 		body, err := ioutil.ReadAll(r.Body)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		if !reflect.DeepEqual(body, tc.Want.RequestBody) {
// 			t.Errorf("RequestBody %s; want %s", body, tc.Want.RequestBody)
// 		}
// 		w.WriteHeader(tc.ResponseCode)
// 		w.Write(tc.Response)
// 	}))
// 	defer server.Close()
// 	client, err := mockClient(server)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	for i, tc := range testCases {
// 		currentTestIdx = i
// 		res, err := client.GetUsersList()
// 		if tc.Want.Error != nil {
// 			if !reflect.DeepEqual(err, tc.Want.Error) {
// 				t.Errorf("Error %d %v; want %v", i, err, tc.Want.Error)
// 			}
// 		} else {
// 			if err != nil {
// 				t.Error(err)
// 			}
// 		}
// 		if !reflect.DeepEqual(res, tc.Want.Response) {
// 			t.Errorf("Response %d %v; want %v", i, res, tc.Want.Response)
// 		}
// 	}
// }
