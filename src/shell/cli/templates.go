package cli

import (
	"strings"
	"text/template"
)

const versionTemplate = `
{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
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
}

func NewTemplate() *Template {
	return &Template{
		version: updateTemplate(versionTemplate),
		help:    updateTemplate(helpTemplate),
	}
}
func (t *Template) SetHelp(s string) {
	t.help = updateTemplate(s)
}
func (t *Template) Help() string {
	return t.help
}

func (t *Template) SetVersion(s string) {
	t.version = updateTemplate(s)
}

func (t *Template) Version() string {
	return t.version
}

func (t *Template) Exec(c *Command, text string) (string, error) {
	type Component struct {
		Name                    string
		Version                 string
		LongHelp                string
		ShortHelp               string
		HasSubCommands          bool
		UsageString             string
		HasAvailableSubCommands bool
		IsAvailableCommand      bool
		Commands                []*Command
		CommandPath             string
	}
	z := Component{
		Name:           c.Name(),
		LongHelp:       c.longHelp,
		ShortHelp:      c.shortHelp,
		HasSubCommands: c.HasSubCommands(),
		Commands:       c.commands,
		CommandPath:    c.CommandPath(),
	}
	tpl := template.New("top")
	tpl.Funcs(_templateFuncs)
	template.Must(tpl.Parse(text))
	var b strings.Builder
	if err := tpl.Execute(&b, z); err != nil {
		return "", nil
	}
	return b.String(), nil
}
