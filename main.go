package main

import (
	"lockchain/blockchain"
)

func main() {
	bc := blockchain.NewBlockchain()
	defer bc.DB.Close()

	cli := blockchain.CLI{bc}
	cli.Run()
}
