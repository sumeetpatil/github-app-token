package apptoken

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type RepoInfo struct {
	HostName string
	Repo     string
	Owner    string
	ApiUrl   string
	PRNumber int
}

type ClientOptions struct {
	PrivateKeyPath string
	RepoURL        string
	GithubAppID    string
}

type Client struct {
	client   *github.Client
	repoInfo *RepoInfo
}

func NewClient(ctx context.Context, options ClientOptions) (*Client, error) {
	signedToken, err := getGitHubAppToken(options.GithubAppID, options.PrivateKeyPath)
	if err != nil {
		return nil, err
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: signedToken})

	httpClient := oauth2.NewClient(ctx, src)

	repoInfo, err := getGitRepoInfo(options.RepoURL)
	if err != nil {
		return nil, err
	}

	client := &github.Client{}
	if repoInfo.HostName != "github.com" {
		client, err = newEnterpriseClient(repoInfo.ApiUrl, httpClient)
		if err != nil {
			return nil, err
		}
	} else {
		client = github.NewClient(httpClient)
	}

	return &Client{client: client, repoInfo: &repoInfo}, nil
}

func (g *Client) GetInstallationID(ctx context.Context) (*github.Installation, error) {
	installation, _, err := g.client.Apps.FindRepositoryInstallation(ctx, g.repoInfo.Owner, g.repoInfo.Repo)
	if err != nil {
		fmt.Println("Error getting installationID", err)
		return &github.Installation{}, err
	}

	return installation, err
}

func (g *Client) GetToken(ctx context.Context, installationID int64) (string, error) {
	installationToken, _, err := g.client.Apps.CreateInstallationToken(ctx, installationID, &github.InstallationTokenOptions{})
	return installationToken.GetToken(), err
}

func getGitRepoInfo(repoUri string) (RepoInfo, error) {
	repoInfo := RepoInfo{}
	pat := regexp.MustCompile(`^(https:\/\/|git@)([\S]+:[\S]+@)?([^\/:]+)[\/:]([^\/:]+\/[\S]+)$`)
	matches := pat.FindAllStringSubmatch(repoUri, -1)
	if len(matches) > 0 {
		match := matches[0]
		repoInfo.HostName = match[3]
		repoData := strings.Split(strings.TrimSuffix(match[4], ".git"), "/")
		if len(repoData) != 2 {
			return repoInfo, fmt.Errorf("Invalid repository %s", repoUri)
		}

		repoInfo.Owner = repoData[0]
		repoInfo.Repo = repoData[1]
		repoInfo.ApiUrl = "https://" + match[3] + "/api/v3"
		return repoInfo, nil
	}

	return repoInfo, fmt.Errorf("Invalid repository %s", repoUri)
}

func newEnterpriseClient(apiUrl string, httpClient *http.Client) (*github.Client, error) {
	return github.NewEnterpriseClient(
		apiUrl,
		apiUrl,
		httpClient)
}

func getGitHubAppToken(appID, privateKeyPath string) (string, error) {
	// Read the private key from the file
	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return "", err
	}

	// Create the JWT token using the app's private key
	token := jwt.New(jwt.SigningMethodRS256)
	token.Claims = jwt.MapClaims{
		"iss": appID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 10).Unix(), // Token expiration time (10 minutes)
	}

	// Parse the private key
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return "", err
	}

	// Sign the token with the private key
	signedToken, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
