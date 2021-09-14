package shot

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type metadata struct {
	Name     string          `yaml:"Name"`
	Requires []string        `yaml:"Requires"`
	Params   []string        `yaml:"Params"`
	Tests    map[string]test `yaml:"Tests"`
}

type test struct {
	Args   []argument `yaml:"Args"`
	Result result     `yaml:"Result"`
}

type result struct {
	Source string     `yaml:"Source"`
	Args   []argument `yaml:"Args"`
}

type argument struct {
	Name   string     `yaml:"Name"`
	Source string     `yaml:"Source"`
	Const  string     `yaml:"Const"`
	Table  string     `yaml:"Table"`
	Args   []argument `yaml:"Args"`
}

func getMetadata(fsys fs.FS, templateName string) (metadata, error) {
	var metadata metadata

	fileName := templateName + ".yaml"
	pattern := fileName
	for i := 0; i < depth; i++ {
		list, err := fs.Glob(fsys, pattern)
		if err != nil {
			return metadata, err
		}
		if len(list) > 0 {
			return openMetadata(fsys, list[0])
		}
		pattern = "*/" + pattern
	}

	return metadata, fmt.Errorf("couldn't find definition file %q", fileName)
}

func openMetadata(fsys fs.FS, path string) (metadata, error) {

	var metadata metadata

	file, err := fsys.Open(path)
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
