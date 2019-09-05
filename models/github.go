package models

import "time"

type GitHubPRPayload struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		URL      string `json:"url"`
		ID       int    `json:"id"`
		NodeID   string `json:"node_id"`
		HTMLURL  string `json:"html_url"`
		DiffURL  string `json:"diff_url"`
		PatchURL string `json:"patch_url"`
		IssueURL string `json:"issue_url"`
		Number   int    `json:"number"`
		State    string `json:"state"`
		Locked   bool   `json:"locked"`
		Title    string `json:"title"`
		User     struct {
			Login             string `json:"login"`
			ID                int    `json:"id"`
			NodeID            string `json:"node_id"`
			AvatarURL         string `json:"avatar_url"`
			GravatarID        string `json:"gravatar_id"`
			URL               string `json:"url"`
			HTMLURL           string `json:"html_url"`
			FollowersURL      string `json:"followers_url"`
			FollowingURL      string `json:"following_url"`
			GistsURL          string `json:"gists_url"`
			StarredURL        string `json:"starred_url"`
			SubscriptionsURL  string `json:"subscriptions_url"`
			OrganizationsURL  string `json:"organizations_url"`
			ReposURL          string `json:"repos_url"`
			EventsURL         string `json:"events_url"`
			ReceivedEventsURL string `json:"received_events_url"`
			Type              string `json:"type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"user"`
		Body                interface{}   `json:"body"`
		CreatedAt           time.Time     `json:"created_at"`
		UpdatedAt           time.Time     `json:"updated_at"`
		ClosedAt            interface{}   `json:"closed_at"`
		MergedAt            interface{}   `json:"merged_at"`
		MergeCommitSha      string        `json:"merge_commit_sha"`
		Assignee            interface{}   `json:"assignee"`
		Assignees           []interface{} `json:"assignees"`
		RequestedReviewers  []interface{} `json:"requested_reviewers"`
		RequestedTeams      []interface{} `json:"requested_teams"`
		CommitsURL          string        `json:"commits_url"`
		ReviewCommentsURL   string        `json:"review_comments_url"`
		ReviewCommentURL    string        `json:"review_comment_url"`
		CommentsURL         string        `json:"comments_url"`
		StatusesURL         string        `json:"statuses_url"`
		AuthorAssociation   string        `json:"author_association"`
		Draft               bool          `json:"draft"`
		Merged              bool          `json:"merged"`
		Mergeable           interface{}   `json:"mergeable"`
		Rebaseable          interface{}   `json:"rebaseable"`
		MergeableState      string        `json:"mergeable_state"`
		MergedBy            interface{}   `json:"merged_by"`
		Comments            int           `json:"comments"`
		ReviewComments      int           `json:"review_comments"`
		MaintainerCanModify bool          `json:"maintainer_can_modify"`
		Commits             int           `json:"commits"`
		Additions           int           `json:"additions"`
		Deletions           int           `json:"deletions"`
		ChangedFiles        int           `json:"changed_files"`
	} `json:"pull_request"`
}

type GitHubPushPayload struct {
	Commits []struct {
		ID        string    `json:"id"`
		TreeID    string    `json:"tree_id"`
		Distinct  bool      `json:"distinct"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		URL       string    `json:"url"`
		Author    struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"author"`
		Committer struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Username string `json:"username"`
		} `json:"committer"`
		Added    []interface{} `json:"added"`
		Removed  []interface{} `json:"removed"`
		Modified []string      `json:"modified"`
	} `json:"commits"`
	Ref     string      `json:"ref"`
	Before  string      `json:"before"`
	After   string      `json:"after"`
	Created bool        `json:"created"`
	Deleted bool        `json:"deleted"`
	Forced  bool        `json:"forced"`
	BaseRef interface{} `json:"base_ref"`
	Compare string      `json:"compare"`
	Pusher  struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
}
