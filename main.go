package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

func main()  {
	var types string
	//var typeName string
	var typeField []map[string]string

	flag.StringVar(&types, "types", "", "类型名称")
	flag.Parse()
	typeSlice := strings.Split(types, ",")

	//fmt.Println(os.Args[0],os.Getenv("GOPACKAGE"), os.Getenv("GOFILE"))
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rootPath := strings.TrimRight(dir, "/")
	//filename := os.Getenv("GOFILE")
	filePath := rootPath + "/" +  os.Getenv("GOFILE")

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", fileData, 0)
	if err != nil {
		panic(err)
	}
	//ast.Print(fset, f)

	//fmt.Println(typeSlice)
	//obj := f.Scope.Lookup(typeSlice[0])
	//fmt.Println(obj.Name, obj.Decl)

	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			//fmt.Println("not GenDecl")
			continue
		}
		spec, ok := gen.Specs[0].(*ast.TypeSpec)
		if !ok {
			//fmt.Println("not TypeSpec")
			continue
		}
		if spec.Name.Name == typeSlice[0] {
			//fmt.Println(spec.Name.Name)
			//typeName = spec.Name.Name
			//解析类型的字段和字段类型
			structType := spec.Type.(*ast.StructType)
			for _, field := range structType.Fields.List {
				var fieldInfo = map[string]string{}
				fieldInfo["name"] = field.Names[0].Name
				if _, ok := field.Type.(*ast.Ident); ok {
					fieldInfo["type"] = field.Type.(*ast.Ident).Name
				}
				if _, ok := field.Type.(*ast.ArrayType); ok {
					fieldInfo["type"] = "[]" + field.Type.(*ast.ArrayType).Elt.(*ast.Ident).Name
				}
				typeField = append(typeField, fieldInfo)
			}
		}
	}
	result := struct {
		PkgName string
		TypeName string
		Fields []map[string]string
	}{
		PkgName: os.Getenv("GOPACKAGE"),
		TypeName: typeSlice[0],
		Fields: typeField,
	}

	tpl := tpl()
	tmpl, err := template.New("template").Parse(tpl)
	if err != nil {
		panic(err)
	}

	err = tmpl.Execute(os.Stdout, result)
	if err != nil {
		panic(err)
	}
	//fmt.Println(result)
}

func tpl() string {
	return `
package {{ .PkgName }}
{{ $typename := .TypeName }}
{{- range $index, $item := .Fields }}
func (s {{ $typename }}) Set{{ $item.name }}(value {{ $item.type }}) {
	s.{{ $item.name }} = value
}

func (s {{ $typename }}) Get{{ $item.name }}() {{ $item.type }} {
	return s.{{ $item.name }}
}
{{ end }}
`
}
