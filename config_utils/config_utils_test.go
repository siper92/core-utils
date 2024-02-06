package config_utils

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_ExtractImports(t *testing.T) {
	yamlContent := `
module:
  name: pfg/test
  rootPath: ${cwd}

outputType: file

entity:
  package:
    path: generated
  db:
	@import(db)
  cache:
	@import(cache)
  model:
	path: model

@import(./apiConf.yaml)
@import(../apiConf.yaml)
@import(http://apiConf.com/test.yaml)
@import(https://apiConf.com/test.yaml)
`

	imports := extractImports(yamlContent)

	testCases := []struct {
		key   string
		value string
	}{
		{"db", "${def}/db.yaml"},
		{"cache", "${def}/cache.yaml"},
		{"./apiConf.yaml", "./apiConf.yaml"},
		{"../apiConf.yaml", "../apiConf.yaml"},
		{"http://apiConf.com/test.yaml", "http://apiConf.com/test.yaml"},
		{"https://apiConf.com/test.yaml", "https://apiConf.com/test.yaml"},
	}
	for _, tc := range testCases {
		if val, ok := imports[tc.key]; !ok {
			t.Errorf("%s not found", tc.key)
		} else if val.value != tc.value {
			t.Errorf("Expected %s, got %s", tc.value, val)
		}
	}
}

func Test_LoadConfig(t *testing.T) {
	type TestConf struct {
		Module     Module
		OutputType string `yaml:"output_type"`
		Entity     struct {
			Package PackageImport
		}
	}
	testContent := `module:
  name: pfg/test
  root_path: ${cwd}

output_type: file

entity:
  package:
    path: "generated"`
	var conf *TestConf
	var err error

	conf, err = LoadConfigFromString(&TestConf{}, testContent)
	if err != nil {
		t.Errorf("Error loading config: %s", err)
	}

	if conf.Module.Name != "pfg/test" {
		t.Errorf("Expected pfg/test, got %s", conf.Module.Name)
	} else if conf.OutputType != "file" {
		t.Errorf("Expected file, got %s", conf.OutputType)
	} else if conf.Entity.Package.Path != "generated" {
		t.Errorf("Expected generated, got %s", conf.Entity.Package.Path)
	}
}

func Test_GetDefaultConfigPath(t *testing.T) {
	var err error

	testCase := []struct {
		expected      string
		filesToCreate []string
	}{{
		expected:      ".conf.yaml",
		filesToCreate: []string{},
	}, {
		expected:      ".conf.yaml",
		filesToCreate: []string{".conf.yaml"},
	}, {
		expected:      ".conf.local.yaml",
		filesToCreate: []string{".conf.local.yaml"},
	}, {
		expected:      ".conf.prod.yaml",
		filesToCreate: []string{".conf.prod.yaml"},
	}, {
		expected:      ".conf.demo.yaml",
		filesToCreate: []string{".conf.demo.yaml", ".conf.prod.yaml"},
	}, {
		expected:      ".conf.local.yaml",
		filesToCreate: []string{".conf.local.yaml", ".conf.prod.yaml", ".conf.demo.yaml", ".conf.dev.yaml"},
	}}

	for _, tc := range testCase {
		for _, file := range tc.filesToCreate {
			if _, err = os.Create(file); err != nil {
				t.Fatal(err)
			}
		}

		path := GetDefaultConfigPath()
		var expected string
		expected, err = filepath.Abs(tc.expected)

		if path != expected {
			t.Errorf("Expected %s, got %s", expected, path)
		}

		for _, file := range tc.filesToCreate {
			if err = os.Remove(file); err != nil {
				t.Fatal(err)
			}
		}
	}
}
