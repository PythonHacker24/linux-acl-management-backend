package transprocessor

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
)

/* maintains locks on file which are actively under ACL modifications */
var pathLocks sync.Map

/* locks a given file path */
func getPathLock(path string) *sync.Mutex {
	mtx, _ := pathLocks.LoadOrStore(path, &sync.Mutex{})
	return mtx.(*sync.Mutex)
}

/* handles local transaction execution (change permissions via mounts) */
func (p *PermProcessor) HandleLocalTransaction(txn *types.Transaction, absolutePath string) error {
	aclEntry := BuildACLEntry(txn.Entries)

	/* REMOVE THIS */
	zap.L().Info("ACL Request recieved",
		zap.String("Transaction ID", txn.ID.String()),
		zap.String("Action", txn.Entries.Action),
		zap.String("Entry", aclEntry),
		zap.String("Path", absolutePath),
	)

	/* lock the file path for thread safety (ensure unlock even on panic) */
	lock := getPathLock(absolutePath)
	lock.Lock()
	defer lock.Unlock()

	/* execute the ACL modifications with acl commands */
	var cmd *exec.Cmd
	switch txn.Entries.Action {
	case "add", "modify":
		cmd = exec.Command("setfacl", "-m", aclEntry, absolutePath)
	case "remove":
		cmd = exec.Command("setfacl", "-x", aclEntry, absolutePath)
	default:
		// sendResponse(conn, false, "Unsupported action: "+req.Action)
		txn.ErrorMsg = fmt.Sprintf("unsupported ACL action: %s", txn.Entries.Action)
	}

	start := time.Now()

	output, err := cmd.CombinedOutput()

	duration := time.Since(start).Milliseconds()

	txn.Output = string(output)
	txn.DurationMs = duration

	if err != nil {
		/* status of transaction is successful but execution failed */
		txn.Status = types.StatusSuccess
		txn.ExecStatus = false
		txn.ErrorMsg = err.Error()

		/* REMOVE THIS */
		zap.L().Error("ACL command failed",
			zap.String("Transaction ID", txn.ID.String()),
			zap.String("Output", string(output)),
			zap.Error(err),
		)

		txn.ErrorMsg = fmt.Sprintf("setfacl failed: %w, output: %s", err, output)
	}

	txn.Status = types.StatusSuccess
	txn.ExecStatus = true

	/* REMOVE THIS */
	zap.L().Info("ACL command executed successfully",
		zap.String("Transaction ID", txn.ID.String()),
		zap.String("Output", string(output)),
	)

	return nil
}

/* builds the ACL entry string for setfacl */
func BuildACLEntry(entry types.ACLEntry) string {
	var sb strings.Builder

	if entry.IsDefault {
		sb.WriteString("default:")
	}

	sb.WriteString(entry.EntityType)
	sb.WriteString(":")
	sb.WriteString(entry.Entity)
	sb.WriteString(":")
	sb.WriteString(entry.Permissions)

	return sb.String()
}
