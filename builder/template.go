package builder

import (
	"bytes"
	"text/template"
)

const GoModTemplate string = `
module github.com/roadrunner-server/roadrunner/v2

go 1.19

require (
        github.com/buger/goterm v1.0.4
        github.com/dustin/go-humanize v1.0.1
        github.com/joho/godotenv v1.4.0
        github.com/olekukonko/tablewriter v0.0.5
        github.com/spf13/cobra v1.6.1
		github.com/spf13/viper v1.14.0
        github.com/stretchr/testify v1.8.1
		go.uber.org/automaxprocs v1.5.1
)

replace (
	{{range $v := .Entries}}{{if (ne $v.Replace "")}}{{$v.Module}} => {{$v.Replace}}
	{{end}}{{end}}
)
`

const PluginsTemplate string = `
package container

import (
	"github.com/roadrunner-server/informer/v3"
	"github.com/roadrunner-server/resetter/v3"
	{{range $v := .Entries}}{{$v.Prefix}} "{{$v.Module}}"
	{{end}}
)

func Plugins() []any {
		return []any {
		// bundled
		// informer plugin (./rr workers, ./rr workers -i)
		&informer.Plugin{},
		// resetter plugin (./rr reset)
		&resetter.Plugin{},
	
		// std and custom plugins
		{{range $v := .Entries}}&{{$v.Prefix}}.{{$v.Structure}},
		{{end}}
	}
}
`

// Entry represents all info about module
type Entry struct {
	Module    string
	Structure string
	Prefix    string
	Version   string
	// Replace directive, should include path
	Replace string
}

type Template struct {
	Entries []*Entry
}

func compileTemplate(buf *bytes.Buffer, data *Template) error {
	tmplt, err := template.New("plugins.go").Parse(PluginsTemplate)
	if err != nil {
		return err
	}

	err = tmplt.Execute(buf, data)
	if err != nil {
		return err
	}

	return nil
}

func compileGoModTemplate(buf *bytes.Buffer, data *Template) error {
	tmplt, err := template.New("go.mod").Parse(GoModTemplate)
	if err != nil {
		return err
	}

	err = tmplt.Execute(buf, data)
	if err != nil {
		return err
	}

	return nil
}
