/*
 * Copyright 2018 The openwallet Authors
 * This file is part of the openwallet library.
 *
 * The openwallet library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The openwallet library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 */

package avalanche

import (
	"fmt"

	"github.com/blocktree/go-owcdrivers/addressEncoder/bech32"
	"github.com/blocktree/openwallet/v2/openwallet"

	"github.com/blocktree/go-owcdrivers/addressEncoder"
	"github.com/blocktree/go-owcrypt"
)

func init() {

}

var(
	AVAX_mainnetAddressP2PKH = addressEncoder.AddressType{
		"bech32",
		addressEncoder.BTCBech32Alphabet,
		"sha3_256_ripemd160",
		"sha3_256",
		20,
		[]byte("avax"),
		nil}

	AVAX_testnetAddressP2PKH = addressEncoder.AddressType{
		"bech32",
		addressEncoder.BTCBech32Alphabet,
		"sha3_256_ripemd160",
		"sha3_256",
		20,
		[]byte("local"),
		nil}

	//AVAX_testnetAddressP2PKH = addressEncoder.AddressType{
	//	"bech32",
	//	addressEncoder.BTCBech32Alphabet,
	//	"sha3_256_ripemd160",
	//	"sha3_256",
	//	20,
	//	[]byte("fuji"),
	//	nil}
)

type AddressDecoder interface {
	openwallet.AddressDecoder
	ScriptPubKeyToBech32Address(scriptPubKey []byte) (string, error)
}

type addressDecoder struct {
	wm *WalletManager //钱包管理者
}

//NewAddressDecoder 地址解析器
func NewAddressDecoder(wm *WalletManager) *addressDecoder {
	decoder := addressDecoder{}
	decoder.wm = wm
	return &decoder
}

//PrivateKeyToWIF 私钥转WIF
func (decoder *addressDecoder) PrivateKeyToWIF(priv []byte, isTestnet bool) (string, error) {

	cfg := addressEncoder.BTC_mainnetPrivateWIFCompressed
	if decoder.wm.Config.IsTestNet {
		cfg = addressEncoder.BTC_testnetPrivateWIFCompressed
	}

	wif := addressEncoder.AddressEncode(priv, cfg)

	return wif, nil

}

//PublicKeyToAddress 公钥转地址
func (decoder *addressDecoder) PublicKeyToAddress(pub []byte, isTestnet bool) (string, error) {

	cfg := AVAX_mainnetAddressP2PKH
	if decoder.wm.Config.IsTestNet {
		cfg = AVAX_testnetAddressP2PKH
	}

	hash := owcrypt.Hash(pub, 0, owcrypt.HASH_ALG_SHA3_256)
	rip := owcrypt.Hash(hash, 20, owcrypt.HASH_ALG_RIPEMD160)

	address := bech32.Encode(string(cfg.Prefix), cfg.Alphabet, rip, nil)

	//if decoder.wm.Config.RPCServerType == RPCServerCore {
	//	//如果使用core钱包作为全节点，需要导入地址到core，这样才能查询地址余额和utxo
	//	err := decoder.wm.ImportAddress(address, "")
	//	if err != nil {
	//		return "", err
	//	}
	//}

	// 加上chainId，只实现X链的
	// X chain = X
	// P chain = P
	// C chain = C
	return "X-"+address, nil

}

//RedeemScriptToAddress 多重签名赎回脚本转地址
func (decoder *addressDecoder) RedeemScriptToAddress(pubs [][]byte, required uint64, isTestnet bool) (string, error) {

	cfg := addressEncoder.BTC_mainnetAddressP2SH
	if decoder.wm.Config.IsTestNet {
		cfg = addressEncoder.BTC_testnetAddressP2SH
	}

	redeemScript := make([]byte, 0)

	for _, pub := range pubs {
		redeemScript = append(redeemScript, pub...)
	}

	pkHash := owcrypt.Hash(redeemScript, 0, owcrypt.HASH_ALG_HASH160)

	address := addressEncoder.AddressEncode(pkHash, cfg)

	return address, nil

}

//WIFToPrivateKey WIF转私钥
func (decoder *addressDecoder) WIFToPrivateKey(wif string, isTestnet bool) ([]byte, error) {

	cfg := addressEncoder.BTC_mainnetPrivateWIFCompressed
	if decoder.wm.Config.IsTestNet {
		cfg = addressEncoder.BTC_testnetPrivateWIFCompressed
	}

	priv, err := addressEncoder.AddressDecode(wif, cfg)
	if err != nil {
		return nil, err
	}

	return priv, err

}

//ScriptPubKeyToBech32Address scriptPubKey转Bech32地址
func (decoder *addressDecoder) ScriptPubKeyToBech32Address(scriptPubKey []byte) (string, error) {
	return scriptPubKeyToBech32Address(scriptPubKey, decoder.wm.Config.IsTestNet)

}

//ScriptPubKeyToBech32Address scriptPubKey转Bech32地址
func scriptPubKeyToBech32Address(scriptPubKey []byte, isTestNet bool) (string, error) {
	var (
		hash []byte
	)

	cfg := addressEncoder.BTC_mainnetAddressBech32V0
	if isTestNet {
		cfg = addressEncoder.BTC_testnetAddressBech32V0
	}

	if len(scriptPubKey) == 22 || len(scriptPubKey) == 34 {

		hash = scriptPubKey[2:]

		address := addressEncoder.AddressEncode(hash, cfg)

		return address, nil

	} else {
		return "", fmt.Errorf("scriptPubKey length is invalid")
	}

}
