package git

import (
	"fmt"
	"os"
	"os/exec"
)

func CloneRepository(gitURL, destination string, accessToken string) (string, error) {
	// Ensure the destination directory exists, if exists remove it
	dir := fmt.Sprintf("./tmp/%s", destination)
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		if err := os.RemoveAll(dir); err != nil {
			return "", err
		}
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}

	cloneURL := gitURL
	var err error
	if accessToken != "" {
		cloneURL, err = injectTokenIntoURL(gitURL, accessToken)
		if err != nil {
			return "", err
		}
	}

	cmd := exec.Command("git", "clone", cloneURL, dir)
	return dir, cmd.Run()
}

// Supports both https://github.com/owner/repo.git and https://gitlab.com/owner/repo.git
// Insert token after https://
// e.g., https://<token>@github.com/owner/repo.git
func injectTokenIntoURL(gitURL, accessToken string) (string, error) {
	if accessToken == "" {
		return gitURL, nil
	}

	if gitURL == "" {
		return "", fmt.Errorf("git URL cannot be empty")
	}

	if len(gitURL) > 0 && gitURL[:8] == "https://" {
		return fmt.Sprintf("https://%s@%s", accessToken, gitURL[8:]), nil
	}
	return "", fmt.Errorf("unsupported git URL format: %s", gitURL)
}
