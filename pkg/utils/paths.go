package utils

import (
    "os"
    "path/filepath"
    "runtime"
)

// ProjectRoot returns the absolute path to the project root directory
func ProjectRoot() string {
    _, filename, _, _ := runtime.Caller(0)
    dir := filepath.Dir(filename)
    return filepath.Join(dir, "../..")
}

// EnsureDirExists creates a directory if it doesn't exist
func EnsureDirExists(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return os.MkdirAll(path, 0755)
    }
    return nil
}