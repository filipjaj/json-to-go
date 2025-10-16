package jsontogo

import (
	"strings"
	"testing"
)

func TestJSONToGo_SimpleObject(t *testing.T) {
	json := `{"name": "John", "age": 30}`
	result := JSONToGo(json, "Person", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	expected := []string{
		"type Person struct {",
		"Name string `json:\"name\"`",
		"Age int `json:\"age\"`",
	}
	
	for _, exp := range expected {
		if !strings.Contains(result.Go, exp) {
			t.Errorf("Expected output to contain %q, got:\n%s", exp, result.Go)
		}
	}
}

func TestJSONToGo_NestedObject(t *testing.T) {
	json := `{"name": "John", "address": {"street": "123 Main St", "city": "NYC"}}`
	result := JSONToGo(json, "Person", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	expected := []string{
		"type Person struct {",
		"Name string `json:\"name\"`",
		"Address Address `json:\"address\"`",
		"type Address struct {",
		"Street string `json:\"street\"`",
		"City string `json:\"city\"`",
	}
	
	for _, exp := range expected {
		if !strings.Contains(result.Go, exp) {
			t.Errorf("Expected output to contain %q, got:\n%s", exp, result.Go)
		}
	}
}

func TestJSONToGo_Array(t *testing.T) {
	json := `{"items": [1, 2, 3]}`
	result := JSONToGo(json, "Data", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	if !strings.Contains(result.Go, "Items []int `json:\"items\"`") {
		t.Errorf("Expected array of int, got:\n%s", result.Go)
	}
}

func TestJSONToGo_ArrayOfObjects(t *testing.T) {
	json := `{"users": [{"name": "John", "age": 30}, {"name": "Jane", "age": 25}]}`
	result := JSONToGo(json, "Response", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	expected := []string{
		"type Response struct {",
		"Users []Users `json:\"users\"`",
		"type Users struct {",
		"Name string `json:\"name\"`",
		"Age int `json:\"age\"`",
	}
	
	for _, exp := range expected {
		if !strings.Contains(result.Go, exp) {
			t.Errorf("Expected output to contain %q, got:\n%s", exp, result.Go)
		}
	}
}

func TestJSONToGo_MixedTypes(t *testing.T) {
	json := `{"name": "John", "age": 30, "active": true, "score": 3.14}`
	result := JSONToGo(json, "User", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	expected := []string{
		"Name string `json:\"name\"`",
		"Age int `json:\"age\"`",
		"Active bool `json:\"active\"`",
		"Score float64 `json:\"score\"`",
	}
	
	for _, exp := range expected {
		if !strings.Contains(result.Go, exp) {
			t.Errorf("Expected output to contain %q, got:\n%s", exp, result.Go)
		}
	}
}

func TestJSONToGo_Timestamp(t *testing.T) {
	json := `{"created_at": "2023-10-15T12:30:45Z"}`
	result := JSONToGo(json, "Event", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	if !strings.Contains(result.Go, "time.Time") {
		t.Errorf("Expected time.Time for timestamp, got:\n%s", result.Go)
	}
}

func TestJSONToGo_NullValue(t *testing.T) {
	json := `{"value": null}`
	result := JSONToGo(json, "Data", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	if !strings.Contains(result.Go, "any") {
		t.Errorf("Expected any for null value, got:\n%s", result.Go)
	}
}

func TestJSONToGo_NumberTypes(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{
			name:     "small integer",
			json:     `{"value": 42}`,
			expected: "int",
		},
		{
			name:     "large integer",
			json:     `{"value": 9223372036854775807}`,
			expected: "int64",
		},
		{
			name:     "float",
			json:     `{"value": 3.14}`,
			expected: "float64",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JSONToGo(tt.json, "Data", true, false, false)
			if result.Error != "" {
				t.Fatalf("Unexpected error: %s", result.Error)
			}
			if !strings.Contains(result.Go, tt.expected) {
				t.Errorf("Expected %q in output, got:\n%s", tt.expected, result.Go)
			}
		})
	}
}

func TestJSONToGo_SpecialCharacters(t *testing.T) {
	json := `{"user-name": "John", "email@address": "john@example.com"}`
	result := JSONToGo(json, "Data", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	// Should convert to valid Go identifiers (special chars are removed, resulting in single words)
	expected := []string{
		"Username",
		"Emailaddress",
	}
	
	for _, exp := range expected {
		if !strings.Contains(result.Go, exp) {
			t.Errorf("Expected output to contain %q, got:\n%s", exp, result.Go)
		}
	}
}

func TestJSONToGo_NumberAsFieldName(t *testing.T) {
	json := `{"123": "value", "1name": "test"}`
	result := JSONToGo(json, "Data", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	// Should prefix numbers to make valid Go identifiers
	if !strings.Contains(result.Go, "Num123") {
		t.Errorf("Expected Num123 in output, got:\n%s", result.Go)
	}
	// "1name" becomes "One_name" after formatNumber, then "OneName" after format()
	if !strings.Contains(result.Go, "OneName") {
		t.Errorf("Expected OneName in output, got:\n%s", result.Go)
	}
}

func TestJSONToGo_CommonInitialisms(t *testing.T) {
	json := `{"user_id": 1, "api_key": "abc", "http_url": "http://example.com"}`
	result := JSONToGo(json, "Config", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	// Should use proper Go initialisms (per Go conventions, APIKey not APIKEY since "Key" isn't an initialism)
	expected := []string{
		"UserID",
		"APIKey",
		"HTTPURL",
	}
	
	for _, exp := range expected {
		if !strings.Contains(result.Go, exp) {
			t.Errorf("Expected output to contain %q, got:\n%s", exp, result.Go)
		}
	}
}

func TestJSONToGo_Omitempty(t *testing.T) {
	json := `{"required": "value", "optional": null}`
	result := JSONToGo(json, "Data", true, false, true)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	// With allOmitempty=true, all fields should have omitempty
	if !strings.Contains(result.Go, ",omitempty") {
		t.Errorf("Expected omitempty in output, got:\n%s", result.Go)
	}
}

func TestJSONToGo_OptionalFieldsInArray(t *testing.T) {
	json := `{"items": [{"id": 1, "name": "A"}, {"id": 2}]}`
	result := JSONToGo(json, "Data", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	// Name should be optional (omitempty)
	lines := strings.Split(result.Go, "\n")
	foundNameWithOmitempty := false
	for _, line := range lines {
		if strings.Contains(line, "Name") && strings.Contains(line, "omitempty") {
			foundNameWithOmitempty = true
			break
		}
	}
	
	if !foundNameWithOmitempty {
		t.Errorf("Expected Name field to have omitempty, got:\n%s", result.Go)
	}
}

func TestJSONToGo_InvalidJSON(t *testing.T) {
	json := `{"name": invalid}`
	result := JSONToGo(json, "Data", true, false, false)
	
	if result.Error == "" {
		t.Error("Expected error for invalid JSON")
	}
}

func TestJSONToGo_EmptyObject(t *testing.T) {
	json := `{}`
	result := JSONToGo(json, "Empty", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	if !strings.Contains(result.Go, "type Empty struct {") {
		t.Errorf("Expected empty struct, got:\n%s", result.Go)
	}
}

func TestJSONToGo_FloatHack(t *testing.T) {
	// Test the .0 to .1 conversion hack
	json := `{"value": 42.0}`
	result := JSONToGo(json, "Data", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	// Should detect as float, not int
	if !strings.Contains(result.Go, "float64") {
		t.Errorf("Expected float64 for value with .0, got:\n%s", result.Go)
	}
}

func TestJSONToGo_ComplexNested(t *testing.T) {
	json := `{
		"user": {
			"name": "John",
			"contacts": [
				{
					"type": "email",
					"value": "john@example.com"
				},
				{
					"type": "phone",
					"value": "555-1234"
				}
			]
		}
	}`
	result := JSONToGo(json, "Response", true, false, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	expected := []string{
		"type Response struct {",
		"User User `json:\"user\"`",
		"type User struct {",
		"Name string `json:\"name\"`",
		"Contacts []Contacts `json:\"contacts\"`",
		"type Contacts struct {",
		"Type string `json:\"type\"`",
		"Value string `json:\"value\"`",
	}
	
	for _, exp := range expected {
		if !strings.Contains(result.Go, exp) {
			t.Errorf("Expected output to contain %q, got:\n%s", exp, result.Go)
		}
	}
}

func TestJSONToGo_Example(t *testing.T) {
	json := `{"name": "John", "age": 30}`
	result := JSONToGo(json, "User", true, true, false)
	
	if result.Error != "" {
		t.Fatalf("Unexpected error: %s", result.Error)
	}
	
	// With example=true, should include example tags
	if !strings.Contains(result.Go, "example:") {
		t.Errorf("Expected example tag in output, got:\n%s", result.Go)
	}
}

func TestFormatNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"123", "Num123"},
		{"1name", "One_name"},
		{"2things", "Two_things"},
		{"name", "name"},
		{"", ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := formatNumber(tt.input)
			if result != tt.expected {
				t.Errorf("formatNumber(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestToProperCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user_id", "UserID"},
		{"api_key", "APIKey"},
		{"http_url", "HTTPURL"},
		{"simple_name", "SimpleName"},
		{"SCREAMING_CASE", "ScreamingCase"},
	}
	
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := toProperCase(tt.input)
			if result != tt.expected {
				t.Errorf("toProperCase(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGoType(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"nil", nil, "any"},
		{"string", "hello", "string"},
		{"int", float64(42), "int"},
		{"int64", float64(9223372036854775807), "int64"},
		{"float64", 3.14, "float64"},
		{"bool", true, "bool"},
		{"timestamp", "2023-10-15T12:30:45Z", "time.Time"},
		{"array", []interface{}{1, 2, 3}, "slice"},
		{"object", map[string]interface{}{"key": "value"}, "struct"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := goType(tt.value)
			if result != tt.expected {
				t.Errorf("goType(%v) = %q, want %q", tt.value, result, tt.expected)
			}
		})
	}
}

func TestUniqueTypeName(t *testing.T) {
	seen := []string{"Name", "Value"}
	
	tests := []struct {
		name     string
		seen     []string
		prefix   string
		expected string
	}{
		{"NewName", []string{}, "", "NewName"},
		{"Name", seen, "", "Name0"},
		{"Name", seen, "User", "UserName"},
		{"Value", seen, "", "Value0"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := uniqueTypeName(tt.name, tt.seen, tt.prefix)
			if result != tt.expected {
				t.Errorf("uniqueTypeName(%q, %v, %q) = %q, want %q", 
					tt.name, tt.seen, tt.prefix, result, tt.expected)
			}
		})
	}
}
