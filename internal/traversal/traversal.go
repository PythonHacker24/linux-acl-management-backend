package traversal

import (
	"os"
	"path/filepath"
)

func ListFiles(path string, userID string) ([]FileEntry, error) {
    var entries []FileEntry

	/* list all the files in the given directory */
    files, err := os.ReadDir(path)
    if err != nil {
        return nil, err
    }

	/* retrive information for each file in the directory */
    for _, f := range files {
        fullPath := filepath.Join(path, f.Name())

        /* check ACL access first */
        hasAccess := checkACLAccess(fullPath, userID)
        if !hasAccess {
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
