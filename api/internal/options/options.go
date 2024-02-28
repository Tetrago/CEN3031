package options

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Options struct {
	TokenSecret      string
	Hostname         string
	Hostport         int
	SslEnabled       bool
	BasePath         string
	DatabaseHostname string
	DatabasePort     int
	DatabaseName     string
	DatabaseUsername string
	DatabasePassword string
	Port             int
	ImageFolderPath  string
}

func require[T any](v T, err error) T {
	if err != nil {
		panic(err.Error())
	} else {
		return v
	}
}

func otherwise[T any](def T) func(T, error) T {
	return func(v T, err error) T {
		if err != nil {
			return def
		} else {
			return v
		}
	}
}

func getSecret(name string) (string, error) {
	value := require(getString(name))

	if !strings.HasPrefix(value, "file://") {
		return value, nil
	}

	var path string
	fmt.Sscanf(value, "file://%s", &path)

	if data, err := os.ReadFile(path); err != nil {
		return "", err
	} else {
		return string(data), nil
	}
}

func getString(name string) (string, error) {
	if value, ok := os.LookupEnv(name); ok {
		return value, nil
	} else {
		return value, fmt.Errorf("environment variable `%s` not set", name)
	}
}

func getInt(name string) (int, error) {
	str, err := getString(name)
	if err != nil {
		return 0, err
	}

	if value, err := strconv.Atoi(str); err != nil {
		return 0, err
	} else {
		return value, nil
	}
}

func LoadFromEnvironment() Options {
	match := regexp.MustCompile(`^(https?)://([^:]+):(\d+)(/[^:]+)/?$`).FindStringSubmatch(require(getString("API_ENDPOINT_FQDN")))
	if match == nil {
		panic("Invalid endpoint FQDN!")
	}

	port, err := strconv.Atoi(match[3])
	if err != nil {
		panic("Invalid FQDN")
	}

	return Options{
		TokenSecret:      require(getSecret("API_TOKEN_SECRET")),
		Hostname:         match[2],
		Hostport:         port,
		SslEnabled:       match[1] == "https",
		BasePath:         match[4],
		DatabaseHostname: require(getString("API_DATABASE_HOSTNAME")),
		DatabasePort:     otherwise(5432)(getInt("API_DATABASE_PORT")),
		DatabaseName:     otherwise("motmot")(getString("API_DATABASE_NAME")),
		DatabaseUsername: otherwise("motmot")(getString("API_DATABASE_USERNAME")),
		DatabasePassword: require(getSecret("API_DATABASE_PASSWORD")),
		Port:             otherwise(8080)(getInt("API_PORT")),
		ImageFolderPath:  require(getSecret("API_IMAGE_FOLDER")),
	}
}
