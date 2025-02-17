package builder

import (
	"bytes"
	"testing"

	"github.com/roadrunner-server/velox/builder/templates"
	"github.com/stretchr/testify/require"
)

const res string = `
package container

import (
	"github.com/roadrunner-server/informer/v4"
	"github.com/roadrunner-server/resetter/v4"
	ab "github.com/roadrunner-server/rpc/v4"
	cd "github.com/roadrunner-server/http/v4"
	ef "github.com/roadrunner-server/grpc/v4"
	jk "github.com/roadrunner-server/logger/v4"
	
)

func Plugins() []any {
		return []any {
		// bundled
		// informer plugin (./rr workers, ./rr workers -i)
		&informer.Plugin{},
		// resetter plugin (./rr reset)
		&resetter.Plugin{},
	
		// std and custom plugins
		&ab.Plugin{},
		&cd.Plugin{},
		&ef.Plugin{},
		&jk.Plugin{},
		
	}
}
`

func TestCompile(t *testing.T) {
	tt := &templates.Template{
		Entries: make([]*templates.Entry, 0, 10),
	}

	tt.Entries = append(tt.Entries, &templates.Entry{
		Module:    "github.com/roadrunner-server/rpc/v4",
		Structure: "Plugin{}",
		Prefix:    "ab",
	})
	tt.Entries = append(tt.Entries, &templates.Entry{
		Module:    "github.com/roadrunner-server/http/v4",
		Structure: "Plugin{}",
		Prefix:    "cd",
	})
	tt.Entries = append(tt.Entries, &templates.Entry{
		Module:    "github.com/roadrunner-server/grpc/v4",
		Structure: "Plugin{}",
		Prefix:    "ef",
	})
	tt.Entries = append(tt.Entries, &templates.Entry{
		Module:    "github.com/roadrunner-server/logger/v4",
		Structure: "Plugin{}",
		Prefix:    "jk",
	})

	buf := new(bytes.Buffer)
	err := templates.CompileTemplateV2023(buf, tt)
	require.NoError(t, err)

	require.Equal(t, res, buf.String())
}
