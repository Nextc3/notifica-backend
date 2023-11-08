package web

import (
	"encoding/json"
	"fmt"
	notifica_model "github.com/Nextc3/notifica-model"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func formatJSONError(mensagem string) []byte {
	appErro := struct {
		Mensagem string `json:"mensagem"`
	}{
		mensagem,
	}
	response, err := json.Marshal(appErro)
	if err != nil {
		return []byte(err.Error())
	}
	return response
}

func (setup *OrgSetup) AssetHandler(w http.ResponseWriter, r *http.Request) {
	//enableCors(&w)
	//chainCodeName := r.FormValue("chaincodeid")
	chainCodeName := "notifica-chaincode"
	//channelID := r.FormValue("channelid")
	channelID := "mychannel"
	//function := r.FormValue("function")
	function := ""
	//usando uma a variável function fica mais fácil manipular qual metódo será executado
	//e evita repetição de código
	asset := notifica_model.Asset{}
	_ = json.NewDecoder(r.Body).Decode(&asset)
	//args := []string{"dataNascimento", asset.DataNascimento, "dataDiagnostico", asset.DataDiagnostico, "dataNotificacao", asset.DataNotificacao, "dataInicioSintomas", asset.DataInicioSintomas, "bairro", asset.Bairro, "cidade", asset.Cidade, "endereco", asset.Endereco, "estado", asset.Estado, "pais", asset.Pais, "doenca", asset.Doenca, "informacoesClinicas", asset.InformacoesClinicas, "sexo", asset.Sexo}
	//args := r.Form["args"]
	fmt.Printf("channel: %s, chaincode: %s, function: %s, args: %s\n", channelID, chainCodeName, function, asset)
	network := setup.Gateway.GetNetwork(channelID)
	contract := network.GetContract(chainCodeName)

	header := w.Header()
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Headers", "*")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)

	}

	//-------
	//

	sid := strings.TrimPrefix(r.URL.Path, "/notificacao/")
	id, _ := strconv.Atoi(sid)

	switch {
	case r.Method == "GET" && id > 0:
		obterNotificacao(w, r, sid, contract)
	case r.Method == "GET":
		obterTodasNotificacoes(w, r, contract)
	case r.Method == "POST":
		salvarNotificacao(w, r, contract)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Desculpa... :(")
	}
}

func obterNotificacao(w http.ResponseWriter, r *http.Request, id string, contract *client.Contract) {
	//enableCors(&w)
	existe, err := contract.SubmitTransaction("assetExists", id)
	deuCerto, err := strconv.ParseBool(string(existe))
	a := notifica_model.Asset{}
	if err != nil {
		log.Fatalf("Erro em converter String para Bool e saber se existe notificação")
	}
	if !deuCerto {

		json, _ := json.Marshal(a)
		fmt.Println("Não encontrada asset na ledger. Retornando vazia")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(json))

	}

	aEmBytes, err := contract.SubmitTransaction("readAsset", id)
	if err != nil {
		log.Fatalf("Falhou em Transação consultarNotificacao : %v\n", err)
	}
	err = json.Unmarshal(aEmBytes, &a)
	if err != nil {
		w.Write([]byte("Não encontrada Notificação"))
		w.WriteHeader(http.StatusNotFound)
		w.Write(formatJSONError(err.Error()))
		return
	}

	json, _ := json.Marshal(a)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))

}
func obterTodasNotificacoes(w http.ResponseWriter, r *http.Request, contract *client.Contract) {

	log.Println("--> Contrato: Transação ObterTodasNotificacoes, função que retorna todos os ativos na ledger")
	var resultEmBytes []byte
	var assets []*notifica_model.Asset
	var err error

	log.Println("Consultando contratointeligente")
	fmt.Println(contract.ChaincodeName())
	resultEmBytes, err = contract.SubmitTransaction("getAllAssets")

	log.Println("Obteve resultado do contrato inteligente")

	if err != nil {
		log.Fatalf("Falhou em obter todos os assets transação: %v", err)
	}

	_ = json.Unmarshal(resultEmBytes, &assets)
	if err != nil {
		//passa um erro como resposta. Sinaliza também no cabeçalho
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))

		return
	}

	json, _ := json.Marshal(assets)

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(json))

}
func salvarNotificacao(w http.ResponseWriter, r *http.Request, contract *client.Contract) {
	w.Header().Set("Content-Type", "application/json")

	//vamos pegar os dados enviados pelo usuário via body
	var asset notifica_model.Asset

	err := json.NewDecoder(r.Body).Decode(&asset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(formatJSONError(err.Error()))
		return
	}
	log.Println("--> Transação de Submit: CreateAsset, cria ativos do tipo Asset")

	resultado, err := contract.EvaluateTransaction("getUltimoId")
	if err != nil {
		log.Println("erro em buscar o último Id")
	}
	aux, _ := strconv.Atoi(string(resultado))
	aux++
	asset.Id = aux

	aEmBytes, _ := json.Marshal(asset)
	aString := string(aEmBytes)
	result, err := contract.SubmitTransaction("createAsset", aString)
	if err != nil {
		log.Fatalf("Falhou a transação de Criar Asset SUBMIT (altera estado da ledger) transação: %v", err)
	}
	log.Println(string(result))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatJSONError(err.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
}
