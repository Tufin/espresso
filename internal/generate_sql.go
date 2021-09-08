package internal

import (
	"bytes"
	"fmt"
	"io/fs"
	"text/template"
)

func generateSQL(fs fs.FS, templateName string, params map[string]string) (string, error) {

	// how many directories to enter until giving up
	const depth = 5

	fileName := templateName + ".sql"
	pattern := fileName
	for i := 0; i < depth; i++ {
		result, err := generateSQLInternal(fs, pattern, templateName, params)
		if err == nil {
			return result, nil
		}
		pattern = "*/" + pattern
	}

	return "", fmt.Errorf("couldn't find template file %q", fileName)
}

func generateSQLInternal(fs fs.FS, glob string, templateName string, params map[string]string) (string, error) {

	t, err := template.ParseFS(fs, glob)
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
