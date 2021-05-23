package sqltemplate

type DynamicDriver struct {
	sqls []Sql
}

func NewDynamicDriver() *DynamicDriver {
	return &DynamicDriver{}
}

func (dd *DynamicDriver) DriverName() string {
	return "Dynamic"
}

func (dd *DynamicDriver) Register(name, sql string) {
	dd.RegisterWithDescs(name, "", sql)
}

func (dd *DynamicDriver) RegisterWithDescs(name, description, sql string) {
	dd.sqls = append(dd.sqls, Sql{
		Name:        name,
		Description: description,
		Script:      sql,
	})
}

func (dd *DynamicDriver) Load() ([]Sql, error) {
	return dd.sqls, nil
}
