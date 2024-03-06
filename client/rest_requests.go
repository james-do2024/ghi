package client

import (
	"context"

	"github.com/google/go-github/v60/github"
)

type RestRequest struct {
	GitClient *github.Client
	Ctx       context.Context
}

func (r *RestRequest) GetContent(owner, repo, path string) (*string, []*github.RepositoryContent, error) {
	opt := &github.RepositoryContentGetOptions{}
	fileContent, directoryContent, _, err := r.GitClient.Repositories.GetContents(r.Ctx, owner, repo, path, opt)
	if err != nil {
		return nil, nil, err
	}

	// fileContent can only be non-nil if we have requested a file
	// We return a pointer to the file's contents if this is the case
	if fileContent != nil {
		rawFile, err := getFileContent(fileContent)
		if err != nil {
			return nil, nil, err
		}
		return rawFile, nil, nil
	}
	return nil, directoryContent, nil
}

func getFileContent(content *github.RepositoryContent) (*string, error) {
	rawFile, err := content.GetContent()
	if err != nil {
		return nil, err
	}

	return &rawFile, nil
}
