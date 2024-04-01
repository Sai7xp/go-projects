package models

// build request received via kafka
type BuildRequestDetails struct {
	BuildId          string                 `json:"build_id" bson:"build_id"`
	ProjectGithubUrl string                 `json:"project_github_url" bson:"project_github_url"`
	BuildCommand     string                 `json:"build_command" bson:"build_command"`
	BuildOutDir      string                 `json:"build_out_dir" bson:"build_out_dir"`
}
