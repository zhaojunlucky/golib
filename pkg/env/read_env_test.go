package env

import (
	"os"
	"reflect"
	"testing"
)

func TestNewOSEnv(t *testing.T) {
	env := NewOSEnv()
	
	if env == nil {
		t.Fatal("NewOSEnv() returned nil")
	}
	
	// Test that it contains some OS environment variables
	allEnvs := env.GetAll()
	if len(allEnvs) == 0 {
		t.Fatal("NewOSEnv() should contain OS environment variables")
	}
}

func TestNewEmptyReadEnv(t *testing.T) {
	env := NewEmptyReadEnv()
	
	if env == nil {
		t.Fatal("NewEmptyReadEnv() returned nil")
	}
	
	// Should be empty
	allEnvs := env.GetAll()
	if len(allEnvs) != 0 {
		t.Fatalf("NewEmptyReadEnv() should be empty, got %d variables", len(allEnvs))
	}
}

func TestNewReadEnv(t *testing.T) {
	// Create a parent env
	parentEnvs := map[string]string{
		"PARENT_KEY": "parent_value",
	}
	parent := NewReadEnv(nil, parentEnvs)
	
	// Create child env
	childEnvs := map[string]string{
		"CHILD_KEY": "child_value",
	}
	child := NewReadEnv(parent, childEnvs)
	
	if child == nil {
		t.Fatal("NewReadEnv() returned nil")
	}
	
	// Test that child can access both parent and child keys
	if child.Get("PARENT_KEY") != "parent_value" {
		t.Errorf("Expected 'parent_value', got '%s'", child.Get("PARENT_KEY"))
	}
	
	if child.Get("CHILD_KEY") != "child_value" {
		t.Errorf("Expected 'child_value', got '%s'", child.Get("CHILD_KEY"))
	}
}

func TestNewReadEnv_WithNilParent(t *testing.T) {
	envs := map[string]string{
		"TEST_KEY": "test_value",
	}
	env := NewReadEnv(nil, envs)
	
	if env == nil {
		t.Fatal("NewReadEnv() returned nil")
	}
	
	// Should use OSEnv as parent when nil is passed
	if env.Parent == nil {
		t.Error("Expected parent to be set to OSEnv when nil is passed")
	}
}

func TestReadEnv_Get(t *testing.T) {
	envs := map[string]string{
		"TEST_KEY1": "value1",
		"TEST_KEY2": "value2",
	}
	env := NewReadEnv(NewEmptyReadEnv(), envs)
	
	// Test existing key
	if got := env.Get("TEST_KEY1"); got != "value1" {
		t.Errorf("Get('TEST_KEY1') = %s, want 'value1'", got)
	}
	
	// Test non-existing key
	if got := env.Get("NON_EXISTING"); got != "" {
		t.Errorf("Get('NON_EXISTING') = %s, want ''", got)
	}
}

