package client

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
)

type RestRequest struct {
	GitClient *github.Client
	Ctx       context.Context
}

func (r *RestRequest) GetDirectoryContents(owner, repo, path string) ([]*github.RepositoryContent, error) {
	opt := &github.RepositoryContentGetOptions{}
	_, directoryContent, _, err := r.GitClient.Repositories.GetContents(r.Ctx, owner, repo, path, opt)
	if err != nil {
		return nil, fmt.Errorf("fetching directory contents failed: %v", err)
	}
	return directoryContent, nil
}

func (r *RestRequest) GetFileContent(owner, repo, path string) (*github.RepositoryContent, error) {
	opt := &github.RepositoryContentGetOptions{}
	fileContent, _, _, err := r.GitClient.Repositories.GetContents(r.Ctx, owner, repo, path, opt)
	if err != nil {
		return nil, fmt.Errorf("fetching file content failed: %v", err)
	}
	return fileContent, nil
}
