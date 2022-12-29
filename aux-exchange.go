package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/stat"
	"github.com/spf13/cobra"
)

func auxCmd(s *protocolStat) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "aux",
		Short: "show stats for aux exchange",
		Args:  cobra.NoArgs,
	}

	cmd.Run = func(cmd *cobra.Command, args []string) {
		s.client = aptos.MustNewClient(aptos.Mainnet, "")
		s.protocol = stat.NewStatForConstantProductPool()

		auxConfig, _ := aptos.GetAuxClientConfig(aptos.Mainnet)

		resp, err := s.client.GetAccountResources(context.Background(), &aptos.GetAccountResourcesRequest{
			Address: auxConfig.Address,
		})
		if err != nil {
			panic(err)
		}

		for _, resource := range *resp.Parsed {
			resourceType := resource.Type

			if resourceType.Module == "amm" && resourceType.Name == "Pool" && aptos.IsAddressEqual(&auxConfig.Address, &resourceType.Address) {
				var amm aptos.AuxAmmPool
				if err := json.Unmarshal(resource.Data, &amm); err != nil {
					fmt.Printf("failed to parse %s due to %v\n", string(resource.Data), err)
					continue
				}
				s.protocol.AddSinglePool(resourceType.GenericTypeParameters[0].Struct, uint64(amm.XReserve.Value), resourceType.GenericTypeParameters[1].Struct, uint64(amm.YReserve.Value))
			}
		}

		s.after()
	}

	return cmd
}
