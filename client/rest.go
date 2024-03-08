/*
 *   Because we are leveraging the GitHub client module, much
 *   of the heavy lifting has already been done. However, it
 *   would be nice to be able to override some defaults with
 *   environment variables and set the auth token that way.
 *
 *   A thin wrapper can achieve this.
 */

package client

import (
	"fmt"
	"log"
	"net/url"
	"regexp"

	"github.com/google/go-github/v60/github"
	"github.com/james-do2024/ghi/config"
)

func NewRestClient() (*github.Client, error) {
	client := github.NewClient(nil)
	if err := applyEnvVars(client); err != nil {
		return nil, err
	}

	return client, nil
}

func applyEnvVars(c *github.Client) error {
	if env, base, err := config.GetEnvIfSet("defaultBaseURL"); err != nil {
		return err
	} else if base != "" {
		// It is possible for the URL in the env var to be invalid - check it here
		parsedURL, err := url.ParseRequestURI(base)
		if err != nil {
			// Log the problem, but carry on with module's default value
			log.Printf("malformed url in env var %s: %s\n", env, base)
		} else {
			c.BaseURL = parsedURL
		}
	}

	if _, agent, err := config.GetEnvIfSet("defaultUserAgent"); err != nil {
		return err
	} else if agent != "" {
		c.UserAgent = agent
	}

	// For simplicity, require token and treat its absence as fatal
	env, token, err := config.GetEnvIfSet("authToken")
	if err != nil {
		return err
	} else if token == "" {
		return fmt.Errorf("auth token not set or empty in %s", env)
	}

	if !isPATValid(token) {
		return fmt.Errorf("supplied token does not appear to be a valid GitHub PAT")
	}

	c = c.WithAuthToken(token)
	return nil
}

// Validates structure and size only
func isPATValid(token string) bool {
	tokenRx := regexp.MustCompile(`^github_pat_[a-zA-Z0-9_]{29,244}$`)
	return tokenRx.MatchString(token)
}
