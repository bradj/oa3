package goserver

import (
	"bytes"
	"fmt"
	"go/format"
	"sort"
	"strings"
	"text/template"
	"unicode"

	"github.com/aarondl/oa3/generator"
	"github.com/aarondl/oa3/openapi3spec"
	"github.com/aarondl/oa3/templates"
)

const (
	// DefaultPackage name for go
	DefaultPackage = "oa3gen"
	// Disclaimer printed to the top of Go files
	Disclaimer = `// Code generated by oa3 (https://github.com/aarondl/oa3). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.`
)

// Constants for keys recognized in the parameters for the Go server
const (
	PackageKey = "package"
)

// templates for generation
var tpls = []string{
	"api_interface.tpl",
	"api_methods.tpl",
	"schema.tpl",
	"schema_top.tpl",

	"validate_schema.tpl",
	"validate_field.tpl",
}

// funcs to use for generation
var funcs = map[string]interface{}{
	"camelSnake":        camelSnake,
	"newData":           newData,
	"recurseData":       recurseData,
	"primitive":         primitive,
	"isInlinePrimitive": isInlinePrimitive,
	"taggedPaths":       tagPaths,
	"responseKind":      responseKind,
}

// generator generates templates for Go
type gen struct {
	tpl *template.Template
}

// New go generator
func New() generator.Interface {
	return &gen{}
}

// Load templates
func (g *gen) Load(dir string) error {
	var err error
	g.tpl, err = templates.Load(funcs, dir, tpls...)
	return err
}

// Do generation for Go.
func (g *gen) Do(spec *openapi3spec.OpenAPI3, params map[string]string) ([]generator.File, error) {
	if params == nil {
		params = make(map[string]string)
	}
	if pkg, ok := params[PackageKey]; !ok || len(pkg) == 0 {
		params[PackageKey] = DefaultPackage
	}

	var files []generator.File
	f, err := generateTopLevelSchemas(spec, params, g.tpl)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schemas: %w", err)
	}

	files = append(files, f...)

	f, err = generateAPIInterface(spec, params, g.tpl)
	if err != nil {
		return nil, fmt.Errorf("failed to api interface: %w", err)
	}

	files = append(files, f...)

	f, err = generateAPIMethods(spec, params, g.tpl)
	if err != nil {
		return nil, fmt.Errorf("failed to api methods: %w", err)
	}

	files = append(files, f...)

	for i, f := range files {
		formatted, err := format.Source(f.Contents)
		if err != nil {
			return nil, fmt.Errorf("failed to format file(%s): %w\n%s", f.Name, err, f.Contents)
		}

		files[i].Contents = formatted
	}

	return files, nil
}

func generateAPIMethods(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template) ([]generator.File, error) {
	if spec.Paths == nil {
		return nil, nil
	}

	return nil, nil
}

