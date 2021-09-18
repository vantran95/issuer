package issuer

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

const RPCURL = "https://volta-rpc.energyweb.org"
const key = `{"address":"c292eb3833a2506713d315e587d1521f04efbb16","crypto":{"cipher":"aes-128-ctr","ciphertext":"4be208667540e735aa0b4d6b42251db557520d207eb274b804df481745fb499b","cipherparams":{"iv":"4c5d2021efd58f68e55f663f6c62fa87"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"1e98822fed6cf3c9f80fca3e02bf587e2ac457381d49d55759d5737c05ddf013"},"mac":"44782c35cabcccb32ba4b4751372b5316decce96f7265cafd6fe95f5444dddb0"},"id":"ccaaaa04-9d99-409f-ac67-1ae17be009d8","version":3}`
const pwd = "1234" //password from Geth
const StrPrivateKey = "0xcdb00958437a1a2aca132945f7ee72cd5b23aafdeae79933c92a8d6be557483e"

func TestIssuerDeployContract(t *testing.T) {
	//const address = "http://127.0.0.1:8545"
	//const key = `{"address":"c292eb3833a2506713d315e587d1521f04efbb16","crypto":{"cipher":"aes-128-ctr","ciphertext":"4be208667540e735aa0b4d6b42251db557520d207eb274b804df481745fb499b","cipherparams":{"iv":"4c5d2021efd58f68e55f663f6c62fa87"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"1e98822fed6cf3c9f80fca3e02bf587e2ac457381d49d55759d5737c05ddf013"},"mac":"44782c35cabcccb32ba4b4751372b5316decce96f7265cafd6fe95f5444dddb0"},"id":"ccaaaa04-9d99-409f-ac67-1ae17be009d8","version":3}`
	//const key = `{"address":"9f0b91204fb9daa8c4cf0b3289c267b3c33fae4d","crypto":{"cipher":"aes-128-ctr","ciphertext":"8cf7ac663144d1caa405bd70314d0e397ab7ba3b5f160dc7b716788e980bb9b7","cipherparams":{"iv":"61bbf8f0fb01cda61af708c9fa5157b3"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"391caaecfffed1895c5e05d58cec2e176f5464ee799e9e31dea067095efb30c9"},"mac":"8e0cb285ced87b376b34db5c112dbb8bf0700f3c3a4d030e0a192a4902fb8829"},"id":"9b1ed4d8-21a7-4584-95c1-0414f2297073","version":3}`

	const address = "https://volta-rpc.energyweb.org"
	const key = `{"address":"c292eb3833a2506713d315e587d1521f04efbb16","crypto":{"cipher":"aes-128-ctr","ciphertext":"4be208667540e735aa0b4d6b42251db557520d207eb274b804df481745fb499b","cipherparams":{"iv":"4c5d2021efd58f68e55f663f6c62fa87"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"1e98822fed6cf3c9f80fca3e02bf587e2ac457381d49d55759d5737c05ddf013"},"mac":"44782c35cabcccb32ba4b4751372b5316decce96f7265cafd6fe95f5444dddb0"},"id":"ccaaaa04-9d99-409f-ac67-1ae17be009d8","version":3}`
	const pwd = "1234" //password from Geth
	conn, err := ethclient.Dial(address)
	assert.NoError(t, err)
	chainID := big.NewInt(73799)
	auth, err := bind.NewTransactorWithChainID(strings.NewReader(key), pwd, chainID)
	assert.NoError(t, err)

	rs, _, _, _ := DeployIssuer(auth, conn)
	strAddress := fmt.Sprintf("Contract deploy address: 0x%x\n", rs)
	fmt.Println(strAddress)
	assert.True(t, len(strAddress) > 0)
}

