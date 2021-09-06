package internal

import (
	"bytes"
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

func loadQueryRecursive(fs fs.FS, templateName string, args []Argument) (string, error) {

	params := map[string]string{}

	for _, arg := range args {
		query, err := loadQueryRecursive(fs, arg.Source, arg.Args)
		if err != nil {
			return "", err
		}
		params[arg.Name] = query
	}

	return generateSQL(fs, templateName, params)
}

func generateSQL(fs fs.FS, templateName string, params map[string]string) (string, error) {

	t, err := template.New("").ParseFS(fs, "**.sql")
	if err != nil {
		t, err = template.New("").ParseFS(fs, "**/*.sql")
		if err != nil {
			log.Error(err)
			return "", err
		}
	}

	buf := bytes.Buffer{}
	err = t.ExecuteTemplate(&buf, templateName, params)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return buf.String(), nil
}
