package config

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tomochain/backend-matching-engine/swap/ethereum"
)

type EthereumConfig struct {
	NetworkID       string `mapstructure:"network_id"`
	MasterPublicKey string `mapstructure:"master_public_key"`
	// Minimum value of transaction accepted by Bifrost in ETH.
	// Everything below will be ignored.
	MinimumValueEth string `mapstructure:"minimum_value_eth"`
	// TokenPrice is a price of one token in ETH
	TokenPrice string `mapstructure:"token_price"`
	// Host only
	RpcServer string `mapstructure:"rpc_server"`

	// Block number to confirm
	ConfirmedBlockNumber uint64 `mapstructure:"confirmed_block_number"`
}

type TomochainConfig struct {
	// TokenAssetCode is asset code of token that will be purchased using ETH.
	TokenAssetCode string `mapstructure:"token_asset_code"`
	// IssuerPublicKey is public key of the assets issuer.
	IssuerPublicKey string `mapstructure:"issuer_public_key"`
	// DistributionPublicKey is public key of the distribution account.
	// Distribution account can be the same account as issuer account however it's recommended
	// to use a separate account.
	// Distribution account is also used to fund new accounts, this is via smart contract.
	DistributionPublicKey string `mapstructure:"distribution_public_key"`
	// SignerPrivateKey is:
	// * Distribution's secret key if only one instance of Bifrost is deployed.
	// Signer's sequence number will be consumed in transaction's sequence number.
	SignerPrivateKey string `mapstructure:"signer_private_key"`
	// StartingBalance is the starting amount of TOKEN for newly created accounts.
	// Default value is 41. Increase it if you need Data records / other custom entities on new account.
	StartingBalance string `mapstructure:"starting_balance"`
	// LockUnixTimestamp defines unix timestamp when user account will be unlocked.
	LockUnixTimestamp uint64 `mapstructure:"lock_unix_timestamp"`

	signerPrivateKey *ecdsa.PrivateKey
}

type Config struct {
	Ethereum  *EthereumConfig  `mapstructure:"ethereum"`
	Tomochain *TomochainConfig `mapstructure:"tomochain"`
}

func (c *Config) SignerPrivateKey() *ecdsa.PrivateKey {
	if c.Tomochain.signerPrivateKey == nil {
		// from private key to sign smart contract
		// may contain 0x must use FromHex instead of HexString to bytes directly
		keyBytes := common.FromHex(c.Tomochain.SignerPrivateKey)
		c.Tomochain.signerPrivateKey, _ = crypto.ToECDSA(keyBytes)

		// fmt.Printf("address key:%s, err: %v", privkey, err)

	}

	return c.Tomochain.signerPrivateKey
}

func (c *Config) SignerPublicKey() common.Address {
	if c.Tomochain == nil {
		return ethereum.EmptyAddress
	}
	privkey := c.SignerPrivateKey()
	if privkey == nil {
		return ethereum.EmptyAddress
	}

	address := crypto.PubkeyToAddress(privkey.PublicKey)
	return address
}
