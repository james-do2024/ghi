package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v60/github"
)

type RestRequest struct {
	Owner       string
	Repo        string
	CurrentPath string
	GitClient   *github.Client
	Ctx         context.Context
}

func (r *RestRequest) NavigateUp() {
	if r.CurrentPath == "" {
		// Already at the root, do nothing
		return
	}
	lastSlash := strings.LastIndex(r.CurrentPath, "/")
	if lastSlash > 0 {
		r.CurrentPath = r.CurrentPath[:lastSlash]
	} else {
		// Move to root
		r.CurrentPath = ""
	}
}

func (r *RestRequest) NavigateRoot() {
	r.CurrentPath = "" // Simply reset to root
}

func (r *RestRequest) NavigateIndex(idx int, dirMap map[int]string) {
	// Check if the index is valid for the current directory listing
	if name, ok := dirMap[idx]; ok {
		// If valid, update the path
		if r.CurrentPath == "" {
			r.CurrentPath = name
		} else {
			r.CurrentPath = strings.TrimRight(r.CurrentPath, "/") + "/" + name
		}
	} else {
		fmt.Printf("invalid index: %d\n", idx)
	}
}

func (r *RestRequest) GetContent() (*string, []*github.RepositoryContent, error) {
	var allContent []*github.RepositoryContent

	opt := &github.RepositoryContentGetOptions{}
	for {
		fileContent, directoryContent, resp, err := r.GitClient.Repositories.GetContents(r.Ctx, r.Owner, r.Repo, r.CurrentPath, opt)
		if err != nil {
			if errResp, ok := err.(*github.ErrorResponse); ok {
				if errResp.Response.StatusCode == 404 {
					// Specific handling for 404 errors
					return nil, nil, fmt.Errorf("content not found at path '%s' in repository '%s/%s'", r.CurrentPath, r.Owner, r.Repo)
				}
			}
			// Otherwise, return whatever kind of error it is.
			return nil, nil, err
		}

		// fileContent can only be non-nil if we have requested a file
		// We return a pointer to the file's contents if this is the case
		if fileContent != nil {
			rawFile, err := getFileContent(fileContent)
			if err != nil {
				return nil, nil, err
			}
			return rawFile, nil, nil // Pagination not required for files
		}

		// Append as we go, break if there is no further pagination to be done
		allContent = append(allContent, directoryContent...)
		if resp.NextPage == 0 {
			break
		}
	}
	return nil, allContent, nil
}

func getFileContent(content *github.RepositoryContent) (*string, error) {
	rawFile, err := content.GetContent()
	if err != nil {
		return nil, err
	}

	return &rawFile, nil
}
