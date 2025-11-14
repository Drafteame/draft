package constants

import "os"

// File and directory permissions
const (
	DefaultFileMode os.FileMode = 0644 // rw-r--r-- (standard file permissions)
	DefaultDirMode  os.FileMode = 0755 // rwxr-xr-x (directories need execute bit for traversal)
)
