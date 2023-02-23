package smtp

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// ParseTemplateDir 👇 Email template parser
func ParseTemplateDir(dir string) (*template.Template, error) {
	var paths []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	fmt.Println(paths)

	return template.ParseFiles(paths...)
}
