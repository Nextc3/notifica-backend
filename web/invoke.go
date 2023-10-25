package web

import (
	"encoding/json"
	"fmt"
	notifica_model "github.com/Nextc3/notifica-model"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// Invoke handles chaincode invoke requests.
func (setup *OrgSetup) Invoke(w http.ResponseWriter, r *http.Request) {
	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	//header.Set("Access-Control-Allow-Origin", "*")

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
	asset := notifica_model.Asset{}

	_ = json.NewDecoder(r.Body).Decode(&asset)
	//args := []string{"dataNascimento", asset.DataNascimento, "dataDiagnostico", asset.DataDiagnostico, "dataNotificacao", asset.DataNotificacao, "dataInicioSintomas", asset.DataInicioSintomas, "bairro", asset.Bairro, "cidade", asset.Cidade, "endereco", asset.Endereco, "estado", asset.Estado, "pais", asset.Pais, "doenca", asset.Doenca, "informacoesClinicas", asset.InformacoesClinicas, "sexo", asset.Sexo}
	//args := r.Form["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, asset)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)
	resultado, err := contract.EvaluateTransaction("getUltimoId")
	if err != nil {
		log.Println("erro em buscar o último Id")
	}
	aux, _ := strconv.Atoi(string(resultado))
	aux++
	asset.Id = aux

	nEmBytes, _ := json.Marshal(asset)
	aString := string(nEmBytes)
	txn_proposal, err := contract.NewProposal(function, client.WithArguments(aString))
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
