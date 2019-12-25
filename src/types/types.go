package types

type ConfigYAML struct {
	Routes   Routes `yaml:"routes"`
	Root     string `yaml:"root"`
	Wildcard string `yaml:"wildcard"`
}

type Routes map[string]string
