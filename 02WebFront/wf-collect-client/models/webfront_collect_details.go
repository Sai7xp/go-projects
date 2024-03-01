package models

// Request Body data for /api/v1/collect endpoint
type WebfrontCollectDetails struct {
	// A public github repo url
	ProjectGithubUrl string `json:"project_github_url" validate:"required"`

	// Command to generate the build when executed in the root folder
	// of the cloned repository. Usually its `npm run build` or `yarn build`
	BuildCommand string `json:"build_command" validate:"required"`

	// Name of the folder generated after the build
	BuildOutDir string `json:"build_out_dir" validate:"required"`
}
