package strings

import (
	"bytes"
	"text/template"
)

var (
	tmplCache = map[string]*template.Template{}
)

// Replace 替换字符串内容
func Replace(s string, data any) (string, error) {
	tmpl, ok := tmplCache[s]
	var err error
	if !ok {
		tmpl, err = template.New(s).Parse(s)
		if err != nil {
			return "", err
		}
		tmplCache[s] = tmpl
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, data)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