//func TestIssuerCertificateRequest(t *testing.T) {
//	const address = "https://volta-rpc.energyweb.org"
//	const key = `{"address":"c292eb3833a2506713d315e587d1521f04efbb16","crypto":{"cipher":"aes-128-ctr","ciphertext":"4be208667540e735aa0b4d6b42251db557520d207eb274b804df481745fb499b","cipherparams":{"iv":"4c5d2021efd58f68e55f663f6c62fa87"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"1e98822fed6cf3c9f80fca3e02bf587e2ac457381d49d55759d5737c05ddf013"},"mac":"44782c35cabcccb32ba4b4751372b5316decce96f7265cafd6fe95f5444dddb0"},"id":"ccaaaa04-9d99-409f-ac67-1ae17be009d8","version":3}`
//	const pwd = "1234" //password from Geth
//	conn, err := ethclient.Dial(address)
//	assert.NoError(t, err)
//	chainID := big.NewInt(73799)
//	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
//	fmt.Println("Owner Address: ", ownerAddress)
//
//	contractAddress := common.HexToAddress("0xff101748b0f90d7f27e0c8f90a70728d6a4027b6")
//	issuer, err := NewIssuer(contractAddress, conn)
//	assert.NoError(t, err)
//
//	auth, err := bind.NewTransactorWithChainID(strings.NewReader(key), pwd, chainID)
//	assert.NoError(t, err)
//	//var data []byte
//	//copy(data[:], "Draft data")
//	data := []byte{'a', 'b', 'c', 'd', 'e', 'f'}
//	fmt.Println("data: ", data)
//
//	requestCertificate, err := issuer.RequestCertificationFor(auth, data, ownerAddress)
//	assert.NoError(t, err)
//	fmt.Println("id: ", requestCertificate.Data())
//	fmt.Println("tx sent: ", requestCertificate.Hash().Hex())
//}

//func TestApproveCertificateRequest(t *testing.T) {
//	const address = "https://volta-rpc.energyweb.org"
//	const key = `{"address":"c292eb3833a2506713d315e587d1521f04efbb16","crypto":{"cipher":"aes-128-ctr","ciphertext":"4be208667540e735aa0b4d6b42251db557520d207eb274b804df481745fb499b","cipherparams":{"iv":"4c5d2021efd58f68e55f663f6c62fa87"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"1e98822fed6cf3c9f80fca3e02bf587e2ac457381d49d55759d5737c05ddf013"},"mac":"44782c35cabcccb32ba4b4751372b5316decce96f7265cafd6fe95f5444dddb0"},"id":"ccaaaa04-9d99-409f-ac67-1ae17be009d8","version":3}`
//	const pwd = "1234" //password from Geth
//	conn, err := ethclient.Dial(address)
//	assert.NoError(t, err)
//	chainID := big.NewInt(73799)
//	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
//	fmt.Println("Owner Address: ", ownerAddress)
//
//	//contractAddress := common.HexToAddress("0xff101748b0f90d7f27e0c8f90a70728d6a4027b6")
//	issuer, err := NewIssuer(ownerAddress, conn)
//	assert.NoError(t, err)
//
//	auth, err := bind.NewTransactorWithChainID(strings.NewReader(key), pwd, chainID)
//	assert.NoError(t, err)
//	//var data []byte
//	//copy(data[:], "Draft data")
//	gasLimit, err := estimateGas(conn, "0xc292eb3833a2506713d315e587d1521f04efbb16")
//	assert.NoError(t, err)
//	fmt.Println("Gas")
//	auth.GasLimit = gasLimit
//
//	approveCertificate, err := issuer.ApproveCertificationRequest(auth, big.NewInt(9), big.NewInt(200))
//	assert.NoError(t, err)
//	fmt.Println("id: ", approveCertificate.Data())
//	fmt.Println("tx sent: ", approveCertificate.Hash().Hex())
//}

func Test(t *testing.T) {
	const address = "https://volta-rpc.energyweb.org"
	const key = `{"address":"c292eb3833a2506713d315e587d1521f04efbb16","crypto":{"cipher":"aes-128-ctr","ciphertext":"4be208667540e735aa0b4d6b42251db557520d207eb274b804df481745fb499b","cipherparams":{"iv":"4c5d2021efd58f68e55f663f6c62fa87"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"1e98822fed6cf3c9f80fca3e02bf587e2ac457381d49d55759d5737c05ddf013"},"mac":"44782c35cabcccb32ba4b4751372b5316decce96f7265cafd6fe95f5444dddb0"},"id":"ccaaaa04-9d99-409f-ac67-1ae17be009d8","version":3}`
	const pwd = "1234" //password from Geth
	conn, err := ethclient.Dial(address)
	assert.NoError(t, err)
	//chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0xb3506aa9f04656761914f2c34c879c35db98cb98")
	issuer, err := NewIssuer(contractAddress, conn)
	assert.NoError(t, err)

	bytecode, err := conn.CodeAt(context.Background(), contractAddress, nil) // nil is latest block
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("block id: ", hex.EncodeToString(bytecode))

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(13579961),
		ToBlock:   big.NewInt(13579961),
		Addresses: []common.Address{
			contractAddress,
		},
	}

	logs, err := conn.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("logs", logs)

	for _, vLog := range logs {
		fmt.Println(vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
		fmt.Println(vLog.BlockNumber)     // 2394201
		fmt.Println(vLog.TxHash.Hex())    // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6
		fmt.Println("data: ", vLog.Data)

		rs, err := issuer.ParseCertificationRequested(vLog)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("results: ", rs.Raw)
	}
}

