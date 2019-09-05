package models

import "time"

type TravisConfig struct {
	Config struct {
		Notifications struct {
			Webhook struct {
				PublicKey string `json:"public_key"`
			} `json:"webhook"`
		} `json:"notifications"`
	} `json:"config"`
}

type TravisHookPayload struct {
	ID                int         `json:"id"`
	Number            string      `json:"number"`
	Type              string      `json:"type"`
	State             string      `json:"state"`
	Status            int         `json:"status"`
	Result            int         `json:"result"`
	StatusMessage     string      `json:"status_message"`
	ResultMessage     string      `json:"result_message"`
	StartedAt         time.Time   `json:"started_at"`
	FinishedAt        time.Time   `json:"finished_at"`
	Duration          int         `json:"duration"`
	BuildURL          string      `json:"build_url"`
	CommitID          int         `json:"commit_id"`
	Commit            string      `json:"commit"`
	BaseCommit        interface{} `json:"base_commit"`
	HeadCommit        interface{} `json:"head_commit"`
	Branch            string      `json:"branch"`
	Message           string      `json:"message"`
	CompareURL        string      `json:"compare_url"`
	CommittedAt       time.Time   `json:"committed_at"`
	AuthorName        string      `json:"author_name"`
	AuthorEmail       string      `json:"author_email"`
	CommitterName     string      `json:"committer_name"`
	CommitterEmail    string      `json:"committer_email"`
	PullRequest       bool        `json:"pull_request"`
	PullRequestNumber interface{} `json:"pull_request_number"`
	PullRequestTitle  interface{} `json:"pull_request_title"`
	Tag               interface{} `json:"tag"`
	Repository        struct {
		ID        int         `json:"id"`
		Name      string      `json:"name"`
		OwnerName string      `json:"owner_name"`
		URL       interface{} `json:"url"`
	} `json:"repository"`
}
