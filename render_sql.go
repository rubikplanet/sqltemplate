package sqltemplate

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"
)

type bindings struct {
	values []interface{}
}

func (b *bindings) bind(value interface{}) string {
	b.values = append(b.values, value)
	return "?"
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