func TestGetCertificateRequest(t *testing.T) {
	const address = "https://volta-rpc.energyweb.org"
	const key = `{"address":"c292eb3833a2506713d315e587d1521f04efbb16","crypto":{"cipher":"aes-128-ctr","ciphertext":"4be208667540e735aa0b4d6b42251db557520d207eb274b804df481745fb499b","cipherparams":{"iv":"4c5d2021efd58f68e55f663f6c62fa87"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"1e98822fed6cf3c9f80fca3e02bf587e2ac457381d49d55759d5737c05ddf013"},"mac":"44782c35cabcccb32ba4b4751372b5316decce96f7265cafd6fe95f5444dddb0"},"id":"ccaaaa04-9d99-409f-ac67-1ae17be009d8","version":3}`
	const pwd = "1234" //password from Geth
	conn, err := ethclient.Dial(address)
	assert.NoError(t, err)
	//chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0xff101748b0f90d7f27e0c8f90a70728d6a4027b6")
	issuer, err := NewIssuer(contractAddress, conn)
	assert.NoError(t, err)

	//auth, err := bind.NewTransactorWithChainID(strings.NewReader(key), pwd, chainID)
	//assert.NoError(t, err)

	certificate, err := issuer.GetCertificationRequest(nil, big.NewInt(4))
	assert.NoError(t, err)

	fmt.Println("certificate owner: ", certificate.Owner)
	fmt.Println("certificate data: ", string(certificate.Data))

}

func TestGetCertificationRequestedLog(t *testing.T) {
	const address = "https://volta-rpc.energyweb.org"
	conn, err := ethclient.Dial(address)
	assert.NoError(t, err)

	//contractAddress := common.HexToAddress("0x2DcBD21CDc98aa449a4E5fd74c24C0c89F26a08b")
	//bytecode, err := conn.CodeAt(context.Background(), contractAddress, nil) // nil is latest block
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("block id: ", hex.EncodeToString(bytecode))
	//blockId := hex.EncodeToString(bytecode)

	// Get the latest block number
	header, err := conn.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(header.Number.String()) // 5671744
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(13597324),
		ToBlock:   big.NewInt(13597324),
		//Addresses: []common.Address{
		//	contractAddress,
		//},
	}

	logs, err := conn.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("logs len", len(logs))

	for _, vLog := range logs {
		//fmt.Println(vLog.BlockHash.Hex()) // 0x3404b8c050aa0aacd0223e91b5c32fee6400f357764771d0684fa7b3f448f1a8
		fmt.Println("Block num: ", vLog.BlockNumber) // 2394201
		fmt.Println("tnx hash: ", vLog.TxHash.Hex()) // 0x280201eda63c9ff6f305fcee51d5eb86167fab40ca3108ec784e8652a0e2b1a6
		//fmt.Println("data: ",vLog.Data)

		var topics [4]string
		for i := range vLog.Topics {
			topics[i] = vLog.Topics[i].Hex()
		}
		//var fillEvent LogFill

		//err := contractAbi.UnpackIntoInterface(&fillEvent, "CertificationRequested", vLog.Data)
		//if err != nil {
		//	log.Fatal(err)
		//}

		//fmt.Println("topic: ", vLog.Topics[3])

		//fillEvent.Maker = common.HexToAddress(vLog.Topics[1].Hex())
		//fillEvent.FeeRecipient = common.HexToAddress(vLog.Topics[2].Hex())
		//fillEvent.Tokens = vLog.Topics[3]
		//
		//fmt.Printf("Maker: %s\n", fillEvent.Maker.Hex())
		//fmt.Printf("Fee Recipient: %s\n", fillEvent.FeeRecipient.Hex())
		//fmt.Printf("Tokens: %s\n", hexutil.Encode(fillEvent.Tokens[:]))
		//fmt.Printf("Order Data: %s\n", hexutil.Encode(fillEvent.data[:]))
		//fmt.Printf("Owner address: %s\n", fillEvent.owner.Hex())

		//ownerAddress := common.HexToAddress(vLog.Topics[1].Hex())
		////address := crypto.Keccak256Hash(ownerAddress)
		//fmt.Println("address: ", ownerAddress)
		//
		//id, err := strconv.ParseInt(hexaNumberToInteger(vLog.Topics[2].Hex()), 16, 64)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//fmt.Println("id: ", id)
		fmt.Println("data 0 ", topics[0]) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
		eventSignature := []byte("CertificationRequestApproved(address,uint256,uint256)")
		hash := crypto.Keccak256Hash(eventSignature)
		fmt.Println("hash: ", hash.Hex())
		if hash.Hex() == topics[0] {
			fmt.Println("data 1 ", topics[1])
			fmt.Println("data 2 ", topics[2])
			fmt.Println("data 3 ", topics[3])
		}

		//fmt.Println("data 0 ", topics[0]) // 0xe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d4
		//fmt.Println("data 1 ", topics[1])
		//fmt.Println("data 2 ", topics[2])
		//fmt.Println("data 3 ", topics[3])
	}
}

