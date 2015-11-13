package genmain

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/raphael/goa/design"
	"github.com/raphael/goa/goagen/codegen"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Generator is the application code generator.
type Generator struct {
	genfiles []string
}

// Generate is the generator entry point called by the meta generator.
func Generate(api *design.APIDefinition) ([]string, error) {
	g, err := NewGenerator()
	if err != nil {
		return nil, err
	}
	return g.Generate(api)
}

// NewGenerator returns the application code generator.
func NewGenerator() (*Generator, error) {
	app := kingpin.New("Main generator", "application main generator")
	codegen.RegisterFlags(app)
	NewCommand().RegisterFlags(app)
	_, err := app.Parse(os.Args[1:])
	if err != nil {
		return nil, fmt.Errorf(`invalid command line: %s. Command line was "%s"`,
			err, strings.Join(os.Args, " "))
	}
	return new(Generator), nil
}

// Generate produces the skeleton main.
func (g *Generator) Generate(api *design.APIDefinition) ([]string, error) {
	mainFile := filepath.Join(codegen.OutputDir, "main.go")
	if Force {
		os.Remove(mainFile)
	}
	_, err := os.Stat(mainFile)
	funcs := template.FuncMap{
		"tempvar":            tempvar,
		"generateJSONSchema": generateJSONSchema,
		"goify":              codegen.Goify,
		"okResp":             okResp,
	}
	if err != nil {
		tmpl, err := template.New("main").Funcs(funcs).Parse(mainTmpl)
		if err != nil {
			panic(err.Error()) // bug
		}
		gg := codegen.NewGoGenerator(mainFile)
		g.genfiles = []string{mainFile}
		outPkg, err := filepath.Rel(os.Getenv("GOPATH"), codegen.OutputDir)
		if err != nil {
			return nil, err
		}
		outPkg = strings.TrimPrefix(outPkg, "src/")
		appPkg := filepath.Join(outPkg, "app")
		imports := []*codegen.ImportSpec{
			codegen.SimpleImport("github.com/raphael/goa"),
			codegen.SimpleImport(appPkg),
			codegen.NewImport("log", "gopkg.in/inconshreveable/log15.v2"),
		}
		if generateJSONSchema() {
			jsonSchemaPkg := filepath.Join(outPkg, "schema")
			imports = append(imports, codegen.SimpleImport(jsonSchemaPkg))
		}
		gg.WriteHeader("", "main", imports)
		data := map[string]interface{}{
			"Name":      AppName,
			"Resources": api.Resources,
		}
		err = tmpl.Execute(gg, data)
		if err != nil {
			g.Cleanup()
			return nil, err
		}
		if err := gg.FormatCode(); err != nil {
			g.Cleanup()
			return nil, err
		}
	}
	tmpl, err := template.New("ctrl").Funcs(funcs).Parse(ctrlTmpl)
	if err != nil {
		panic(err.Error()) // bug
	}
	imp, err := filepath.Rel(filepath.Join(os.Getenv("GOPATH"), "src"), codegen.OutputDir)
	if err != nil {
		return nil, err
	}
	imp = filepath.Join(imp, "app")
	imports := []*codegen.ImportSpec{codegen.SimpleImport(imp)}
	err = api.IterateResources(func(r *design.ResourceDefinition) error {
		filename := filepath.Join(codegen.OutputDir, snakeCase(r.Name)+".go")
		if Force {
			if err := os.Remove(filename); err != nil {
				return err
			}
		}
		if _, err := os.Stat(filename); err != nil {
			resGen := codegen.NewGoGenerator(filename)
			g.genfiles = append(g.genfiles, filename)
			resGen.WriteHeader("", "main", imports)
			err := tmpl.Execute(resGen, r)
			if err != nil {
				g.Cleanup()
				return err
			}
			if err := resGen.FormatCode(); err != nil {
				g.Cleanup()
				return err
			}
		}
		return nil
	})
	if err != nil {
		g.Cleanup()
		return nil, err
	}

	return g.genfiles, nil
}

