package rpc

import (
	"errors"
)

func (c *Client) GetAccountBalance(address string) (string,string, error) {
	params := make(map[string]interface{})
	params["request_type"] = "view_account"
	params["account_id"] = address
	params["finality"] = "final"
	var res struct {
		GeneralResponse
		Result AcccountResult `json:"result"`
	}
	err := c.request("query",params,&res)
	if err != nil {
		return "","", err
	}
	if res.Error != nil {
		return "","", res.Error
	}
	if res.Result.Amount == "" {
		return "","",errors.New("amount is null")
	}
	if res.Result.Locked == ""  {
		return "","",errors.New("locked amount is null")
	}
	return res.Result.Amount,res.Result.Locked, nil
}

type AcccountResult struct {
	Amount string `json:"amount"`
	BlockHash string `json:"block_hash"`
	BlockHeight int64 `json:"block_height"`
	Locked string `json:"locked"`
	StoragePaidAt int64 `json:"storage_paid_at"`
	StorageUsage int64 `json:"storage_usage"`
}
