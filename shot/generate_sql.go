package shot

import (
	"bytes"
	"fmt"
	"io/fs"
	"text/template"
)

// how many directories to descend when looking for templates and definitions
const depth = 5

func generateSQL(fsys fs.FS, templateName string, params map[string]string) (string, error) {

	fileName := templateName + ".sql"
	pattern := fileName
	for i := 0; i < depth; i++ {
		result, err := generateSQLInternal(fsys, pattern, templateName, params)
		if err == nil {
			return result, nil
		}
		pattern = "*/" + pattern
	}

	return "", fmt.Errorf("couldn't find template file %q", fileName)
}

func generateSQLInternal(fsys fs.FS, glob string, templateName string, params map[string]string) (string, error) {

	t, err := template.ParseFS(fsys, glob)
	if err != nil {
		return "", err
	}

	buf := bytes.Buffer{}
	err = t.ExecuteTemplate(&buf, templateName, params)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
