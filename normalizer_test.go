package normalizer

import (
	"fmt"
	"strings"
	"testing"
)

func TestNormalizerSetScheme(t *testing.T) {
	testSchemes := map[string]string{
		"http":  "HTTP://www.Example.com/",
		"https": "HTTPs://www.Example.com/"}
	for proto, testUrl := range testSchemes {
		n, err := NewNormalizer(testUrl)
		if err != nil {
			t.Error("Error: Normalizer `NewNormalizer` test failed!")
		}
		if n.Protocol != proto && n.RawUrl != strings.ToLower(testUrl) {
			t.Error("Error: Normalizer `NewNormalizer` test failed!")
		}
	}
}

func TestNormalizerRemoveDefaultPort(t *testing.T) {
	testUrl := "http://www.example.com:8080/"
	n, err := NewNormalizer(testUrl)
	if err != nil {
		t.Error("Error: Normalizer `removeDefaultPort` test failed!")
	}
	n.removeDefaultPort()
	if n.transformedUrl.Host != "www.example.com" {
		fmt.Printf("host is %s \n", n.transformedUrl.Host)
		t.Error("Error: Normalizer `removeDefaultPort` test failed!")
	}
}

func TestNormalizerRemoveDotSegments(t *testing.T) {
	testUrl := "http://www.example.com/../a/b/../c/./d.html"
	n, err := NewNormalizer(testUrl)
	if err != nil {
		t.Error("Error: Normalizer `removeDotSegments` test failed!")
	}
	n.removeDefaultPort()
	n.removeDotSegments()
	if n.transformedUrl.Path != "/a/c/d.html" {
		fmt.Printf("path is %s \n", n.transformedUrl.Path)
		t.Error("Error: Normalizer `removeDotSegments` test failed!")
	}
}

func TestNormalizerRemoveTrailingSlash(t *testing.T) {
	testUrl := "http://www.example.com/alice/"
	n, err := NewNormalizer(testUrl)
	if err != nil {
		t.Error("Error: Normalizer `removeTrailingSlash` test failed!")
	}
	n.removeDefaultPort()
	n.removeTrailingSlash()
	if n.transformedUrl.Path != "/alice" {
		fmt.Printf("path is %s \n", n.transformedUrl.Path)
		t.Error("Error: Normalizer `removeTrailingSlash` test failed!")
	}
}

// This tests for removal of http/https, www and directory
func TestNormalizerRemoveDirectoryIndexAndProtocol(t *testing.T) {
	testUrls := map[string]string{
		"/":   "http://www.example.com/default.asp",
		"/a/": "http://www.example.com/a/index.html",
		"/b/": "http://https://www.example.com/b/default.html",
	}
	for path, testUrl := range testUrls {
		n, err := NewNormalizer(testUrl)
		if err != nil {
			t.Error("Error: Normalizer `RemoveDirectoryIndex` test failed!")
		}
		n.removeDefaultPort()
		n.removeDirectoryIndex()
		n.removeProtocol()
		if n.transformedUrl.Path != path &&
			n.transformedUrl.Scheme != "" &&
			n.transformedUrl.Host != "example.com" {
			fmt.Printf("path is %s, protocol is %s, host is %s \n",
				n.transformedUrl.Path,
				n.transformedUrl.Scheme,
				n.transformedUrl.Host)
			t.Error("Error: Normalizer `RemoveDirectoryIndex` test failed!")
		}
	}
}

func TestNormalizeHandleQueryParams(t *testing.T) {
	testUrls := map[string]string{
		"www.example.com?id=1234":    "http://www.example.com?id=1234&utm=a",
		"www.example.com?libid=1234": "http://www.example.com?libid=1234&utm=b",
		"www.example.com?p=1234":     "http://www.example.com?p=1234&utm=b",
		"www.example.com":            "http://www.example.com?cid=1234&id=345&libid=67890",
	}
	for result, testUrl := range testUrls {
		n, err := NewNormalizer(testUrl)
		if err != nil {
			t.Error("Error: Normalizer `HandleQueryParams` test failed!")
		}
		n.Normalize()
		if n.UrlIdentifier != result {
			fmt.Println(result, n.UrlIdentifier)
			t.Error("Error: Normalizer `HandleQueryParams` test failed!")
		}
	}
}

func TestNormalizerRemoveDuplicateSlashes(t *testing.T) {
	testUrl := "http://www.example.com/foo//bar.html"
	n, err := NewNormalizer(testUrl)
	if err != nil {
		t.Error("Error: Normalizer `RemoveDuplicateSlashes` test failed!")
	}
	n.removeDefaultPort()
	n.removeDirectoryIndex()
	n.removeProtocol()
	n.removeDuplicateSlashes()
	if n.transformedUrl.Path != "/foo/bar.html" && n.transformedUrl.Scheme != "" {
		fmt.Printf("path is %s, protocol is %s \n",
			n.transformedUrl.Path, n.transformedUrl.Scheme)
		t.Error("Error: Normalizer `RemoveDuplicateSlashes` test failed!")
	}
}

func TestNormalizerNormalize(t *testing.T) {
	testUrl := "Http%3A%2F%2FSomeUrl.com%3A8080%2Fa%2F+%2F..%2F.%2Fc%2F%2F%2Findex.html%3Fc%3D3%26a%3D1%26b%3D9%26c%3D0%23target"
	n, err := NewNormalizer(testUrl)
	if err != nil {
		t.Error("Error: Normalizer `Normalize` test failed!")
	}
	n.Normalize()
	if n.UrlIdentifier != "someurl.com/a/c" &&
		n.UrlIdentifierHash != "358af25b35e40bc8c376c4e35a7474b5" {
		fmt.Printf("normalized url: %s, urlhash: %s \n", n.UrlIdentifier, n.UrlIdentifierHash)
		t.Error("Error: Normalizer `Normalize` test failed!")
	}
}