func generateAPIInterface(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template) ([]generator.File, error) {
	if spec.Paths == nil {
		return nil, nil
	}

	files := make([]generator.File, 0)

	apiName := strings.Title(strings.ReplaceAll(spec.Info.Title, " ", ""))

	tData := templates.NewTemplateData(spec, params)
	data := templateData{
		TemplateData: tData,
		Name:         apiName,
		Object:       nil,
	}

	filename := generator.FilenameFromTitle(spec.Info.Title) + ".go"

	buf := new(bytes.Buffer)
	if err := tpl.ExecuteTemplate(buf, "api_interface", data); err != nil {
		return nil, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	fileBytes := new(bytes.Buffer)
	pkg := params["package"]

	fileBytes.WriteString(Disclaimer)
	fmt.Fprintf(fileBytes, "\npackage %s\n", pkg)
	if imps := imports(data.Imports); len(imps) != 0 {
		fileBytes.WriteByte('\n')
		fileBytes.WriteString(imports(data.Imports))
		fileBytes.WriteByte('\n')
	}
	fileBytes.WriteByte('\n')
	fileBytes.Write(buf.Bytes())

	content := make([]byte, len(fileBytes.Bytes()))
	copy(content, fileBytes.Bytes())
	files = append(files, generator.File{Name: filename, Contents: content})

	tData = templates.NewTemplateData(spec, params)
	data = templateData{
		TemplateData: tData,
		Name:         apiName,
		Object:       nil,
	}

	filename = generator.FilenameFromTitle(spec.Info.Title) + "_methods.go"

	buf.Reset()
	fileBytes.Reset()
	if err := tpl.ExecuteTemplate(buf, "api_methods", data); err != nil {
		return nil, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	fileBytes.WriteString(Disclaimer)
	fmt.Fprintf(fileBytes, "\npackage %s\n", pkg)
	if imps := imports(data.Imports); len(imps) != 0 {
		fileBytes.WriteByte('\n')
		fileBytes.WriteString(imports(data.Imports))
		fileBytes.WriteByte('\n')
	}
	fileBytes.WriteByte('\n')
	fileBytes.Write(buf.Bytes())

	files = append(files, generator.File{Name: filename, Contents: fileBytes.Bytes()})

	return files, nil
}

// generateSchemas creates files for the topLevel-level referenceable types
//
// Some supported Inline are also generated.
// Prefixed with their recursive names and Inline.
// components.responses[name].headers[headername].schema
// components.responses[name].content[mime-type].schema
// components.responses[name].content[mime-type].encoding[propname].headers[headername].schema
// components.parameters[name].schema
// components.requestBodies[name].content[mime-type].schema
// components.requestBodies[name].content[mime-type].encoding[propname].headers[headername].schema
// components.headers[name].schema
// paths.parameters[0].schema
// paths.(get|put...).parameters[0].schema
// paths.(get|put...).requestBody.content[mime-type].schema
// paths.(get|put...).responses[name].headers[headername].schema
// paths.(get|put...).responses[name].content[mime-type].schema
// paths.(get|put...).responses[name].content[mime-type].encoding[propname].headers[headername].schema
func generateTopLevelSchemas(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template) ([]generator.File, error) {
	if spec.Components == nil {
		return nil, nil
	}

	keys := make([]string, 0, len(spec.Components.Schemas))
	for k := range spec.Components.Schemas {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	topLevelStructs := make([]generator.File, 0, len(keys))

	for _, k := range keys {
		v := spec.Components.Schemas[k]
		filename := "schema_" + camelSnake(k) + ".go"

		generated, err := makePseudoFile(spec, params, tpl, filename, k, v)
		if err != nil {
			return nil, err
		}
		topLevelStructs = append(topLevelStructs, generated)
	}

	type opMap struct {
		Verb string
		Op   *openapi3spec.Operation
	}

	for _, p := range spec.Paths {
		opMaps := []opMap{
			{"GET", p.Get}, {"POST", p.Post}, {"PUT", p.Put},
			{"PATCH", p.Patch}, {"TRACE", p.Trace}, {"HEAD", p.Head},
			{"DELETE", p.Delete},
		}
		for _, o := range opMaps {
			if o.Op == nil {
				continue
			}

			// If we have no request body ignore this op
			if o.Op.RequestBody != nil {
				if len(o.Op.RequestBody.Ref) != 0 {
					continue
				}

				schema := o.Op.RequestBody.Content["application/json"].Schema
				// Refs are taken care of already
				if len(schema.Ref) != 0 {
					continue
				}

				filename := "schema_" + camelSnake(o.Op.OperationID) + "_reqbody.go"
				generated, err := makePseudoFile(spec, params, tpl, filename, strings.Title(o.Op.OperationID)+"Inline", &schema)
				if err != nil {
					return nil, err
				}

				topLevelStructs = append(topLevelStructs, generated)
			}
		}

		for _, o := range opMaps {
			if o.Op == nil {
				continue
			}

			for code, resp := range o.Op.Responses {
				if len(resp.Ref) != 0 {
					continue
				}
				if len(resp.Content) == 0 {
					continue
				}

				schema := resp.Content["application/json"].Schema
				// Refs are taken care of already
				if len(schema.Ref) != 0 {
					continue
				}

				filename := "schema_" + camelSnake(o.Op.OperationID) + "_" + code + "_respbody.go"
				generated, err := makePseudoFile(spec, params, tpl, filename, strings.Title(o.Op.OperationID)+strings.Title(code)+"Inline", &schema)
				if err != nil {
					return nil, err
				}

				topLevelStructs = append(topLevelStructs, generated)
			}
		}
	}

	for name, req := range spec.Components.RequestBodies {
		schema := req.Content["application/json"].Schema
		if len(schema.Ref) != 0 {
			continue
		}

		filename := "schema_" + camelSnake(name) + "_reqbody.go"
		generated, err := makePseudoFile(spec, params, tpl, filename, name+"Inline", &schema)
		if err != nil {
			return nil, err
		}

		topLevelStructs = append(topLevelStructs, generated)
	}

	for name, resp := range spec.Components.Responses {
		if len(resp.Content) == 0 {
			continue
		}
		schema := resp.Content["application/json"].Schema
		if len(schema.Ref) != 0 {
			continue
		}

		filename := "schema_" + camelSnake(name) + "_respbody.go"
		generated, err := makePseudoFile(spec, params, tpl, filename, name+"Inline", &schema)
		if err != nil {
			return nil, err
		}
		topLevelStructs = append(topLevelStructs, generated)
	}

	return topLevelStructs, nil
}

var (
	fileBuf   = new(bytes.Buffer)
	headerBuf = new(bytes.Buffer)
)

func makePseudoFile(spec *openapi3spec.OpenAPI3, params map[string]string, tpl *template.Template, filename string, name string, schema *openapi3spec.SchemaRef) (generator.File, error) {
	fileBuf.Reset()
	headerBuf.Reset()

	tData := templates.NewTemplateData(spec, params)
	data := templateData{
		TemplateData: tData,
		Name:         name,
		Object:       schema,
	}

	if err := tpl.ExecuteTemplate(fileBuf, "schema_top", data); err != nil {
		return generator.File{}, fmt.Errorf("failed rendering template %q: %w", "schema", err)
	}

	pkg := DefaultPackage
	if pkgParam := params["package"]; len(pkgParam) > 0 {
		pkg = pkgParam
	}

	headerBuf.WriteString(Disclaimer)
	fmt.Fprintf(headerBuf, "\npackage %s\n", pkg)
	if imps := imports(data.Imports); len(imps) != 0 {
		headerBuf.WriteByte('\n')
		headerBuf.WriteString(imports(data.Imports))
		headerBuf.WriteByte('\n')
	}
	headerBuf.WriteByte('\n')

	headerLen, fileLen := headerBuf.Len(), fileBuf.Len()
	contents := make([]byte, headerLen+fileLen)
	copy(contents, headerBuf.Bytes())
	copy(contents[headerLen:], fileBuf.Bytes())

	return generator.File{Name: filename, Contents: contents}, nil
}

func isInlinePrimitive(schema *openapi3spec.Schema) bool {
	if schema.Type == "object" {
		return schema.AdditionalProperties != nil
	}

	return true
}

func primitive(tdata templateData, schema *openapi3spec.Schema) (string, error) {
	if schema.Nullable {
		return primitiveNil(tdata, schema)
	}

	return primitiveNonNil(tdata, schema)
}

func primitiveNonNil(tdata templateData, schema *openapi3spec.Schema) (string, error) {
	switch schema.Type {
	case "integer":
		if schema.Format != nil {
			switch *schema.Format {
			case "int32":
				return "int32", nil
			case "int64":
				return "int64", nil
			}
		}

		return "int", nil
	case "number":
		if schema.Format != nil {
			switch *schema.Format {
			case "float":
				return "float32", nil
			case "double":
				return "float64", nil
			}
		}

		return "float64", nil
	case "string":
		return "string", nil
	case "boolean":
		return "bool", nil
	}

	return "", fmt.Errorf("schema expected primitive type (integer, number, string, boolean) but got: %s", schema.Type)
}

func primitiveNil(tdata templateData, schema *openapi3spec.Schema) (string, error) {
	switch schema.Type {
	case "integer":
		tdata.Import("github.com/volatiletech/null/v8")

		if schema.Format != nil {
			switch *schema.Format {
			case "int32":
				return "null.Int32", nil
			case "int64":
				return "null.Int64", nil
			}
		}

		return "null.Int", nil
	case "number":
		tdata.Import("github.com/volatiletech/null/v8")

		if schema.Format != nil {
			switch *schema.Format {
			case "float":
				return "null.Float32", nil
			case "double":
				return "null.Float64", nil
			}
		}

		return "null.Float64", nil
	case "string":
		tdata.Import("github.com/volatiletech/null/v8")
		return "null.String", nil
	case "boolean":
		tdata.Import("github.com/volatiletech/null/v8")
		return "null.Bool", nil
	}

	return "", fmt.Errorf("schema had expected primitive nil type (integer, number, string, boolean) but got: %s", schema.Type)
}

func imports(imps map[string]struct{}) string {
	if len(imps) == 0 {
		return ""
	}

	var std, third []string
	for imp := range imps {
		splits := strings.Split(imp, "/")
		if len(splits) > 0 && strings.ContainsRune(splits[0], '.') {
			third = append(third, imp)
			continue
		}

		std = append(std, imp)
	}

	sort.Strings(std)
	sort.Strings(third)

	buf := new(bytes.Buffer)
	buf.WriteString("import (")
	for _, imp := range std {
		fmt.Fprintf(buf, "\n\t\"%s\"", imp)
	}
	if len(std) != 0 && len(third) != 0 {
		buf.WriteByte('\n')
	}
	for _, imp := range third {
		fmt.Fprintf(buf, "\n\t\"%s\"", imp)
	}
	buf.WriteString("\n)")

	return buf.String()
}

// schema_UserIDProfile -> schema_user_id_profile
// ID -> id
func camelSnake(filename string) string {
	build := new(strings.Builder)

	var upper bool

	in := []rune(filename)
	for i, r := range []rune(in) {
		if !unicode.IsLetter(r) {
			upper = false
			build.WriteRune(r)
			continue
		}

		if !unicode.IsUpper(r) {
			upper = false
			build.WriteRune(r)
			continue
		}

		addUnderscore := false
		if upper {
			if i+1 < len(in) && unicode.IsLower(in[i+1]) {
				addUnderscore = true
			}
		} else {
			if i-1 > 0 && unicode.IsLetter(in[i-1]) {
				addUnderscore = true
			}
		}

		if addUnderscore {
			build.WriteByte('_')
		}

		upper = true
		build.WriteRune(unicode.ToLower(r))
	}

	return build.String()
}

// responseKind returns the type of abstraction we need for a specific response
// code in an operation.
//
// The return can be one of three values: "wrapped" "empty" or ""
//
// Wrapped indicates it must be wrapped in a struct because it either has
// headers or it is a duplicate response type (say two strings) and we need
// to differentiate code.
//
// Empty means there is no response body, and an empty response code type
// must be used in its place.
//
// An empty string means that no special handling is required and the type
// response type can be used directly.
func responseKind(op *openapi3spec.Operation, code string) string {
	r := op.Responses[code]
	if len(r.Headers) != 0 {
		return "wrapped"
	}

	// Return here since there's no point continuing if we can't find bodies to
	// collide with
	if len(r.Content) == 0 {
		return "empty"
	}

	body := r.Content["application/json"].Schema

	for respCode, resp := range op.Responses {
		if respCode == code {
			continue
		}

		// Don't compare bodies if the other one doesn't have one
		if len(resp.Content) == 0 {
			continue
		}

		otherBody := resp.Content["application/json"].Schema
		if len(body.Ref) != 0 && len(otherBody.Ref) != 0 && body.Ref == otherBody.Ref {
			return "wrapped"
		} else if len(body.Ref) == 0 && len(otherBody.Ref) == 0 && body.Type == otherBody.Type {
			if isInlinePrimitive(body.Schema) {
				return "wrapped"
			}
		}
	}

	return ""
}
