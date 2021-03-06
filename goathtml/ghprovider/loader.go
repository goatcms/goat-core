package ghprovider

import (
	"html/template"
	"sync"

	"github.com/goatcms/goatcore/filesystem"
	"github.com/goatcms/goatcore/varutil/goaterr"
)

// TemplateLoader provide method to load templates from filesystem
type TemplateLoader struct {
	template   *template.Template
	muTemplate *sync.Mutex
}

// NewTemplateLoader create TemplateLoader instance
func NewTemplateLoader(template *template.Template) *TemplateLoader {
	return &TemplateLoader{
		template:   template,
		muTemplate: &sync.Mutex{},
	}
}

// Load get all templates code form files in subPath and add it to template
func (loader *TemplateLoader) Load(fs filesystem.Filespace, subPath string) error {
	bytes, err := fs.ReadFile(subPath)
	if err != nil {
		return err
	}
	loader.muTemplate.Lock()
	defer loader.muTemplate.Unlock()
	if len(bytes) == 0 {
		return goaterr.Errorf("empty file")
	}
	if bytes == nil || len(bytes) == 0 {
		return goaterr.Errorf("file data is empty")
	}
	if _, err := loader.template.Parse(string(bytes)); err != nil {
		return goaterr.Errorf("%v: %v", subPath, err)
	}
	return nil
}

// Template return loaded template
func (loader *TemplateLoader) Template() *template.Template {
	return loader.template
}
