package filepathconstants

import (
	"os"
	"path/filepath"
	"runtime"
)

// FilePathConstants defines constants related to file paths.
var _, b, _, _ = runtime.Caller(0)
var FileSeparator string = string(filepath.Separator)
var ProjectPath string = filepath.Dir(filepath.Join(filepath.Dir(b), ".."))
var ProjectName string = os.Getenv("PROJECT_NAME")
var PresentDir string = filepath.Join(filepath.Dir(ProjectPath), "")
var UploadDir string = filepath.Join(PresentDir, "?same-document", "?same", FileSeparator)
