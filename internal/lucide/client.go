package lucide

import (
	"context"
	"fmt"

	"github.com/google/go-github/v78/github"
)

const (
	owner       = "kaugesaar"
	repo        = "lucide-go"
	lucideOwner = "lucide-icons"
	lucideRepo  = "lucide"
)

// Client provides access to the Lucide icon repository via GitHub API.
type Client struct {
	gh *github.Client
}

// Release represents a Lucide release with its metadata.
type Release struct {
	TagName string
	Name    string
	URL     string
	Body    string
	Assets  []*github.ReleaseAsset
}

// NewClient creates a new Lucide client.
// Pass an empty string for token to access public data without authentication.
func NewClient(token string) *Client {
	var gh *github.Client
	if token != "" {
		gh = github.NewClient(nil).WithAuthToken(token)
	} else {
		gh = github.NewClient(nil)
	}

	return &Client{gh: gh}
}

// GetLatestRelease fetches the latest release from the Lucide repository.
func (c *Client) GetLatestRelease(ctx context.Context) (*Release, error) {
	release, _, err := c.gh.Repositories.GetLatestRelease(ctx, lucideOwner, lucideRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest release: %w", err)
	}

	return &Release{
		TagName: release.GetTagName(),
		Name:    release.GetName(),
		URL:     release.GetHTMLURL(),
		Body:    release.GetBody(),
		Assets:  release.Assets,
	}, nil
}

// FindIconsAsset finds the lucide-icons zip asset in a release.
func (r *Release) FindIconsAsset() (*github.ReleaseAsset, error) {
	for _, asset := range r.Assets {
		name := asset.GetName()
		if len(name) > 13 && name[:13] == "lucide-icons-" && name[len(name)-4:] == ".zip" {
			return asset, nil
		}
	}
	return nil, fmt.Errorf("icons asset not found in release %s", r.TagName)
}

// CreateRelease creates a GitHub release for the lucide-go repository.
// Returns the HTML URL of the created release.
func (c *Client) CreateRelease(ctx context.Context, version, releaseNotes string) (string, error) {
	release := &github.RepositoryRelease{
		TagName: github.Ptr(version),
		Name:    github.Ptr(version),
		Body:    github.Ptr(releaseNotes),
	}

	created, _, err := c.gh.Repositories.CreateRelease(ctx, owner, repo, release)
	if err != nil {
		return "", fmt.Errorf("failed to create release: %w", err)
	}

	return created.GetHTMLURL(), nil
}
