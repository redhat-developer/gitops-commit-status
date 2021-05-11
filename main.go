package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/go-scm/scm/factory"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	authTokenEnvName = "GITHOSTACCESSTOKEN"
)

type commitStatus struct {
	status  string
	SHA     string
	URL     string
	path    string
	desc    string
	token   string
	context string
}

func main() {
	cs := &commitStatus{}

	cs.registerFlags()

	err := cs.checkMandatoryFlags()
	if err != nil {
		log.Fatalln(err)
	}

	if cs.token == "" {
		envToken := os.Getenv(authTokenEnvName)
		if envToken == "" {
			log.Fatalf("token not found in env variable: %s", authTokenEnvName)
		}
		cs.token = envToken
	}

	if cs.path == "" {
		cs.path, err = getRepoPath(cs.URL)
		if err != nil {
			log.Fatal(err)
		}
	}

	client, err := createSCMClient(cs.URL, cs.token)
	if err != nil {
		log.Fatalf("failed to create scm client: %v", err)
	}

	input := &scm.StatusInput{
		State: getStatus(cs.status),
		Desc:  cs.desc,
		Label: cs.context,
	}

	log.Infof("creating commit status for driver: %s repo: %s SHA: %s status: %s", client.Driver.String(), cs.path, cs.SHA, cs.status)

	_, _, err = client.Repositories.CreateStatus(context.Background(), cs.path, cs.SHA, input)
	if err != nil {
		log.Fatalf("failed to create commit status: %v", err)
	}
}

func createSCMClient(url, token string) (*scm.Client, error) {
	newURL, err := addTokenToURL(url, token)
	if err != nil {
		return nil, err
	}
	return factory.FromRepoURL(newURL)
}

func addTokenToURL(s, token string) (string, error) {
	parsedURL, err := url.Parse(s)
	if err != nil {
		return "", err
	}
	parsedURL.User = url.UserPassword("", token)
	return parsedURL.String(), nil
}

func getRepoPath(repoURL string) (string, error) {
	parsedURL, err := url.Parse(repoURL)
	if err != nil {
		return "", err
	}

	parts := strings.Split(parsedURL.Path, "/")
	components := []string{}
	for _, part := range parts {
		if part != "" {
			components = append(components, part)
		}
	}

	if len(components) < 2 {
		return "", fmt.Errorf("failed to determine repo path from URL: %s", repoURL)
	}
	return strings.Join(components, "/"), nil
}

func getStatus(state string) scm.State {
	switch strings.ToLower(state) {
	case "succeeded":
		return scm.StateSuccess
	case "failed":
		return scm.StateFailure
	default:
		return scm.StatePending
	}
}

func (cs *commitStatus) registerFlags() {
	pflag.StringVarP(&cs.status, "status", "p", "pending", "The status for a given commit")
	pflag.StringVar(&cs.SHA, "sha", "", "SHA to the set the commit status for")
	pflag.StringVar(&cs.URL, "url", "", "URL of your git repository")
	pflag.StringVar(&cs.path, "path", "", "Repository path ex: org/repo")
	pflag.StringVar(&cs.desc, "description", "", "Description of the status")
	pflag.StringVar(&cs.token, "token", "", "Personal access token")
	pflag.StringVar(&cs.context, "context", "", "Context for the status to be set")

	pflag.Parse()
}

func (cs *commitStatus) checkMandatoryFlags() error {
	missingFlags := []string{}

	if cs.URL == "" {
		missingFlags = append(missingFlags, "url")
	}

	if cs.context == "" {
		missingFlags = append(missingFlags, "context")
	}

	if cs.SHA == "" {
		missingFlags = append(missingFlags, "sha")
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("required flag(s) missing: %v", missingFlags)
	}

	return nil
}
