package web

import (
	"fmt"
	"net/http"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// OrgSetup contains organization's config to interact with the network.
type OrgSetup struct {
	OrgName      string
	MSPID        string
	CryptoPath   string
	CertPath     string
	KeyPath      string
	TLSCertPath  string
	PeerEndpoint string
	GatewayPeer  string
	Gateway      client.Gateway
}

// Serve starts http web server.
func Serve(setups OrgSetup) {
	//Fica pra depois fazer uma solução que implemente com query também
	//http.HandleFunc("/query", setups.Query)
	//http.HandleFunc("/invoke", setups.Invoke)
	http.HandleFunc("/notificacao/", setups.AssetHandler)
	fmt.Println("Escutando (http://localhost:8080/)...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
