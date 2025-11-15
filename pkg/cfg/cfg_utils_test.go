package cfg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetCfgPath(t *testing.T) {
	appName := "testapp"
	cfgFile := "config.yaml"

	// Get expected values
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get user home dir: %v", err)
	}
	curDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get current dir: %v", err)
	}

	// Call the function
	cfgPaths := GetCfgPath(appName, cfgFile)

	// Expected paths
	expectedPaths := []string{
		filepath.Join(curDir, cfgFile),
		filepath.Join(homeDir, ".config", appName, cfgFile),
		filepath.Join("/etc", appName, cfgFile),
	}

	// Verify the number of paths
	if len(cfgPaths) != len(expectedPaths) {
		t.Errorf("expected %d paths, got %d", len(expectedPaths), len(cfgPaths))
	}

	// Verify each path
	for i, expected := range expectedPaths {
		if i >= len(cfgPaths) {
			t.Errorf("missing path at index %d: expected %s", i, expected)
			continue
		}
		if cfgPaths[i] != expected {
			t.Errorf("path at index %d: expected %s, got %s", i, expected, cfgPaths[i])
		}
	}
}

func TestGetCfgPathWithDifferentInputs(t *testing.T) {
	testCases := []struct {
		name     string
		appName  string
		cfgFile  string
	}{
		{
			name:    "standard config",
			appName: "myapp",
			cfgFile: "config.yaml",
		},
		{
			name:    "json config",
			appName: "webapp",
			cfgFile: "settings.json",
		},
		{
			name:    "nested config file",
			appName: "service",
			cfgFile: "conf/app.toml",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfgPaths := GetCfgPath(tc.appName, tc.cfgFile)

			// Verify we get 3 paths
			if len(cfgPaths) != 3 {
				t.Errorf("expected 3 paths, got %d", len(cfgPaths))
			}

			// Verify last path starts with /etc
			if !filepath.IsAbs(cfgPaths[2]) || cfgPaths[2][:4] != "/etc" {
				t.Errorf("last path should start with /etc, got %s", cfgPaths[2])
			}

			// Verify paths contain the app name and config file
			for i, path := range cfgPaths {
				if i > 0 { // Last two paths should contain appName
					if filepath.Base(filepath.Dir(path)) != tc.appName && filepath.Base(filepath.Dir(filepath.Dir(path))) != tc.appName {
						t.Errorf("path %d should contain app name %s, got %s", i, tc.appName, path)
					}
				}
				if filepath.Base(path) != filepath.Base(tc.cfgFile) {
					t.Errorf("path %d should end with config file %s, got %s", i, filepath.Base(tc.cfgFile), path)
				}
			}
		})
	}
}
