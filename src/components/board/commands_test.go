package board

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCommand(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		desc        string
		command     interface{}
		expectPanic bool
		expectedErr string
	}{
		{"valid function", "add", "add two numbers", func(a, b int) int { return a + b }, false, ""},
		{"missing function", "nil", "nil command", nil, true, "command is nil"},
		{"non-function command", "nonfunc", "command is not a function", 42, true, "command isn't a function"},
		{"invalid return type", "badret", "command returns multiple values", func() (int, error) { return 0, nil }, true, "command must return a single value"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectPanic {
				assert.Panics(t, func() { NewCommand(tt.id, tt.desc, tt.command) })
			} else {
				cmd := NewCommand(tt.id, tt.desc, tt.command)
				assert.NotNil(t, cmd)
				assert.Equal(t, tt.id, cmd.id)
				assert.Equal(t, tt.desc, cmd.description)
			}
		})
	}
}

func TestCommand_Exec(t *testing.T) {
	addCmd := NewCommand("add", "add integers", func(a, b int) int { return a + b })
	strCmd := NewCommand("join", "join strings", func(a, b string) string { return a + b })
	noArgCmd := NewCommand("noargs", "no arguments", func() int { return 0 })

	tests := []struct {
		name      string
		cmd       *Command
		args      []interface{}
		expect    interface{}
		expectErr bool
		errString string
	}{
		{"Valid Execution", addCmd, []interface{}{1, 2}, 3, false, ""},
		{"Wrong Arg Count", addCmd, []interface{}{1}, nil, true, "wrong number of arguments"},
		{"Wrong Arg Type", addCmd, []interface{}{1, "two"}, nil, true, "wrong argument type"},
		{"Valid String Concat", strCmd, []interface{}{"hello", " world"}, "hello world", false, ""},
		{"No arguments command", noArgCmd, []interface{}{}, 0, false, ""},
		{"No arguments command wrong", noArgCmd, []interface{}{1}, nil, true, "wrong number of arguments"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.cmd.Exec(tt.args)
			if tt.expectErr {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tt.errString, err.Error())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expect, result)
			}
		})
	}
}

func TestCommands_Add(t *testing.T) {
	mockAdd := func(a int, b int) int {
		return a + b
	}
	cmds := NewCommands()
	//Caso di successo
	err := cmds.Add("sum", "Add two numbers", mockAdd)
	assert.NoError(t, err)
	assert.True(t, cmds.Exists("sum"))
}

func TestCommands_Remove(t *testing.T) {
	mockAdd := func(a int, b int) int {
		return a + b
	}
	cmds := NewCommands()
	err := cmds.Add("sum", "Add two numbers", mockAdd)
	assert.NoError(t, err)
	assert.True(t, cmds.Exists("sum"))
	cmds.Remove("sum")
	assert.False(t, cmds.Exists("sum"))
}

func TestCommands_Exec(t *testing.T) {
	mockJoin := func(a string, b string) string {
		return a + b
	}
	cmds := NewCommands()
	_ = cmds.Add("concat", "Concatenate two strings", mockJoin)
	tests := []struct {
		name      string
		commandID string
		args      []interface{}
		expected  interface{}
		expectErr bool
	}{
		{
			name:      "Successful execution",
			commandID: "concat",
			args:      []interface{}{"Hello, ", "World!"},
			expected:  "Hello, World!",
			expectErr: false,
		},
		{
			name:      "Command not found",
			commandID: "nonexistent",
			args:      []interface{}{},
			expected:  nil,
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := cmds.Exec(tt.commandID, tt.args)
			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestCommands_PrintHelp(t *testing.T) {
	cmds := NewCommands()
	_ = cmds.Add("greet", "Greets the user", func(name string) string {
		return fmt.Sprintf("Hello, %s!", name)
	})
	_ = cmds.Add("add", "Adds two numbers", func(a, b int) int {
		return a + b
	})

	help := cmds.Documentation()
	//fmt.Printf("%v\n", help)
	assert.NotEmpty(t, help)
	assert.Len(t, help, 2)
	assert.Contains(t, help[0], "add(int, int) int : Adds two numbers")
	assert.Contains(t, help[1], "greet(string) string : Greets the user")
}
