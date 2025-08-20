package fcfs

import (
	"strings"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

func FindServerFromPath(servers []config.FileSystemServers, filepath string) (isRemote bool, host string, port int, found bool) {
	/* search through all the servers */
	for _, server := range servers {
		/* check if the server path has the prefix for filepath */
		if strings.HasPrefix(filepath, server.Path) {
			/* check if it's remote */
			if server.Remote != nil {
				return true, server.Remote.Host, server.Remote.Port, true
			}
			/* local filesystem */
			return false, "", 0, true
		}
	}

	/* filesystem not found */
	return false, "", 0, false
}