func TestRequestToIssue(t *testing.T) {
	client, err := ethclient.Dial(RPCURL)
	assert.NoError(t, err)
	chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0xff101748b0f90d7f27e0c8f90a70728d6a4027b6")
	issuer, err := NewIssuer(contractAddress, client)
	assert.NoError(t, err)

	privateKey, err := RetrievePrivateKey("")
	assert.NoError(t, err)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	assert.NoError(t, err)
	const deviceOwnerAddr = "0xde0bed9f459169d438a4018c888d0ec27baa7b0b"

	tnx, err := issuer.RequestCertificationFor(auth, nil, common.HexToAddress(deviceOwnerAddr))
	assert.NoError(t, err)

	fmt.Println("gas price: ", tnx.GasPrice())
	fmt.Println("gas limit: ", tnx.Gas())
	fmt.Println("address string: ", tnx.Hash().String())
	fmt.Println("address hex: ", tnx.Hash().Hex())
}

func TestInitialize(t *testing.T) {
	client, err := ethclient.Dial(RPCURL)
	assert.NoError(t, err)
	//chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0x92218f951dFa608bD88607042fa6C760dd2d5213")
	issuer, err := NewIssuer(contractAddress, client)
	assert.NoError(t, err)

	privateKey, err := RetrievePrivateKey("")
	if err != nil {
		log.Fatal(err)
	}
	chainID := big.NewInt(73799)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	assert.NoError(t, err)

	registryAddress := common.HexToAddress("0xD3F6182F019Ef97522D1A0d07fb5165d91C00A25")

	tnx, err := issuer.Initialize(auth, big.NewInt(01), registryAddress)
	assert.NoError(t, err)
	fmt.Println("Retrieve: ", tnx.Data())
}

func TestGetRegistryAddress(t *testing.T) {
	client, err := ethclient.Dial(RPCURL)
	assert.NoError(t, err)
	//chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0x92218f951dFa608bD88607042fa6C760dd2d5213")
	issuer, err := NewIssuer(contractAddress, client)
	assert.NoError(t, err)

	//privateKey, err := RetrievePrivateKey("")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//chainID := big.NewInt(73799)
	//
	//auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	//assert.NoError(t, err)
	//
	//registryAddress := common.HexToAddress("0xD3F6182F019Ef97522D1A0d07fb5165d91C00A25")

	tnx, err := issuer.GetRegistryAddress(nil)
	assert.NoError(t, err)
	fmt.Println("Retrieve: ", tnx.Hex())
}

