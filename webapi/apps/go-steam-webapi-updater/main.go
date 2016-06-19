// Copyright 2016 A.W. Stanley All rights reserved.
// Use of this source code is governed by a BSD-style
// licence that can be found in the LICENCE.md file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Param struct {
	name      string
	paramType string
	required  bool
}

var repository string = "github.com/awstanley/GoSteam/webapi"

// Public endpoint
var publicAPI string = "api.steampowered.com"

// Partner endpoint
var partnerAPI string = "partner.steam-api.com"

func Usage() {
	fmt.Println("Usage:")
	fmt.Printf("%s --key <key>\n\n", os.Args[0])
	fmt.Println("  --file <file>")
	fmt.Println("  --partner")
	fmt.Println("  --insecure")
}

var tmplRoot string

func toPrettyGoName(old string) string {
	// Make it marginally prettier
	out := strings.Replace(old, "_", " ", -1)
	out = strings.Title(out)
	out = strings.Replace(out, " ", "", -1)
	return out
}

func loadTemplate(name string) *template.Template {
	return template.Must(template.New(name).ParseFiles(filepath.Clean(fmt.Sprintf("%s/%s", tmplRoot, name))))
}

func main() {

	partner := flag.Bool("partner", false, "If true the partner api endpoint is used.")
	insecure := flag.Bool("insecure", false, "If true HTTP is used instead of HTTPS.")
	key := flag.String("key", "", "Steam API Key")
	localJSON := flag.String("file", "", "JSON file (for local load)")

	flag.Usage = Usage

	flag.Parse()

	base := publicAPI
	if *partner {
		*insecure = true
		base = partnerAPI
	}

	proto := "https://"
	if *insecure {
		proto = "http://"
	}

	// Get the output directory
	dst := fmt.Sprintf("%s/src/%s", os.Getenv("GOPATH"), repository)

	// Template root
	tmplRoot = fmt.Sprintf("%s/apps/go-steam-webapi-updater/tmpl", dst)

	tmplHeader := loadTemplate("header.txt")
	tmplFunc := loadTemplate("func.txt")
	tmplStruct := loadTemplate("struct.txt")
	tmplFuncGet := loadTemplate("funcGet.txt")
	tmplFuncPost := loadTemplate("funcPost.txt")

	var api apiSteam

	if *localJSON == "" {

		uri := fmt.Sprintf("%s%s/ISteamWebAPIUtil/GetSupportedAPIList/v1/", proto, base)
		if *key != "" {
			uri = fmt.Sprintf("%s?key=%s", uri, *key)
		}

		println(uri)

		response, err := http.Get(uri)
		if err != nil {
			fmt.Printf("Error encountered getting supported list: %s\n", err)
			return
		}
		defer response.Body.Close()

		// This should be safe, given the SteamAPI doesn't typically
		// return 2GB+ files...
		//content = make([]byte, response.ContentLength)
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("Error getting API: %s\n\n", err)
			return
		}

		var root jsonSteamSupportedRoot

		// Load the JSON
		err = json.Unmarshal(content, &root)

		if err != nil {
			fmt.Printf("Error arose unmarshalling: %s\n", err)
			return
		}

		// Now process it using the json
		api.load(&root)
	} else {
		content, err := ioutil.ReadFile(*localJSON)
		if err != nil {
			fmt.Printf("Failure in loading local file\n\terr: %s\n", err)
			return
		}

		var root jsonSteamSupportedRoot

		// Load the JSON
		err = json.Unmarshal(content, &root)

		if err != nil {
			fmt.Printf("Error arose unmarshalling: %s\n", err)
			return
		}

		// Now process it using the json
		api.load(&root)
	}

	tmplData := make(map[string]interface{})
	tmplData["webapi"] = fmt.Sprintf("%s/core", repository)

	// Now that we have that, we can build GoLang files ...
	for interfaceName, interfaceObj := range api.interfaces {
		goInterfaceName := toPrettyGoName(interfaceName)

		tmplData["interface"] = interfaceName
		folder := fmt.Sprintf("%s/%s", dst, goInterfaceName)
		os.MkdirAll(folder, 0644)
		for methodName, methodMap := range interfaceObj.methods {
			file := filepath.Clean(fmt.Sprintf("%s/%sRequest.go", folder, methodName))
			fp, err := os.Create(file)
			if err != nil {
				fmt.Printf("error opening '%s'\nerr: %s\n\n", file, err)
				return
			}

			// File header
			err = tmplHeader.Execute(fp, tmplData)
			if err != nil {
				fmt.Printf("failed to execute template (header)\nerr: %s\n\n", err)
				return
			}

			for version, versionObj := range methodMap.methods {
				tmplMethodName := fmt.Sprintf("%sV%d", methodName, version)
				tmplData["uri"] = fmt.Sprintf("%s/%s/v%d/", interfaceName, methodName, version)
				tmplData["version"] = fmt.Sprintf("%d", version)
				tmplData["method"] = tmplMethodName
				tmplData["verb"] = versionObj.verb

				// Build the struct first.
				err = tmplStruct.Execute(fp, tmplData)
				if err != nil {
					fmt.Printf("failed to execute template (struct)\nerr: %s\n\n", err)
					return
				}

				fmt.Fprintf(fp, "\ntype %s struct {\n", tmplMethodName)

				reqParams := make(map[string]*Param)
				optParams := make(map[string]*Param)

				typeString := ""
				requiresKey := false
				for _, p := range versionObj.params {
					if p.name == "key" {
						requiresKey = true
					} else {
						switch p.varType {
						case "string":
							typeString = "string"
							break
						case "{message}":
							typeString = "string"
							break
						case "int32":
							typeString = "int32"
							break
						case "bool":
							typeString = "bool"
							break
						case "float":
							typeString = "float32"
							break
						case "rawbinary":
							typeString = "[]byte"
							break
						case "uint32":
							typeString = "uint32"
							break
						case "uint64":
							typeString = "uint64"
							break
						}

						name := toPrettyGoName(p.name)

						if p.optional {
							optParams[p.name] = &Param{
								name:      name,
								paramType: typeString,
								required:  !p.optional,
							}
						} else {
							reqParams[p.name] = &Param{
								name:      name,
								paramType: typeString,
								required:  !p.optional,
							}
						}

						if p.description == "" {
							fmt.Fprintln(fp, "// No description provided by Valve")
						} else {
							fmt.Fprintf(fp, "// %s\n", p.description)
						}

						fmt.Fprintf(fp, "%s %s\n", name, typeString)
					}
				}
				fmt.Fprintf(fp, "}\n\n")

				tmplData["requiresKey"] = requiresKey

				// Func
				err = tmplFunc.Execute(fp, tmplData)
				if err != nil {
					fmt.Printf("failed to execute template (func)\nerr: %s\n\n", err)
					return
				}

				// This is where I do a walkthrough of the variables, building a query based on type.
				// Required = easy
				for n, v := range reqParams {
					// This is sort of nasty.
					switch v.paramType {
					case "string":
						fmt.Fprintf(fp, "params.AddString(\"%s\", method.%s)\n", n, v.name)
						break
					case "int32":
						fmt.Fprintf(fp, "params.AddInt32(\"%s\", method.%s)\n", n, v.name)
						break
					case "bool":
						fmt.Fprintf(fp, "params.AddBoolean(\"%s\",  method.%s)\n", n, v.name)
						break
					case "float32":
						fmt.Fprintf(fp, "params.AddFloat32(\"%s\", method.%s)\n", n, v.name)
						break
					case "[]byte":
						fmt.Fprintf(fp, "params.AddBytes(\"%s\", method.%s)\n", n, v.name)
						break
					case "uint32":
						fmt.Fprintf(fp, "params.AddUInt32(\"%s\", method.%s)\n", n, v.name)
						break
					case "uint64":
						fmt.Fprintf(fp, "params.AddUInt64(\"%s\", method.%s)\n", n, v.name)
						break
					}
				}

				for n, v := range optParams {
					// This is sort of nasty.
					switch v.paramType {
					case "string":
						fmt.Fprintf(fp, "if method.%s != \"\" {\nparams.AddString(\"%s\", method.%s)\n}\n", v.name, n, v.name)
						break
					case "int32":
						fmt.Fprintf(fp, "if method.%s != 0 {\nparams.AddInt32(\"%s\", method.%s)\n}\n", v.name, n, v.name)
						break
					case "bool":
						fmt.Fprintf(fp, "if method.%s != false {\nparams.AddBoolean(\"%s\",  method.%s)\n}\n", v.name, n, v.name)
						break
					case "float32":
						fmt.Fprintf(fp, "if method.%s != 0 {\nparams.AddFloat32(\"%s\", method.%s)\n}\n", v.name, n, v.name)
						break
					case "[]byte":
						fmt.Fprintf(fp, "if method.%s != nil {\nparams.AddBytes(\"%s\", method.%s)\n}\n", v.name, n, v.name)
						break
					case "uint32":
						fmt.Fprintf(fp, "if method.%s != 0 {\nparams.AddUInt32(\"%s\", method.%s)\n}\n", v.name, n, v.name)
						break
					case "uint64":
						fmt.Fprintf(fp, "if method.%s != 0 {\nparams.AddUInt64(\"%s\", method.%s)\n}\n", v.name, n, v.name)
						break
					}
				}

				// This does absolutely nothing special, in reality.
				switch versionObj.verb {
				case "GET":

					err = tmplFuncGet.Execute(fp, tmplData)
					if err != nil {
						fmt.Printf("failed to execute template (funcGet)\nerr: %s\n\n", err)
						return
					}
					break
				case "POST":
					err = tmplFuncPost.Execute(fp, tmplData)
					if err != nil {
						fmt.Printf("failed to execute template (funcPost)\nerr: %s\n\n", err)
						return
					}
					break

				}

				// End of function -- safety first, pad it with newlines
				fmt.Fprintf(fp, "\n}\n")
			}

			//
			//
			// What happens below here is ugly.
			//
			//

			fp.Close()
			fp, err = os.Open(file)
			if err != nil {
				fmt.Printf("failed to reopen file for reading\n\tfile: %s\n\terr: %s\n", file, err)
				return
			}

			raw, err := ioutil.ReadAll(fp)
			if err != nil {
				fmt.Printf("failed to read file content\n\tfile: %s\n\terr: %s\n", file, err)
				return
			}
			fp.Close()

			src, err := format.Source(raw)
			if err != nil {
				fmt.Printf("failed to format source (%s)\n\terr: %s\n", file, err)
				return
			}

			fp, err = os.Create(file)
			if err != nil {
				fmt.Printf("error re-opening '%s'\nerr: %s\n\n", file, err)
				return
			}

			fp.Write(src)
			fp.Close()
		}
	}
}
