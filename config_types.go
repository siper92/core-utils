package core_utils

const DefaultConfigFile = ".conf.yaml"

type Module struct {
	Name     string `yaml:"name" validate:"required"`
	RootPath string `yaml:"root_path" validate:"required"`
}

type BaseCliTool struct {
	Module Module `yaml:"module" validate:"required"`
}

type BaseApp struct {
	Module Module `yaml:"module" validate:"required"`
}

type BaseAPI struct {
	Version        string   `yaml:"version"`
	Port           string   `yaml:"port" validate:"required"`
	Mode           string   `yaml:"mode" validate:"omitempty,oneof=dev debug prod"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

func (c *BaseAPI) IsDev() bool {
	return c.Mode == "dev"
}

func (c *BaseAPI) IsDebug() bool {
	return c.Mode == "debug"
}

func (c *BaseAPI) IsProd() bool {
	return c.Mode == "prod"
}

type PackageImport struct {
	Alias string `yaml:"alias"`
	Path  string `yaml:"path" validate:"required"`
}

type UseFunc struct {
	Package PackageImport `yaml:"package" validate:"required"`
	UseFunc string        `yaml:"use_func" validate:"required"`
}
