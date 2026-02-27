package analyzer

import (
	"testing"
)

func TestIsLoggerFunc(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"simple name", "debug", true},
		{"uppercase", "INFO", true},
		{"mixed case", "Error", true},
		{"part of word", "InfoContext", true},

		{"no match ", "hello", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLoggerFunc(tt.input)
			if result != tt.expected {
				t.Errorf("isLoggerFunc(%s): expected {%v}, got {%v}", tt.input, tt.expected, result)
			}
		})
	}
}

func TestHasFirstCapital(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"has first capital letter", "Test", true},

		{"capial letter is not first", "tesT", false},
		{"has no first capital letter", "test", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := hasFirstCapital(tt.input)
			if result != tt.expected {
				t.Errorf("hasFirstCapital(%s): expected {%v}, got {%v}", tt.input, tt.expected, result)
			}
		})
	}
}

func TestHasInvalidSymbol(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedBool bool
		expectedPos  int
	}{
		{"valid text", "server started on port 8080", false, 0},
		{"empty string", "", false, 0},

		{"invalid text: wrong language", "успешное подключение к базе данных", true, 0},
		{"invalid text: special symbols", "key: value", true, 3},
		{"invalid text: emoji", "test🔥", true, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resBool, resPos := hasInvalidSymbol(tt.input)
			if resBool != tt.expectedBool || resPos != tt.expectedPos {
				t.Errorf("hasInvalidSymbol(%s): expected {%v, %d}, got {%v, %d}",
					tt.input, tt.expectedBool, tt.expectedPos, resBool, resPos,
				)
			}
		})
	}
}

func TestHasSensetiveData(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedBool bool
		expectedPos  int
	}{
		{"valid text", "server started on port 8080", false, 0},
		{"empty string", "", false, 0},

		{"sensetive data: variable name", "userPassword", true, 4},
		{"sensetive data: variable name", "apiKey", true, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resBool, resPos := hasSensetiveData(tt.input)
			if resBool != tt.expectedBool || resPos != tt.expectedPos {
				t.Errorf("hasInvalidSymbol(%s): expected {%v, %d}, got {%v, %d}",
					tt.input, tt.expectedBool, tt.expectedPos, resBool, resPos,
				)
			}
		})
	}
}