// Cleanup removes all the files generated by this generator during the last invokation of Generate.
func (g *Generator) Cleanup() {
	for _, f := range g.genfiles {
		os.Remove(f)
	}
	g.genfiles = nil
}

// tempCount is the counter used to create unique temporary variable names.
var tempCount int

// tempvar generates a unique temp var name.
func tempvar() string {
	tempCount++
	if tempCount == 1 {
		return "c"
	}
	return fmt.Sprintf("c%d", tempCount)
}

// generateJSONSchema returns true if the API JSON schema should be generated.
func generateJSONSchema() bool {
	return codegen.CommandName == "" || codegen.CommandName == "schema"
}

func okResp(a *design.ActionDefinition) map[string]interface{} {
	var ok *design.ResponseDefinition
	for _, resp := range a.Responses {
		if resp.Status == 200 {
			ok = resp
			break
		}
	}
	if ok == nil {
		return nil
	}
	var mt *design.MediaTypeDefinition
	var ok2 bool
	if mt, ok2 = design.Design.MediaTypes[ok.MediaType]; !ok2 {
		return nil
	}
	typeref := codegen.GoTypeRef(mt, 1)
	if strings.HasPrefix(typeref, "*") {
		typeref = "&app." + typeref[1:]
	} else {
		typeref = "app." + typeref
	}
	return map[string]interface{}{
		"Name":             ok.Name,
		"HasMultipleViews": len(mt.Views) > 1,
		"GoType":           codegen.GoNativeType(mt),
		"TypeRef":          typeref,
	}
}

// snakeCase produces the snake_case version of the given CamelCase string.
func snakeCase(name string) string {
	var b bytes.Buffer
	var lastUnderscore bool
	ln := len(name)
	if ln == 0 {
		return ""
	}
	b.WriteRune(unicode.ToLower(rune(name[0])))
	for i := 1; i < ln; i++ {
		r := rune(name[i])
		if unicode.IsUpper(r) {
			if !lastUnderscore {
				b.WriteRune('_')
				lastUnderscore = true
			}
			b.WriteRune(unicode.ToLower(r))
		} else {
			b.WriteRune(r)
			lastUnderscore = false
		}
	}
	return b.String()
}

const mainTmpl = `
func main() {
	// Create service
	api := goa.New("{{.Name}}")

	// Setup middleware
	api.Use(goa.RequestID())
	api.Use(goa.LogRequest())
	api.Use(goa.Recover())

{{range $name, $res := .Resources}}	// Mount "{{$res.Name}}" controller
	{{$tmp := tempvar}}{{$tmp}} := New{{goify $res.Name true}}Controller()
	app.Mount{{goify $res.Name true}}Controller(api, {{$tmp}})
{{end}}{{if generateJSONSchema}}
	// Mount Swagger spec provider controller
	swagger.MountController(api)
{{end}}
	// Start service, listen on port 8080
	api.ListenAndServe(":8080")
}
`
const ctrlTmpl = `// {{$ctrlName := printf "%s%s" (goify .Name true) "Controller"}}{{$ctrlName}} implements the {{.Name}} resource.
type {{$ctrlName}} struct {}

// New{{$ctrlName}} creates a {{.Name}} controller.
func New{{$ctrlName}}() *{{$ctrlName}} {
	return &{{$ctrlName}}{}
}
{{$ctrl := .}}{{range .Actions}}
// {{goify .Name true}} runs the {{.Name}} action.
func (c *{{$ctrlName}}) {{goify .Name true}}(ctx *app.{{goify .Name true}}{{goify $ctrl.Name true}}Context) error {
{{$ok := okResp .}}{{if $ok}}	res := {{$ok.TypeRef}}{}
{{end}}	return {{if $ok}}ctx.{{$ok.Name}}(res{{if $ok.HasMultipleViews}}, "default"{{end}}){{else}}nil{{end}}
}
{{end}}
`
