package client

import (
	"context"
	"fmt"

	"github.com/google/go-github/v60/github"
)

type RestRequest struct {
	gitClient *github.Client
	ctx       context.Context
}

type RestResponse struct {
	currentPath string
	responseMap map[int]string
	fileContent *github.RepositoryContent
	dirContent  []*github.RepositoryContent
}

func (r *RestRequest) GetRepoContents(owner, repo, path string) (*RestResponse, error) {
	resp := &RestResponse{
		currentPath: path,
		responseMap: make(map[int]string),
	}

	opt := &github.RepositoryContentGetOptions{}
	fileContent, entries, _, err := r.gitClient.Repositories.GetContents(r.ctx, owner, repo, path, opt)
	if err != nil {
		return nil, fmt.Errorf("fetching repo at toplevel failed: %v", err)
	}
	for i, entry := range entries {
		if *entry.Type == "dir" {
			*entry.Name += "/"
		}
		resp.responseMap[i] = *entry.Name
	}
	resp.fileContent = fileContent
	resp.dirContent = entries

	return resp, nil
}
