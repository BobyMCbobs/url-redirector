package types

type ConfigYAML struct {
	Routes Routes `yaml:"routes"`
}

type Routes map[string]string
