package sqltemplate

import (
	"fmt"
	"strings"
	"testing"
	"text/template"
)

func TestDynamicLoad(t *testing.T) {
	st := New()
	store := NewDynamicDriver()
	store.Register("rest1", `select * from table where id = {{. }} or {{ block "rest" . }} {{ end }}`)
	store.Register("rest", `select * from table where id = "{{ test . }}"`)
	st.Use(store)
	st.RegisterFunc(template.FuncMap{
		"test": func(v string) string {
			return strings.ToUpper(v)
		},
	})
	st.Load()
	sql, err := st.RenderTPL("rest1", "test")
	if err != nil {
		panic(err)
	}
	fmt.Println(sql)
}
