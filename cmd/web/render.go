package main

import (
	"github.com/smarulanda97/app-stripe/internal/utils"

	"embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var functions = template.FuncMap{
	"formatCurrency": utils.FormatCurrency,
}

//go:embed theme/templates
var templateFS embed.FS

// Add default data to all given templates
func (app *application) addDefaultData(templateData *utils.TemplateData, r *http.Request) *utils.TemplateData {
	templateData.ApiUrl = app.kernel.ApiUrl
	templateData.StripePublicKey = app.kernel.Stripe.PublicKey

	return templateData
}

// check if the given template exist in render cache
// render cache only works in production environment
func (app *application) isTemplateCached(templateName string) bool {
	renderCache := app.renderCache
	_, isCached := renderCache[templateName]

	return isCached && app.kernel.Environment == "production"
}

// Execute and render the template with the writer
func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, pageName string, templateData *utils.TemplateData, partials ...string) error {
	var t *template.Template

	templateName := fmt.Sprintf("%s/%s.page.gohtml", app.kernel.ThemePath, pageName)

	if app.isTemplateCached(templateName) {
		t = app.renderCache[templateName]
	} else {
		t, _ = app.parseTemplate(partials, pageName, templateName)
	}

	if templateData == nil {
		templateData = &utils.TemplateData{}
	}

	templateData = app.addDefaultData(templateData, r)
	if err := t.Execute(w, templateData); err != nil {
		return err
	}

	return nil
}

// Build partials for each template
func (app *application) buildPartials(partials []string) []string {
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("%s/%s.partial.gohtml", app.kernel.ThemePath, x)
		}
	}

	return partials
}

// Parse each template
func (app *application) parseTemplate(partials []string, page, templateName string) (*template.Template, error) {
	var t *template.Template

	partials = app.buildPartials(partials)
	templateFile := fmt.Sprintf("%s.page.gohtml", page)
	templateBaseLayout := fmt.Sprintf("%s/base.layout.gohtml", app.kernel.ThemePath)

	if len(partials) > 0 {
		t, _ = template.New(templateFile).Funcs(functions).ParseFS(templateFS, templateBaseLayout, strings.Join(partials, ","), templateName)
	} else {
		t, _ = template.New(templateFile).Funcs(functions).ParseFS(templateFS, templateBaseLayout, templateName)
	}

	app.renderCache[templateName] = t

	return t, nil
}
