package controllers

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/QuarkChain/goquarkchain/account"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ybbus/jsonrpc"
	"io/ioutil"
	"math/big"
	"os"
	"sync"
)

type LocalConfig struct {
	Private   string `json:"PrivateKey"`
	Host      string `json:"Host"`
	NetWorkID uint32 `json:"NetWorkID"`
}

func LoadConfig() *LocalConfig {
	conf := new(LocalConfig)
	f, err := os.Open("./localConfig.json")
	if err != nil {
		panic(err)
	}
	buffer, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buffer, conf)
	return conf
}

var (
	localConfig  = LoadConfig()
	sdk          = NewQKCSDK()
	gasLimit     = uint64(6000000)
	gasPrice     = uint64(1000000000)
	MaxPostLen   = (int(gasLimit) - 21000) / 68
	emptyAddress = common.Address{}
	token        = TokenIDEncode("QKC")
	fullShardID  = 262145
)

type QKCSDK struct {
	jrpcHost    jsonrpc.RPCClient
	signAccount account.Account
	mu          sync.Mutex

	nonce uint64
}

func NewQKCSDK() *QKCSDK {
	acc, err := account.NewAccountWithKey(account.BytesToIdentityKey(common.FromHex(localConfig.Private)))
	acc.QKCAddress = acc.QKCAddress.AddressInBranch(account.Branch{Value: uint32(fullShardID)})
	if err != nil {
		panic(err)
	}

	q := &QKCSDK{
		jrpcHost:    jsonrpc.NewClient(localConfig.Host),
		signAccount: acc,
	}
	q.resetNonce()
	return q
}

func (q *QKCSDK) resetNonce() {
	var err error
	tryCnt := 0
	for tryCnt <= 10 {
		tryCnt++
		q.nonce, err = q.GetNonceFromJRPC()
		if err == nil {
			return
		}
	}

}

func (q *QKCSDK) SendFormData(nonce uint64, payLoad []byte) (string, error) {
	tx := newEvmTransaction(nonce, &emptyAddress, new(big.Int), gasLimit, gasPrice, uint32(fullShardID), uint32(fullShardID), token, token, localConfig.NetWorkID, 0, payLoad)
	prvKey, err := crypto.ToECDSA(common.FromHex(q.signAccount.PrivateKey()))
	if err != nil {
		return "", err
	}

	tx, err = SignTx(tx, prvKey)
	if err != nil {
		return "", err
	}
	hash, err := q.SendTransaction(tx)
	if err != nil {
		return "", err
	}
	return hash, err
}

func (q *QKCSDK) SendTransaction(tx *EvmTransaction) (string, error) {
	data, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return "", err
	}

	resp, err := q.jrpcHost.Call("sendRawTransaction", common.ToHex(data))
	if err != nil {
		return "", err
	}
	if resp.Error != nil {
		return "", resp.Error
	}
	return resp.Result.(string), nil
}

func (q *QKCSDK) GetNonceFromJRPC() (nonce uint64, err error) {
	shrd, err := q.GetAccountData(true)
	if err != nil {
		return 0, err
	}
	return hexutil.DecodeUint64(shrd["transactionCount"].(string))
}

func (q *QKCSDK) GetNonce() uint64 {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.nonce++
	return q.nonce - 1
}
func GetFullShardIdByFullShardKey(fullShardKey uint32) uint32 {
	chainID := fullShardKey >> 16
	shardsize := uint32(1)
	shardID := fullShardKey & (shardsize - 1)
	return (chainID << 16) | shardsize | shardID
}

func (q *QKCSDK) GetAccountData(includeShards bool) (map[string]interface{}, error) {
	resp, err := q.jrpcHost.Call("getAccountData", q.signAccount.QKCAddress.ToHex(), nil, includeShards)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	fullShardId := GetFullShardIdByFullShardKey(q.signAccount.QKCAddress.FullShardKey)
	shards := resp.Result.(map[string]interface{})["shards"]
	for _, val := range shards.([]interface{}) {
		shrd := val.(map[string]interface{})
		id, err := hexutil.DecodeUint64(shrd["fullShardId"].(string))
		if err != nil {
			return nil, err
		}
		if id == uint64(fullShardId) {
			return shrd, nil
		}
	}
	return nil, errors.New("has no such account")
}

func (q *QKCSDK) GetTransactionById(txid string) (result []byte, err error) {
	resp, err := q.jrpcHost.Call("getTransactionById", []string{txid})
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	if _, ok := resp.Result.(map[string]interface{}); !ok {
		return nil, fmt.Errorf("txid %v not in the chain", txid)
	}
	return hex.DecodeString(resp.Result.(map[string]interface{})["data"].(string)[2:])
}
