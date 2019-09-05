package utils

import (
	"MuseBot/models"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

var ErrNoRedirect = errors.New("no redirect")
var ErrRedirect = errors.New("redirect exist")

func CheckUrlExist(url string) bool {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	if resp.StatusCode >= 400 {
		return false
	}
	return true
}

func GetRedirectUrl(url string) (string, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", err
	}
	client := new(http.Client)
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return ErrRedirect
	}

	response, err := client.Do(req)
	if err != nil && response != nil && (response.StatusCode == http.StatusFound || response.StatusCode == http.StatusMovedPermanently) {
		if redUrl, err := response.Location(); err != nil {
			return "", err
		} else {
			return redUrl.String(), nil
		}
	} else if err != nil {
		return "", err
	} else {
		return "", ErrNoRedirect
	}
}

func TravisGetPubKey(url string) (string, error) {
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	config := models.TravisConfig{}
	if err := json.Unmarshal(data, &config); err != nil {
		return "", err
	}
	return config.Config.Notifications.Webhook.PublicKey, nil
}
