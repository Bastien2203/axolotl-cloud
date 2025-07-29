package utils

func IsAbsolutePath(path string) bool {
	if len(path) < 1 {
		return false
	}
	if path[0] == '/' || (len(path) > 2 && path[1] == ':' && path[2] == '\\') {
		return true
	}
	return false
}
