package main

import (
	"MuseBot/models"
	"MuseBot/utils"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

var SkipMessage = errors.New("skip message")

type Message interface {
	GetText() (string, error)
}

type PRMessage struct {
	ID      int
	BaseUrl string
}

func (m PRMessage) GetText() (string, error) {
	url := m.BaseUrl + strconv.Itoa(m.ID)
	if utils.CheckUrlExist(url) {
		return fmt.Sprintf("<a href='%s'>PR #%d</a>", url, m.ID), nil
	}
	return "", fmt.Errorf("PR does not exist: %s", url)
}

type NodeMessage struct {
	ID      int
	BaseUrl string
}

func (m NodeMessage) GetText() (string, error) {
	url := m.BaseUrl + strconv.Itoa(m.ID)
	if utils.CheckUrlExist(url) {
		return fmt.Sprintf("<a href='%s'>Node #%d</a>", url, m.ID), nil
	}
	return "", fmt.Errorf("PR does not exist: %s", url)
}

type WikiMessage struct {
	Texts []string
}

func (m WikiMessage) GetText() (string, error) {
	i := rand.Intn(len(m.Texts))
	url, err := utils.GetRedirectUrl("https://en.wikipedia.org/wiki/Special:Random")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("<a href='%s'>%s</a>", url, m.Texts[i]), nil
}

type BasicMessage struct {
	Text string
}

func (m BasicMessage) GetText() (string, error) {
	return m.Text, nil
}

type BasicRandMessage struct {
	Texts []string
}

func (m BasicRandMessage) GetText() (string, error) {
	i := rand.Intn(len(m.Texts))
	return m.Texts[i], nil
}

type TravisHookMessage struct {
	Payload models.TravisHookPayload
}

func (m TravisHookMessage) GetText() (string, error) {
	if m.Payload.PullRequest {
		return "", SkipMessage
	}

	if strings.ToLower(m.Payload.Repository.OwnerName) != "musescore" {
		return "", SkipMessage
	}

	status := ""
	switch strings.ToLower(m.Payload.StatusMessage) {
	case "fixed":
		status = "has been fixed"
	case "broken":
		status = "has been broken"
	case "still failing":
		status = "is still failing"
	case "errored":
		status = "has errored"
	default:
		return "", SkipMessage
	}

	commitUrl := fmt.Sprintf("<a href='%s%s'>%s</a>", Config.GitHubCommitUrl, m.Payload.Commit, m.Payload.Commit[0:6])
	buildUrl := fmt.Sprintf("<a href='%s'>build</a>", m.Payload.BuildURL)
	msg := fmt.Sprintf("MuseScore/%s : %s by %s: %s %s", m.Payload.Branch, commitUrl, m.Payload.CommitterName, buildUrl, status)
	return msg, nil
}

type GitHubHookPRMessage struct {
	Payload models.GitHubPRPayload
}

func (m GitHubHookPRMessage) GetText() (string, error) {
	if m.Payload.Action != "opened" {
		return "", SkipMessage
	}
	msg := fmt.Sprintf("New Pull Request: <a href='%s'>#%d - %s</a> by %s",
		m.Payload.PullRequest.HTMLURL, m.Payload.PullRequest.Number,
		utils.SanitizeText(m.Payload.PullRequest.Title), m.Payload.PullRequest.User.Login,
	)
	return msg, nil
}

type GitHubHookPushMessage struct {
	Payload models.GitHubPushPayload
}

func (m GitHubHookPushMessage) GetText() (string, error) {
	if len(m.Payload.Commits) == 0 {
		return "", SkipMessage
	}

	branch := strings.Replace(m.Payload.Ref, "refs/heads/", "", -1)
	commitMsg := utils.StripToLength(utils.SanitizeText(m.Payload.Commits[0].Message), 70, "...")
	commitLink := fmt.Sprintf("<a href='%s'>%s</a> - <i>%s</i>", m.Payload.Commits[0].URL, m.Payload.Commits[0].ID[0:6], commitMsg)
	msg := ""
	if len(m.Payload.Commits) > 1 {
		msg = fmt.Sprintf("%s pushed %d commits to %s, including %s", m.Payload.Pusher.Name, len(m.Payload.Commits), branch, commitLink)
	} else {
		msg = fmt.Sprintf("%s pushed commit to %s: %s", m.Payload.Pusher.Name, branch, commitLink)
	}
	return msg, nil
}
