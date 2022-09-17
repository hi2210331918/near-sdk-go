package rpc

import (
	"encoding/hex"
	"fmt"
	"github.com/MixinNetwork/near-sdk-go/account"
	"github.com/btcsuite/btcutil/base58"
	"strings"
)

func (c *Client) GetNonce(address, publicKey, finality string) (int64, error) {
	var err error
	var pkStr string
	if !strings.HasPrefix(publicKey, "ed25519:") {
		var pk []byte
		if len(publicKey) == 64 {
			pk, err = hex.DecodeString(publicKey)
			if err != nil {
				return -1, err
			}
			pkStr = account.PublicKeyToString(pk)
		} else {
			pk = base58.Decode(publicKey)
			if len(pk) == 0 {
				return -1, fmt.Errorf("b58 decode public key error, %s", publicKey)
			}
			pkStr = "ed25519:" + publicKey
		}
	} else {
		pkStr = publicKey
	}
	params := make(map[string]interface{})
	params["request_type"] = "view_access_key"
	params["account_id"] = address
	params["public_key"] = pkStr
	params["finality"] = finality
	if finality == "" {
		params["finality"] = "optimistic"
	}
	var res struct {
		GeneralResponse
		Result ViewAccessKey `json:"result"`
	}
	err = c.request("query", params, &res)
	if err != nil {
		return -1, err
	}
	nonce:=res.Result.Nonce
	if nonce <= 0 {
		return -1, fmt.Errorf("res nonce is null,resp=%v", res)
	}
	return nonce + 1, nil
}



func (c *Client) GetTx(transactionHash, senderId string) (*TransactionResult, error) {
	var res struct {
		GeneralResponse
		Result TransactionResult `json:"result"`
	}
	err := c.request("tx", []string{transactionHash, senderId}, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &res.Result, nil
}
func (c *Client) BroadcastTxCommit(raw string) (string, error) {
	var res struct {
		GeneralResponse
		Result NearTxStatus `json:"result"`
	}
	err := c.request("broadcast_tx_commit", []string{raw}, &res)
	if err != nil {
		return "", err
	}

	if res.Error != nil {
		//return "", fmt.Errorf("res=%v", res.Error.Data.TxExecutionError.InvalidTxError)
		return "", res.Error
	}
	//只关心txid，不关心是否发送成功
	txid := res.Result.Transaction.Txid
	return txid, nil
}

func (c *Client) BroadcastTxAsync(raw string) (string, error) {
	var res struct {
		GeneralResponse
		Result string `json:"result"`
	}
	err := c.request("broadcast_tx_async", []string{raw}, &res)
	if err != nil {
		return "", err
	}
	if res.Error != nil {
		return "", res.Error
	}
	return res.Result, nil
}

type ViewAccessKey struct {
	BlockHash  string `json:"block_hash"`
	BlockHeight int64 `json:"block_height"`
	Nonce      int64  `json:"nonce"`
	Permission string  `json:"permission"`
}

type NearTxStatus struct {
	Status             interface{}            `json:"status"`
	Transaction        NearTransaction        `json:"transaction"`
	TransactionOutcome NearTransactionOutcome `json:"transaction_outcome"`
}
type NearTransaction struct {
	SignerId   string        `json:"signer_id"`
	PublicKey  string        `json:"public_key"`
	Nonce      int           `json:"nonce"`
	ReceiverId string        `json:"receiver_id"`
	Actions    []interface{} `json:"actions"`
	Signature  string        `json:"signature"`
	Txid       string        `json:"hash"`
}
type NearTransactionOutcome struct {
	Proof     []interface{} `json:"proof"`
	BlockHash string        `json:"block_hash"`
	Id        string        `json:"id"`
	Outcome   NearOutcome   `json:"outcome"`
}
type NearOutcome struct {
	ReceiptIds  []string          `json:"receipt_ids"`
	GasBurnt    int64             `json:"gas_burnt"`
	TokensBurnt string            `json:"tokens_burnt"`
	ExecutorId  string            `json:"executor_id"`
	Status      map[string]string `json:"status"` // "status": {"SuccessReceiptId": "4xjmyr15T4UqRjdd5YERpfS8QRdCWTH392sMKdQWaJzM"}
}

type TransactionResult struct {
	ReceiptsOutcome []struct {
		BlockHash string `json:"block_hash"`
		ID        string `json:"id"`
		Outcome   struct {
			ExecutorID string        `json:"executor_id"`
			GasBurnt   int64         `json:"gas_burnt"`
			Logs       []interface{} `json:"logs"`
			ReceiptIds []string      `json:"receipt_ids"`
			Status     struct {
				SuccessValue string `json:"SuccessValue"`
			} `json:"status"`
			TokensBurnt string `json:"tokens_burnt"`
		} `json:"outcome"`
		Proof []struct {
			Direction string `json:"direction"`
			Hash      string `json:"hash"`
		} `json:"proof"`
	} `json:"receipts_outcome"`
	Status struct {
		SuccessValue     string `json:"SuccessValue,omitempty"`
		SuccessReceiptId string `json:"SuccessReceiptId,omitempty"`
		Failure          string `json:"Failure,omitempty"`
		Unknown          string `json:"Unknown,omitempty"`
	} `json:"status"`
	Transaction struct {
		Actions []struct {
			Transfer struct {
				Deposit string `json:"deposit"`
			} `json:"Transfer"`
		} `json:"actions"`
		Hash       string `json:"hash"`
		Nonce      int    `json:"nonce"`
		PublicKey  string `json:"public_key"`
		ReceiverID string `json:"receiver_id"`
		Signature  string `json:"signature"`
		SignerID   string `json:"signer_id"`
	} `json:"transaction"`
	TransactionOutcome struct {
		BlockHash string `json:"block_hash"`
		ID        string `json:"id"`
		Outcome   struct {
			ExecutorID string        `json:"executor_id"`
			GasBurnt   int64         `json:"gas_burnt"`
			Logs       []interface{} `json:"logs"`
			ReceiptIds []string      `json:"receipt_ids"`
			Status     struct {
				SuccessReceiptID string `json:"SuccessReceiptId"`
			} `json:"status"`
			TokensBurnt string `json:"tokens_burnt"`
		} `json:"outcome"`
		Proof []struct {
			Direction string `json:"direction"`
			Hash      string `json:"hash"`
		} `json:"proof"`
	} `json:"transaction_outcome"`
}
