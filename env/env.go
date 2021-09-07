package env

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func GetEnvWithDefault(variable, defaultValue string) string {

	ret := os.Getenv(variable)
	if ret == "" {
		ret = defaultValue
	}
	log.Infof("%q: %q", variable, ret)

	return ret
}

func GetEnvOrExit(variable string) string {

	ret := os.Getenv(variable)
	if ret == "" {
		log.Fatalf("Please, set %q", variable)
	}
	log.Infof("%q: %q", variable, ret)

	return ret
}

func GetSensitive(variable string) string {

	ret := os.Getenv(variable)
	if ret != "" {
		log.Infof("%q: [sensitive]", variable)
	}

	return ret
}

func GetEnvSensitiveOrExit(variable string) string {

	ret := GetSensitive(variable)
	if ret == "" {
		log.Fatalf("Please, set %q", variable)
	}

	return ret
}
