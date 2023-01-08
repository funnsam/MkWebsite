package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
)

var GHAPIToken = os.Getenv("GHAPIToken")
var DefaultJSONValue MarkdownDescription
var DefaultJSONValueMap map[string]string

type MarkdownFile struct {
	Content string `json:"text"`
	Mode    string `json:"mode"`
	info    MarkdownDescription
}

type MarkdownDescription struct {
	Title        string
	TemplateName string
	Description  string
	Image        string
}

func CompNavBar() string {
	tmpf, _ := os.ReadFile("./templates/navbaritems.html")
	tmpm := MarkdownFile{string(tmpf), "gfm", MarkdownDescription{}}
	tmpj, _ := json.Marshal(tmpm)
	tmpc, _ := Markdown2HTML(tmpj)

	return string(tmpc)
}

func main() {
	tmp, _ := os.ReadFile("./templates/default.html")
	var defaulttemplate, _ = template.New("defaulttemplate").Parse(string(tmp))
	var metadatatemplate, _ = template.New("metadatatemplate").Parse(`<meta name="description" content="{{.Description}}">
	<meta property="og:title" content="{{.Title}}" />
	<meta property="og:image" content="{{.Image}}" />
	<meta property="og:description" content="{{.Description}}" />`)
	var navbar = CompNavBar()

	tmp, _ = os.ReadFile("./markdown/default.json")
	json.Unmarshal(tmp, &DefaultJSONValue)
	DefaultJSONValueMap = ToMap(&DefaultJSONValue)

	var wg sync.WaitGroup

	if filepath.Walk("./markdown",
		func(path string, info os.FileInfo, err error) error {
			wg.Add(1)
			go func(path string, info os.FileInfo, err error) {
				defer wg.Done()

				if info.IsDir() {
					switch filepath.Base(path) {
					case "res":
						filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
							if info.IsDir() {
								return nil
							}

							a, _ := os.ReadFile(path)

							os.MkdirAll(filepath.Join(".", "website", filepath.Dir(path[9:])), 0664)
							os.WriteFile(filepath.Join(".", "website", path[9:]), a, 0664)

							return nil
						})
					}
				} else {
					if strings.Contains(filepath.Ext(path), "_") {
						a, _ := os.ReadFile(path)
						os.MkdirAll(filepath.Join(".", "website", filepath.Dir(path[9:])), 0664)
						os.WriteFile(filepath.Join(".", "website", path[9:]), a, 0664)
					} else {
						switch filepath.Ext(path) {
						case ".md":
							tmpf, _ := os.ReadFile(path)
							tmpm := MarkdownFile{string(tmpf), "gfm", GetMarkdownInfo(path)}
							tmpj, _ := json.Marshal(tmpm)
							tmpc, _ := Markdown2HTML(tmpj)

							prehtml := make(map[string]string, 20)
							prehtml["Content"] = string(tmpc)
							prehtml["NavBar"] = navbar

							if tmpm.info.Title == "" {
								prehtml["Title"] = strings.TrimSuffix(filepath.Base(path), ".md")
							} else {
								prehtml["Title"] = tmpm.info.Title
							}

							prehtml["Metadata"] = FormatToTemplate(metadatatemplate, ToMap(&tmpm.info))

							var html string

							if tmpm.info.TemplateName == "" {
								html = FormatToTemplate(defaulttemplate, prehtml)
							} else {
								t, _ := os.ReadFile(filepath.Join(".", "templates", tmpm.info.TemplateName+".html"))
								html = FormatToTemplate(template.Must(template.New(path).Parse(string(t))), prehtml)
							}

							os.MkdirAll(filepath.Join(".", "website", filepath.Dir(path[9:])), 0664)
							os.WriteFile(filepath.Join(".", "website", strings.TrimSuffix(path, filepath.Ext(path))[9:]+".html"), []byte(html), 0664)

							fmt.Println("Compiled ", path)

						case ".html":
							tmpf, _ := os.ReadFile(path)
							tmpm := MarkdownFile{"", "gfm", GetMarkdownInfo(path)}

							prehtml := make(map[string]string, 20)
							prehtml["Content"] = string(tmpf)
							prehtml["NavBar"] = navbar

							if tmpm.info.Title == "" {
								prehtml["Title"] = strings.TrimSuffix(filepath.Base(path), ".md")
							} else {
								prehtml["Title"] = tmpm.info.Title
							}

							prehtml["Metadata"] = FormatToTemplate(metadatatemplate, ToMap(&tmpm.info))

							var html string

							if tmpm.info.TemplateName == "" {
								html = FormatToTemplate(defaulttemplate, prehtml)
							} else {
								t, _ := os.ReadFile(filepath.Join(".", "templates", tmpm.info.TemplateName+".html"))
								html = FormatToTemplate(template.Must(template.New(path).Parse(string(t))), prehtml)
							}

							os.MkdirAll(filepath.Join(".", "website", filepath.Dir(path[9:])), 0664)
							os.WriteFile(filepath.Join(".", "website", strings.TrimSuffix(path, filepath.Ext(path))[9:]+".html"), []byte(html), 0664)

							fmt.Println("Compiled ", path)
						}
					}
				}

			}(path, info, err)
			return nil
		}) != nil {
		panic("Error while getting all files")
	}

	wg.Wait()
}

func GetMarkdownInfo(path string) MarkdownDescription {
	jso, err := os.ReadFile(strings.TrimSuffix(path, filepath.Ext(path)) + ".json")
	if err != nil {
		return DefaultJSONValue
	}
	r := MarkdownDescription{}
	if json.Unmarshal(jso, &r) != nil {
		return DefaultJSONValue
	}
	return r
}

func FormatToTemplate(templ *template.Template, data map[string]string) string {
	builder := &strings.Builder{}
	if err := templ.Execute(builder, data); err != nil {
		panic(err)
	}
	return builder.String()
}

func ToMap(a *MarkdownDescription) map[string]string {
	var b map[string]interface{}
	c, _ := json.Marshal(a)
	json.Unmarshal(c, &b)
	return MapStringInterface2String(b)
}

func MapStringInterface2String(a map[string]interface{}) map[string]string {
	result := map[string]string{}

	for k := range a {
		if v, ok := a[k].(string); ok {
			result[k] = v
		} else {
			result[k] = DefaultJSONValueMap[k]
		}
	}

	return result
}
