package sqltemplate

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

type bindings struct {
	values []interface{}
}

func (b *bindings) bind(value interface{}) string {
	var paramPlaceHolders []string
	v := reflect.ValueOf(value)
	if v.Type().Kind() != reflect.Slice {
		b.values = append(b.values, value)
		paramPlaceHolders = append(paramPlaceHolders, "?")
	} else {
		length := v.Len()
		for i := 0; i < length; i++ {
			b.values = append(b.values, v.Index(i).Interface())
			paramPlaceHolders = append(paramPlaceHolders, "?")
		}
	}
	return strings.Join(paramPlaceHolders, ",")
}

func (st *SqlTemplate) RenderTPL(name string, data interface{}) (string, []interface{}, error) {
	values := &bindings{values: []interface{}{}}
	clonedTmpl, err := st.tpl.Clone()
	if err != nil {
		return "", nil, fmt.Errorf("unable to parse template %w", err)
	}
	clonedTmpl.Funcs(template.FuncMap{"bind": values.bind})
	var buff bytes.Buffer
	err = clonedTmpl.ExecuteTemplate(&buff, name, data)
	if err != nil {
		sql, has := st.findTpl(name)
		if has {
			return "", nil, fmt.Errorf("sqltemplate - ERROR: %s[%s] %w", name, sql.Description, err)
		}
		return "", nil, fmt.Errorf("sqltemplate - ERROR: %s %w", name, errors.New(fmt.Sprintf("template: %s no found", name)))
	}
	return buff.String(), values.values, nil
}
