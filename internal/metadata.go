package internal

import (
	"bufio"
	"io"
	"io/fs"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type Metadata struct {
	Name     string          `yaml:"Name"`
	Requires []string        `yaml:"Requires"`
	Params   []string        `yaml:"Params"`
	Tests    map[string]Test `yaml:"Tests"`
}

type Test struct {
	Args   []Argument `yaml:"Args"`
	Result Result     `yaml:"Result"`
}

type Result struct {
	Source string     `yaml:"Source"`
	Args   []Argument `yaml:"Args"`
}

type Argument struct {
	Name   string     `yaml:"Name"`
	Source string     `yaml:"Source"`
	Const  string     `yaml:"Const"`
	Args   []Argument `yaml:"Args"`
}

func GetMetadata(fs fs.FS, path string) (Metadata, error) {
	var metadata Metadata

	file, err := fs.Open(path)
	if err != nil {
		log.Error(err)
		return metadata, err
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	buf := make([]byte, 1024)
	bytes := []byte{}

	for {
		n, err := reader.Read(buf)
		if err != nil {
			if err != io.EOF {
				log.Error(err)
				return metadata, err
			}
			break
		}
		bytes = append(bytes, buf[0:n]...)
	}

	if err := yaml.Unmarshal(bytes, &metadata); err != nil {
		log.Error(err)
		return metadata, err
	}

	return metadata, nil
}
