package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	fcfg = &FileConfig{
		Clusters: map[string]*Config{},
	}
	cmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%+v\n", fcfg)
		},
	}
)

func init() {
	cmd.Flags().Var(&fcfg.Defaults, "cluster-defaults", "default values for all clusters")
	cmd.Flags().VarP(ConfigMap(fcfg.Clusters), "cluster", "c", "per-cluster config and overrides")
}

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
