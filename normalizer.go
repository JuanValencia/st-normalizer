package normalizer

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
)

const (
	DecodeUrlError     = "DecodeUrlError"
	EmptyInputUrl      = "EmptyInputUrl"
	ParseError         = "ParseError"
	PathUnescapeError  = "PathUnescapeError"
	QueryUnescapeError = "QueryUnescapeError"
)

var rxDirIndex = regexp.MustCompile(`(^|/)((?:default|index)\.\w{1,4})$`)
var rxDupSlashes = regexp.MustCompile(`/{2,}`)
var rxQueryParamId = regexp.MustCompile(`(^(id|libid|p)=*)`)

type Normalizer struct {
	RawUrl            string `json:"raw_url"`
	Protocol          string `json:"protocol"`
	CanonicalUrl      string `json:"canonical_url"`
	CanonicalUrlHash  string `json:"canonical_url_hash"`
	UrlIdentifier     string `json:"url_identifier"`
	UrlIdentifierHash string `json:"url_identifier_hash"`
	RawQueryParams    string `json:"raw_query_params"`
	parsedRawUrl      *url.URL
	transformedUrl    *url.URL
}

func decodeRawUrl(rawUrl string) (string, error) {
	rawUrl, err := url.QueryUnescape(rawUrl)
	if err != nil {
		log.Printf("%s: %s \n", QueryUnescapeError, err.Error())
		return "", err
	}
	return rawUrl, nil
}

func NewNormalizer(rawUrl string) (*Normalizer, error) {

	if rawUrl == "" {
		return nil, errors.New(EmptyInputUrl)
	}

	// lower case
	rawUrl = strings.ToLower(rawUrl)

	// decode
	rawUrl, err := decodeRawUrl(rawUrl)
	if err != nil {
		return nil, errors.New(DecodeUrlError)
	}

	// parse
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Printf("%s: %s \n", ParseError, err.Error())
		return nil, errors.New(ParseError)
	}

	// Remove fragment
	u.Fragment = ""

	return &Normalizer{
		RawUrl:         rawUrl,
		Protocol:       u.Scheme,
		RawQueryParams: u.RawQuery,
		parsedRawUrl:   u,
		transformedUrl: u,
	}, nil
}

func (n *Normalizer) String() string {
	rawUrl := fmt.Sprintf("raw_url: %s \n", n.RawUrl)
	protocol := fmt.Sprintf("protocol: %s \n", n.Protocol)
	rawQueryParams := fmt.Sprintf("raw_query_params: %s \n", n.RawQueryParams)
	canonicalUrl := fmt.Sprintf("canonical_url: %s \n", n.CanonicalUrl)
	canonicalUrlHash := fmt.Sprintf("canonical_url_hash: %s \n", n.CanonicalUrlHash)
	UrlIdentifier := fmt.Sprintf("url_identifier: %s \n", n.UrlIdentifier)
	UrlIdentifierHash := fmt.Sprintf("url_identifier_hash: %s \n", n.UrlIdentifierHash)
	return (rawUrl + protocol + rawQueryParams + canonicalUrl + canonicalUrlHash +
		UrlIdentifier + UrlIdentifierHash)
}

func (n *Normalizer) removeTrailingSlash() {
	u := n.transformedUrl
	if l := len(u.Path); l > 0 {
		if strings.HasSuffix(u.Path, "/") {
			u.Path = u.Path[:l-1]
		}
	} else if l = len(u.Host); l > 0 {
		if strings.HasSuffix(u.Host, "/") {
			u.Host = u.Host[:l-1]
		}
	}
	n.transformedUrl = u
}

func (n *Normalizer) removeDirectoryIndex() {
	u := n.transformedUrl
	if len(u.Path) > 0 {
		u.Path = rxDirIndex.ReplaceAllString(u.Path, "$1")
	}
	n.transformedUrl = u
}

