package config_utils

import "fmt"

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
	Port           int      `yaml:"port" validate:"required"`
	Mode           string   `yaml:"mode" validate:"omitempty,oneof=dev debug prod"`
	AllowedOrigins []string `yaml:"allowed_origins"`
}

func (c BaseAPI) IsDev() bool {
	return c.Mode == "dev"
}

func (c BaseAPI) IsDebug() bool {
	return c.Mode == "debug"
}

func (c BaseAPI) IsProd() bool {
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

type ApiSecurity struct {
	JWTSecret     string `yaml:"jwt_secret" validate:"required"`
	CaptchaSecret string `yaml:"captcha_secret" validate:"required"`
}

type MySQLConfig struct {
	User     string `yaml:"user"  validate:"required"`
	Pass     string `yaml:"pass"  validate:"required"`
	Host     string `yaml:"host"  validate:"required"`
	Post     int    `yaml:"port"  validate:"required"`
	Database string `yaml:"database"  validate:"required"`
}

func (dbConf MySQLConfig) GetDNS() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConf.User,
		dbConf.Pass,
		dbConf.Host,
		dbConf.Post,
		dbConf.Database,
	)
}

type RedisConfig struct {
	Pass     string `yaml:"pass"`
	Host     string `yaml:"host" validate:"required"`
	Port     int    `yaml:"port" validate:"required"`
	Database int    `yaml:"database"`
}

func (dbConf RedisConfig) GetAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		dbConf.Host,
		dbConf.Port,
	)
}
