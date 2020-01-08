package types

type RouteHost struct {
	Routes   Routes `yaml:"routes"`
	Root     string `yaml:"root"`
	Wildcard string `yaml:"wildcard"`
}

type RouteHosts map[string]RouteHost

type Routes map[string]string
