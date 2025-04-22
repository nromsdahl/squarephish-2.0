package templates

import (
	"html/template"
	"sync"
)

var (
	tmplCache = make(map[string]*template.Template)
	mu        sync.Mutex
)

// GetTemplate retrieves a template from the cache or parses it if not cached.
// It locks the cache to prevent concurrent access issues.
// Parameters:
//   - tmpl: The name of the template file to retrieve.
//
// It returns the template and an error if the template is not found or executed
// correctly.
func GetTemplate(tmpl string) (*template.Template, error) {
	mu.Lock()
	defer mu.Unlock()

	if cached, ok := tmplCache[tmpl]; ok {
		return cached, nil
	}

	templates := template.New("template")
	_, err := templates.ParseFiles("templates/base.html", "templates/"+tmpl)
	if err != nil {
		return nil, err
	}

	tmplCache[tmpl] = templates
	return templates, nil
}
