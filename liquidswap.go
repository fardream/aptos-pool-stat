package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/stat"
	"github.com/spf13/cobra"
)

type LiquidSwapLiquidityPool struct {
	CoinXReserve aptos.Coin `json:"coin_x_reserve"`
	CoinYReserve aptos.Coin `json:"coin_y_reserve"`
}

func liquidswapCmd(s *protocolStat) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "liquidswap",
		Short: "display stat for liquidswap",
		Args:  cobra.NoArgs,
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		s.client = aptos.MustNewClient(aptos.Mainnet, "")
		s.protocol = stat.NewStatForConstantProductPool()

		resp, err := s.client.GetAccountResources(context.Background(), &aptos.GetAccountResourcesRequest{
			Address: aptos.MustParseAddress("0x05a97986a9d031c4567e15b797be516910cfcb4156312482efc6a19c0a30c948"),
		})
		if err != nil {
			panic(err)
		}

		for _, resource := range *resp.Parsed {
			resourceType := resource.Type

			if resourceType.Module == "liquidity_pool" && resourceType.Name == "LiquidityPool" {
				var amm LiquidSwapLiquidityPool
				if err := json.Unmarshal(resource.Data, &amm); err != nil {
					fmt.Printf("failed to parse %s: %s\n", string(resource.Data), err.Error())
					continue
				}
				s.protocol.AddSinglePool(resourceType.GenericTypeParameters[0].Struct, uint64(amm.CoinXReserve.Value), resourceType.GenericTypeParameters[1].Struct, uint64(amm.CoinYReserve.Value))
			}
		}

		s.after()
	}

	return cmd
}
