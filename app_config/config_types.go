package app_config

const DefaultConfigFile = ".conf.yaml"

type Module struct {
	Name     string `yaml:"name" validate:"required"`
	RootPath string `yaml:"root_path" validate:"required"`
}

type PackageImport struct {
	Alias string `yaml:"alias"`
	Path  string `yaml:"path" validate:"required"`
}

type UsedPackages struct {
	Package PackageImport `yaml:"package" validate:"required"`
	UseStr  string        `yaml:"useStr" validate:"required"`
}

type APIConfig struct {
	Version        string   `yaml:"version"`
	Port           string   `yaml:"port" validate:"required"`
	Mode           string   `yaml:"mode" validate:"omitempty,oneof=dev debug prod"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}
