package sqltemplate

import (
	"fmt"
	"testing"
)

func TestLoad(t *testing.T) {
	st := New()
	st.Use(NewMarkdownDriverWithDir("./test-sql"))
	st.Load()
	sql, args, err := st.RenderTPL("GetStudentByID2", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(sql, args)
}
