package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"log"
	"math"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go_web/store"
	"github.com/go_web/token"
	"github.com/go_web/util"
	"golang.org/x/crypto/sha3"
)

const (
	testUrl          = "ethereum-sepolia-rpc.publicnode.com"
	contractBytecode = "608060405234801561000f575f5ffd5b5060405161087838038061087883398181016040528101906100319190610193565b805f908161003f91906103ea565b50506104b9565b5f604051905090565b5f5ffd5b5f5ffd5b5f5ffd5b5f5ffd5b5f601f19601f8301169050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6100a58261005f565b810181811067ffffffffffffffff821117156100c4576100c361006f565b5b80604052505050565b5f6100d6610046565b90506100e2828261009c565b919050565b5f67ffffffffffffffff8211156101015761010061006f565b5b61010a8261005f565b9050602081019050919050565b8281835e5f83830152505050565b5f610137610132846100e7565b6100cd565b9050828152602081018484840111156101535761015261005b565b5b61015e848285610117565b509392505050565b5f82601f83011261017a57610179610057565b5b815161018a848260208601610125565b91505092915050565b5f602082840312156101a8576101a761004f565b5b5f82015167ffffffffffffffff8111156101c5576101c4610053565b5b6101d184828501610166565b91505092915050565b5f81519050919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061022857607f821691505b60208210810361023b5761023a6101e4565b5b50919050565b5f819050815f5260205f209050919050565b5f6020601f8301049050919050565b5f82821b905092915050565b5f6008830261029d7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82610262565b6102a78683610262565b95508019841693508086168417925050509392505050565b5f819050919050565b5f819050919050565b5f6102eb6102e66102e1846102bf565b6102c8565b6102bf565b9050919050565b5f819050919050565b610304836102d1565b610318610310826102f2565b84845461026e565b825550505050565b5f5f905090565b61032f610320565b61033a8184846102fb565b505050565b5b8181101561035d576103525f82610327565b600181019050610340565b5050565b601f8211156103a25761037381610241565b61037c84610253565b8101602085101561038b578190505b61039f61039785610253565b83018261033f565b50505b505050565b5f82821c905092915050565b5f6103c25f19846008026103a7565b1980831691505092915050565b5f6103da83836103b3565b9150826002028217905092915050565b6103f3826101da565b67ffffffffffffffff81111561040c5761040b61006f565b5b6104168254610211565b610421828285610361565b5f60209050601f831160018114610452575f8415610440578287015190505b61044a85826103cf565b8655506104b1565b601f19841661046086610241565b5f5b8281101561048757848901518255600182019150602085019450602081019050610462565b868310156104a457848901516104a0601f8916826103b3565b8355505b6001600288020188555050505b505050505050565b6103b2806104c65f395ff3fe608060405234801561000f575f5ffd5b506004361061003f575f3560e01c806348f343f31461004357806354fd4d5014610073578063f56256c714610091575b5f5ffd5b61005d600480360381019061005891906101d7565b6100ad565b60405161006a9190610211565b60405180910390f35b61007b6100c2565b604051610088919061029a565b60405180910390f35b6100ab60048036038101906100a691906102ba565b61014d565b005b6001602052805f5260405f205f915090505481565b5f80546100ce90610325565b80601f01602080910402602001604051908101604052809291908181526020018280546100fa90610325565b80156101455780601f1061011c57610100808354040283529160200191610145565b820191905f5260205f20905b81548152906001019060200180831161012857829003601f168201915b505050505081565b8060015f8481526020019081526020015f20819055507fe79e73da417710ae99aa2088575580a60415d359acfad9cdd3382d59c80281d48282604051610194929190610355565b60405180910390a15050565b5f5ffd5b5f819050919050565b6101b6816101a4565b81146101c0575f5ffd5b50565b5f813590506101d1816101ad565b92915050565b5f602082840312156101ec576101eb6101a0565b5b5f6101f9848285016101c3565b91505092915050565b61020b816101a4565b82525050565b5f6020820190506102245f830184610202565b92915050565b5f81519050919050565b5f82825260208201905092915050565b8281835e5f83830152505050565b5f601f19601f8301169050919050565b5f61026c8261022a565b6102768185610234565b9350610286818560208601610244565b61028f81610252565b840191505092915050565b5f6020820190508181035f8301526102b28184610262565b905092915050565b5f5f604083850312156102d0576102cf6101a0565b5b5f6102dd858286016101c3565b92505060206102ee858286016101c3565b9150509250929050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b5f600282049050600182168061033c57607f821691505b60208210810361034f5761034e6102f8565b5b50919050565b5f6040820190506103685f830185610202565b6103756020830184610202565b939250505056fea26469706673582212200130a39a0045c96b27ab38d5ad599b97445c46843cd214811a764799cdb321d364736f6c634300081d0033"
	contractAddress  = "0xCAF8146F3de921Fa76B869a5cC5A6Fd0309B8085"
	tokenAddress     = "0x787275d81d57c431c2b26b909e2ce69d8b189f59"
	keyfilepath      = "E:\\vscode\\UTC--2025-04-11T03-25-33.208480829Z--22f2efae7f8155120b11973a3315977540f6f58b"
)

