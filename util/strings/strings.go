package strings

import (
	"bytes"
	"sort"
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

// InArray 字符串数组中是否存在目标串
func InArray(target string, str_array []string) bool {
	// 快排字符串数组
	sort.Strings(str_array)
	// 二分法搜索
	index := sort.SearchStrings(str_array, target)
	if index < len(str_array) && str_array[index] == target {
		return true
	}
	return false
}