func TestIssuerCertificateRequest(t *testing.T) {
	client, err := ethclient.Dial(RPCURL)
	assert.NoError(t, err)
	//chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0x92218f951dFa608bD88607042fa6C760dd2d5213")
	issuer, err := NewIssuer(contractAddress, client)
	assert.NoError(t, err)

	privateKey, err := RetrievePrivateKey("")
	if err != nil {
		log.Fatal(err)
	}
	chainID := big.NewInt(73799)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	assert.NoError(t, err)
	data := []byte{'a', 'b', 'c', 'd', 'e', 'f'}
	fmt.Println("data: ", data)

	requestCertificate, err := issuer.RequestCertificationFor(auth, data, ownerAddress)
	assert.NoError(t, err)
	fmt.Println("id: ", requestCertificate.Data())
	fmt.Println("tx sent: ", requestCertificate.Hash().Hex())
}

func TestApproveCertificateRequest(t *testing.T) {
	client, err := ethclient.Dial(RPCURL)
	assert.NoError(t, err)
	//chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0x92218f951dFa608bD88607042fa6C760dd2d5213")
	issuer, err := NewIssuer(contractAddress, client)
	assert.NoError(t, err)

	privateKey, err := RetrievePrivateKey("")
	if err != nil {
		log.Fatal(err)
	}
	chainID := big.NewInt(73799)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	assert.NoError(t, err)
	//var data []byte
	//copy(data[:], "Draft data")
	gasLimit, err := estimateGas(client, "0xc292eb3833a2506713d315e587d1521f04efbb16")
	assert.NoError(t, err)
	fmt.Println("Gas: ", gasLimit)
	auth.GasLimit = gasLimit

	approveCertificate, err := issuer.ApproveCertificationRequest(auth, big.NewInt(2), big.NewInt(200))
	assert.NoError(t, err)
	fmt.Println("id: ", approveCertificate.Data())
	fmt.Println("tx sent: ", approveCertificate.Hash().Hex())
}

func TestIssue(t *testing.T) {
	client, err := ethclient.Dial(RPCURL)
	assert.NoError(t, err)
	//chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0x92218f951dFa608bD88607042fa6C760dd2d5213")
	issuer, err := NewIssuer(contractAddress, client)
	assert.NoError(t, err)

	privateKey, err := RetrievePrivateKey("")
	if err != nil {
		log.Fatal(err)
	}
	chainID := big.NewInt(73799)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	assert.NoError(t, err)
	//var data []byte
	//copy(data[:], "Draft data")
	gasLimit, err := estimateGas(client, "0xc292eb3833a2506713d315e587d1521f04efbb16")
	assert.NoError(t, err)
	fmt.Println("Gas: ", gasLimit)
	auth.GasLimit = gasLimit
	data := []byte{'a', 'b', 'c', 'd', 'e', 'f'}
	issue, err := issuer.Issue(auth, ownerAddress, big.NewInt(100), data)
	assert.NoError(t, err)
	fmt.Println("id: ", issue.Data())
	fmt.Println("tx sent: ", issue.Hash().Hex())
	fmt.Println("Certificate id: ", issue.Value())
}

func TestGetCertificateID(t *testing.T) {
	client, err := ethclient.Dial(RPCURL)
	assert.NoError(t, err)
	//chainID := big.NewInt(73799)
	ownerAddress := common.HexToAddress("0xc292eb3833a2506713d315e587d1521f04efbb16")
	fmt.Println("Owner Address: ", ownerAddress)

	contractAddress := common.HexToAddress("0x44055D7beb0e2e2Ca8B342F8b1D8aCd311441E2e")
	issuer, err := NewIssuer(contractAddress, client)
	assert.NoError(t, err)

	tnx, err := issuer.GetCertificationRequest(nil, big.NewInt(2))
	assert.NoError(t, err)
	fmt.Println("Retrieve: ", tnx.Owner)
}

func hexaNumberToInteger(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

func estimateGas(client *ethclient.Client, address string) (uint64, error) {
	tokenAddr := common.HexToAddress(address)
	estimatedGas, err := client.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddr,
		Data: []byte{0},
	})
	if err != nil {
		return 0, err
	}

	return uint64(float64(estimatedGas) * 100), nil
}

// RetrievePrivateKey retrieve private key by iamId
func RetrievePrivateKey(iamId string) (*ecdsa.PrivateKey, error) {
	// TODO: need encrypt StrPrivateKey before inserted to DB
	b, err := hexutil.Decode(StrPrivateKey)
	if err != nil {
		return nil, err
	}

	return crypto.ToECDSA(b)
}
