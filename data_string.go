package core_utils

import (
	"crypto/sha256"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math/rand"
	"strconv"
	"strings"
)

func ToStringPointer(s string) *string {
	return &s
}

func ToString(s interface{}) string {
	switch val := s.(type) {
	case string:
		return val
	case *string:
		return *val
	case []byte:
		return string(val)
	case *[]byte:
		return string(*val)
	case int:
		return IntToString(val)
	case *int:
		return IntToString(*val)
	case int32, int64, uint, uint32, uint64:
		return fmt.Sprintf("%d", val)
	case *int32:
		return fmt.Sprintf("%d", *val)
	case *int64:
		return fmt.Sprintf("%d", *val)
	case *uint:
		return fmt.Sprintf("%d", *val)
	case *uint32:
		return fmt.Sprintf("%d", *val)
	case float32, float64:
		return fmt.Sprintf("%f", val)
	case *float32:
		return fmt.Sprintf("%f", *val)
	case *float64:
		return fmt.Sprintf("%f", *val)
	case bool:
		return fmt.Sprintf("%t", val)
	case *bool:
		return fmt.Sprintf("%t", *val)
	case fmt.Stringer:
		return val.String()
	}

	PrintWarningMessage("ToString: unknown type %T", s)

	return ""
}

type UrlString string

func (u UrlString) ToPath() string {
	url := string(u)
	if !strings.HasPrefix(url, "/") {
		url = "/" + url
	}

	return url
}

func StringToInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func IntToString(i int) string {
	return strconv.Itoa(i)
}

func IF[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func Sha256(data string) string {
	return RawSha256([]byte(data))
}

func RawSha256(data []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(data))
}

func HashPassword(pass string) string {
	return Sha256(pass + Sha256(pass))
}

func FirstLetterToUpper(s string) string {
	if len(s) == 0 {
		return ""
	}

	return strings.ToUpper(s[0:1]) + s[1:]
}

func CamelCaseToSnakeCase(s string) string {
	var result string
	for i, c := range s {
		if i > 0 && c >= 'A' && c <= 'Z' {
			result += "_"
		}
		result += string(c)
	}

	return strings.Trim(strings.ToLower(result), "_")
}

func SnakeCaseToCamelCase(s string) string {
	if s == "" {
		return ""
	}

	s = strings.Trim(s, "_")
	parts := strings.Split(s, "_")
	result := ""
	for _, part := range parts {
		result += FirstLetterToUpper(part)
	}
	result = strings.ToLower(result[0:1]) + result[1:]

	return result
}

func GetContentBetween(start string, end string, content string) string {
	startIndex := strings.Index(content, start)
	if startIndex == -1 {
		return ""
	}

	endIndex := strings.Index(content, end)
	if endIndex == -1 {
		return ""
	}

	return content[startIndex+len(start) : endIndex]
}

func Title(s string) string {
	return cases.Title(language.Und).String(s)
}