func TestReadEnv_Get_WithParent(t *testing.T) {
	parentEnvs := map[string]string{
		"PARENT_KEY": "parent_value",
		"OVERRIDE_KEY": "parent_override",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	childEnvs := map[string]string{
		"CHILD_KEY": "child_value",
		"OVERRIDE_KEY": "child_override",
	}
	child := NewReadEnv(parent, childEnvs)
	
	// Test parent key access
	if got := child.Get("PARENT_KEY"); got != "parent_value" {
		t.Errorf("Get('PARENT_KEY') = %s, want 'parent_value'", got)
	}
	
	// Test child key access
	if got := child.Get("CHILD_KEY"); got != "child_value" {
		t.Errorf("Get('CHILD_KEY') = %s, want 'child_value'", got)
	}
	
	// Test override (child should override parent)
	if got := child.Get("OVERRIDE_KEY"); got != "child_override" {
		t.Errorf("Get('OVERRIDE_KEY') = %s, want 'child_override'", got)
	}
}

func TestReadEnv_Contains(t *testing.T) {
	envs := map[string]string{
		"TEST_KEY": "value",
	}
	env := NewReadEnv(NewEmptyReadEnv(), envs)
	
	// Test existing key
	if !env.Contains("TEST_KEY") {
		t.Error("Contains('TEST_KEY') should return true")
	}
	
	// Test non-existing key
	if env.Contains("NON_EXISTING") {
		t.Error("Contains('NON_EXISTING') should return false")
	}
}

func TestReadEnv_Contains_WithParent(t *testing.T) {
	parentEnvs := map[string]string{
		"PARENT_KEY": "parent_value",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	childEnvs := map[string]string{
		"CHILD_KEY": "child_value",
	}
	child := NewReadEnv(parent, childEnvs)
	
	// Test parent key
	if !child.Contains("PARENT_KEY") {
		t.Error("Contains('PARENT_KEY') should return true")
	}
	
	// Test child key
	if !child.Contains("CHILD_KEY") {
		t.Error("Contains('CHILD_KEY') should return true")
	}
	
	// Test non-existing key
	if child.Contains("NON_EXISTING") {
		t.Error("Contains('NON_EXISTING') should return false")
	}
}

func TestReadEnv_GetAll(t *testing.T) {
	envs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	env := NewReadEnv(NewEmptyReadEnv(), envs)
	
	all := env.GetAll()
	
	if len(all) != 2 {
		t.Errorf("GetAll() returned %d items, want 2", len(all))
	}
	
	if all["KEY1"] != "value1" {
		t.Errorf("GetAll()['KEY1'] = %s, want 'value1'", all["KEY1"])
	}
	
	if all["KEY2"] != "value2" {
		t.Errorf("GetAll()['KEY2'] = %s, want 'value2'", all["KEY2"])
	}
}

func TestReadEnv_GetAll_WithParent(t *testing.T) {
	parentEnvs := map[string]string{
		"PARENT_KEY": "parent_value",
		"OVERRIDE_KEY": "parent_override",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	childEnvs := map[string]string{
		"CHILD_KEY": "child_value",
		"OVERRIDE_KEY": "child_override",
	}
	child := NewReadEnv(parent, childEnvs)
	
	all := child.GetAll()
	
	expected := map[string]string{
		"PARENT_KEY": "parent_value",
		"CHILD_KEY": "child_value",
		"OVERRIDE_KEY": "child_override", // Child should override parent
	}
	
	if !reflect.DeepEqual(all, expected) {
		t.Errorf("GetAll() = %v, want %v", all, expected)
	}
}

func TestReadEnv_Expand(t *testing.T) {
	envs := map[string]string{
		"NAME": "World",
		"GREETING": "Hello",
	}
	env := NewReadEnv(NewEmptyReadEnv(), envs)
	
	// Test simple expansion
	result := env.Expand("${GREETING}, ${NAME}!")
	expected := "Hello, World!"
	if result != expected {
		t.Errorf("Expand('${GREETING}, ${NAME}!') = %s, want %s", result, expected)
	}
	
	// Test expansion with non-existing variable
	result = env.Expand("${GREETING}, ${UNKNOWN}!")
	expected = "Hello, !"
	if result != expected {
		t.Errorf("Expand('${GREETING}, ${UNKNOWN}!') = %s, want %s", result, expected)
	}
}

func TestReadEnv_Expand_WithParent(t *testing.T) {
	parentEnvs := map[string]string{
		"PARENT_VAR": "from_parent",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	childEnvs := map[string]string{
		"CHILD_VAR": "from_child",
	}
	child := NewReadEnv(parent, childEnvs)
	
	// Test expansion accessing parent variable
	result := child.Expand("${PARENT_VAR} and ${CHILD_VAR}")
	expected := "from_parent and from_child"
	if result != expected {
		t.Errorf("Expand('${PARENT_VAR} and ${CHILD_VAR}') = %s, want %s", result, expected)
	}
}

func TestReadEnv_Expand_InConstructor(t *testing.T) {
	// Set up a parent with a base value
	parentEnvs := map[string]string{
		"BASE_PATH": "/usr/local",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	// Create child with expansion in constructor
	childEnvs := map[string]string{
		"FULL_PATH": "${BASE_PATH}/bin",
	}
	child := NewReadEnv(parent, childEnvs)
	
	// The expansion should have happened during construction
	result := child.Get("FULL_PATH")
	expected := "/usr/local/bin"
	if result != expected {
		t.Errorf("Get('FULL_PATH') = %s, want %s", result, expected)
	}
}

func TestReadEnv_Set(t *testing.T) {
	env := NewReadEnv(NewEmptyReadEnv(), nil)
	
	// Set method is empty in ReadEnv, so this should not panic
	env.Set("KEY", "value")
	
	// The key should not be set since Set is not implemented
	if env.Get("KEY") != "" {
		t.Error("Set() should not actually set values in ReadEnv")
	}
}

func TestReadEnv_SetAll(t *testing.T) {
	env := NewReadEnv(NewEmptyReadEnv(), nil)
	
	envs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	
	// SetAll method is empty in ReadEnv, so this should not panic
	env.SetAll(envs)
	
	// The keys should not be set since SetAll is not implemented
	if env.Get("KEY1") != "" || env.Get("KEY2") != "" {
		t.Error("SetAll() should not actually set values in ReadEnv")
	}
}

func TestReadEnv_initOSEnv(t *testing.T) {
	// Set a test environment variable
	testKey := "TEST_INIT_OS_ENV"
	testValue := "test_value"
	os.Setenv(testKey, testValue)
	defer os.Unsetenv(testKey)
	
	// Create a new OS env
	env := NewOSEnv()
	
	// Should contain our test variable
	if !env.Contains(testKey) {
		t.Errorf("OS env should contain %s", testKey)
	}
	
	if got := env.Get(testKey); got != testValue {
		t.Errorf("Get('%s') = %s, want %s", testKey, got, testValue)
	}
}
