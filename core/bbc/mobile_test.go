package bbc

import (
	"math/rand"
	"testing"

	"github.com/dabankio/wallet-core/bip39"
	"github.com/dabankio/wallet-core/bip44"
	"github.com/stretchr/testify/require"
)

// 兼容性测试，确保代码改动不会破坏已有的数据, 该测试的测试数据不要改动
func TestDeriveAddressCompatible(t *testing.T) {
	for _, tt := range []struct {
		name                            string
		mnemonic, salt, symbol, address string
		// Apply                     func(*Wallet)
	}{
		{
			name:     "PMine",
			mnemonic: "connect auto goose panda extend ozone absent climb abstract doll west crazy",
			salt:     "bbc_keys",
			symbol:   "BBC",
			address:  "18y6y4mdkt2c3q0cypsevs8xt2q08b5ng8j4ztgbdhfpfthpgtp9484re",
		},
		{
			name:     "ColdLegacy",
			mnemonic: "connect auto goose panda extend ozone absent climb abstract doll west crazy",
			salt:     bip44.Password,
			symbol:   "BBC",
			address:  "18wq4j7gewb4zkg0h51jk72f9rhp995g0baxxcky7nbz2zyevjjdypp2g",
		},
		{
			name:     "MKFLegacy",
			mnemonic: "lecture leg select like delay limit spread retire toward west grape bachelor",
			salt:     bip44.Password,
			symbol:   "MKF",
			address:  "1vx6bd4d0jvhte4qndwgcf0hdc4cstmz3zqg8eh2bfsrarewv65xezpdz",
		},
		// TODO BBC 使用新的ID, MKF 使用和BBC一样的ID
	} {
		t.Run(tt.name, func(t *testing.T) {
			seed := bip39.NewSeed(tt.mnemonic, tt.salt)
			d, err := NewSymbolCoin(tt.symbol, bip44.PathFormat, "", seed)
			require.NoError(t, err)
			add, err := d.DeriveAddress()
			require.NoError(t, err)
			require.Equal(t, tt.address, add)
		})
	}
}

func TestNewBip44Deriver(t *testing.T) {
	entropy, err := bip39.NewEntropy(128)
	require.NoError(t, err)

	err = bip39.SetWordListLang(bip39.LangChineseSimplified)
	require.NoError(t, err)
	mnemonic, err := bip39.NewMnemonic(entropy)
	require.NoError(t, err)
	t.Log("mnemonic:", mnemonic)

	seed := bip39.NewSeed(mnemonic, "")

	d, err := NewSymbolCoin("BBC", bip44.PathFormat, "", seed)
	require.NoError(t, err)

	t.Log(d.DeriveAddress())
	t.Log(d.DerivePrivateKey())
	t.Log(d.DerivePublicKey())
}

// TestDeriveConsistent 该测试确保api的稳定性，代码改动过程中确保同样的助记词始终推导出一样的私钥/地址
func TestDeriveConsistent(t *testing.T) {
	mnemonic := "注 笼 伍 叹 纶 林 尸 售 招 愤 勒 熙"
	r := require.New(t)
	r.NoError(bip39.SetWordListLang(bip39.LangChineseSimplified))

	d, err := NewSymbolCoin("BBC", bip44.PathFormat, "", bip39.NewSeed(mnemonic, ""))
	require.NoError(t, err)
	r.NoError(err)
	add, err := d.DeriveAddress()
	require.NoError(t, err)
	r.Equal("11qy08xpjwhv1y012n7c3zv74ww7vy4hnrgz3esv1hzaq115xvdfttng6", add)
}

func TestDecodeTX(t *testing.T) {
	raw := "01000000f345785e00000000701af4705c5e6fcb04efc3ca3c851c1e4d8948e10923025f54bea9b000000000026c2ffa7c6fce7b535aa06b436b7d239c18ec033bb886f689e0a0094beef0775e005a5e2804636414cacc577351e542ff4bb81afa23e45317d298d401fcf345785e010174bc27dc9bfdced95b9b01be398ddd1820350115024fcdb4afc23c3d36bd83bb9c64cd1d00000000640000000000000000816578f1ccb4309f9238db2538b8727952e917cbd3b9ee4dc54cbb8876e072a1e801d5748bcbd807c3c18c0120e88e1e592b339eff523b0fbd353182fe65a3a05ede4bac3d4622e8478ec542aabed3223b9862965289b1d35279ebb2e5b754c21cbc7d8fa7f5c23e4d246065cf12a5c4e29aa2be6b37c70cf8f0927536faa75ac303"
	de, err := DecodeSymbolTX(SymbolBBC, raw)
	require.NoError(t, err)
	t.Log("de tx", de)

	/**
	{
		"transaction" : {
			"txid" : "5e7845f31913abc371132db237f5b44379b9e718876a626230e5223d9e4780ab",
			"version" : 1,
			"type" : "token",
			"time" : 1584940531,
			"lockuntil" : 0,
			"anchor" : "00000000b0a9be545f022309e148894d1e1c853ccac3ef04cb6f5e5c70f41a70",
			"vin" : [
				{
					"txid" : "5e77f0ee4b09a0e089f686b83b03ec189c237d6b436ba05a537bce6f7cfa2f6c",
					"vout" : 0
				},
				{
					"txid" : "5e7845f3fc01d498d21753e423fa1ab84bff42e5517357ccca14646304285e5a",
					"vout" : 1
				}
			],
			"sendfrom" : "20g03dfrhttamxxs3ca4fx7f7h1336h9hw9rnza1nb0e2666aq2e9sm9m",
			"sendto" : "1ejy2fq4vzq7djpwv06z3k3ex30g3a08n097wvd5fr8y3tdnxgexybjvm",
			"amount" : 499.999900,
			"txfee" : 0.000100,
			"data" : "",
			"sig" : "6578f1ccb4309f9238db2538b8727952e917cbd3b9ee4dc54cbb8876e072a1e801d5748bcbd807c3c18c0120e88e1e592b339eff523b0fbd353182fe65a3a05ede4bac3d4622e8478ec542aabed3223b9862965289b1d35279ebb2e5b754c21cbc7d8fa7f5c23e4d246065cf12a5c4e29aa2be6b37c70cf8f0927536faa75ac303",
			"fork" : "00000000b0a9be545f022309e148894d1e1c853ccac3ef04cb6f5e5c70f41a70",
			"confirmations" : 43
		}
	}
	*/
}

func TestParsePrivateKey(t *testing.T) {
	priv, pub, add := "be692b83d565a862933906605eb3ff2816bfbcb1ca51c8066bedbdb677228def", "d27e4e9b76041045e31fcc656def79d36df3a0d5ac2a1bab93649a7b98f6fb08", "113xzd63vk9j97arv5apdb87kdq9qkvvdcq61zrt52027d6tefv9e16md"

	info, err := ParsePrivateKey(priv)
	require.NoError(t, err)
	require.Equal(t, priv, info.PrivateKey)
	require.Equal(t, pub, info.PublicKey)
	require.Equal(t, add, info.Address)
}

func TestNewSymbolBip44Deriver(t *testing.T) {
	b := make([]byte, 64)
	rand.Read(b)

	w, err := NewSymbolBip44Deriver("BBC", bip44.FullPathFormat, "BBC", b)
	require.NoError(t, err)

	t.Log(w.DeriveAddress())
	t.Log(w.DerivePrivateKey())
	t.Log(w.DerivePublicKey())
}
