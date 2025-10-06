// Package sprig provides template functions compatible with Masterminds/sprig.
package sprig

import (
	htmltemplate "html/template"
	"text/template"
)

// FuncMap returns the standard function map for templates.
// This is an alias for TxtFuncMap.
func FuncMap() template.FuncMap {
	return TxtFuncMap()
}

// TxtFuncMap returns a function map for text templates.
func TxtFuncMap() template.FuncMap {
	return template.FuncMap(genericFuncMap())
}

// HtmlFuncMap returns a function map for HTML templates.
func HtmlFuncMap() htmltemplate.FuncMap {
	return htmltemplate.FuncMap(genericFuncMap())
}

// HermeticTxtFuncMap returns a function map with only repeatable text template functions.
// Functions that depend on time, randomness, or environment are excluded.
func HermeticTxtFuncMap() template.FuncMap {
	return template.FuncMap(hermeticFuncMap())
}

// HermeticHtmlFuncMap returns a function map with only repeatable HTML template functions.
// Functions that depend on time, randomness, or environment are excluded.
func HermeticHtmlFuncMap() htmltemplate.FuncMap {
	return htmltemplate.FuncMap(hermeticFuncMap())
}

// GenericFuncMap returns a copy of the basic function map as a map[string]interface{}.
func GenericFuncMap() map[string]interface{} {
	return genericFuncMap()
}
