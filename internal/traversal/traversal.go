package traversal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"go.uber.org/zap"
)

/* comprehensive list of dangerous characters */
var (
	dangerousChars = []string{";", "|", "&", "`", "$", "(", ")", "<", ">", "{", "}", "[", "]", "\\", "'", "\""}
)

/* list files in a given directory with some basic information */
func ListFiles(path string, userID string) ([]FileEntry, error) {
	var entries []FileEntry

	/* combine basePath with the requested path */
	fullPath := filepath.Join(config.BackendConfig.AppInfo.BasePath, path)

	/* clean the path to prevent directory traversal */
	fullPath = filepath.Clean(fullPath)

	/* evaluate symlinks to get the real path */
	realPath, err := filepath.EvalSymlinks(fullPath)
	if err != nil {
		zap.L().Warn("Failed to evaluate symlinks",
			zap.String("path", fullPath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("invalid path or broken symlink: %w", err)
	}

	/* ensure the resulting path is still within the basePath (prevent directory traversal) */
	if !strings.HasPrefix(realPath, filepath.Clean(config.BackendConfig.AppInfo.BasePath)) {
		zap.L().Warn("Path traversal attempt detected",
			zap.String("path", path),
			zap.String("resolved_path", realPath),
		)
		return nil, fmt.Errorf("access denied: path outside allowed directory")
	}

	/* list all the files in the given directory */
	files, err := os.ReadDir(realPath)
	if err != nil {
		zap.L().Error("Failed to read directory",
			zap.String("path", realPath),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	/* retrieve information for each file in the directory */
	for _, f := range files {
		entryPath := filepath.Join(path, f.Name())
		fullEntryPath := filepath.Join(realPath, f.Name())

		/* evaluate symlinks for each entry */
		realEntryPath, err := filepath.EvalSymlinks(fullEntryPath)
		if err != nil {
			zap.L().Warn("Failed to evaluate symlinks for entry",
				zap.String("entry", f.Name()),
				zap.Error(err),
			)
			continue
		}

		/* verify the entry is still within allowed directory */
		if !strings.HasPrefix(realEntryPath, filepath.Clean(config.BackendConfig.AppInfo.BasePath)) {
			zap.L().Warn("Entry symlink points outside allowed directory",
				zap.String("entry", f.Name()),
				zap.String("resolved_path", realEntryPath),
			)
			continue
		}

		/* Open the file with O_NOFOLLOW to prevent symlink races */
		file, err := os.OpenFile(realEntryPath, os.O_RDONLY|syscall.O_NOFOLLOW, 0)
		if err != nil {
			zap.L().Warn("Failed to open file",
				zap.String("path", realEntryPath),
				zap.Error(err),
			)
			continue
		}
		defer file.Close()

		/* Get file descriptor for further operations */
		fd := file.Fd()

		/* check ACL access using the file descriptor */
		isOwner, err := isOwnerFd(fd, realEntryPath, userID)
		if err != nil {
			zap.L().Error("Failed to check ownership",
				zap.String("path", realEntryPath),
				zap.Error(err),
			)
			continue
		}

		if !isOwner {
			continue
		}

		/* get file information using the same file descriptor */
		info, err := file.Stat()
		if err != nil {
			zap.L().Warn("Error while getting file information",
				zap.String("path", realEntryPath),
				zap.Error(err),
			)
			continue
		}

		entries = append(entries, FileEntry{
			Name:    f.Name(),
			Path:    entryPath,
			IsDir:   info.IsDir(),
			Size:    info.Size(),
			ModTime: info.ModTime().Unix(),
		})
	}

	return entries, nil
}

/*
checks if the user is the owner of the file using a file descriptor
this reduces race conditions by operating on an already-open file
*/
func isOwnerFd(fd uintptr, filePath string, userCN string) (bool, error) {
	cleanPath := filepath.Clean(filePath)

	/* validation to ensure that the path doesn't contain dangerous characters */
	for _, char := range dangerousChars {
		if strings.Contains(cleanPath, char) {
			zap.L().Warn("Illegal character detected in file path",
				zap.String("path", cleanPath),
				zap.String("character", char),
			)
			return false, fmt.Errorf("invalid character in file path")
		}
	}

	/* get the file's ACL using getfacl with the file descriptor */
	cmd := exec.Command("getfacl", fmt.Sprintf("/proc/self/fd/%d", fd))
	output, err := cmd.Output()
	if err != nil {
		zap.L().Error("Failed to execute getfacl",
			zap.String("path", cleanPath),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to check file permissions: %w", err)
	}

	/* parse the getfacl output to check ownership */
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "# owner:") {
			owner := strings.TrimSpace(strings.TrimPrefix(line, "# owner:"))
			if strings.EqualFold(owner, userCN) {
				return true, nil
			}
		}

		if strings.HasPrefix(line, "user:") && !strings.HasPrefix(line, "user::") {
			parts := strings.Split(line, ":")
			if len(parts) >= 3 {
				aclUser := parts[1]
				permissions := parts[2]
				if strings.EqualFold(aclUser, userCN) && strings.Contains(permissions, "w") {
					return true, nil
				}
			}
		}
	}

	return false, nil
}