package internal

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	log "github.com/sirupsen/logrus"
)

func GetQuery(fs embed.FS, queryName string, testName string) (string, error) {

	metadata, err := getMetadata(fs, queryName)
	if err != nil {
		log.Errorf("failed to get metadata with '%v'", err)
		return "", err
	}

	test, ok := metadata.Tests[testName]
	if !ok {
		err := fmt.Errorf("test '%s' undefined", testName)
		log.Error(err)
		return "", err
	}

	query, err := loadQueryRecursive(fs, test.Source, test.Args)
	if err != nil {
		log.Errorf("failed to load query with '%v'", err)
		return "", err
	}

	return query, nil
}

func generateSQL(fs embed.FS, templateName string, params map[string]string) (string, error) {

	t, err := template.New("").ParseFS(fs, templateDir+"/"+templateName+".sql")
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

func loadQueryRecursive(fs embed.FS, source string, args []Argument) (string, error) {

	params := map[string]string{}

	for _, arg := range args {
		query, err := loadQueryRecursive(fs, arg.Source, arg.Args)
		if err != nil {
			return "", err
		}
		params[arg.Name] = query
	}

	return generateSQL(fs, source, params)
}
