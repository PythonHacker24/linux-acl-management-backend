package transprocessor

import (
	"path"
	"strings"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

func FindServerFromPath(filepath string) (isRemote bool, host string, port int, found bool, absolutePath string) {
	/* search through all the servers */
	for _, server := range config.BackendConfig.FileSystemServers {
		/* check if the server path has the prefix for filepath */
		if strings.HasPrefix(filepath, server.Path) {
			absolutePath := strings.TrimPrefix(filepath, server.Path)
			/* check if it's remote */
			if server.Remote != nil {
				return true, server.Remote.Host, server.Remote.Port, true, absolutePath
			}
			/* local filesystem */
			return false, "", 0, true, path.Join(config.BackendConfig.AppInfo.BasePath, filepath)
		}
	}

	/* filesystem not found */
	return false, "", 0, false, ""
}
