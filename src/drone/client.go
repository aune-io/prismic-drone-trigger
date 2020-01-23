package drone

import (
	"github.com/drone/drone-go/drone"
	"golang.org/x/oauth2"
)

// Client wraps the Drone client objects
type Client struct {
	droneClient drone.Client
}

// NewClient creates a new instance of the Drone API client
func NewClient(droneHost, droneToken string) *Client {
	config := new(oauth2.Config)
	auth := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: droneToken,
		},
	)
	droneClient := drone.NewClient(droneHost, auth)

	return &Client{droneClient}
}

// GetLastBuildID returns the last build for the given owner/repo:branch
func (c *Client) GetLastBuildID(repoOwner string, repoName string, repoBranch string) (int64, error) {
	build, err := c.droneClient.BuildLast(repoOwner, repoName, repoBranch)
	if err != nil {
		return 0, err
	}

	return build.ID, err
}

// TriggerBuild fetches the last build for the given owner/repo:branch, and triggers a restart
func (c *Client) TriggerBuild(repoOwner string, repoName string, repoBranch string) (int64, error) {
	build, err := c.droneClient.BuildLast(repoOwner, repoName, repoBranch)
	if err != nil {
		return 0, err
	}

	nBuild, err := c.droneClient.BuildRestart(repoOwner, repoName, int(build.ID), nil)
	if err != nil {
		return 0, err
	}

	return nBuild.ID, err
}
