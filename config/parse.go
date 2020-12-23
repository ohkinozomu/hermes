package config

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

func ParseConfig(configFilePath string) (*Config, error) {
	fileBytes, err := fetchFileBytes(configFilePath)
	if err != nil {
		return nil, err
	}

	extension := filepath.Ext(configFilePath)
	return parseConfig(fileBytes, extension)
}

func fetchFileBytes(configFilePath string) ([]byte, error) {
	remoteUrl, err := url.Parse(configFilePath)
	if err != nil {
		return nil, err
	}

	var fileBytesFetcherFn func(string) ([]byte, error)
	scheme := remoteUrl.Scheme
	if scheme == "http" {
		fileBytesFetcherFn = downloadFile
	} else {
		// Local file
		fileBytesFetcherFn = readLocalFile
	}

	return fileBytesFetcherFn(configFilePath)
}

func downloadFile(url string) ([]byte, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(response.Body)
}

func readLocalFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(file)
}

func parseConfig(fileBytes []byte, extension string) (*Config, error) {
	var configParserFn func([]byte) (*Config, error)
	switch extension {
	case ".json":
		configParserFn = jsonParser
		break
	case ".yaml":
	case ".yml":
		configParserFn = yamlParser
	}

	if configParserFn == nil {
		return nil, errors.New("unsupported config type")
	}

	return configParserFn(fileBytes)
}

func jsonParser(fileBytes []byte) (*Config, error) {
	var config Config
	err := json.Unmarshal(fileBytes, &config)
	return &config, err
}

func yamlParser(fileBytes []byte) (*Config, error) {
	var config Config
	err := yaml.Unmarshal(fileBytes, &config)
	return &config, err
}
