package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type Google0AuthToken struct {
	AccessToken string `json:"access_token"`
	IdToken     string `json:"id_token"`
}

func GetGoogleOAuthToken(code string) (*Google0AuthToken, error) {
	const rootURI = "https://oauth2.googleapis.com/token"

	values := url.Values{}
	values.Add("grant_type", "authorization_code")
	values.Add("code", code)
	values.Add("client_id", os.Getenv("GOOGLE_OAUTH_CLIENT_ID"))
	values.Add("client_secret", os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"))
	values.Add("redirect_uri", os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"))

	query := values.Encode()

	req, err := http.NewRequest("POST", rootURI, strings.NewReader(query))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{
		Timeout: time.Second * 120,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get Google OAuth token")
	}

	var resBody bytes.Buffer

	_, err = resBody.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	var GoogleOAuthTokenRes map[string]interface{}

	err = json.Unmarshal(resBody.Bytes(), &GoogleOAuthTokenRes)
	if err != nil {
		return nil, err
	}

	tokenBody := Google0AuthToken{
		AccessToken: GoogleOAuthTokenRes["access_token"].(string),
		IdToken:     GoogleOAuthTokenRes["id_token"].(string),
	}

	return &tokenBody, nil
}

type GoogleUserInfo struct {
	Id 		string `json:"id"`
	Email 	string `json:"email"`
	Verified_email bool `json:"verified_email"`
	Name 	string `json:"name"`
	Given_name 	string `json:"given_name"`
	Family_name 	string `json:"family_name"`
	Picture 	string `json:"picture"`
	Locale 	string `json:"locale"`
}

func GetGoogleUserInfo(accessToken string, id_token string) (*GoogleUserInfo, error) {
	rootUrl := fmt.Sprintf("https://www.googleapis.com/oauth2/v3/userinfo?access_token=%s", accessToken)

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get Google user info")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, resp.Body)
	if err != nil {
		return nil, err
	}

	var GoogleUserInfoRes map[string]interface{}

	err = json.Unmarshal(resBody.Bytes(), &GoogleUserInfoRes)
	if err != nil {
		return nil, err
	}

	userBody := &GoogleUserInfo{
		Id: GoogleUserInfoRes["sub"].(string),
		Email: GoogleUserInfoRes["email"].(string),
		Verified_email: GoogleUserInfoRes["verified_email"].(bool),
		Name: GoogleUserInfoRes["name"].(string),
		Given_name: GoogleUserInfoRes["given_name"].(string),
		Family_name: GoogleUserInfoRes["family_name"].(string),
		Picture: GoogleUserInfoRes["picture"].(string),
		Locale: GoogleUserInfoRes["locale"].(string),
	}

	return userBody, nil
}