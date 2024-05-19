package pkg

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	Production bool
	HttpPort   string
	GinDebug   bool
	KafkaUrls  []string
	O11Y       O11Y
}

type O11Y struct {
	Enable     bool
	MetricPort string
}

//

type Unmarshal func(bData []byte, v any) error

type Marshal func(v any) ([]byte, error)

func LoadJsonConfingByLocal(path string) (Config, error) {
	return LoadConfigByLocal[Config](
		json.Unmarshal,
		path,
		"CONF_PATH",
		"GARMIN2024.conf",
	)
}

func LoadConfigByLocal[T any](decode Unmarshal, path, byEnv, byHomeDirFileName string) (obj T, err error) {
	const (
		byNormal int = iota + 1
		byEnvironmentVariable
		byHomeDir
		stop
	)

	getPath := func(search int) (string, error) {
		switch search {
		case byNormal:
			return path, nil

		case byEnvironmentVariable:
			byEnvPath, _ := os.LookupEnv(byEnv)
			return byEnvPath, nil

		case byHomeDir:
			osUser, err := user.Current()
			if err != nil {
				return "", err
			}
			return filepath.Join(osUser.HomeDir, byHomeDirFileName), nil

		default:
			return "", nil
		}
	}

	var targetPath string
	var file *os.File

	for searchKind := byNormal; searchKind < stop; searchKind++ {
		targetPath, err = getPath(searchKind)
		if err != nil {
			continue
		}

		if targetPath == "" {
			err = errors.New("path is empty")
			continue
		}

		log.Printf("open config from %v", targetPath)
		file, err = os.Open(targetPath)
		if err == nil {
			break
		}
	}

	if err != nil {
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return
	}

	body := new(T)
	err = decode(data, body)
	if err != nil {
		return
	}

	return *body, err
}
