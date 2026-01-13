package gen

import "testing"

func TestIsImportable(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{"User", true},
		{"#CustomType", true},
		{"string", false},
		{"number", false},
		{"boolean", false},
		{"any", false},
	}

	for _, test := range tests {
		result := isImportable(test.typeName)
		if result != test.expected {
			t.Errorf("isImportable(%q) = %v; want %v", test.typeName, result, test.expected)
		}
	}
}

func TestIsConstant(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{"constants.UserStatus", true},
		{"constants.SomeOtherConstant", true},
		{"UserStatus", false},
		{"string", false},
		{"constants", false},
	}

	for _, test := range tests {
		result := isConstant(test.typeName)
		if result != test.expected {
			t.Errorf("isConstant(%q) = %v; want %v", test.typeName, result, test.expected)
		}
	}
}

func TestGetObjectNameAndImportPath(t *testing.T) {
	tests := []struct {
		baseType   string
		objName    string
		importPath string
	}{
		{"models.User", "User", "gen/models/User"},
		{"dtos.UserDTO", "UserDTO", "gen/dtos/UserDTO"},
		{"CustomType", "CustomType", "./CustomType"},
	}

	for _, test := range tests {
		objName, importPath := getObjectNameAndImportPath(test.baseType)
		if objName != test.objName || importPath != test.importPath {
			t.Errorf("getObjectNameAndImportPath(%q) = (%q, %q); want (%q, %q)", test.baseType, objName, importPath, test.objName, test.importPath)
		}
	}
}

func TestRequireConstructor(t *testing.T) {
	tests := []struct {
		typeName string
		expected bool
	}{
		{"User", true},
		{"[]User", false},
		{"map[string]User", true},
		{"[]*User", false},
	}

	for _, test := range tests {
		result := requireConstructor(test.typeName)
		if result != test.expected {
			t.Errorf("requireConstructor(%q) = %v; want %v", test.typeName, result, test.expected)
		}
	}
}
