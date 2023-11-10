(O notifica-backend, notifica-models, notifica-chaincode e o Hyperledger Fabric foram desenvolvidos em plataforma Linux Ubuntu 20.04)
Instalar os seguintes softwares e versões mais recentes(com exceção do Go):
•	git
•	curl
•	docker
•	docker compose
•	go 1.19
•	jq
Caso queira interagir usando terminal Git Bash no Windows para interação com API REST do notifica-backend faça os seguintes comandos para reconhecimento de quebra de linha em Linux e tratamento e de endereços:
git config --global core.autocrlf false
git config --global core.longpaths true

No terminal do Linux onde estarão todos componentes, exceto notifica-frontend, execute comandos para configuração do golang:
export GOPATH=$HOME/go

export PATH=$PATH:$GOPATH/bin

mkdir -p $HOME/go/src/github.com/<usuário_no_github>
cd $HOME/go/src/github.com/<usuário_no_github>

Para instalar Hyperledger Fabric:

curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh
./install-fabric.sh

Para colocar visível ao sistema os binários do hyperledger, no arquivo /home/seuusuario/.bashrc, no final coloque:
 
export FABRIC_RAIZ=~/go/src/github.com/<usuariogithub>/fabric-samples/

export FABRIC_CFG_PATH=$FABRIC_RAIZ/config/

export export PATH=$FABRIC_RAIZ/bin:$PATH

export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="Org1MSP"
export CORE_PEER_TLS_ROOTCERT_FILE=$FABRIC_RAIZ/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=$FABRIC_RAIZ/organizations/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
export CORE_PEER_ADDRESS=localhost:7051

saia do seu usuário e faça login novamente.

No terminal execute:
cd $FABRIC_RAIZ
cd ..

faça um git clone do model:

git clone https://github.com/Nextc3/notifica-model
cd fabric-samples
cd asset-transfer-basic

faça git clone do chaincode e models
git clone https://github.com/Nextc3/notifica-backend
git clone https://github.com/Nextc3/notifica-chaincode

entre no diretório da rede de test:
cd $FABRIC_RAIZ/test-network

Levante a rede com couchdb, crie um canal com nome padrão:
./network.sh up -s couchdb
./network.sh createChannel 

Faça deploy do chaincode

./network.sh deployCC -ccn notifica-chaincode -ccp ../asset-transfer-basic/notifica-chaincode -ccl go

vá na pasta ../fabric-samples/test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/ e copie o arquivo existente com extensão .pem e coloque o nome cert.pem sem apagar o original

vá na pasta do backend:
cd $FABRIC_RAIZ/asset-transfer-basic/notifica-backend

execute

go run main.go

se tudo estiver ok um servidor web será criado escutando na porta 8080.

Teste se o chaincode funciona e execute dentro da pasta $FABRIC_RAIZ/test-network:



peer chaincode invoke -o localhost:7050 --ordererTLSHostnameOverride orderer.example.com --tls --cafile "${PWD}/organizations/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem" -C mychannel -n notifica-frontend --peerAddresses localhost:7051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt" --peerAddresses localhost:9051 --tlsRootCertFiles "${PWD}/organizations/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt" -c '{"function":"InitLedger","Args":[]}'

Se ocorrer certo retornará:
INFO 001 Chaincode invoke successful. result: status:200

Faça uma consulta ao chaincode com:
peer chaincode query -C mychannel -n notifica-chaincode -c '{"Args":["GetAllAssets"]}'

retornará asset padrão criado com método InitLedger

Em outro computador baixe o frontend:

git clone https://github.com/Nextc3/notifica-frontend 

tenha instalado Node.JS, npm, React.JS

Com todos instalados entre dentro de notifica-frontend e execute:
npm start

No arquivo: notifica-frontend\src\main\App.js na linha 27 coloque o ip onde o notifica-backend está executando:
const baseUrl = “http://192.168.1.80:8080/notificacao/”

