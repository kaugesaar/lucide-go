package lucide

import (
	"context"
	"testing"

	"github.com/google/go-github/v78/github"
)

func TestNewClient(t *testing.T) {
	client := NewClient("")
	if client == nil {
		t.Fatal("NewClient(\"\") returned nil")
	}
	if client.gh == nil {
		t.Error("NewClient(\"\").gh is nil")
	}

	client = NewClient("test-token")
	if client == nil {
		t.Fatal("NewClient(\"test-token\") returned nil")
	}
	if client.gh == nil {
		t.Error("NewClient(\"test-token\").gh is nil")
	}
}

func TestGetSourceArchiveURL(t *testing.T) {
	client := NewClient("")
	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	ctx := context.Background()
	archiveURL, err := client.GetSourceArchiveURL(ctx, "0.460.0")
	if err != nil {
		t.Fatalf("GetSourceArchiveURL() error = %v", err)
	}
	if archiveURL == nil {
		t.Fatal("GetSourceArchiveURL() returned nil URL")
	}
	if archiveURL.String() == "" {
		t.Error("GetSourceArchiveURL() returned empty URL")
	}
}

func TestReleaseFindIconsAsset(t *testing.T) {
	tests := []struct {
		name      string
		assets    []*github.ReleaseAsset
		wantFound bool
		wantName  string
	}{
		{
			name: "finds correct asset",
			assets: []*github.ReleaseAsset{
				{Name: github.Ptr("lucide-0.553.0.tar.gz")},
				{Name: github.Ptr("lucide-icons-0.553.0.zip")},
				{Name: github.Ptr("other-file.txt")},
			},
			wantFound: true,
			wantName:  "lucide-icons-0.553.0.zip",
		},
		{
			name: "no matching asset",
			assets: []*github.ReleaseAsset{
				{Name: github.Ptr("lucide-0.553.0.tar.gz")},
				{Name: github.Ptr("other-file.txt")},
			},
			wantFound: false,
		},
		{
			name:      "empty assets",
			assets:    []*github.ReleaseAsset{},
			wantFound: false,
		},
		{
			name: "wrong extension",
			assets: []*github.ReleaseAsset{
				{Name: github.Ptr("lucide-icons-0.553.0.tar.gz")},
			},
			wantFound: false,
		},
		{
			name: "wrong prefix",
			assets: []*github.ReleaseAsset{
				{Name: github.Ptr("icons-0.553.0.zip")},
			},
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			release := &Release{
				TagName: "0.553.0",
				Assets:  tt.assets,
			}

			asset, err := release.FindIconsAsset()

			if tt.wantFound {
				if err != nil {
					t.Errorf("FindIconsAsset() error = %v, want nil", err)
				}
				if asset == nil {
					t.Fatal("FindIconsAsset() returned nil asset")
				}
				if asset.GetName() != tt.wantName {
					t.Errorf("FindIconsAsset() asset name = %q, want %q", asset.GetName(), tt.wantName)
				}
			} else {
				if err == nil {
					t.Error("FindIconsAsset() should return error when asset not found")
				}
				if asset != nil {
					t.Errorf("FindIconsAsset() returned asset %v, want nil", asset)
				}
			}
		})
	}
}
