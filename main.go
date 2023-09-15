package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ezbuy/ezorm-gen-flink-cdc/internal/flink"
	"github.com/ezbuy/ezorm/v2/pkg/plugin"
)

func main() {
	r := bufio.NewScanner(os.Stdin)
	ctx := context.Background()
	for r.Scan() {
		req, err := plugin.Decode(r.Bytes())
		if err != nil {
			fmt.Fprintf(os.Stderr, "decode error: %s\n", err)
			return
		}
		exes := make(map[plugin.TemplateName]flink.CDCExecution)
		if err := req.Each(func(t plugin.TemplateName, s plugin.Schema) error {
			d, err := s.GetDriver()
			if err != nil {
				return err
			}
			if d != "mysqlr" {
				fmt.Fprintln(os.Stdout, "driver is not mysqlr , skip")
				return nil
			}
			csm := &flink.CDCSchema{
				From: flink.Schema{
					Args: make(map[string]interface{}),
				},
				To: flink.Schema{
					Args: make(map[string]interface{}),
				},
			}
			if err := csm.Load(ctx, s, req.GetArgs()); err != nil {
				return err
			}
			if err := csm.Write(ctx, req.GetOutputPath()); err != nil {
				return err
			}
			if csm.From.Args["connector"] == nil {
				return errors.New("from.connector must be declared")
			}
			if csm.To.Args["connector"] == nil {
				return errors.New("to.connector must be declared")
			}
			exes[t] = flink.CDCExecution{
				From: flink.Execution{
					Connector: csm.From.Args["connector"].(string),
					Table:     csm.From.Table,
					Database:  csm.From.Database,
				},
				To: flink.Execution{
					Connector: csm.To.Args["connector"].(string),
					Table:     csm.To.Table,
					Database:  csm.To.Database,
				},
			}
			return nil
		}); err != nil {
			fmt.Fprintf(os.Stderr, "each error: %s\n", err)
		}
		if err := flink.WriteExecutionTemplate(ctx, exes, req.GetOutputPath()); err != nil {
			fmt.Fprintf(os.Stderr, "write execution template error: %q\n", err)
		}
	}
}
