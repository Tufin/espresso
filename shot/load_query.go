package shot

import (
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

		if arg.Const != "" {
			params[arg.Name] = arg.Const
			continue
		}
		if arg.Table != "" {
			params[arg.Name] = shot.getTableName(arg.Table)
			continue
		}
		query, err := shot.loadQueryRecursive(arg.Source, arg.Args)
		if err != nil {
			return "", err
		}
		params[arg.Name] = query
	}

	return generateSQL(shot.sqlTemplates, templateName, params)
}

func (shot Shot) getTableName(table string) string {
	return shot.projectID + "." + shot.dataset + "." + table
}
