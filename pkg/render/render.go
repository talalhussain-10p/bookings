package render

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"path/filepath"

	"github.com/bookings/pkg/model"
	"github.com/pkg/errors"
)

var cacheTemplate = make(map[string]*template.Template)

type Manager struct {
	useCache     bool
	template     map[string]*template.Template
	name         string
	templateData *model.TemplateData
}

func init() {
	c, err := createTemplateCache()
	if err != nil {
		log.Fatal(err)
	}
	cacheTemplate = c
}

func (m *Manager) RenderTemplate() (io.ReadCloser, error) {
	var err error
	if !m.useCache {
		m.template, err = createTemplateCache()
		if err != nil {
			return nil, err
		}
	} else {
		m.template = cacheTemplate
	}
	// get requested template from cache
	t, ok := m.template[m.name]
	if !ok {
		return nil, errors.New("no cache template found")
	}
	var out bytes.Buffer

	m.templateData = addDefaultData(m.templateData)

	if err = t.Execute(&out, m.templateData); err != nil {
		log.Println(err)
	}

	return io.NopCloser(&out), nil
}

func createTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	// get all of the files name *page.tmpl from templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return nil, err
	}

	// loop through all the files ending with page.tmpl
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return nil, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return nil, err
			}
		}
		cache[name] = ts
	}
	return cache, nil
}

func addDefaultData(td *model.TemplateData) *model.TemplateData {
	return td
}

func Client(tmpl string, td *model.TemplateData, useCache bool) *Manager {
	return &Manager{
		useCache:     useCache,
		name:         tmpl,
		templateData: td,
	}
}

// func Client(tmpl string) *Manager {
// 	// parsedTemplate, err := template.ParseFiles(fmt.Sprintf("./templates/%s", tmpl), fmt.Sprintf("./templates/%s", "base.layout.tmpl"))
// 	// if err != nil {
// 	// 	log.Fatal("couldn't load the template:", err)
// 	// }
// 	// err = parsedTemplate.Execute(w, nil)
// 	// if err != nil {
// 	// 	log.Fatal("couldn't execute the template:")
// 	// }
// 	var out bytes.Buffer
// 	return &Manager{
// 		useCache:      false,
// 		templateCache: cacheMap,
// 		out:           out,
// 		name:          tmpl,
// 	}
// }

// func (c *Manager) RenderTemplate() error {
// 	var t *template.Template
// 	var err error
// 	if _, ok := c.templateCache[c.name]; !ok {
// 		log.Println("creating template and adding to cache")
// 		err = c.cacheTemplate(c.name)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		log.Println("using cache template")
// 	}
// 	t = c.templateCache[c.name]
// 	return t.Execute(c.w, nil)
// }

// func (c *Manager) cacheTemplate(name string) error {
// 	templates := []string{
// 		fmt.Sprintf("./templates/%s", name),
// 		"./templates/base.layout.tmpl",
// 	}
// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}
// 	c.templateCache[name] = tmpl
// 	return nil
// }
