package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/farishan-labs/dist-plugin/internal/goatcounter"
	"github.com/farishan-labs/dist-plugin/internal/gumroad"
)

func main() {
	root := &cobra.Command{Use: "dist-plugin", Short: "Metric provider plugins for dist"}
	root.AddCommand(&cobra.Command{
		Use:   "gumroad",
		Short: "Collect Gumroad metrics as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := gumroad.CollectJSON()
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	})
	root.AddCommand(&cobra.Command{
		Use:   "goatcounter",
		Short: "Collect GoatCounter metrics as JSON",
		RunE: func(cmd *cobra.Command, args []string) error {
			b, err := goatcounter.CollectJSON()
			if err != nil {
				return err
			}
			fmt.Println(string(b))
			return nil
		},
	})
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
