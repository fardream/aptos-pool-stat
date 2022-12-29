package main

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/fardream/go-aptos/aptos"
	"github.com/fardream/go-aptos/aptos/known"
	"github.com/fardream/go-aptos/aptos/stat"
)

type protocolStat struct {
	coinListOutput string
	poolListOutput string
	protocol       *stat.ConstantProductPoolProtocol
	client         *aptos.Client
}

func (s *protocolStat) after() {
	known.ReloadHippoCoinRegistry(known.HippoCoinRegistryUrl)

	usdSymbols := []string{
		"USDC",
		"USDCso",
		"USDT",
		"ceUSDC",
		"ceDAI",
		"ceUSDT",
		"zUSDC",
		"zUSDT",
		"BUSD",
		"ceBUSD",
		"zBUSD",
	}

	for _, usdSymbol := range usdSymbols {
		stable := known.GetCoinInfoBySymbol(aptos.Mainnet, usdSymbol)
		if stable != nil {
			s.protocol.AddStableCoins(stable.TokenType.Type)
		}
	}

	s.protocol.FillCoinInfo(context.Background(), aptos.Mainnet, s.client)

	s.protocol.FillStat()

	var coinBuf bytes.Buffer

	fmt.Fprintln(&coinBuf, "Coin Type, Coin Symbol, Coin Name, Coin Decimal, Total Reserve, Price, Total Value, IsHippo")

	for _, coin := range s.protocol.Coins {
		isHippo := 0
		if coin.IsHippo {
			isHippo = 1
		}

		fmt.Fprintf(&coinBuf, "%s,%s,%s,%d,%d,%f,%f,%d\n", coin.MoveTypeTag.String(), coin.Symbol, coin.Name, coin.Decimals, coin.TotalQuantity, coin.Price, coin.TotalValue, isHippo)
	}

	if s.coinListOutput == "-" {
		fmt.Print(coinBuf.String())
	} else if s.coinListOutput != "" {
		os.WriteFile(s.coinListOutput, coinBuf.Bytes(), 0o666)
	}

	var poolBuf bytes.Buffer

	fmt.Fprint(&poolBuf, "Coin 0 Type, Coin 0 Symbol, Coin 0 Is Hippo, Coin 0 Decimal, Coin 0 Reserve, Coin 0 Price, Coin 0 Value,")
	fmt.Fprint(&poolBuf, "Coin 1 Type, Coin 1 Symbol, Coin 1 Is Hippo, Coin 1 Decimal, Coin 1 Reserve, Coin 1 Price, Coin 1 Value,")
	fmt.Fprintln(&poolBuf, "Total Value")

	for _, pool := range s.protocol.Pools {
		coin0Name := pool.Coin0.String()
		coin0Info := s.protocol.Coins[coin0Name]
		coin0IsHippo := 0
		if coin0Info.IsHippo {
			coin0IsHippo = 1
		}
		fmt.Fprintf(&poolBuf, "%s,%s,%d,%d,%d,%f,%f,", coin0Name, coin0Info.Symbol, coin0IsHippo, coin0Info.Decimals, pool.Coin0Reserve, coin0Info.Price, pool.Coin0Value)

		coin1Name := pool.Coin1.String()
		coin1Info := s.protocol.Coins[coin1Name]
		coin1IsHippo := 0
		if coin1Info.IsHippo {
			coin1IsHippo = 1
		}
		fmt.Fprintf(&poolBuf, "%s,%s,%d,%d,%d,%f,%f,", coin1Name, coin1Info.Symbol, coin1IsHippo, coin1Info.Decimals, pool.Coin1Reserve, coin1Info.Price, pool.Coin1Value)

		fmt.Fprintf(&poolBuf, "%f\n", pool.TotalValueLocked)
	}

	if s.poolListOutput == "-" {
		fmt.Print(poolBuf.String())
	} else if s.poolListOutput != "" {
		os.WriteFile(s.poolListOutput, poolBuf.Bytes(), 0o666)
	}
}
