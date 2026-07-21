package prom

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"go.yaml.in/yaml/v3"
)

type Config struct {
	Global GlobalConfig
}

type GlobalConfig struct {
	ScrapeInterval     time.Duration     `yaml:"scrape_interval"`
	ScrapeTimeout      time.Duration     `yaml:"scrape_timeout"`
	EvaluationInterval time.Duration     `yaml:"evaluation_interval"`
	ExternalLabels     map[string]string `yaml:"external_labels"`
}

func NewPrometheusConfig() *Config {
	config := &Config{
		Global: GlobalConfig{
			ScrapeTimeout: time.Second * 10, // Default from comment
		},
	}
	return config
}

func ConfigFromYAML(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	config := NewPrometheusConfig()
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return Config{}, err
	}
	// Return copy of data only.
	return *config, nil
}

func ToYAML(a any) string {
	var buffer bytes.Buffer

	encoder := yaml.NewEncoder(&buffer)
	encoder.SetIndent(2)

	err := encoder.Encode(a)
	if err != nil {
		panic(err)
	}

	return buffer.String()
}

func main() {
	path := "test/data/sample.yaml"
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("cannot read %q because %s", path, err)
	}

	yml := ToYAML(string(data))
	fmt.Print(yml)

}
