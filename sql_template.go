package sqltemplate

import (
	"log"
	"reflect"
	"runtime"
	"strings"
	"text/template"
)

type SqlTemplate struct {
	sqls    []Sql
	drivers map[string]Driver
	funcs   template.FuncMap
	tpl     *template.Template
}

func New() *SqlTemplate {
	st := &SqlTemplate{
		drivers: make(map[string]Driver),
		funcs:   template.FuncMap{},
	}
	return st
}

func (st *SqlTemplate) Use(plugin Driver) {
	if _, ok := st.drivers[plugin.DriverName()]; ok {
		log.Printf("sqltemplate - WARN: %s already used\n", plugin.DriverName())
	}
	st.drivers[plugin.DriverName()] = plugin
}

func (st *SqlTemplate) Load() {
	st.tpl = nil
	for _, driver := range st.drivers {
		sqls, err := driver.Load()
		if err != nil {
			log.Printf("sqltemplate - ERROR: %s load failed: ", sqls)
			log.Panicln(err)
		}
		for _, sql := range sqls {
			d, has := st.findTpl(sql.Name)
			if has {
				log.Printf("sqltemplate - WARN: %s Has duplicate sql: It will be cover [%s] with [ %s ]\n", sql.Name, strings.ReplaceAll(d.Script, "\n", ""), strings.ReplaceAll(sql.Script, "\n", ""))
			}
			st.sqls = append(st.sqls, sql)
			if st.tpl == nil {
				st.tpl = template.New(sql.Name)
				st.tpl.Funcs(st.funcs)
			} else {
				st.tpl = st.tpl.New(sql.Name)
			}
			st.tpl, err = st.tpl.Parse(sql.Script)
			if err != nil {
				panic(err)
			}
		}
		log.Printf("sqltemplate - INFO: %s loaded %d sqls.\n", driver.DriverName(), len(sqls))
	}
}

func (st *SqlTemplate) findTpl(name string) (*Sql, bool) {
	for _, tpl := range st.sqls {
		if tpl.Name == name {
			return &tpl, true
		}
	}
	return nil, false
}

func (st *SqlTemplate) RegisterFunc(funcs template.FuncMap) {
	for k, v := range funcs {
		if temp, ok := st.funcs[k]; ok {
			log.Printf("sqltemplate - WARN: %s Has duplicate func: It will be cover [%s] with [%s]\n", k, getFunctionName(temp), getFunctionName(v))
		}
		st.funcs[k] = v
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type Driver interface {
	Load() ([]Sql, error)
	DriverName() string
}

type Sql struct {
	Name        string // 名称
	Description string // 描述
	Script      string // sql
}
