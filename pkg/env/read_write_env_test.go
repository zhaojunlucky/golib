package env

import (
	"reflect"
	"testing"
)

func TestNewReadWriteEnv(t *testing.T) {
	// Create a parent env
	parentEnvs := map[string]string{
		"PARENT_KEY": "parent_value",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	// Create ReadWriteEnv with initial values
	initialEnvs := map[string]string{
		"INITIAL_KEY": "initial_value",
	}
	env := NewReadWriteEnv(parent, initialEnvs)
	
	if env == nil {
		t.Fatal("NewReadWriteEnv() returned nil")
	}
	
	// Test that initial values are set
	if env.Get("INITIAL_KEY") != "initial_value" {
		t.Errorf("Expected 'initial_value', got '%s'", env.Get("INITIAL_KEY"))
	}
	
	// Test that parent is accessible
	if env.Get("PARENT_KEY") != "parent_value" {
		t.Errorf("Expected 'parent_value', got '%s'", env.Get("PARENT_KEY"))
	}
}

func TestNewReadWriteEnv_WithNilParent(t *testing.T) {
	envs := map[string]string{
		"TEST_KEY": "test_value",
	}
	env := NewReadWriteEnv(nil, envs)
	
	if env == nil {
		t.Fatal("NewReadWriteEnv() returned nil")
	}
	
	// Should use OSEnv as parent when nil is passed
	if env.Parent == nil {
		t.Error("Expected parent to be set to OSEnv when nil is passed")
	}
}

func TestNewReadWriteEnv_WithEmptyEnvs(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	if env == nil {
		t.Fatal("NewReadWriteEnv() returned nil")
	}
	
	// Should handle nil envs gracefully
	allEnvs := env.GetAll()
	if len(allEnvs) != 0 {
		t.Errorf("Expected empty environment, got %d variables", len(allEnvs))
	}
}

func TestNewEmptyRWEnv(t *testing.T) {
	env := NewEmptyRWEnv()
	
	if env == nil {
		t.Fatal("NewEmptyRWEnv() returned nil")
	}
	
	// Should be empty with no parent
	allEnvs := env.GetAll()
	if len(allEnvs) != 0 {
		t.Errorf("Expected empty environment, got %d variables", len(allEnvs))
	}
}

func TestReadWriteEnv_Get(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	// Set a value
	env.Set("TEST_KEY", "test_value")
	
	// Test getting existing key
	if got := env.Get("TEST_KEY"); got != "test_value" {
		t.Errorf("Get('TEST_KEY') = %s, want 'test_value'", got)
	}
	
	// Test getting non-existing key
	if got := env.Get("NON_EXISTING"); got != "" {
		t.Errorf("Get('NON_EXISTING') = %s, want ''", got)
	}
}

func TestReadWriteEnv_Get_WithParent(t *testing.T) {
	parentEnvs := map[string]string{
		"PARENT_KEY": "parent_value",
		"OVERRIDE_KEY": "parent_override",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	env := NewReadWriteEnv(parent, nil)
	env.Set("CHILD_KEY", "child_value")
	env.Set("OVERRIDE_KEY", "child_override")
	
	// Test parent key access
	if got := env.Get("PARENT_KEY"); got != "parent_value" {
		t.Errorf("Get('PARENT_KEY') = %s, want 'parent_value'", got)
	}
	
	// Test child key access
	if got := env.Get("CHILD_KEY"); got != "child_value" {
		t.Errorf("Get('CHILD_KEY') = %s, want 'child_value'", got)
	}
	
	// Test override (child should override parent)
	if got := env.Get("OVERRIDE_KEY"); got != "child_override" {
		t.Errorf("Get('OVERRIDE_KEY') = %s, want 'child_override'", got)
	}
}

func TestReadWriteEnv_Contains(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	env.Set("TEST_KEY", "test_value")
	
	// Test existing key
	if !env.Contains("TEST_KEY") {
		t.Error("Contains('TEST_KEY') should return true")
	}
	
	// Test non-existing key
	if env.Contains("NON_EXISTING") {
		t.Error("Contains('NON_EXISTING') should return false")
	}
}

func TestReadWriteEnv_Contains_WithParent(t *testing.T) {
	parentEnvs := map[string]string{
		"PARENT_KEY": "parent_value",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	env := NewReadWriteEnv(parent, nil)
	env.Set("CHILD_KEY", "child_value")
	
	// Test parent key
	if !env.Contains("PARENT_KEY") {
		t.Error("Contains('PARENT_KEY') should return true")
	}
	
	// Test child key
	if !env.Contains("CHILD_KEY") {
		t.Error("Contains('CHILD_KEY') should return true")
	}
	
	// Test non-existing key
	if env.Contains("NON_EXISTING") {
		t.Error("Contains('NON_EXISTING') should return false")
	}
}

func TestReadWriteEnv_Set(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	// Test setting a value
	env.Set("NEW_KEY", "new_value")
	
	if got := env.Get("NEW_KEY"); got != "new_value" {
		t.Errorf("After Set('NEW_KEY', 'new_value'), Get('NEW_KEY') = %s, want 'new_value'", got)
	}
	
	// Test overwriting a value
	env.Set("NEW_KEY", "updated_value")
	
	if got := env.Get("NEW_KEY"); got != "updated_value" {
		t.Errorf("After Set('NEW_KEY', 'updated_value'), Get('NEW_KEY') = %s, want 'updated_value'", got)
	}
}

func TestReadWriteEnv_Set_WithExpansion(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	// Set base value
	env.Set("BASE", "hello")
	
	// Set value with expansion
	env.Set("EXPANDED", "${BASE} world")
	
	// The expansion should happen during Set
	if got := env.Get("EXPANDED"); got != "hello world" {
		t.Errorf("Get('EXPANDED') = %s, want 'hello world'", got)
	}
}

func TestReadWriteEnv_SetAll(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	envs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
		"KEY3": "value3",
	}
	
	env.SetAll(envs)
	
	// Test all values are set
	for key, expectedValue := range envs {
		if got := env.Get(key); got != expectedValue {
			t.Errorf("Get('%s') = %s, want %s", key, got, expectedValue)
		}
	}
}

func TestReadWriteEnv_SetAll_WithExpansion(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	// Set a base value first
	env.Set("PREFIX", "test")
	
	envs := map[string]string{
		"PATH1": "${PREFIX}_path1",
		"PATH2": "${PREFIX}_path2",
	}
	
	env.SetAll(envs)
	
	// Test expansion happened during SetAll
	if got := env.Get("PATH1"); got != "test_path1" {
		t.Errorf("Get('PATH1') = %s, want 'test_path1'", got)
	}
	
	if got := env.Get("PATH2"); got != "test_path2" {
		t.Errorf("Get('PATH2') = %s, want 'test_path2'", got)
	}
}

func TestReadWriteEnv_GetAll(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	envs := map[string]string{
		"KEY1": "value1",
		"KEY2": "value2",
	}
	
	env.SetAll(envs)
	
	all := env.GetAll()
	
	if len(all) != 2 {
		t.Errorf("GetAll() returned %d items, want 2", len(all))
	}
	
	for key, expectedValue := range envs {
		if got, ok := all[key]; !ok || got != expectedValue {
			t.Errorf("GetAll()[%s] = %s, want %s", key, got, expectedValue)
		}
	}
}

func TestReadWriteEnv_GetAll_WithParent(t *testing.T) {
	parentEnvs := map[string]string{
		"PARENT_KEY": "parent_value",
		"OVERRIDE_KEY": "parent_override",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	env := NewReadWriteEnv(parent, nil)
	env.Set("CHILD_KEY", "child_value")
	env.Set("OVERRIDE_KEY", "child_override")
	
	all := env.GetAll()
	
	expected := map[string]string{
		"PARENT_KEY": "parent_value",
		"CHILD_KEY": "child_value",
		"OVERRIDE_KEY": "child_override", // Child should override parent
	}
	
	if !reflect.DeepEqual(all, expected) {
		t.Errorf("GetAll() = %v, want %v", all, expected)
	}
}

func TestReadWriteEnv_Expand(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	env.Set("NAME", "World")
	env.Set("GREETING", "Hello")
	
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

func TestReadWriteEnv_Expand_WithParent(t *testing.T) {
	parentEnvs := map[string]string{
		"PARENT_VAR": "from_parent",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	env := NewReadWriteEnv(parent, nil)
	env.Set("CHILD_VAR", "from_child")
	
	// Test expansion accessing parent variable
	result := env.Expand("${PARENT_VAR} and ${CHILD_VAR}")
	expected := "from_parent and from_child"
	if result != expected {
		t.Errorf("Expand('${PARENT_VAR} and ${CHILD_VAR}') = %s, want %s", result, expected)
	}
}

func TestReadWriteEnv_ChainedExpansion(t *testing.T) {
	env := NewReadWriteEnv(NewEmptyReadEnv(), nil)
	
	env.Set("BASE", "hello")
	env.Set("MIDDLE", "${BASE}_middle")
	env.Set("FINAL", "${MIDDLE}_final")
	
	// Test chained expansion
	if got := env.Get("FINAL"); got != "hello_middle_final" {
		t.Errorf("Get('FINAL') = %s, want 'hello_middle_final'", got)
	}
}

func TestReadWriteEnv_InitialEnvsWithExpansion(t *testing.T) {
	// Create parent with base values
	parentEnvs := map[string]string{
		"ROOT": "/usr/local",
	}
	parent := NewReadEnv(NewEmptyReadEnv(), parentEnvs)
	
	// Create ReadWriteEnv with initial values that use expansion
	initialEnvs := map[string]string{
		"BIN_PATH": "${ROOT}/bin",
		"LIB_PATH": "${ROOT}/lib",
	}
	env := NewReadWriteEnv(parent, initialEnvs)
	
	// Test that expansion happened during construction
	if got := env.Get("BIN_PATH"); got != "/usr/local/bin" {
		t.Errorf("Get('BIN_PATH') = %s, want '/usr/local/bin'", got)
	}
	
	if got := env.Get("LIB_PATH"); got != "/usr/local/lib" {
		t.Errorf("Get('LIB_PATH') = %s, want '/usr/local/lib'", got)
	}
}