func Connect() *ethclient.Client {
	client, err := ethclient.Dial("https://" + testUrl)
	if err != nil {
		log.Fatal(err)
	}
	return client
}

var _privateKey string

func PrivateKey() {
	// Read key from file.
	keyjson, err := os.ReadFile(keyfilepath)
	if err != nil {
		utils.Fatalf("Failed to read the keyfile at '%s': %v", keyfilepath, err)
	}

	// Decrypt key with passphrase.
	passphrase := "12345"
	key, err := keystore.DecryptKey(keyjson, passphrase)
	if err != nil {
		utils.Fatalf("Error decrypting key: %v", err)
	}

	address := key.Address.Hex()
	privateKey := hex.EncodeToString(crypto.FromECDSA(key.PrivateKey))
	publicKey := hex.EncodeToString(crypto.FromECDSAPub(&key.PrivateKey.PublicKey))
	_privateKey = privateKey

	util.Log("PrivateKey: ", privateKey)
	util.Log("PublicKey: ", publicKey)
	util.Log("Address: ", address)
	util.Separator()
}

func key(client *ethclient.Client) (*ecdsa.PrivateKey, uint64) {
	// util.Log("_privateKey: ", _privateKey)
	privateKey, err := crypto.HexToECDSA(_privateKey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	util.Log("From Address: ", fromAddress.Hex())

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	util.Log("PendingNonce: ", nonce)

	return privateKey, nonce
}

func DepolyContract() {
	client := Connect()
	privateKey, nonce := key(client)

	// 获取建议的gas价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasPrice = big.NewInt(int64(10000000000))
	gasLimit := uint64(500000)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice
	address, tx, instance, err := store.DeployStore(auth, client, "1.0")
	if err != nil {
		log.Fatal(err)
	}
	_ = instance

	util.Log("Contract Address: ", address.Hex())

	/* // 解码合约字节码
	data, err := hex.DecodeString(contractBytecode)
	if err != nil {
		log.Fatal(err)
	}
	// 创建交易
	auth := types.NewContractCreation(nonce, big.NewInt(0), gasLimit, gasPrice, data)
	tx, err := types.SignTx(auth, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	// 发送交易
	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatal(err)
	}
	*/
	// 等待交易被挖矿
	receipt, err := waitForReceipt(client, tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	util.Log("Contract deployed at: ", receipt.ContractAddress.Hex())

	util.Log("Transaction Hash: ", tx.Hash().Hex())
	util.Separator()
}

func LoadContract() {
	client := Connect()
	privateKey, err := crypto.HexToECDSA(_privateKey)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	opt, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	storeContract, err := store.NewStore(common.HexToAddress(contractAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	var key [32]byte
	var value [32]byte
	copy(key[:], "store_key3")
	copy(value[:], "store_value3")
	// 调用合约方法
	tx, err := storeContract.SetItem(opt, key, value)
	if err != nil {
		log.Fatal(err)
	}
	util.Log("Transaction Hash: ", tx.Hash().Hex())

	receipt, err := waitForReceipt(client, tx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	util.Log("SetItem result: ", receipt.Status)

	//查询数据
	call := &bind.CallOpts{Context: context.Background()}
	value, err = storeContract.Items(call, key)
	if err != nil {
		log.Fatal(err)
	}
	util.Log("Value: ", string(value[:]))
	util.Separator()

}

func QueryEvent() {
	client := Connect()
	ctx := context.Background()

	contractAddr := common.HexToAddress(contractAddress)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(8120130),
		// ToBlock:   big.NewInt(0),
		Addresses: []common.Address{
			contractAddr,
			common.HexToAddress(tokenAddress),
		},
		Topics: [][]common.Hash{
			// {common.BytesToHash(crypto.Keccak256([]byte("ItemSet(bytes32,bytes32)")))},
		},
	}

	logs, err := client.FilterLogs(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	util.Log("logs count: ", len(logs))
	contractAbi, err := abi.JSON(strings.NewReader(store.StoreMetaData.ABI))
	if err != nil {
		log.Fatal(err)
	}
	// 示例：ERC20 Token 的 ABI（仅包含 Transfer 事件）
	erc20ABI, err := abi.JSON(strings.NewReader(`[
		{
			"anonymous": false,
			"inputs": [
				{"indexed": true, "name": "from", "type": "address"},
				{"indexed": true, "name": "to", "type": "address"},
				{"indexed": false, "name": "value", "type": "uint256"}
			],
			"name": "Transfer",
			"type": "event"
		}
	]`))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
	}

	for _, vlog := range logs {
		util.Log("vlog: ", vlog.BlockHash.Hex(), vlog.BlockNumber, vlog.TxHash.Hex())
		if vlog.Address.Hex() == contractAddress {
			// util.Log("contract Data: ", string(vlog.Data))
			var event store.StoreItemSet
			// event := struct {
			// 	Key   [32]byte
			// 	Value [32]byte
			// }{}

			err := contractAbi.UnpackIntoInterface(&event, "ItemSet", vlog.Data)
			if err != nil {
				log.Fatal(err)
			}

			util.Log("contract Data key: ", string(event.Key[:]), "value: ", string(event.Value[:]))
		} else {
			// util.Log("token Data: ", vlog.Data)
			// token固定签名 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
			parsedLog, err := erc20ABI.Unpack("Transfer", vlog.Data)
			if err != nil {
				log.Fatalf("Failed to unpack log: %v", err)
			}

			// 从解析结果中提取参数
			from := common.HexToAddress(vlog.Topics[1].Hex())
			to := common.HexToAddress(vlog.Topics[2].Hex())
			value := parsedLog[0].(*big.Int)
			fbal := new(big.Float)
			fbal.SetString(value.String())
			tokenNum := fbal.Quo(fbal, big.NewFloat(1e18))
			util.Log("Token data:", from.Hex(), "to", to.Hex(), "transfer", tokenNum)
		}

		/* var topics []string
		for i := range vlog.Topics {
			topics = append(topics, vlog.Topics[i].Hex())
		}
		util.Log("topics[0]=", topics[0])
		if len(topics) > 1 {
			util.Log("index topic:", topics[1:])
		} */
	}

	util.Log(crypto.Keccak256Hash([]byte("ItemSet(bytes32,bytes32)")))
	util.Separator()

}

func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	second := 0
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		// util.Log("waitForReceipt second: ", second+1)
		second++
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(1 * time.Second)
	}

}

func TestQueryBlock() {
	client := Connect()
	ctx := context.Background()
	//查询区块
	// block, err := client.HeaderByNumber(ctx, big.NewInt(0))
	block, err := client.BlockByNumber(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	chainId, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	util.Log("Hash: ", block.Hash().Hex())
	util.Log("Number: ", block.Number().Uint64())
	util.Log("Time: ", block.Time())
	util.Log("Difficulty: ", block.Difficulty().Uint64())
	util.Log("Transactions: ", len(block.Transactions()))

	count, err := client.TransactionCount(ctx, block.Hash())
	if err != nil {
		log.Fatal(err)
	}
	util.Log("TransactionCount: ", count)

	//查询交易
	for _, tx := range block.Transactions() {
		util.Log("Hash: ", tx.Hash().Hex())
		util.Log("Value: ", tx.Value().String())
		util.Log("Gas: ", tx.Gas())
		util.Log("GasPrice: ", tx.GasPrice().Uint64())
		util.Log("Nonce: ", tx.Nonce())
		util.Log("Data: ", tx.Data())
		util.Log("To: ", tx.To().Hex())
		sender, err := types.Sender(types.NewEIP155Signer(chainId), tx)
		if err == nil {
			util.Log("Sender: ", sender.Hex())
		} else {
			log.Fatal(err)
		}

		//查询交易回执
		receipt, err := client.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			log.Fatal(err)
		}
		util.Log("receipt Status: ", receipt.Status)
		util.Log("receipt log: ", receipt.Logs)

		txhash := common.HexToHash(tx.Hash().Hex())
		tx, isPending, errr := client.TransactionByHash(ctx, txhash)
		if errr != nil {
			log.Fatal(errr)
		}
		util.Log("isPending: ", isPending)
		util.Log("tx: ", tx.Hash().Hex())
		break
	}

	//查询区块收据
	receipts, err := client.BlockReceipts(ctx, rpc.BlockNumberOrHashWithNumber(rpc.BlockNumber(block.Number().Uint64())))
	if err != nil {
		log.Fatal(err)
	}
	util.Log("receipts count: ", len(receipts))
	for _, receipt := range receipts {
		util.Log("receipts info: ", receipt.Status, receipt.Logs, receipt.TxHash.Hex(), receipt.TransactionIndex, receipt.ContractAddress.Hex(), receipt.BlockNumber)
		break
	}
	util.Separator()
}

func SubscribeNewHead() {
	client, err := ethclient.Dial("wss://" + testUrl)

	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	ctx2, _ := context.WithTimeout(ctx, 5*time.Second)

	headers := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(ctx2, headers)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}
			util.Log("Hash: ", block.Hash().Hex())
			util.Log("Number: ", block.Number().Uint64())
			util.Log("Time: ", block.Time())
			util.Log("Nonce: ", block.Nonce())
			util.Log("Difficulty: ", block.Difficulty().Uint64())
			util.Log("Transactions: ", len(block.Transactions()))
		case <-ctx.Done():
			sub.Unsubscribe()
			return
		}
	}
}

