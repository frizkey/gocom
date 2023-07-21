package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
)

var ifaceName string

func init() {

	flag.StringVar(&ifaceName, "src", "", "input interface name")
}

func main() {

	flag.Parse()

	if ifaceName == "" {
		flag.PrintDefaults()
	}

	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println("Unable to get working directory")
		os.Exit(-1)
	}

	fmt.Println("distobj gen :", ifaceName, cwd)
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, cwd, nil, parser.ParseComments)

	if err != nil {
		fmt.Println("Error when parsing package :", err)
		return
	}

	for _, pkg := range pkgs {

		for _, fl := range pkg.Files {

			for _, dcl := range fl.Decls {

				genDecl, ok := dcl.(*ast.GenDecl)

				if !ok {
					continue
				}

				if genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {

					astSpec, ok := spec.(*ast.TypeSpec)

					if !ok {
						continue
					}

					if astSpec.Name.Name != ifaceName {
						continue
					}

					// start generate proxy --------------

					proxyName := ifaceName + "Proxy"

					strOut := `package ` + pkg.Name + `

import "github.com/adlindo/gocom/distobj"

type ` + proxyName + ` struct {
	prefix string
}

func Get` + proxyName + `(prefix string) *` + proxyName + ` {

	ret = &` + proxyName + `{}
	ret.prefix = prefix

	return ret
}

					`

					ast.Print(fset, fl)

					astIFace, ok := astSpec.Type.(*ast.InterfaceType)

					if !ok {
						continue
					}

					for _, field := range astIFace.Methods.List {

						if len(field.Names) == 0 {
							continue
						}

						mtdName := field.Names[0].Name
						paramList := map[string]string{}
						resultList := []string{}

						astMtd, ok := field.Type.(*ast.FuncType)

						if !ok {
							continue
						}

						//params
						for _, param := range astMtd.Params.List {

							paramType, ok := param.Type.(*ast.Ident)

							if !ok {
								continue
							}

							for _, name := range param.Names {

								paramList[name.Name] = paramType.Name
							}
						}

						//result
						for _, result := range astMtd.Results.List {

							typeIdent, ok := result.Type.(*ast.Ident)

							if !ok {
								continue
							}

							resultList = append(resultList, typeIdent.Name)
						}

						// start generate method ----------------------------

						str := "func (o *" + proxyName + ") " + mtdName + "("

						paramDecl := ""
						paramInvoke := ""

						for pName, pType := range paramList {

							paramDecl += ", " + pName + " " + pType
							paramInvoke += ", " + pName
						}

						if len(paramList) > 0 {
							paramDecl = paramDecl[2:]
							paramInvoke = paramInvoke[2:]
						}

						str += paramDecl + ")"

						resultDecl := ""

						for _, result := range resultList {

							resultDecl += "," + result
						}

						if len(resultList) > 1 {
							resultDecl = "(" + resultDecl[1:] + ") "
						} else if len(resultList) > 0 {
							resultDecl = resultDecl[1:] + " "
						}

						str += " " + resultDecl + `{

	ret, err := distobj.Invoke(o.prefix, "` + ifaceName + `", "` + mtdName + `", ` + paramInvoke + `)
	
	return ret, err
}`

						strOut += `
` + str + `
						`
					}

					fmt.Println("====================================>")
					fmt.Println(strOut)
				}
			}
		}
	}
}
