package internal

import (
	"bytes"
	"fmt"
	"io/fs"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func GetQuery(fs fs.FS, templateName string, args []Argument) (string, error) {

	query, err := loadQueryRecursive(fs, templateName, args)
	if err != nil {
		log.Errorf("failed to load query with '%v'", err)
		return "", err
	}

	return query, nil
}

// how many directories to enter until giving up
const depth = 5

func loadQueryRecursive(fs fs.FS, templateName string, args []Argument) (string, error) {

	params := map[string]string{}

	for _, arg := range args {
		query, err := loadQueryRecursive(fs, arg.Source, arg.Args)
		if err != nil {
			return "", err
		}
		params[arg.Name] = query
	}

	fileName := templateName + ".sql"
	pattern := fileName
	for i := 0; i < depth; i++ {
		result, err := generateSQL(fs, pattern, templateName, params)
		if err == nil {
			return result, nil
		}
		pattern = "*/" + pattern
	}

	return "", fmt.Errorf("couldn't find template file '%s'", fileName)
}

func generateSQL(fs fs.FS, glob string, templateName string, params map[string]string) (string, error) {

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

// var filenames []string
// fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
// 	if d.IsDir() {
// 		log.Infof("dir %s: %v", path, d)
// 		//fs.WalkDir(fsys, "", visit)
// 	} else {
// 		log.Infof("appending %s: %v", path, d)
// 		filenames = append(filenames, d.Name())
// 	}
// 	return nil
// })

// return template.ParseFiles(filenames...)
