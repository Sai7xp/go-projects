package models

// events
type Events struct{}

// Response Body endpoint for /api/v1/build/{build_id}
type WebfronBuildDetails struct {
	BuildID          string `json:"build_id"`
	ProjectGitHubURL string `json:"project_github_url"`
	Events           Events `json:"events"`
}
