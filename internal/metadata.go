package internal

import (
	"embed"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const templateDir = "queries"

type Metadata struct {
	Name     string          `yaml:"Name"`
	Requires []string        `yaml:"Requires"`
	Params   []string        `yaml:"Params"`
	Tests    map[string]Test `yaml:"Tests"`
}

type Test struct {
	Source string     `yaml:"Source"`
	Args   []Argument `yaml:"Args"`
	Output string     `yaml:"Output"`
}

type Argument struct {
	Name   string     `yaml:"Name"`
	Source string     `yaml:"Source"`
	Args   []Argument `yaml:"Args"`
}

func getMetadata(fs embed.FS, templateName string) (Metadata, error) {
	var metadata Metadata

	bytes, err := fs.ReadFile(templateDir + "/" + templateName + ".yaml")
	if err != nil {
		log.Error(err)
		return metadata, err
	}

	if err := yaml.Unmarshal(bytes, &metadata); err != nil {
		log.Error(err)
		return metadata, err
	}

	return metadata, nil
}
