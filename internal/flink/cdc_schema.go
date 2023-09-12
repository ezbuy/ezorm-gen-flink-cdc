package flink

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/ezbuy/ezorm/v2/pkg/plugin"
	"github.com/iancoleman/strcase"
)

//go:embed templates/*.tpl
var cdcTemplate embed.FS

var files = []string{
	"templates/flink_cdc_schema.tpl",
}

var templateFuncs = template.FuncMap{
	"formatArgs": formatArgs,
}

type Schema struct {
	Args     map[string]interface{}
	Database string
	Table    string
	Fields   []Field
	PKs      []string
	Output   io.Writer
}

func (s *Schema) PrimaryKey() string {
	return strings.Join(s.PKs, ",")
}

func (s *Schema) Load(ctx context.Context, psm plugin.Schema) error {
	s.Database = psm["dbname"].(string)
	s.Table = psm["dbtable"].(string)
	rawFields := psm["fields"].([]any)
	var pks []any
	if _, ok := psm["primary"]; ok {
		pks = psm["primary"].([]any)
	}
	for _, pk := range pks {
		s.PKs = append(s.PKs, strcase.ToSnake(pk.(string)))
	}
	fs := make([]Field, 0, len(rawFields))
	for _, rawField := range rawFields {
		var f Field
		rf, ok := rawField.(map[string]any)
		if !ok {
			return fmt.Errorf("rawField:invalid field: %v", rawField)
		}

		var fname string
		var ftype string

		for k, v := range rf {
			if k[:1] == strings.ToUpper(k[:1]) {
				fname = strcase.ToSnake(k)
				ftype = toFlinkSQLType(v.(string))
			}
		}

		_, ok = rf["sqltype"]
		if ok {
			ftype = toFlinkSQLType(rf["sqltype"].(string))
		}

		_, ok = rf["sqlname"]
		if ok {
			fname = rf["sqlname"].(string)
		}

		f.name = fname
		f.t = ftype

		fl, ok := rf["flags"]
		if ok {
			fg, ok := fl.([]any)
			if ok {
				for _, flag := range fg {
					flag := flag.(string)
					if flag == "primary" {
						s.PKs = append(s.PKs, strcase.ToSnake(f.name))
					}
				}
			}
		}

		fs = append(fs, f)
	}
	s.Fields = fs
	return nil
}

func toFlinkSQLType(raw string) string {
	switch raw {
	case "uint8", "int8":
		return "TINYINT"
	case "uint16", "int16":
		return "SMALLINT"
	case "uint32", "int32", "int":
		return "INTEGER"
	case "uint64", "int64":
		return "BIGINT"
	case "float32":
		return "FLOAT"
	case "float64":
		return "DOUBLE"
	case "time.Time", "*time.Time", "timestamp", "timeint":
		return "TIMESTAMP"
	case "datetime":
		return "DATE"
	case "bool":
		return "BOOLEAN"
	default:
		return "VARCHAR"
	}
}

type Field struct {
	name string
	t    string
}

func (f Field) GetName() string {
	return f.name
}

func (f Field) GetType() string {
	return f.t
}

type CDCSchema struct {
	From Schema
	To   Schema
}

func (cs *CDCSchema) parseArgs(args map[string]string) {
	for k, v := range args {
		switch {
		case strings.HasPrefix(k, "from."):
			cs.From.Args[strings.TrimPrefix(k, "from.")] = v
		case strings.HasPrefix(k, "to."):
			cs.To.Args[strings.TrimPrefix(k, "to.")] = v
		}
	}
}

func (cs *CDCSchema) Load(ctx context.Context,
	psm plugin.Schema,
	args map[string]string,
) error {
	cs.parseArgs(args)
	if err := cs.From.Load(ctx, psm); err != nil {
		return err
	}
	if err := cs.To.Load(ctx, psm); err != nil {
		return err
	}
	return nil
}

func (cs *CDCSchema) Write(ctx context.Context, dir string) error {
	for k, f := range sprig.GenericFuncMap() {
		templateFuncs[k] = f
	}
	t, err := template.New("flink_cdc_schema").Funcs(
		templateFuncs,
	).ParseFS(cdcTemplate, files...)
	if err != nil {
		return err
	}
	if cs.From.Output == nil {
		fp := filepath.Join(dir, fmt.Sprintf("flink_cdc_%s_%s.sql", cs.From.Database, cs.From.Table))
		f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		cs.From.Output = f
	}
	return t.ExecuteTemplate(cs.From.Output, "flink_cdc_schema", cs)
}

func formatArgs(args map[string]interface{}) string {
	var pairs []string
	for k, v := range args {
		pairs = append(pairs, fmt.Sprintf(`'%s' = '%s'`, k, v))
	}
	return strings.Join(pairs, ",")
}
