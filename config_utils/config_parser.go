package config_utils

type ConfigValueType interface {
}

type ConfigValueParser[T any] interface {
	Parse(T) ConfigValueType
}
