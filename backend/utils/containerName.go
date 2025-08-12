package utils

import "regexp"

func FormatContainerName(projectName string, containerName string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	cleanProjectName := reg.ReplaceAllString(projectName, "")
	cleanContainerName := reg.ReplaceAllString(containerName, "")

	return cleanProjectName + "_" + cleanContainerName
}
