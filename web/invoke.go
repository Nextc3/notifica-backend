package web

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Invoke handles chaincode invoke requests.
func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received Invoke request")
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %s", err)
		return
	}
	//chainCodeName := r.FormValue("chaincodeid")
	chainCodeName := "notifica-chaincode"
	//channelID := r.FormValue("channelid")
	channelID := "mychannel"
	//function := r.FormValue("function")
	function := ""
	//usando uma a variável function fica mais fácil manipular qual metódo será executado
	//e evita repetição de código
	sid := strings.TrimPrefix(r.URL.Path, "/notificacao/")
	id, _ := strconv.Atoi(sid)
	switch {
	case r.Method == "GET" && id > 0:
		function = "ReadAsset"
	case r.Method == "GET":
		function = "GetAllAssets"
	case r.Method == "POST":
		function = "CreateAsset"
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Desculpa... :(")
	}
	args := r.Form["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, args)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(args...))
	if err != nil {
		fmt.Fprintf(w, "Error creating txn proposal: %s", err)
		return
	}
	txn_endorsed, err := txn_proposal.Endorse()
	if err != nil {
		fmt.Fprintf(w, "Error endorsing txn: %s", err)
		return
	}
	txn_committed, err := txn_endorsed.Submit()
	if err != nil {
		fmt.Fprintf(w, "Error submitting transaction: %s", err)
		return
	}
	fmt.Fprintf(w, "Transaction ID : %s Response: %s", txn_committed.TransactionID(), txn_endorsed.Result())
}
