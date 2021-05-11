package main

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMandatoryFlags(t *testing.T) {
	t.Run("missing flags", func(t *testing.T) {
		cs := &commitStatus{}
		err := cs.checkMandatoryFlags()
		assert.EqualError(t, err, "required flag(s) missing: [url context sha]")
	})
	t.Run("no missing flags", func(t *testing.T) {
		cs := &commitStatus{
			URL:     "https://github.com/org/test",
			SHA:     "123",
			context: "CI",
		}
		assert.Nil(t, cs.checkMandatoryFlags())
	})
}

func TestGetRepoPath(t *testing.T) {
	tests := []struct {
		name   string
		url    string
		repo   string
		errMsg string
	}{
		{
			"valid repo URL",
			"https://github.com/org/test",
			"org/test",
			"",
		},
		{
			"invalid repo URL",
			"https://github.com/test",
			"org/test",
			"failed to determine repo path from URL: https://github.com/test",
		},
		{
			"custom driver path",
			"https://gitlab.us.com/proj/own/repo",
			"proj/own/repo",
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := getRepoPath(test.url)
			if test.errMsg != "" {
				assert.EqualError(t, err, test.errMsg)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, test.repo, got)
			}
		})
	}
}

func TestAddTokenToURL(t *testing.T) {
	got, err := addTokenToURL("https://github.com/org/test", "123")
	assert.Nil(t, err)

	parsedURL, err := url.Parse(got)
	assert.Nil(t, err)
	pass, _ := parsedURL.User.Password()
	assert.Equal(t, "123", pass)
}
