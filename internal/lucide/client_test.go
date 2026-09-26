package lucide

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
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

func TestCreateRelease(t *testing.T) {
	var got map[string]any
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/repos/kaugesaar/lucide-go/releases" {
			t.Errorf("request = %s %s", r.Method, r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"html_url":"https://github.com/kaugesaar/lucide-go/releases/tag/v1.2.3"}`))
	}))
	defer server.Close()

	gh := github.NewClient(nil)
	gh.BaseURL, _ = url.Parse(server.URL + "/")
	client := &Client{gh: gh}

	releaseURL, err := client.CreateRelease(context.Background(), "v1.2.3", "abc123", "notes")
	if err != nil {
		t.Fatalf("CreateRelease() error = %v", err)
	}
	if releaseURL != "https://github.com/kaugesaar/lucide-go/releases/tag/v1.2.3" {
		t.Errorf("CreateRelease() = %q", releaseURL)
	}
	want := map[string]any{"tag_name": "v1.2.3", "target_commitish": "abc123", "name": "v1.2.3", "body": "notes"}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("request %s = %v, want %v", k, got[k], v)
		}
	}
}
