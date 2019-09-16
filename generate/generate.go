package generate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"strings"
	"text/template"

	packr "github.com/gobuffalo/packr/v2"
	"github.com/gregdhill/go-openrpc/types"
	"github.com/gregdhill/go-openrpc/util"
)

const (
	params    = "Params"
	result    = "Result"
	goExt     = "go"
	goTmplExt = goExt + "tmpl"
)

func maybeMethodParams(method types.Method) string {
	if len(method.Params) > 0 {
		return fmt.Sprintf("%s%s", util.CamelCase(method.Name), params)
	}
	return ""
}

func maybeMethodResult(method types.Method) string {
	if method.Result != nil {
		return fmt.Sprintf("%s%s", util.CamelCase(method.Name), result)
	}
	return ""
}

func maybeMethodComment(method types.Method) string {
	if comment := util.FirstOf(method.Description, method.Summary); comment != "" {
		return fmt.Sprintf("// %s", comment)
	}
	return ""
}

func maybeFieldComment(desc string) string {
	if desc != "" {
		return fmt.Sprintf("// %s", desc)
	}
	return ""
}

type object struct {
	Name   string
	Fields *types.FieldMap
}

func funcMap(openrpc *types.OpenRPCSpec1) template.FuncMap {
	return template.FuncMap{
		"camelCase":          util.CamelCase,
		"lowerFirst":         util.LowerFirst,
		"maybeMethodComment": maybeMethodComment,
		"maybeMethodParams":  maybeMethodParams,
		"maybeMethodResult":  maybeMethodResult,
		"maybeFieldComment":  maybeFieldComment,
		"getObjects": func(om *types.ObjectMap) []object {
			keys := om.GetKeys()
			objects := make([]object, 0, len(keys))
			for _, k := range keys {
				objects = append(objects, object{k, om.Get(k)})
			}
			return objects
		},
		"getFields": func(fm *types.FieldMap) []types.BasicType {
			keys := fm.GetKeys()
			fields := make([]types.BasicType, 0, len(keys))
			for _, k := range keys {
				fields = append(fields, fm.Get(k))
			}
			return fields
		},
		"indent": func(spaces int, v string) string {
			pad := strings.Repeat(" ", spaces)
			return pad + strings.Replace(v, "\n", "\n"+pad, -1)
		},
	}
}

func WriteFile(box *packr.Box, name, pkg string, openrpc *types.OpenRPCSpec1) error {
	data, err := box.Find(fmt.Sprintf("%s.%s", name, goTmplExt))
	if err != nil {
		return err
	}

	tmp, err := template.New(name).Funcs(funcMap(openrpc)).Parse(string(data))
	if err != nil {
		return err
	}

	tmpl := new(bytes.Buffer)
	err = tmp.Execute(tmpl, openrpc)
	if err != nil {
		return err
	}

	fset := new(token.FileSet)
	root, err := parser.ParseFile(fset, "", tmpl.Bytes(), parser.ParseComments)
	if err != nil {
		return err
	}
	ast.SortImports(fset, root)
	cfg := printer.Config{Mode: printer.UseSpaces | printer.TabIndent, Tabwidth: 8}

	root.Name.Name = path.Base(pkg)

	err = os.MkdirAll(pkg, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(path.Join(pkg, fmt.Sprintf("%s.%s", name, goExt)), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return cfg.Fprint(file, fset, root)
}