func (n *Normalizer) removeUnusedQueryParams() {

	// keep first query param and remove rest
	u := n.transformedUrl
	if len(u.RawQuery) > 0 {
		u.RawQuery = strings.SplitN(u.RawQuery, "&", 2)[0]
	}
	n.transformedUrl = u
}

func (n *Normalizer) removeProtocol() {
	u := n.transformedUrl
	if len(u.Scheme) > 0 {
		u.Scheme = ""
	}
	n.transformedUrl = u
}

func (n *Normalizer) removeDuplicateSlashes() {
	u := n.transformedUrl
	if len(u.Path) > 0 {
		u.Path = rxDupSlashes.ReplaceAllString(u.Path, "/")
	}
	n.transformedUrl = u
}

func (n *Normalizer) removeQueryParams() {
	n.transformedUrl.RawQuery = ""
}

func (n *Normalizer) removeDefaultPort() {
	host := n.transformedUrl.Host
	if strings.Contains(host, ":") {
		host = strings.SplitN(host, ":", 2)[0]
	}
	n.transformedUrl.Host = host
}

func (n *Normalizer) removeDotSegments() {
	u := n.transformedUrl
	if len(u.Path) > 0 {
		var dotFree []string
		var lastIsDot bool

		sections := strings.Split(u.Path, "/")
		for _, s := range sections {
			if s == ".." {
				if len(dotFree) > 0 {
					dotFree = dotFree[:len(dotFree)-1]
				}
			} else if s != "." {
				dotFree = append(dotFree, s)
			}
			lastIsDot = (s == "." || s == "..")
		}

		// Special case if host does not end with / and new path does not begin with /
		u.Path = strings.Join(dotFree, "/")
		if u.Host != "" &&
			!strings.HasSuffix(u.Host, "/") && !strings.HasPrefix(u.Path, "/") {
			u.Path = "/" + u.Path
		}

		// Special case if the last segment was a dot, make sure the path ends with a slash
		if lastIsDot && !strings.HasSuffix(u.Path, "/") {
			u.Path += "/"
		}
	}
	n.transformedUrl = u
}

func (n *Normalizer) setUrlIdentifier() {
	normalizedUrlString := fmt.Sprintf(
		"%s%s",
		n.transformedUrl.Host,
		n.transformedUrl.Path)
	if len(n.transformedUrl.RawQuery) > 0 {
		normalizedUrlString = fmt.Sprintf(
			"%s?%s", normalizedUrlString, n.transformedUrl.RawQuery)
	}
	if normalizedUrlString != "" {
		n.UrlIdentifier = normalizedUrlString
		n.UrlIdentifierHash = GetMD5Hash(n.UrlIdentifier)
	}
}

func (n *Normalizer) setCanonicalUrl() {
	var canonicalUrl string
	if n.Protocol != "" {
		canonicalUrl = fmt.Sprintf(
			"%s://%s%s",
			n.Protocol,
			n.transformedUrl.Host,
			n.transformedUrl.Path)
	} else {
		canonicalUrl = n.UrlIdentifier
	}
	n.CanonicalUrl = canonicalUrl
	n.CanonicalUrlHash = GetMD5Hash(n.CanonicalUrl)
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (n *Normalizer) HandleQueryParams() {
	if rxQueryParamId.MatchString(n.transformedUrl.RawQuery) {

		// remove unused query parama
		n.removeUnusedQueryParams()
	} else {

		// remove query params
		n.removeQueryParams()
	}
}

func (n *Normalizer) Normalize() {

	// NewNormalizer does following
	// - converts url to lowercase
	// - decodes url
	// - removes Fragment

	// remove default port
	n.removeDefaultPort()

	// remove dot segments
	n.removeDotSegments()

	// remove directory index
	n.removeDirectoryIndex()

	// remove protocol
	n.removeProtocol()

	// remove duplicate slashes
	n.removeDuplicateSlashes()

	// remove Trailing Slash
	n.removeTrailingSlash()

	// handle query params
	n.HandleQueryParams()

	// set url identfier
	n.setUrlIdentifier()

	// set canonical url
	n.setCanonicalUrl()
}