func accountBalance(client *ethclient.Client, address common.Address) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	fbal := new(big.Float)
	fbal.SetString(balance.String())
	util.Log("balance: ", new(big.Float).Quo(fbal, big.NewFloat(1e18)))
}

func Wallet() {
	// genPrivateKey, err := crypto.GenerateKey()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// publicKey := genPrivateKey.Public().(*ecdsa.PublicKey)
	// privateKeyBytes := crypto.FromECDSA(genPrivateKey)
	// publicKeyBytes := crypto.FromECDSAPub(publicKey)
	// util.Log("Private Key: ", hex.EncodeToString(privateKeyBytes))
	// util.Log("Public Key: ", hex.EncodeToString(publicKeyBytes))
	// address := crypto.PubkeyToAddress(*publicKey).Hex()
	// util.Log("Address: ", address)

	client := Connect()
	ctx := context.Background()
	privateKey, nonce := key(client)
	accountBalance(client, crypto.PubkeyToAddress(*privateKey.Public().(*ecdsa.PublicKey)))

	// Create a transaction.
	value := big.NewInt(int64(math.Pow10(14))) // in wei
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("0xc2f2f0955802a506f2A2c650fD48aCD01E8De6E0")
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)
	chainId, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	util.Log("chainId: ", chainId)
	//quentially sign the transaction
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Send the transaction.
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal(err)
	}
	util.Log("tx successfully, SignedTx: ", signedTx.Hash().Hex())
	// 等待交易被挖矿
	receipt, err := waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	util.Log("tx successfully, SignedTx: ", signedTx.Hash().Hex(), ", receipt Status: ", receipt.Status)
	util.Separator()
}

