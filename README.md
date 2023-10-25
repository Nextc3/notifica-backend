# Cliente API REST

Este módulo é um servidor REST API escrito em golang com endpoints para acesso ao chaincode da blockchain com invoke(ações que alteram estado) e query(ações que não alteram o estado).

  
## Uso

- Setup fabric test network and deploy the asset transfer chaincode by [following this instructions](https://hyperledger-fabric.readthedocs.io/en/release-2.4/test_network.html).

- Entre no diretório notifica-backend 
- Baixe as dependências com o comando `go mod download`
- Execute `go run main.go` para iniciar o servidor REST
- Caso dẽ o seguinte erro:
-  failed to read certificate file: open ../../test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/cert.pem: no such file or directory
- vá na pasta ../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/ e copie o arquivo existente com extensão .pem e coloque o nome cert.pem sem apagar o original e execute de novo

## Sending Requests
Invoke endpoint aceita solicitações POST com funções e argumentos de chaincode. O endpoint de consulta aceita solicitações get com função e argumentos de chaincode.

Um exemplo simples de usando invoke seria o método "createAsset" que cria uma novo asset Response will contain transaction ID for a successful invoke.

``` sh
curl --request POST \
  --url http://localhost:3000/invoke \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data = \
  --data channelid=mychannel \
  --data chaincodeid=basic \
  --data function=createAsset \
  --data args=Asset123 \
  --data args=yellow \
  --data args=54 \
  --data args=Tom \
  --data args=13005
```
Sample chaincode query for getting asset details.

``` sh
curl --request GET \
  --url 'http://localhost:3000/query?channelid=mychannel&chaincodeid=basic&function=ReadAsset&args=Asset123' 
  ```
