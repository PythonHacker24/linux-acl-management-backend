package traversal

/*
file entry contains basic information about a file
this information is displayed in the traversal view of the frontend
*/
type FileEntry struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"is_dir"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mod_time"`
}

/* request for listing files in a given directory path */
type ListRequest struct {
	FilePath string `json:"file_path"`
}
