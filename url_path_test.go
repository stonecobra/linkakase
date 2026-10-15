package linkquisition_test

import (
	"testing"

	"github.com/strobotti/linkquisition"
)

func TestMatchesPathRule(t *testing.T) {
	const pr = "https://github.com/bonsaitechnologies/bonsai/pull/4534?x=1#top"
	cases := []struct {
		rule string
		want bool
	}{
		{"github.com/*", true},
		{"GitHub.com/*", true},
		{"github.com/bonsaitechnologies/*", true},
		{"github.com/bonsaitechnologies/bonsai/*", true},
		{"github.com/bonsaitechnologies/bonsai/pull/*", true},
		{"github.com/bonsaitechnologies/bonsai/issues/*", false},
		{"github.com/bonsai/*", false},
		{"github.com/bonsaitechnologies/bonsa/*", false},
		{"gist.github.com/*", false},
	}
	for _, c := range cases {
		if got := linkquisition.NewURL(pr).MatchesPathRule(c.rule); got != c.want {
			t.Errorf("MatchesPathRule(%q) = %v, want %v", c.rule, got, c.want)
		}
	}
}

func TestGetMatchingBrowserPrefersDeepestPath(t *testing.T) {
	s := &linkquisition.Settings{Browsers: []linkquisition.BrowserSettings{
		{Name: "site", Command: "a", Matches: []linkquisition.BrowserMatch{{Type: "site", Value: "github.com"}}},
		{Name: "org", Command: "b", Matches: []linkquisition.BrowserMatch{{Type: "path", Value: "github.com/bonsaitechnologies/*"}}},
		{Name: "repo", Command: "c", Matches: []linkquisition.BrowserMatch{{Type: "path", Value: "github.com/bonsaitechnologies/bonsai/*"}}},
	}}
	b, err := s.GetMatchingBrowser("https://github.com/bonsaitechnologies/bonsai/pull/4534")
	if err != nil || b.Name != "repo" {
		t.Fatalf("got %+v, %v; want repo", b, err)
	}
	b, err = s.GetMatchingBrowser("https://github.com/bonsaitechnologies/other")
	if err != nil || b.Name != "org" {
		t.Fatalf("got %+v, %v; want org", b, err)
	}
	b, err = s.GetMatchingBrowser("https://github.com/someone/else")
	if err != nil || b.Name != "site" {
		t.Fatalf("got %+v, %v; want site", b, err)
	}
}
