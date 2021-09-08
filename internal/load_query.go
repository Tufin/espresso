package internal

import (
	"io/fs"

	log "github.com/sirupsen/logrus"
)

func GetQuery(fs fs.FS, templateName string, args []Argument) (string, error) {

	query, err := loadQueryRecursive(fs, templateName, args)
	if err != nil {
		log.Errorf("failed to load query with %v", err)
		return "", err
	}

	return query, nil
}

func loadQueryRecursive(fs fs.FS, templateName string, args []Argument) (string, error) {

	params := map[string]string{}

	for _, arg := range args {
		if arg.Const != "" {
			params[arg.Name] = arg.Const
			continue
		}
		query, err := loadQueryRecursive(fs, arg.Source, arg.Args)
		if err != nil {
			return "", err
		}
		params[arg.Name] = query
	}

	return generateSQL(fs, templateName, params)
}
