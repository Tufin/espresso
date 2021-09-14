package shot

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (shot Shot) getQuery(templateName string, args []argument) (string, error) {

	query, err := shot.loadQueryRecursive(templateName, args)
	if err != nil {
		log.Errorf("failed to load query for %q with %v", templateName, err)
		return "", err
	}

	return query, nil
}

func (shot Shot) loadQueryRecursive(templateName string, args []argument) (string, error) {

	params := map[string]string{}

	for _, arg := range args {
		param, err := shot.processArg(arg)
		if err != nil {
			return "", fmt.Errorf("invalid template %q due to %w", templateName, err)
		}
		params[arg.Name] = param
	}

	return generateSQL(shot.fsys, templateName, params)
}

func (shot Shot) processArg(arg argument) (string, error) {
	if arg.Const != "" {
		return arg.Const, nil
	}

	if arg.Table != "" {
		return shot.getTableName(arg.Table), nil
	}

	if arg.Source != "" {
		query, err := shot.loadQueryRecursive(arg.Source, arg.Args)
		if err != nil {
			return "", err
		}
		return query, nil
	}

	return "", fmt.Errorf("invalid arg %q lacks source, const and table definitions", arg.Name)
}

func (shot Shot) getTableName(table string) string {
	return shot.projectID + "." + shot.dataset + "." + table
}
