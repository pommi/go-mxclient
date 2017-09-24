package mxclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
)

type Client struct {
	Config  Config
	Session Session
}

type Config struct {
	Url       string `json:"url"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Anonymous bool   `json:"anonymous"`
	UserAgent string `json:"user_agent"`
}

func DefaultConfig() *Config {
	return &Config{
		Url:       "",
		Username:  "",
		Password:  "",
		Anonymous: false,
		UserAgent: "Go-MxClient/0.1",
	}
}

type Session struct {
	HttpClient *http.Client
	CsrfToken  string
}

func CreateSession() *Session {
	jar, _ := cookiejar.New(nil)
	return &Session{
		HttpClient: &http.Client{
			Jar: jar,
		},
		CsrfToken: "",
	}
}

func NewClient(config *Config) (client *Client, err error) {
	// bootstrap the config
	defConfig := DefaultConfig()

	if len(config.Url) == 0 {
		config.Url = defConfig.Url
	}

	if len(config.Username) == 0 {
		config.Username = defConfig.Username
	}

	if len(config.Password) == 0 {
		config.Password = defConfig.Password
	}

	if len(config.UserAgent) == 0 {
		config.UserAgent = defConfig.UserAgent
	}

	session := CreateSession()
	client = &Client{
		Config:  *config,
		Session: *session,
	}
	return client, nil
}

type RequestAction struct {
	Action  string                 `json:"action"`
	Params  map[string]interface{} `json:"params"`
	Context []string               `json:"context,omitempty"`
}

func (c *Client) Request(ra RequestAction) (map[string]interface{}, error) {
	ra_json, _ := json.Marshal(ra)
	// debug
	//fmt.Println("JSON:", string(ra_json))

	req, err := http.NewRequest("POST", c.Config.Url, bytes.NewBuffer(ra_json))
	req.Header.Set("User-Agent", c.Config.UserAgent)
	req.Header.Set("Content-Type", "application/json")

	request_id, _ := uuid.NewRandom()
	req.Header.Set("X-Mx-ReqToken", request_id.String())
	if c.Session.CsrfToken != "" {
		req.Header.Set("X-Csrf-Token", c.Session.CsrfToken)
	}

	resp, err := c.Session.HttpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("request:", resp.Request)
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	var response map[string]interface{}
	json.Unmarshal(body, &response)

	if csrf, ok := response["csrftoken"]; ok {
		c.Session.CsrfToken = csrf.(string)
	}

	if resp.StatusCode != 200 {
		return response, errors.New(fmt.Sprintf("Received HTTP response: %s", resp.Status))
	}
	return response, nil
}

func (c *Client) GetSessionData() (map[string]interface{}, error) {
	request := RequestAction{
		Action: "get_session_data",
		Params: map[string]interface{}{
			"profile":        "",
			"timezoneoffset": 0,
		},
	}
	return c.Request(request)
}
