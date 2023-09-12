package main

import (
	"bufio"
	"context"
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
		if err := req.Each(func(_ plugin.TemplateName, s plugin.Schema) error {
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
			return nil
		}); err != nil {
			fmt.Fprintf(os.Stderr, "each error: %s\n", err)
		}
	}
}
