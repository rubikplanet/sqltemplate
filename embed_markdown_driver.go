package sqltemplate

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"path"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type EmbedMarkdownDriver struct {
	fs  embed.FS
	dir string
}

func NewMarkdownDriverWithEmbedDir(fs embed.FS, dir string) *EmbedMarkdownDriver {
	return &EmbedMarkdownDriver{
		fs:  fs,
		dir: dir,
	}
}

func NewMarkdownDriverWithEmbed(fs embed.FS) *EmbedMarkdownDriver {
	return NewMarkdownDriverWithEmbedDir(fs, "sql")
}
func (mdd *EmbedMarkdownDriver) DriverName() string {
	return "embed"
}

func (mdd *EmbedMarkdownDriver) Load() ([]Sql, error) {
	var sqls []Sql
	err := fs.WalkDir(mdd.fs, mdd.dir, func(subpath string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		ext := path.Ext(subpath)
		if ext == ".md" || ext == ".markdown" {
			s, err := mdd.parseMarkdown(subpath)
			if err != nil {
				return err
			}
			if len(s) != 0 {
				sqls = append(sqls, s...)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sqls, nil
}

func (mdd *EmbedMarkdownDriver) parseMarkdown(filename string) ([]Sql, error) {
	var sqls []Sql
	buf, err := mdd.fs.ReadFile(filename)
	if err != nil {
		log.Printf("sqltemplate - ERROR: %s loading failed...\n", filename)
		return nil, err
	}
	if bytes.ContainsRune(buf, '\r') {
		buf = markdown.NormalizeNewlines(buf)
	}
	psr := parser.New()
	node := markdown.Parse(buf, psr)
	list := getAll(node)
	i := 0
	for {
		// 1. text, code
		// 2. text, text, code
		if i >= len(list) {
			break
		}
		var tpl Sql
		if i+1 < len(list) && list[i].Type == "text" && list[i+1].Type == "code" {
			tpl.Name = mdd.getName(filename, list[i].Content)
			tpl.Script = list[i+1].Content
			sqls = append(sqls, tpl)
			i += 2
		} else if i+2 < len(list) && list[i].Type == "text" && list[i+1].Type == "text" && list[i+2].Type == "code" {
			tpl.Name = mdd.getName(filename, list[i].Content)
			tpl.Description = list[i+1].Content
			tpl.Script = list[i+2].Content
			sqls = append(sqls, tpl)
			i += 3
		} else {
			return nil, errors.New(fmt.Sprintf("ERROR: parse markdown failed: %s", filename))
		}
	}
	return sqls, nil
}

func (mdd *EmbedMarkdownDriver) getName(filename, code string) string {
	ext := path.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	base = strings.TrimPrefix(base, mdd.dir)
	base = strings.TrimPrefix(base, "/")
	return path.Join(strings.TrimSuffix(base, ext), code)
}
