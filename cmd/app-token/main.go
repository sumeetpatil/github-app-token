package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	appTokenClient "github.com/sumeetpatil/github-app-token"
)

func main() {
	// Load environment variables from the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	// Read the environment variables for appID and privateKeyPath
	appID := os.Getenv("GITHUB_APP_ID")
	privateKeyPath := os.Getenv("GITHUB_PRIVATE_KEY_PATH")
	gitRepoURL := os.Getenv("GITHUB_REPO_URL")

	if appID == "" || privateKeyPath == "" || gitRepoURL == "" {
		fmt.Println("Please set the GITHUB_APP_ID, GITHUB_PRIVATE_KEY_PATH and GITHUB_REPO_URL environment variables.")
		return
	}

	ctx := context.Background()
	client, err := appTokenClient.NewClient(ctx, appTokenClient.ClientOptions{PrivateKeyPath: privateKeyPath, RepoURL: gitRepoURL, GithubAppID: appID})
	if err != nil {
		fmt.Println("Error getting client:", err)
		return
	}

	installation, err := client.GetInstallationID(ctx)
	if err != nil {
		fmt.Println("Error getting installation id:", err)
		return
	}

	fmt.Printf("Installation ID : %s\n", fmt.Sprint(installation.GetID()))

	token, err := client.GetToken(ctx, installation.GetID())
	if err != nil {
		fmt.Println("Error getting token:", err)
		return
	}

	fmt.Printf("Token : %s\n", token)
}
