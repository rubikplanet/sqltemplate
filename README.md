# sqltemplate

## feature

- manage sql with markdown
- render sql with go template
- support load sql with custom plugin

## usage

create a markdown in `sql/test.md`

```markdown
    ### GetStudentByID
    >get student by id, required id
    ```sql
    select * from student where id = {{. | bind}}
    ```
```

in golang

```go
package main

import (
    "fmt"
    "github.com/freshtech2021/sqltemplate"
)

func main() {
    st := sqltemplate.New()
    st.Use(sqltemplate.NewMarkdownDriver())
    // load sql with custom dir
    // st.Use(sqltemplate.NewMarkdownDriverWithDir("./prod-sql"))
    // register go template func
    // st.RegisterFunc(template.FuncMap{
    //     "test": func(v string) string {
    //         return strings.ToUpper(v)
    //     },
    // })
    st.Load()
    sql, args, err := st.RenderTPL("GetStudentByID2", 1)
    if err != nil {
        panic(err)
    }
    fmt.Println(sql, args)
}
```

## custom puglin

> implement sqltemplate.Driver

```go
type CustomeDriver struct {
}

func NewCustomeDriver() *CustomeDriver {
    return &CustomeDriver{}
}

func (mdd *CustomeDriver ) DriverName() string {
    return "CustomeDriver"
}

func (mdd *CustomeDriver ) Load() ([]Sql, error) {
    var list []Sql
    // db.table("sql_store").Find(&list)
    return list, nil
}
```
