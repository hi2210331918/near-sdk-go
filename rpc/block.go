package rpc

import (
	"errors"
)

func (c *Client) GetLatestBlockHash() (int64,string, error) {
	var res struct {
		GeneralResponse
		Result BlockResult `json:"result"`
	}
	err := c.request("block", map[string]interface{}{"finality": "final"}, &res)
	if err != nil {
		return 0,"", err
	}
	if res.Error != nil {
		return 0,"", res.Error
	}
	return res.Result.Header.Height,res.Result.Header.Hash, nil
}

func (c *Client) GetBlockByNumber(id int64) (*BlockResult, error) {
	var res struct {
		GeneralResponse
		Result BlockResult `json:"result"`
	}
	err := c.request("block", map[string]interface{}{"block_id": id}, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &res.Result, nil
}

func (c *Client) GetChainStatus() (*ChainStatus, error) {
	var res struct {
		GeneralResponse
		Result ChainStatus `json:"result"`
	}
	err := c.request("status", []interface{}{},&res)
	if err != nil {
		return nil, err
	}
	if res.Result.SyncInfo.LatestBlockHash == "" {
		return nil, errors.New("latest block hash is null")
	}
	return &res.Result, nil
}

func (c *Client) GetBlock(id string) (*BlockResult, error) {
	var res struct {
		GeneralResponse
		Result BlockResult `json:"result"`
	}
	err := c.request("block", map[string]interface{}{"block_id": id}, &res)
	if err != nil {
		return nil, err
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &res.Result, nil
}

type BlockResult struct {
	Author string `json:"author"`
	Chunks []struct {
		BalanceBurnt         string        `json:"balance_burnt"`
		ChunkHash            string        `json:"chunk_hash"`
		EncodedLength        int           `json:"encoded_length"`
		EncodedMerkleRoot    string        `json:"encoded_merkle_root"`
		GasLimit             int64         `json:"gas_limit"`
		GasUsed              int           `json:"gas_used"`
		HeightCreated        int           `json:"height_created"`
		HeightIncluded       int           `json:"height_included"`
		OutcomeRoot          string        `json:"outcome_root"`
		OutgoingReceiptsRoot string        `json:"outgoing_receipts_root"`
		PrevBlockHash        string        `json:"prev_block_hash"`
		PrevStateRoot        string        `json:"prev_state_root"`
		RentPaid             string        `json:"rent_paid"`
		ShardID              int           `json:"shard_id"`
		Signature            string        `json:"signature"`
		TxRoot               string        `json:"tx_root"`
		ValidatorProposals   []interface{} `json:"validator_proposals"`
		ValidatorReward      string        `json:"validator_reward"`
	} `json:"chunks"`
	Header struct {
		Approvals             []interface{} `json:"approvals"`
		BlockMerkleRoot       string        `json:"block_merkle_root"`
		ChallengesResult      []interface{} `json:"challenges_result"`
		ChallengesRoot        string        `json:"challenges_root"`
		ChunkHeadersRoot      string        `json:"chunk_headers_root"`
		ChunkMask             []bool        `json:"chunk_mask"`
		ChunkReceiptsRoot     string        `json:"chunk_receipts_root"`
		ChunkTxRoot           string        `json:"chunk_tx_root"`
		ChunksIncluded        int64         `json:"chunks_included"`
		EpochID               string        `json:"epoch_id"`
		GasPrice              string        `json:"gas_price"`
		Hash                  string        `json:"hash"`
		Height                int64         `json:"height"`
		LastDsFinalBlock      string        `json:"last_ds_final_block"`
		LastFinalBlock        string        `json:"last_final_block"`
		LatestProtocolVersion int           `json:"latest_protocol_version"`
		NextBpHash            string        `json:"next_bp_hash"`
		NextEpochID           string        `json:"next_epoch_id"`
		OutcomeRoot           string        `json:"outcome_root"`
		PrevHash              string        `json:"prev_hash"`
		PrevStateRoot         string        `json:"prev_state_root"`
		RandomValue           string        `json:"random_value"`
		RentPaid              string        `json:"rent_paid"`
		Signature             string        `json:"signature"`
		Timestamp             int64         `json:"timestamp"`
		TimestampNanosec      string        `json:"timestamp_nanosec"`
		TotalSupply           string        `json:"total_supply"`
		ValidatorProposals    []interface{} `json:"validator_proposals"`
		ValidatorReward       string        `json:"validator_reward"`
	} `json:"header"`
}

type ChainStatus struct {
	ChainId  string        `json:"chain_id"`
	SyncInfo ChainSyncInfo `json:"sync_info"`
}

type ChainSyncInfo struct {
	LatestBlockHash   string `json:"latest_block_hash"`
	LatestBlockHeight int64  `json:"latest_block_height"`
	LatestBlockTime   string `json:"latest_block_time"`
	Syncing           bool   `json:"syncing"`
}