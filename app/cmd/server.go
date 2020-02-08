package cmd

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/Gklenskiy/vkdigest_bot/app/models"

	log "github.com/go-pkgz/lgr"
	_ "github.com/lib/pq"
)

// ServerCommand with params
type ServerCommand struct {
	Port        string `long:"port" env:"PORT" default:"5000" description:"port for listen"`
	RedirectURL string `long:"redirect" env:"REDIRECT_URL" required:"true" description:"url for redirect to bot"`

	VkAppID     string `long:"vk_app_id" env:"VK_APP_ID" required:"true" description:"Vk Application ID"`
	AuthURL     string `long:"auth_url" env:"AUTH_URL" required:"true" description:"Authentication URL"`
	VkAppSecret string `long:"vk_app_secret" env:"VK_APP_SECRET" required:"true" description:"url for redirect to bot"`

	CommonOpts
}

type authResponse struct {
	Token            string `json:"access_token"`
	Error            string `json:"error"`
	VkUserID         int    `json:"user_id"`
	ErrorDescription string `json:"error_description"`
	ExpiresInSec     int64  `json:"expires_in"`
}

var port string
var redirectURL string
var vkAppID string
var authURL string
var vkAppSecret string

// Execute is the entry point for "server" command, called by flag parser
func (serverCmd *ServerCommand) Execute(args []string) error {
	redirectURL = serverCmd.RedirectURL
	port = serverCmd.Port
	vkAppID = serverCmd.VkAppID
	authURL = serverCmd.AuthURL
	vkAppSecret = serverCmd.VkAppSecret
	log.Printf("[INFO] start web server %s", port)

	http.HandleFunc("/auth", redirectToBot)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("ListenAndServe error: %v", err)
	}

	return nil
}

func redirectToBot(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET params were: %s", r.URL.Query())
	code := r.URL.Query().Get("code")
	if code == "" {
		// return error
	}
	log.Printf("code: %s", code)

	state := r.URL.Query().Get("state")
	if state == "" {
		// return error
	}
	log.Printf("state: %s", state)

	userID, err := strconv.Atoi(state)
	if err != nil {
		log.Printf("%s is not an integer.", state)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	////
	httpClient := &http.Client{Timeout: 500 * time.Millisecond * time.Second}
	log.Printf("[DEBUG] Send request for auth")
	req, err := http.NewRequest("GET", "https://oauth.vk.com/access_token", nil)
	if err != nil {
		log.Printf("[ERROR] While Send request for auth: %s", err)
		return
	}

	q := req.URL.Query()
	q.Add("client_id", vkAppID)
	q.Add("client_secret", vkAppSecret)
	q.Add("redirect_uri", authURL)
	q.Add("code", code)

	req.URL.RawQuery = q.Encode()
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("[ERROR] failed on Get access token, %s", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] failed on read response body, %s", err)
		return
	}

	var result authResponse
	log.Printf("[DEBUG] Unmurshal: %s", body)
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Printf("[ERROR] failed to unmarshal response, %s", err)
		return
	}

	if result.Error != "" {
		log.Printf("[ERROR] failed to get access token, Error: %s Description: %s", result.Error, result.ErrorDescription)
		return
	}

	token := result.Token
	////
	err = models.CreateOrUpdate(userID, token)
	if err != nil {
		log.Printf("[ERROR] While save user token: %s", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}
