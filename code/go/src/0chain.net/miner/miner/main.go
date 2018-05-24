package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime"

	"0chain.net/block"
	"0chain.net/chain"
	"0chain.net/client"
	"0chain.net/common"
	"0chain.net/config"
	"0chain.net/encryption"
	"0chain.net/node"
	"0chain.net/transaction"
)

func initServer() {
	// TODO; when a new server is brought up, it needs to first download all the state before it can start accepting requests
}

func initHandlers() {
	if config.Configuration.TestMode {
		http.HandleFunc("/_hash", encryption.HashHandler)
		http.HandleFunc("/_sign", common.ToJSONResponse(encryption.SignHandler))
	}
	node.SetupHandlers()

	chain.SetupHandlers()
	client.SetupHandlers()
	transaction.SetupHandlers()
	block.SetupHandlers()
}

/*Chain - the chain this miner will be working on */
var Chain string

func main() {
	host := flag.String("host", "", "hostname")
	port := flag.Int("port", 7220, "port")
	chainID := flag.String("chain", "", "chain id")
	testMode := flag.Bool("test", false, "test mode?")
	nodesFile := flag.String("nodes_file", "config/single_node.txt", "nodes_file")
	keysFile := flag.String("keys_file", "config/single_node_keys.txt", "keys_file")
	flag.Parse()

	address := fmt.Sprintf("%v:%v", *host, *port)
	chain.SetServerChainID(*chainID)
	config.Configuration.Host = *host
	config.Configuration.Port = *port
	config.Configuration.ChainID = *chainID
	config.Configuration.TestMode = *testMode

	reader, err := os.Open(*keysFile)
	if err != nil {
		panic(err)
	}
	publicKey, privateKey := encryption.ReadKeys(reader)
	reader.Close()

	if *nodesFile == "" {
		panic("Please specify --node_file file.txt option with a file.txt containing peer nodes")
	}

	reader, err = os.Open(*nodesFile)
	if err != nil {
		panic(err)
	}
	node.ReadNodes(reader, &node.Miners, &node.Sharders, &node.Blobbers)
	reader.Close()
	if node.Self == nil {
		panic("node definition for self node doesn't exist")
	} else {
		if node.Self.PublicKey != publicKey {
			fmt.Printf("self: %v\n", node.Self)
			panic(fmt.Sprintf("Pulbic key from the keys file and nodes file don't match %v %v", publicKey, node.Self.PublicKey))
		}
		node.Self.SetPrivateKey(privateKey)
	}

	go node.Miners.StatusMonitor()
	go node.Sharders.StatusMonitor()
	go node.Blobbers.StatusMonitor()

	mode := "main net"
	if *testMode {
		mode = "test net"
		block.BLOCK_SIZE = 100
	}
	fmt.Printf("Num CPUs available %v\n", runtime.NumCPU())
	fmt.Printf("Starting %v on %v for chain %v in %v mode ...\n", os.Args[0], address, chain.GetServerChainID(), mode)
	initServer()
	initHandlers()
	if err := http.ListenAndServe(address, nil); err != nil {
		panic(err)
	}
}