func DigitalAssetToken() {
	client := Connect()
	ctx := context.Background()
	privateKey, nonce := key(client)

	value := big.NewInt(0) // in wei (0 eth)  token不需要转账金额
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}
	gasPrice = big.NewInt(int64(0.001 * math.Pow10(13))) // in wei (1 gwei)
	util.Log("GasPrice: ", gasPrice)

	toAddress := common.HexToAddress("0xc2f2f0955802a506f2A2c650fD48aCD01E8De6E0")
	tokenAddress := common.HexToAddress(tokenAddress) //代币地址 contract address

	transferFnSignature := []byte("transfer(address,uint256)") //ERC20 transfer函数签名
	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferFnSignature)
	methodID := hash.Sum(nil)[:4]
	util.Log(hexutil.Encode(methodID)) // 0xa9059cbb

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32) //地址填充到32字节
	util.Log(hexutil.Encode(paddedAddress))                     // 0x0000000000000000000000004592d8f8d7b001e72cb26a73e4fa1806a51ac79d

	amount := new(big.Int)
	amount.Add(amount, big.NewInt(1e18)) // 1 token
	paddedAmount := common.LeftPadBytes(amount.Bytes(), 32)
	util.Log(hexutil.Encode(paddedAmount)) // 0x00000000000000000000000000000000000000000000003635c9adc5dea00000

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmount...)

	gasLimit, err := client.EstimateGas(ctx, ethereum.CallMsg{
		To:   &toAddress,
		Data: data,
	})
	if err != nil {
		log.Fatal(err)
	}
	gasLimit = uint64(50000) // 21000 is the default gas limit
	util.Log(gasLimit)

	chainID, err := client.NetworkID(ctx)
	if err != nil {
		log.Fatal(err)
	}
	util.Log("chainID", chainID)
	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		log.Fatal(err)
	}
	// 等待交易被挖矿
	receipt, err := waitForReceipt(client, signedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	util.Log("DAT tx sent: ", signedTx.Hash().Hex(), receipt.Status)
	util.Separator()
}

func QueryTokenBalance() {
	client := Connect()
	// Golem (GNT) Address
	tokenAddress := common.HexToAddress(tokenAddress)
	instance, err := token.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	address := common.HexToAddress("0x22F2efAE7f8155120b11973a3315977540F6f58B")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		log.Fatal(err)
	}
	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	util.Log("name: ", name)
	util.Log("symbol: ", symbol)
	util.Log("decimals: ", decimals)
	util.Log("wei: ", bal)
	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))
	util.Log("token balance: ", value)
	util.Separator()
}
