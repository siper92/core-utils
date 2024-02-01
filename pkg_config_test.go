package core_utils

import "testing"

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
	@import(dbConf)
  cache:
	@import(cacheConf)
  model:
	path: model

@import(./apiConf.yaml)
`

	imports := extractImports(yamlContent)

	if len(imports) != 3 {
		t.Errorf("Expected 3 imports, got %d", len(imports))
	}

	if val, ok := imports["dbConf"]; !ok || val != "${default}/dbConf.yaml" {
		t.Errorf("dbConf not found")
	}

	if val, ok := imports["cacheConf"]; !ok || val != "${default}/cacheConf.yaml" {
		t.Errorf("cacheConf not found")
	}

	if _, ok := imports["./apiConf.yaml"]; !ok {
		t.Errorf("apiConf not found")
	}
}
