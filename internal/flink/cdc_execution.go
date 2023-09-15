package flink

import (
	"context"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/cespare/xxhash/v2"
	"github.com/ezbuy/ezorm/v2/pkg/plugin"
)

//go:embed templates/*.tpl
var cdcExecutionTemplate embed.FS

var execFiles = []string{
	"templates/flink_cdc_execution.tpl",
}

type Execution struct {
	Connector string
	Database  string
	Table     string
	Fields    []Field
}

type CDCExecution struct {
	From Execution
	To   Execution
}

func WriteExecutionTemplate(
	_ context.Context,
	data map[plugin.TemplateName]CDCExecution,
	dir string,
) error {
	t, err := template.New("flink_cdc_execution").Funcs(
		sprig.TxtFuncMap(),
	).
		ParseFS(cdcExecutionTemplate, execFiles...)
	if err != nil {
		return err
	}

	var exes []CDCExecution
	d := xxhash.New()
	for t, exe := range data {
		exes = append(exes, exe)
		d.WriteString(string(t))
	}
	hash := d.Sum64()

	fp := filepath.Join(dir, fmt.Sprintf("flink_cdc_execution_%d.sql", hash))
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.ExecuteTemplate(f, "flink_cdc_execution", exes)
}
