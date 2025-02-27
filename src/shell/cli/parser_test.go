package cli

import (
	"testing"
)

func TestParser_Parse(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		expectError bool
		expected    []string
	}{
		{"empty line", "", false, []string{}},
		{"single word", "hello", false, []string{"hello"}},
		{"multiple words", "hello world cli", false, []string{"hello", "world", "cli"}},
		{"quoted arguments", `"hello world" test`, false, []string{"hello world", "test"}},
		{"escaped spaces", `hello\ world test`, false, []string{"hello world", "test"}},
		{"mixed quotes and escapes", `"hello\"test" cli`, false, []string{`hello"test`, "cli"}},
		{"unterminated quote", `"hello world`, true, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser()
			result, err := p.Parse(tc.line)

			if (err != nil) != tc.expectError {
				t.Fatalf("expected error: %v, got error: %v", tc.expectError, err != nil)
			}

			if !tc.expectError && len(result) != len(tc.expected) {
				t.Fatalf("expected parsed args length: %d, got: %d", len(tc.expected), len(result))
			}

			for i, v := range result {
				if v != tc.expected[i] {
					t.Errorf("expected arg[%d]: %s, got: %s", i, tc.expected[i], v)
				}
			}
		})
	}
}

func TestParser_ParseWithEnvs(t *testing.T) {
	tests := []struct {
		name         string
		line         string
		expectError  bool
		expectedEnvs []string
		expectedArgs []string
	}{
		{"no envs, single arg", "hello", false, []string{}, []string{"hello"}},
		{"single env, no args", "FOO=bar", false, []string{"FOO=bar"}, []string{}},
		{"env followed by arg", "FOO=bar hello", false, []string{"FOO=bar"}, []string{"hello"}},
		{"multiple envs and args", "FOO=bar BAR=baz cli test", false, []string{"FOO=bar", "BAR=baz"}, []string{"cli", "test"}},
		{"invalid env format", "=bar cli", false, []string{}, []string{"=bar", "cli"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			p := NewParser()
			envs, args, err := p.ParseWithEnvs(tc.line)

			if (err != nil) != tc.expectError {
				t.Fatalf("expected error: %v, got error: %v", tc.expectError, err != nil)
			}

			if !tc.expectError {
				if len(envs) != len(tc.expectedEnvs) {
					t.Fatalf("expected envs length: %d, got: %d", len(tc.expectedEnvs), len(envs))
				}
				for i, v := range envs {
					if v != tc.expectedEnvs[i] {
						t.Errorf("expected env[%d]: %s, got: %s", i, tc.expectedEnvs[i], v)
					}
				}

				if len(args) != len(tc.expectedArgs) {
					t.Fatalf("expected args length: %d, got: %d", len(tc.expectedArgs), len(args))
				}
				for i, v := range args {
					if v != tc.expectedArgs[i] {
						t.Errorf("expected arg[%d]: %s, got: %s", i, tc.expectedArgs[i], v)
					}
				}
			}
		})
	}
}

func TestReplaceEnv(t *testing.T) {
	tests := []struct {
		name     string
		envFunc  func(string) string
		input    string
		expected string
	}{
		{"empty input", nil, "", ""},
		{"no variable", nil, "hello", "hello"},
		{"single variable", func(k string) string {
			if k == "FOO" {
				return "bar"
			}
			return ""
		}, "hello $FOO", "hello bar"},
		{"nested variable", func(k string) string {
			if k == "FOO" {
				return "bar"
			}
			if k == "BAR" {
				return "baz"
			}
			return ""
		}, "nested ${FOO}_${BAR}", "nested bar_baz"},
		{"undefined variable", func(k string) string {
			if k == "FOO" {
				return "bar"
			}
			return ""
		}, "hello $UNDEFINED", "hello $UNDEFINED"},
		{"escaped dollar", nil, "price\\$100", "price$100"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := replaceEnv(tc.envFunc, tc.input)
			if result != tc.expected {
				t.Errorf("expected: %s, got: %s", tc.expected, result)
			}
		})
	}
}

func TestIsEnv(t *testing.T) {
	tests := []struct {
		name     string
		arg      string
		expected bool
	}{
		{"valid env", "FOO=bar", true},
		{"invalid env, missing =", "FOO", false},
		{"invalid env, extra =", "FOO=bar=baz", false},
		{"empty value", "FOO=", true},
		{"empty key", "=bar", false},
		{"non-alphanumeric key", "123=bar", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := isEnv(tc.arg)
			if result != tc.expected {
				t.Errorf("expected: %v, got: %v", tc.expected, result)
			}
		})
	}
}
