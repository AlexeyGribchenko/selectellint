package analyzer

import (
	"testing"
)

var testChecker = func() Checker {
	c := Checker{}
	cfg := DefaultConfig()
	c.cfg = &cfg
	return c
}()

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
			result := testChecker.isLoggerFunc(tt.input)
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
			result := testChecker.hasFirstCapital(tt.input)
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
			resBool, resPos := testChecker.hasInvalidSymbol(tt.input)
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
			resBool, resPos := testChecker.hasSensetiveData(tt.input)
			if resBool != tt.expectedBool || resPos != tt.expectedPos {
				t.Errorf("hasInvalidSymbol(%s): expected {%v, %d}, got {%v, %d}",
					tt.input, tt.expectedBool, tt.expectedPos, resBool, resPos,
				)
			}
		})
	}
}

func TestFixCapital(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid text", "server started", "server started"},
		{"valid text", "sERVER STARTED", "sERVER STARTED"},
		{"invalid text", "Server started", "server started"},
		{"invalid text", "SeRVer started", "seRVer started"},
		{"empty string", "", ""},
		{"one symbol 1", "S", "s"},
		{"one symbol 2", "s", "s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := testChecker.fixCapital(tt.input)
			if res != tt.expected {
				t.Errorf("fixInvalid(%s): expected {%s}, got {%s}",
					tt.input, tt.expected, res,
				)
			}
		})
	}
}

func TestFixInvalidSymbols(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid text", "request id is", "request id is"},
		{"invalid text: wrong symbol", "te*&st!.", "test"},
		{"invalid text: wrong symbol", "te🔥st🔥", "test"},
		{"invalid text: underline", "_test_name_", "test name"},
		{"invalid text: dash", "-test-name-", "test name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := testChecker.fixInvalid(tt.input)
			if res != tt.expected {
				t.Errorf("fixInvalid(%s): expected {%s}, got {%s}",
					tt.input, tt.expected, res,
				)
			}
		})
	}
}

func TestIsSymbolValid(t *testing.T) {
	tests := []struct {
		name     string
		input    rune
		expected bool
	}{
		{"latin lowercase", 'f', true},
		{"latin uppercase", 'S', true},
		{"digit", '7', true},
		{"space", ' ', true},

		{"underscore", '_', false},
		{"hyphen", '-', false},
		{"dot", '.', false},
		{"at sign", '@', false},
		{"hash", '#', false},

		{"cyrillic lowercase", 'а', false},
		{"cyrillic uppercase", 'Б', false},
		{"japanese character", 'あ', false},

		{"tab", '\t', false},
		{"newline", '\n', false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := testChecker.isSymbolValid(tt.input)
			if result != tt.expected {
				t.Errorf("isSymbolValid(%q) = %v, want %v",
					tt.input, result, tt.expected)
			}
		})
	}
}
