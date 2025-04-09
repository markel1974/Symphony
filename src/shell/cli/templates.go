package cli

import (
	"io"
	"strings"
	"text/template"
)

const versionTemplate = `
{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
`

const usageTemplate = `
Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rPad .Name .NamePadding }} {{.ShortHelp}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rPad .CommandPath .CommandPathPadding}} {{.ShortHelp}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

const helpTemplate = `
{{with (or .LongHelp .ShortHelp)}}{{. | trimTrailingWhitespaces}}
{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`

func updateTemplate(in string) string {
	out := strings.Replace(in, "\r", "", -1)
	out = strings.Replace(out, "\n", "\r\n", -1)
	return out
}

var _templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"rPad":                    rPad,
}

type Template struct {
	version string
	help    string
	usage   string
}

func NewTemplate() *Template {
	return &Template{
		version: updateTemplate(versionTemplate),
		usage:   updateTemplate(usageTemplate),
		help:    updateTemplate(helpTemplate),
	}
}
func (t *Template) SetHelp(s string) {
	t.help = updateTemplate(s)
}
func (t *Template) Help() string {
	return t.help
}

func (t *Template) SetUsage(s string) {
	t.usage = updateTemplate(s)
}
func (t *Template) Usage() string {
	return t.usage
}

func (t *Template) SetVersion(s string) {
	t.version = updateTemplate(s)
}

func (t *Template) Version() string {
	return t.version
}

func (t *Template) Exec(writer io.Writer, c *Command, text string) error {
	type Component struct {
		Name                    string
		Version                 string
		LongHelp                string
		ShortHelp               string
		Runnable                bool
		HasSubCommands          bool
		UsageString             string
		HasAvailableSubCommands bool
		IsAvailableCommand      bool
		Commands                []*Command
		NamePadding             int
		CommandPath             string
		CommandPathPadding      int
	}
	z := Component{
		Name:           c.Name(),
		LongHelp:       c.LongHelp,
		ShortHelp:      c.ShortHelp,
		Runnable:       c.Runnable(),
		HasSubCommands: c.HasSubCommands(),
		//UsageString:             c.Usage(),
		HasAvailableSubCommands: c.HasAvailableSubCommands(),
		IsAvailableCommand:      c.IsAvailableCommand(),
		Commands:                c.commands,
		NamePadding:             c.NamePadding(),
		CommandPath:             c.CommandPath(),
		CommandPathPadding:      c.CommandPathPadding(),
	}
	tpl := template.New("top")
	tpl.Funcs(_templateFuncs)
	template.Must(tpl.Parse(text))
	return tpl.Execute(writer, z)
}
