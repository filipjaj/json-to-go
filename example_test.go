package jsontogo

import (
	"testing"
)

// Demonstration tests showing various JSON conversions
func TestDemonstrationSimple(t *testing.T) {
	json := `{"name": "John", "age": 30}`
	result := JSONToGo(json, "Person", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	t.Logf("Simple JSON to Go:\n%s", result.Go)
}

func TestDemonstrationNested(t *testing.T) {
	json := `{
		"user": {
			"name": "Alice",
			"email": "alice@example.com"
		},
		"status": "active"
	}`
	result := JSONToGo(json, "Response", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	t.Logf("Nested JSON to Go:\n%s", result.Go)
}

func TestDemonstrationArray(t *testing.T) {
	json := `{
		"items": [
			{"id": 1, "title": "First"},
			{"id": 2, "title": "Second"}
		]
	}`
	result := JSONToGo(json, "Data", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	t.Logf("Array JSON to Go:\n%s", result.Go)
}

func TestRealWorldJSON(t *testing.T) {
	// Test with a realistic API response
	json := `{
		"user_id": 12345,
		"username": "john_doe",
		"email": "john@example.com",
		"created_at": "2023-10-15T12:30:45Z",
		"profile": {
			"first_name": "John",
			"last_name": "Doe",
			"age": 30,
			"active": true
		},
		"posts": [
			{
				"post_id": 1,
				"title": "Hello World",
				"views": 100
			}
		],
		"metadata": {
			"ip_address": "192.168.1.1",
			"user_agent": "Mozilla/5.0"
		}
	}`
	
	result := JSONToGo(json, "APIResponse", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	// Verify key structures are present
	expectedStructs := []string{
		"type APIResponse struct {",
		"type Profile struct {",
		"type Posts struct {",
		"type Metadata struct {",
	}
	
	for _, expected := range expectedStructs {
		if !contains([]string{result.Go}, expected) {
			found := false
			// Manual substring check since we can't use strings.Contains here
			for i := 0; i <= len(result.Go)-len(expected); i++ {
				if result.Go[i:i+len(expected)] == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected to find %q in output", expected)
			}
		}
	}
	
	// Print for visual inspection
	t.Logf("Generated Go code:\n%s", result.Go)
}

func TestEdgeCases(t *testing.T) {
	testCases := []struct {
		name        string
		json        string
		typename    string
		shouldError bool
	}{
		{
			name:        "empty array",
			json:        `{"items": []}`,
			typename:    "Data",
			shouldError: false,
		},
		{
			name:        "nested arrays",
			json:        `{"matrix": [[1, 2], [3, 4]]}`,
			typename:    "Matrix",
			shouldError: false,
		},
		{
			name:        "mixed array types",
			json:        `{"values": [1, "two", 3.0]}`,
			typename:    "Mixed",
			shouldError: false,
		},
		{
			name:        "deep nesting",
			json:        `{"a": {"b": {"c": {"d": "value"}}}}`,
			typename:    "Deep",
			shouldError: false,
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := JSONToGo(tc.json, tc.typename, true, false, false)
			
			if tc.shouldError && result.Error == "" {
				t.Error("Expected error but got none")
			}
			if !tc.shouldError && result.Error != "" {
				t.Errorf("Unexpected error: %s", result.Error)
			}
			
			if !tc.shouldError && result.Go == "" {
				t.Error("Expected Go code but got empty string")
			}
			
			t.Logf("Generated:\n%s", result.Go)
		})
	}
}
