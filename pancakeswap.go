package main

import (
	"context"
	"encoding/json"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/stat"
	"github.com/spf13/cobra"
)

type TokenPairReserve struct {
	ReserveX aptos.JsonUint64 `json:"reserve_x"`
	ReserveY aptos.JsonUint64 `json:"reserve_y"`
}

func pancakeSwapCmd(s *protocolStat) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pancakeswap",
		Short: "show stats for pancakeswap",
		Args:  cobra.NoArgs,
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		s.client = aptos.MustNewClient(aptos.Mainnet, "")
		s.protocol = stat.NewStatForConstantProductPool()

		resp, err := s.client.GetAccountResources(context.Background(), &aptos.GetAccountResourcesRequest{
			Address: aptos.MustParseAddress("0xc7efb4076dbe143cbcd98cfaaa929ecfc8f299203dfff63b95ccb6bfe19850fa"),
		})
		if err != nil {
			panic(err)
		}

		for _, v := range *resp.Parsed {
			if v.Type.Name == "TokenPairReserve" {

				var reserve TokenPairReserve

				if err := json.Unmarshal(v.Data, &reserve); err != nil {
					continue
				}

				coin0 := v.Type.GenericTypeParameters[0].Struct
				coin1 := v.Type.GenericTypeParameters[1].Struct

				s.protocol.AddSinglePool(coin0, uint64(reserve.ReserveX), coin1, uint64(reserve.ReserveY))
			}
		}

		s.after()
	}

	return cmd
}
