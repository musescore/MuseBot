package main

import (
	"MuseBot/models"
	"MuseBot/utils"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	return
}

func GitHubHookHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("token") == "" {
		log.Debugf("GitHub hook with empty token")
		w.WriteHeader(401)
		return
	}
	if r.URL.Query().Get("token") != Config.GitHubHookToken {
		log.Debugf("GitHub hook with wrong token")
		w.WriteHeader(401)
		return
	}
	event := r.Header.Get("X-GitHub-Event")
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Debugf("GitHub hook body read error: %s", err)
		w.WriteHeader(500)
		return
	}

	switch strings.ToLower(event) {
	case "pull_request":
		payload := models.GitHubPRPayload{}
		if err := json.Unmarshal(data, &payload); err != nil {
			log.Debugf("GitHub %s hook unmarshal error: %s", event, err)
			w.WriteHeader(500)
			return
		}
		BotMulticast <- GitHubHookPRMessage{Payload: payload}
	case "push":
		payload := models.GitHubPushPayload{}
		if err := json.Unmarshal(data, &payload); err != nil {
			log.Debugf("GitHub %s hook unmarshal error: %s", event, err)
			w.WriteHeader(500)
			return
		}
		BotMulticast <- GitHubHookPushMessage{Payload: payload}
	default:
		log.Debugf("GitHub hook with unknown event type: %s", event)
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)
	return
}

func TravisHookHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Signature") == "" {
		log.Debugf("Travis hook with empty Signature")
		w.WriteHeader(400)
		return
	}
	sign, err := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))
	if err != nil {
		log.Debugf("Travis hook with wrong Signature, error: %s", err)
		w.WriteHeader(400)
		return
	}
	data := []byte(r.FormValue("payload"))
	digest := utils.DataDigest(data)
	if err := rsa.VerifyPKCS1v15(Config.TravisPublicKey, crypto.SHA1, digest, sign); err != nil {
		log.Debugf("Travis hook sign error: %s", err)
		w.WriteHeader(401)
		return
	}

	payload := models.TravisHookPayload{}
	if err := json.Unmarshal(data, &payload); err != nil {
		log.Debugf("Travis hook unmarshal error: %s", err)
		w.WriteHeader(500)
		return
	}

	BotMulticast <- TravisHookMessage{Payload: payload}
	w.WriteHeader(200)
	return
}
