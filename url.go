package linkquisition

import (
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/publicsuffix"
)

var (
	ipRegex   = regexp.MustCompile(`^\d+\.\d+\.\d+\.\d+$`)
	siteRegex = regexp.MustCompile(`^https?://([^/]+)(/|$)`)
)

type URL struct {
	url string
}

func NewURL(u string) *URL {
	return &URL{url: u}
}

func (u URL) GetDomain() (string, error) {
	parsedUrl, err := url.Parse(u.url)
	if err != nil {
		return "", err
	}

	// If the hostname is an IP address, we return it as is
	if ipRegex.MatchString(parsedUrl.Hostname()) {
		return parsedUrl.Hostname(), nil
	}

	tldPlusOne, err := publicsuffix.EffectiveTLDPlusOne(parsedUrl.Hostname())
	if err != nil {
		return "", err
	}

	return tldPlusOne, nil
}

func (u URL) GetSite() (string, error) {
	match := siteRegex.FindStringSubmatch(u.url)
	if len(match) > 1 {
		return match[1], nil
	}

	return "", nil
}

// GetPathSegments returns the non-empty segments of the URL's path (query and fragment excluded).
func (u URL) GetPathSegments() []string {
	parsedUrl, err := url.Parse(u.url)
	if err != nil {
		return nil
	}

	var segments []string
	for _, segment := range strings.Split(parsedUrl.Path, "/") {
		if segment != "" {
			segments = append(segments, segment)
		}
	}

	return segments
}

// MatchesPathRule reports whether the URL is under the path rule "host/seg1/seg2/*": the same host (host:port,
// ignoring case) and a path that starts with those segments. A rule without segments ("host/*") matches the
// whole host.
func (u URL) MatchesPathRule(rule string) bool {
	parsedUrl, err := url.Parse(u.url)
	if err != nil {
		return false
	}

	rule = strings.TrimSuffix(strings.TrimSuffix(rule, "*"), "/")
	host, path, _ := strings.Cut(rule, "/")
	if !strings.EqualFold(host, parsedUrl.Host) {
		return false
	}

	want := strings.FieldsFunc(path, func(r rune) bool { return r == '/' })
	have := u.GetPathSegments()
	if len(have) < len(want) {
		return false
	}
	for i := range want {
		if have[i] != want[i] {
			return false
		}
	}

	return true
}

// PathRuleDepth is the number of path segments in a path rule, for picking the most specific of several.
func PathRuleDepth(rule string) int {
	rule = strings.TrimSuffix(strings.TrimSuffix(rule, "*"), "/")
	_, path, _ := strings.Cut(rule, "/")

	return len(strings.FieldsFunc(path, func(r rune) bool { return r == '/' }))
}
