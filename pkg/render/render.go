package render

import (
	"bytes"
	"fmt"
	"github.com/Igor-Koniukhov/rest_api/pkg/config"
	"github.com/Igor-Koniukhov/rest_api/pkg/datastruct"
	"html/template"
	"net/http"
	"path/filepath"
)

var err error
var app *config.AppConfig

func NewTemplate(a *config.AppConfig)  {
	app = a
}

func RenderTemplate(w http.ResponseWriter, tmpl string, data datastruct.Data) {
	var tc map[string]*template.Template

	if app.UseCache{
		tc = app.TemplateCache
	}else{
		tc, _ = CreateTemplateCache()
	}



	t, ok := tc[tmpl]
	if !ok {
		fmt.Println("Could not get template from template cache")
	}
	buf := new(bytes.Buffer)
	_ = t.Execute(buf, data)
	_, err = buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}

}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*page.tmpl")
	if err != nil {
		return myCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)

		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
			myCache[name] = ts
		}

	}

	return myCache, err

}
