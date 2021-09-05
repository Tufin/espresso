package internal

import (
	"bytes"
	"embed"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func GetQuery(fs embed.FS, path string, queryName string, args []Argument) (string, error) {

	query, err := loadQueryRecursive(fs, path, queryName, args)
	if err != nil {
		log.Errorf("failed to load query with '%v'", err)
		return "", err
	}

	return query, nil
}

func loadQueryRecursive(fs embed.FS, path string, source string, args []Argument) (string, error) {

	params := map[string]string{}

	for _, arg := range args {
		query, err := loadQueryRecursive(fs, path, arg.Source, arg.Args)
		if err != nil {
			return "", err
		}
		params[arg.Name] = query
	}

	return generateSQL(fs, path, source, params)
}

func generateSQL(fs embed.FS, path string, templateName string, params map[string]string) (string, error) {

	t, err := template.New("").ParseFS(fs, path+"/"+templateName+".sql")
	if err != nil {
		log.Error(err)
		return "", err
	}

	buf := bytes.Buffer{}
	err = t.ExecuteTemplate(&buf, templateName, params)
	if err != nil {
		log.Error(err)
		return "", err
	}

	return buf.String(), nil
}
