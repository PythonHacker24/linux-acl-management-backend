package traversal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"go.uber.org/zap"
)

/* list files in a given directory with some basic information */
func ListFiles(path string, userID string) ([]FileEntry, error) {
	var entries []FileEntry

	/* combine basePath with the requested path */
	fullPath := filepath.Join(config.BackendConfig.AppInfo.BasePath, path)

	/* clean the path to prevent directory traversal */
	fullPath = filepath.Clean(fullPath)

	/* ensure the resulting path is still within the basePath (prevent directory traversal) */
	if !strings.HasPrefix(fullPath, filepath.Clean(config.BackendConfig.AppInfo.BasePath)) {
		return nil, fmt.Errorf("Path traversal attempt detected: %s", path)
	}

	/* list all the files in the given directory */
	files, err := os.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	/* retrive information for each file in the directory */
	for _, f := range files {
		fullPath := filepath.Join(path, f.Name())

		/* check ACL access first */
		isOwner, err := isOwner(fullPath, userID)
		if err != nil {
			return nil, fmt.Errorf("error during listing files: %w", err)
		}

		if !isOwner {
			/* if the user doesn't have right ACL permissions for the file, skip it */
			continue
		}

		/* get information about the file */
		info, err := f.Info()
		if err != nil {
			continue
		}

		/* store it in entries that would be returned */
		entries = append(entries, FileEntry{
			Name:    f.Name(),
			Path:    fullPath,
			IsDir:   f.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		})
	}

	return entries, nil
}

/*
checks if the user is the owner of the file
username is the LDAP CN for the user
uses getfacl to fetch the permissions (usually from filesystems mounted from remote servers)
*/
func isOwner(filePath string, userCN string) (bool, error) {

	cleanPath := filepath.Clean(filePath)

	/* additional validation to ensure that the path doesn't contain dangerous characters */
	if strings.Contains(cleanPath, ";") || strings.Contains(cleanPath, "|") ||
		strings.Contains(cleanPath, "&") || strings.Contains(cleanPath, "`") ||
		strings.Contains(cleanPath, "$") || strings.Contains(cleanPath, "(") ||
		strings.Contains(cleanPath, ")") {
		zap.L().Warn("Illegal method attempted while getting file path by injecting dangerous character in the file path!")
		return false, fmt.Errorf("invalid characters in file path: %s", cleanPath)
	}

	/* get the file's ACL using getfacl with properly escaped arguments */
	cmd := exec.Command("getfacl", "--", cleanPath)
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to execute getfacl on %s: %v", cleanPath, err)
	}

	/* parse the getfacl output to check ownership */
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		/* check for owner line (format: "# owner: username") */
		if strings.HasPrefix(line, "# owner:") {
			owner := strings.TrimSpace(strings.TrimPrefix(line, "# owner:"))

			/* compare with the provided CN (case-insensitive) */
			if strings.EqualFold(owner, userCN) {
				return true, nil
			}
		}

		/* also check user ACL entries (format: "user:username:permissions") */
		if strings.HasPrefix(line, "user:") && !strings.HasPrefix(line, "user::") {
			parts := strings.Split(line, ":")
			if len(parts) >= 3 {
				aclUser := parts[1]
				permissions := parts[2]

				/* check if this user has write permissions (indicating ownership-like access) */
				if strings.EqualFold(aclUser, userCN) && strings.Contains(permissions, "w") {
					return true, nil
				}
			}
		}
	}

	return false, nil
}
