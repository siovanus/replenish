/**
 * Copyright (C) 2021 The poly network Authors
 * This file is part of The poly network library.
 *
 * The poly network is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The poly network is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with the poly network.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/gopass"
	sdk "github.com/polynetwork/poly-go-sdk"
	"github.com/siovanus/replenish/config"
	"github.com/siovanus/replenish/log"
	"strings"
)

var confFile string
var chainId uint64
var txHashes string

func init() {
	flag.StringVar(&confFile, "conf", "./config.json", "configuration file path")
	flag.Uint64Var(&chainId, "chainid", 0, "replenish chain id")
	flag.StringVar(&txHashes, "hashes", "", "tx hash list, sep by ','")
	flag.Parse()
}

func setUpPoly(polySdk *sdk.PolySdk, rpcAddr string) error {
	polySdk.NewRpcClient().SetAddress(rpcAddr)
	hdr, err := polySdk.GetHeaderByHeight(0)
	if err != nil {
		return err
	}
	polySdk.SetChainId(hdr.ChainID)
	return nil
}

func main() {
	log.InitLog(log.InfoLog, "./Logs/", log.Stdout)

	conf, err := config.LoadConfig(confFile)
	if err != nil {
		log.Fatalf("LoadConfig fail:%v", err)
		return
	}

	polySdk := sdk.NewPolySdk()
	err = setUpPoly(polySdk, conf.RestURL)
	if err != nil {
		log.Fatalf("setUpPoly failed: %v", err)
		return
	}
	wallet, err := polySdk.OpenWallet(conf.WalletFile)
	if err != nil {
		log.Fatalf("polySdk.OpenWallet failed: %v", err)
		return
	}
	pass := []byte(conf.WalletPwd)
	if len(pass) == 0 {
		fmt.Print("Enter Password: ")
		pass, err = gopass.GetPasswd()
		if err != nil {
			log.Fatalf("gopass.GetPasswd failed: %v", err)
			return
		}
	}

	signer, err := wallet.GetDefaultAccount(pass)
	if err != nil {
		log.Fatalf("wallet.GetDefaultAccount failed: %v", err)
		return
	}

	txHashList := strings.Split(txHashes, ",")
	tx, err := polySdk.Native.Rp.ReplenishTx(
		chainId,
		txHashList,
		signer)
	if err != nil {
		log.Fatalf("polySdk.Native.Rp.ReplenishTx failed: %v", err)
		return
	} else {
		log.Infof("replenish tx - send transaction to poly chain, poly_txhash: %s", tx.ToHexString())
	}
}
