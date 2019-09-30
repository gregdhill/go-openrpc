package generate

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
	"text/template"

	"github.com/gregdhill/go-openrpc/parse"
	"github.com/gregdhill/go-openrpc/types"
	"github.com/test-go/testify/require"
)

var exp = `package generate
type ServiceMethodParams struct {
// this is a desc
Param1 string
// this is a desc
Param2 string

Param3 string
// this is an object
Param4
}
type Param4 struct {
// this is a foo
Foo string
// this is a bar
Bar string
}
type ServiceMethodResult struct {
// the requested data
Data string
}
`

func tmpl(t *testing.T) string {
	data, err := ioutil.ReadFile("test.json")
	require.NoError(t, err)

	spec := types.NewOpenRPCSpec1()
	err = json.Unmarshal(data, spec)
	require.NoError(t, err)

	parse.GetTypes(spec, spec.Objects)

	data, err = ioutil.ReadFile("test.gotmpl")
	require.NoError(t, err)
	tmp, err := template.New("server").Funcs(funcMap(spec)).Parse(string(data))
	require.NoError(t, err)

	buf := new(bytes.Buffer)
	err = tmp.Execute(buf, spec)
	require.NoError(t, err)

	return buf.String()
}

func TestGenerate(t *testing.T) {
	require.Equal(t, exp, tmpl(t))

	t.Run("Should be deterministic", func(t *testing.T) {
		for i := 0; i < 20; i++ {
			got := tmpl(t)
			if exp == got {
				t.Log("same same")
				continue
			}
			t.Log("but different")
			require.Equal(t, exp, got)
		}
	})
}

func TestMaybeComment(t *testing.T) {
	require.Equal(t, "", maybeFieldComment(""))
	require.Equal(t, "// hello, world", maybeFieldComment("hello, world"))
}
