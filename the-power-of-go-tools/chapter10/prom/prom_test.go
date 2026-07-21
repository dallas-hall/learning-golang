package prom_test

import (
	"os"
	"prom"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestYaml_ParsePrometheusConfigCorrectly(t *testing.T) {
	t.Parallel()

	// Create our test YAML string.
	yml := `# my global config
global:
  scrape_interval: 15s
  evaluation_interval: 30s
  # scrape_timeout is set to the global default (10s).

  external_labels:
    monitor: codelab
    foo: bar
`
	want := prom.ToYAML(yml)

	// Read our sample YAML file and convert to YAML string.
	path := "test/data/sample.yaml"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to open %q because %s", path, err)
	}
	got := prom.ToYAML(string(data))

	// Compare YAML strings.
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestProm_ConfigFromYAMLCorrectlyParsesYAMLDataFromFile(t *testing.T) {
	t.Parallel()

	// Create our test data
	want := prom.Config{
		Global: prom.GlobalConfig{
			ScrapeInterval:     15 * time.Second,
			ScrapeTimeout:      10 * time.Second,
			EvaluationInterval: 30 * time.Second,
			ExternalLabels: map[string]string{
				"monitor": "codelab",
				"foo":     "bar",
			},
		},
	}

	// Get our other test data
	path := "test/data/sample.yaml"
	got, err := prom.ConfigFromYAML(path)
	if err != nil {
		t.Fatalf("failed to read %q because %s", path, err)
	}

	// Compare
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}
