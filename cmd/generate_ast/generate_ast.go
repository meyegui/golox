package main

import (
	"fmt"
	"html/template"
	"os"
	"strings"
)

type field struct {
	Name string
	Type string
}

var TPL_ACCEPT_VISITOR_FUNC = //
`func ({{.r}} *{{.rType}}) Accept({{.v}} {{.vType}}) {
	{{.v}}.Visit{{.rType}}({{.r}})
}`

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "Usage: generate_ast <output directory>")
		os.Exit(64)
	}

	outputDir := args[0]
	defineAst(outputDir, "Expr", map[string][]field{
		"Binary": {
			{"Left", "Expr"},
			{"Operator", "*scanner.Token"},
			{"Right", "Expr"},
		},
		"Grouping": {
			{"Expr", "Expr"},
		},
		"Literal": {
			{"Value", "any"},
		},
		"Unary": {
			{"Operator", "*scanner.Token"},
			{"Right", "Expr"},
		},
	})
}

func defineAst(outputDir string, interfaceName string, definitions map[string][]field) {
	outputDir = strings.TrimSuffix(outputDir, "/")

	packageName := strings.ToLower(outputDir[strings.LastIndex(outputDir, "/")+1:])
	if len(packageName) == 0 {
		panic(fmt.Errorf("outputDir results in an empty package name"))
	}

	// Create the interface that will unite all the structs
	filename := outputDir + "/" + strings.ToLower(interfaceName) + ".go"
	os.WriteFile(filename, []byte(fmt.Sprintf(
		`package %s

type %s interface {
	is%s()
	Accept(v %sVisitor)
}
`,
		packageName,
		interfaceName,
		interfaceName,
		interfaceName,
	)), 0o644)

	// Create the visitor interface
	func() {
		filename = outputDir + "/" + strings.ToLower(interfaceName) + "_" + "visitor.go"
		file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		fmt.Fprintln(file, "package", packageName)
		fmt.Fprintln(file)

		fmt.Fprintln(file, "type", interfaceName+"Visitor", "interface {")

		for typeName := range definitions {
			fmt.Fprintf(
				file,
				"\tVisit%s(%s *%s)\n",
				typeName,
				strings.ToLower(string(typeName[0])),
				typeName,
			)
		}

		fmt.Fprintln(file, "}")
	}()

	// Create the structs
	for typeName, fields := range definitions {
		func() {
			// Create or open file in "overwrite" mode
			filename = outputDir + "/" + strings.ToLower(typeName) + ".go"
			file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			// Write out package name
			fmt.Fprintln(file, "package", packageName)
			fmt.Fprintln(file)

			// Write out imports, if any
			if imports := getImports(fields); len(imports) > 0 {
				fmt.Fprintln(file, "import (")

				for _, importPackageName := range imports {
					fmt.Fprintf(
						file,
						"\t\"github.com/meyegui/golox/%s\"\n",
						importPackageName,
					)
				}

				fmt.Fprintln(file, ")")
				fmt.Fprintln(file)
			}

			// Write out type name
			fmt.Fprintln(file, "type", typeName, "struct {")

			// Calculate alignment (width between field name and type)
			alignment := 0
			for _, f := range fields {
				if len(f.Name) > alignment {
					alignment = len(f.Name)
				}
			}
			alignment++

			// Write out fields
			for _, f := range fields {
				fmt.Fprintf(
					file,
					"\t%s%s%s\n",
					f.Name,
					strings.Repeat(" ", alignment-len(f.Name)),
					f.Type,
				)
			}
			fmt.Fprintln(file, "}")
			fmt.Fprintln(file)

			// Write out tag interface method
			fmt.Fprintf(
				file,
				"func (%s %s) is%s() {}\n",
				strings.ToLower(string(typeName[0])),
				typeName,
				interfaceName,
			)

			// Write out visitor implementation, i.e. accept() method
			fmt.Fprintln(file)
			template.
				Must(template.New("").Parse(TPL_ACCEPT_VISITOR_FUNC)).
				Execute(file, map[string]any{
					"r":     strings.ToLower(string(typeName[0])),
					"rType": typeName,
					"v":     strings.ToLower(string(interfaceName[0])) + "v",
					"vType": interfaceName + "Visitor",
				})
			fmt.Fprintln(file)
		}()
	}
}

func getImports(fields []field) []string {
	imports := make([]string, 0)
	for _, f := range fields {
		star := strings.IndexByte(f.Type, '*')
		dot := strings.IndexByte(f.Type, '.')
		if dot == -1 {
			continue
		}

		imports = append(imports, f.Type[star+1:dot])
	}

	return imports
}
