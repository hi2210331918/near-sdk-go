# near-sdk-go

    官网 https://near.org/
    浏览器 https://explorer.near.org/
    文档 https://docs.near.org/api/rpc/introduction

    钱包
    测试 https://wallet.testnet.near.org/
    主网 https://wallet.mainnet.near.org/

    postmanApi:https://documenter.getpostman.com/view/2821468/2s7YfVZVqi#34ffec9b-1f58-4799-b035-f2bfd081020e

## create address
    priv,pub, err := account.GenerateKey()
	if err != nil {
		log.Println(err.Error())
		return
	}
	pubKey:=account.PublicKeyToString(pub)
	prikey:=account.PublicKeyToString(priv)
	address:=account.PublicKeyToAddress(pub)
	fmt.Println("Priv: ",prikey)
	fmt.Println("pubKey: ",pubKey)
	fmt.Println("Address: ",address)
	return

## query
    c := rpc.NewClient(rpc.TestnetRPCEndpoint)
	//获取新的块高和hash
	_, blockhash, _ := c.GetLatestBlockHash()
	//log.Println(height,blockhash)

	//block,err:=c.GetBlock("71FSEbayqXXWVY62r5ff5uJaJ9XqMVieqYUYvoWxFDyF")
	//log.Println(block,err)
	//
	//block,err:=c.GetBlockByNumber(height)
	//log.Printf("%+v err=%s",block,err)

	//gasPrice,err:=c.GetGasPrice("71FSEbayqXXWVY62r5ff5uJaJ9XqMVieqYUYvoWxFDyF")
	//log.Println(gasPrice,err)

	//chian,err:=c.NodeStatus()
	//log.Println(chian.ChainId,chian.Version,err)

	//chunk,err:=c.GetChunk("B5pYQGPz5WcReWCKZzV6FCnLS532WaTd7PhHAVMGNdBf")
	//log.Printf("%+v err=%s",chunk,err)

	//txinfo,err:=c.GetTx("6XMLFGqs5wppyHaoWPUg8t4UoxMBsxSXfcdRgJQRE5eU","hiunique.testnet")
	//log.Printf("%+v err=%s",txinfo,err)

	//c.BroadcastTxAsync()

	//查询账户余额
	//balance,_,err:=c.GetAccountBalance("hiunique.testnet")
	//fmt.Println("balance: ",balance,err)

	//account_id uniquetestnet.testnet
	//ed25519:4nAnFoeQip2dUsMKRCM1cWDwgsVFxcgq7vzQ4GDft38PRdmfor7Fo3JT3tVvv7KpABQ6NbbiMv1u2pm15jcYYNem

	//account_id heunique.testnet
	//ed25519:64Np1cdt3YJFoUHJh5sQYQH4cJ1ESvhEvYmaEauE9aJ1SJuoqpL9sBtmuwk7U4QBunxgu3JP51TmkbVQ7z2NbzxW

	//from:="f9466811157276085dc1be362a28fdaab928fbb1c71dcb98218efc624265eb77"
	//to:="7df5355dd3fb405b5db800f6a1f21fb58c03003cfc1e28b25ff16aadc1a96b22"

## transfer
    from := "heunique.testnet"
	to := "f9466811157276085dc1be362a28fdaab928fbb1c71dcb98218efc624265eb77"
	privateKeyStr := "64Np1cdt3YJFoUHJh5sQYQH4cJ1ESvhEvYmaEauE9aJ1SJuoqpL9sBtmuwk7U4QBunxgu3JP51TmkbVQ7z2NbzxW"

	private := base58.Decode(privateKeyStr)
	accountId := hex.EncodeToString(private[32:])
	fmt.Println("from: ", from, len(from))
	fmt.Println("to: ", to, len(to))
	fmt.Println("accountId: ", accountId)
	//log.Println("from:",from,len(from),len("f9466811157276085dc1be362a28fdaab928fbb1c71dcb98218efc624265eb77"))
	pub, _ := hex.DecodeString(accountId)
	pubKey := account.PublicKeyToString(pub)
	fmt.Println("PublicKey: ", pubKey)

	privateKey := hex.EncodeToString(private[:32])
	fmt.Println("获取处理后的私钥:", privateKey)

	nonce, err := c.GetNonce(from, pubKey, "")
	if err != nil {
		log.Println("获取nonce失败：", err.Error())
		return
	}
	fmt.Println("Nonce: ", nonce)
	tx, err := transaction.CreateTransaction(
		from,
		to,
		pubKey,
		blockhash,
		nonce,
	)

	amount := decimal.NewFromFloat(0.01).Shift(24)
	fmt.Println("Amount:", amount.String())
	//
	ta, err := serialize.CreateTransfer(amount.String())
	if err != nil {
		log.Println("获取nonce失败：", err.Error())
		return
	}

	tx.SetAction(ta)
	txData, err := tx.Serialize()
	if err != nil {
		log.Println("tx序列化失败：", err.Error())
		return
	}
	tx_hex := hex.EncodeToString(txData)

	sig, err := transaction.SignTransaction(tx_hex, privateKey)
	if err != nil {
		log.Println("tx签名失败：", err.Error())
		return
	}
	fmt.Println("Sig: ", sig)
	stx, err := transaction.CreateSignatureTransaction(tx, sig)
	if err != nil {
		log.Println("tx创建签名交易失败：", err.Error())
		return
	}
	stxData, err := stx.Serialize()
	if err != nil {
		log.Println("stx序列化失败：", err.Error())
		return
	}
	b64Data := base64.StdEncoding.EncodeToString(stxData)

	txid, err := c.BroadcastTxCommit(b64Data)
	if err != nil {
		log.Println("广播交易失败：", err.Error())
		return
	}
	fmt.Println("Txid: ", txid)
