package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/go-openapi/spec"
	"github.com/gobuffalo/packr/v2"
	"github.com/gregdhill/go-openrpc/types"
	"github.com/gregdhill/go-openrpc/util"
)

const (
	params    = "Params"
	result    = "Result"
	goExt     = "go"
	goTmplExt = goExt + "tmpl"
)

func schemaAsJSONPretty(s spec.Schema) string {
	b, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return ""
	}
	b = bytes.ReplaceAll(b, []byte("{"), []byte(""))
	b = bytes.ReplaceAll(b, []byte("}"), []byte(""))
	ss := string(b)
	return ss
}

func maybeLookupComponentsContentDescriptor(cmpnts *types.Components, cd *types.ContentDescriptor) (rootCD *types.ContentDescriptor, err error) {
	rootCD = cd
	if cd == nil || cmpnts == nil {
		return
	}
	if strings.Contains(cd.Schema.Ref.String(), "contentDescriptors") {
		r := filepath.Base(cd.Schema.Ref.String())
		rootCD = cmpnts.ContentDescriptors[r]
		return
	}
	return
}

func schemaHazRef(sch spec.Schema) bool {
	return sch.Ref.String() != ""
}

func derefSchemaRecurse(cts *types.Components, sch spec.Schema) spec.Schema {
	if schemaHazRef(sch) {
		sch = getSchemaFromRef(cts, sch.Ref)
		sch = derefSchemaRecurse(cts, sch)
	}
	for i := range sch.OneOf {
		desc := sch.OneOf[i].Description
		got := derefSchemaRecurse(cts, sch.OneOf[i])
		got.Description = desc
		sch.OneOf[i] = got
	}
	for i := range sch.AnyOf {
		desc := sch.AnyOf[i].Description
		got := derefSchemaRecurse(cts, sch.AnyOf[i])
		got.Description = desc
		sch.AnyOf[i] = got
	}
	for i := range sch.AllOf {
		desc := sch.AllOf[i].Description
		got := derefSchemaRecurse(cts, sch.AllOf[i])
		got.Description = desc
		sch.AllOf[i] = got
	}
	for k, _ := range sch.Properties {
		desc := sch.Properties[k].Description
		got := derefSchemaRecurse(cts, sch.Properties[k])
		got.Description = desc
		sch.Properties[k] = got
	}
	for k, _ := range sch.PatternProperties {
		desc := sch.PatternProperties[k].Description
		got := derefSchemaRecurse(cts, sch.PatternProperties[k])
		got.Description = desc
		sch.PatternProperties[k] = got
	}
	if sch.Items == nil {
		return sch
	}
	if sch.Items.Len() > 1 {
		for i := range sch.Items.Schemas {
			desc := sch.Items.Schemas[i].Description
			got := derefSchemaRecurse(cts, sch.Items.Schemas[i])
			got.Description = desc
			sch.Items.Schemas[i] = got
		}
	} else {
		// Is schema
		desc := sch.Items.Schema.Description
		got := derefSchemaRecurse(cts, *sch.Items.Schema)
		got.Description = desc
		sch.Items.Schema = &got
	}

	return sch
}

func getSchemaFromRef(cmpnts *types.Components, ref spec.Ref) (sch spec.Schema) {
	if cmpnts == nil || ref.String() == "" {
		return
	}
	r := filepath.Base(ref.String())
	sch = cmpnts.Schemas[r] // Trust parser
	return
}

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
		"derefSchema": derefSchemaRecurse,
		"schemaHasRef":       schemaHazRef,
		"schemaAsJSONPretty": schemaAsJSONPretty,
		"lookupContentDescriptor": maybeLookupComponentsContentDescriptor,
		"sanitizeBackticks":       util.SanitizeBackticks,
		"slice":                   util.Slice,
		"inspect":                 util.Inpect,
		"lengthOf":                util.LengthOf,
		"camelCase":               util.CamelCase,
		"lowerFirst":              util.LowerFirst,
		"maybeMethodComment":      maybeMethodComment,
		"maybeMethodParams":       maybeMethodParams,
		"maybeMethodResult":       maybeMethodResult,
		"maybeFieldComment":       maybeFieldComment,
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
