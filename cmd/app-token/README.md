# GitHub App Access Token Example

This is a simple Go program that demonstrates how to obtain an access token for a GitHub App using environment variables for the GitHub App ID, private key file path and repo url.

## Usage 
1. Rename the example.env, to .env .
2. Replace the contents in your file to 'your_github_app_id', 'path/to/your/private_key.pem' and 'repo_url' with the appropriate values for your GitHub App.
3. Compile and run the Go program:
   ```
   go run main.go
   ``` 
4. The program will use the environment variables to obtain an access token for the GitHub App.
5. The access token will be displayed in the output.

## Notes
1. Make sure to keep your private key secure and don't expose it in your code or version control.
2. The github.com/golang-jwt/jwt package provides a simple way to work with JSON Web Tokens in Go.
3. Feel free to use this code as a starting point for your GitHub App authentication in Go!