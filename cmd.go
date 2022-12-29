package main

import "github.com/spf13/cobra"

func main() {
	s := &protocolStat{}
	cmd := &cobra.Command{
		Use:   "aptos-pool-stat",
		Short: "display stats for aptos pool",
		Args:  cobra.NoArgs,
	}

	cmd.PersistentFlags().StringVarP(&s.coinListOutput, "coin-output", "c", s.coinListOutput, "coin list output")
	cmd.MarkPersistentFlagFilename("coin-output", "csv", "txt")
	cmd.PersistentFlags().StringVarP(&s.poolListOutput, "pool-output", "p", s.poolListOutput, "pool list output")
	cmd.MarkPersistentFlagFilename("pool-output", "csv", "txt")

	cmd.AddCommand(
		auxCmd(s),
		pancakeSwapCmd(s),
		liquidswapCmd(s),
	)

	cmd.Execute()
}
