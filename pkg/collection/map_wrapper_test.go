package collection

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMapWrapper_GetMap(t *testing.T) {
	testPlanStr := `
name: rest test plan
depends:
  - xxx
type: plan
enabled: true # default is true
environment:
  API_PREFIX: http://localhost:8080/api
  API_TOKEN: some_token
  GH_SERVER: github.com
global:
  headers:
    Authorization: Bearer ${API_TOKEN}
    Content-Type: application/json
  dataDir: ./data
suites:
  - rest_test_suite
`

	var obj map[string]any

	err := yaml.Unmarshal([]byte(testPlanStr), &obj)
	if err != nil {
		t.Fatal(err)
	}

	m := NewMapWrapper(obj)

	var target map[string]string
	err = m.Get("environment", &target)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(target)

	var target2 map[string]interface{}
	err = m.Get("environment", &target2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(target2)

	var target3 map[string]int
	err = m.Get("environment", &target3)
	if err == nil {
		t.Fatal("expect to fail")
	} else {
		t.Log(err)
	}
}

func TestMapWrapper_GetSlice(t *testing.T) {
	testPlanStr := `
name: rest test plan
depends:
  - xxx
type: plan
enabled: true # default is true
environment:
  API_PREFIX: http://localhost:8080/api
  API_TOKEN: some_token
  GH_SERVER: github.com
global:
  headers:
    Authorization: Bearer ${API_TOKEN}
    Content-Type: application/json
  dataDir: ./data
suites:
  - rest_test_suite
  - rest_test_suite2
`

	var obj map[string]any

	err := yaml.Unmarshal([]byte(testPlanStr), &obj)
	if err != nil {
		t.Fatal(err)
	}

	m := NewMapWrapper(obj)

	var target []string
	err = m.Get("suites", &target)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(target)

	var target2 = make([]string, 2)
	err = m.Get("suites", &target2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(target2)

	var target3 [2]string
	err = m.Get("suites", &target3)
	if err != nil {
		t.Fatal(err)
	}

	var target4 [1]int
	err = m.Get("suites", &target4)
	if err == nil {
		t.Fatal("expect to fail")
	} else {
		t.Log(err)
	}
}

func TestMapWrapper_Get(t *testing.T) {
	testPlanStr := `
name: rest test plan
id: 1
depends:
  - xxx
type: plan
enabled: true # default is true
environment:
  API_PREFIX: http://localhost:8080/api
  API_TOKEN: some_token
  GH_SERVER: github.com
global:
  headers:
    Authorization: Bearer ${API_TOKEN}
    Content-Type: application/json
  dataDir: ./data
suites:
  - rest_test_suite
  - rest_test_suite2
`

	var obj map[string]any

	err := yaml.Unmarshal([]byte(testPlanStr), &obj)
	if err != nil {
		t.Fatal(err)
	}

	m := NewMapWrapper(obj)

	var target string
	err = m.Get("name", &target)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(target)

	var target2 int
	err = m.Get("id", &target2)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(target2)

}
