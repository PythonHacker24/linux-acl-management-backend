package traversal

/*
	file entry contains basic information about a file
	this information is displayed in the traversal view of the frontend
*/
type FileEntry struct {
    Name      string `json:"name"`
    Path      string `json:"path"`
    IsDir     bool   `json:"isDir"`
    Size      int64  `json:"size"`
    ModTime   int64  `json:"modTime"`
}
