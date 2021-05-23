package sqltemplate

import (
	"bytes"
	"errors"
	"fmt"
	"log"
)

func (st *SqlTemplate) RenderTPL(name string, data interface{}) (string, error) {
	var buff bytes.Buffer
	err := st.tpl.ExecuteTemplate(&buff, name, data)
	if err != nil {
		sql, has := st.findTpl(name)
		if has {
			return "", fmt.Errorf("sqltemplate - ERROR: %s[%s] %w", name, sql.Description, err)
		}
		return "", fmt.Errorf("sqltemplate - ERROR: %s %w", name, errors.New(fmt.Sprintf("template: %s no found", name)))
	}
	return buff.String(), nil
}

func (st *SqlTemplate) RenderTPLUnSave(name string, data interface{}) string {
	sql, err := st.RenderTPL(name, data)
	if err != nil {
		log.Printf("sqltemplate - ERROR: %s, render error: %s", name, err.Error())
		return ""
	}
	return sql
}
