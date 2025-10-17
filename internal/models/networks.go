package models

// Network is the model for a network.
type Network struct {
	// The ID of the network.
	Id int64 `json:"id"`

	// Chain ID as defined by the network (e.g., 1 for Ethereum mainnet)
	ChainId int64 `json:"chain_id"`

	// The type of the network.
	// refer to NetworkType enum for the possible values.
	NetworkType NetworkType `json:"network_type"`

	// coin type reffering to SLIP-0044
	CoinType Coin_Type `json:"coin_type"`

	// Human-readable name of the network (e.g., "Ethereum Mainnet")
	Name string `json:"name"`

	// Network symbol/abbreviation (e.g., "ETH", "MATIC")
	Symbol string `json:"symbol"`

	// Number of decimal places for the native token
	Decimals int64 `json:"decimals"`

	// PNG format icon URL for the network
	IconPngUrl string `json:"icon_png_url"`

	// SVG format icon URL for the network
	IconSvgUrl string `json:"icon_svg_url"`

	// Whether the network is a testnet.
	IsTestnet bool `json:"is_testnet"`

	// Whether this network is currently active and supported
	IsActive bool `json:"is_active"`

	// Regular expression to validate addresses for the network
	AddressRegex string `json:"address_regex"`

	// Explorer configuration
	Explorer Explorer `json:"explorer"`

	// The ID of the network on CoinMarketCap.
	CoinMarketCapId int64 `json:"coin_marketcap_id"`
}

// Explorer is the model for a network explorer.
type Explorer struct {
	// The base URL of the explorer.
	BaseUrl string `json:"base_url"`

	// The template of the address.
	AddressTemplate string `json:"address_template"`

	// The template of the transaction.
	TransactionTemplate string `json:"transaction_template"`

	// The template of the token.
	TokenTemplate string `json:"token_template"`

	// The template of the block.
	BlockTemplate string `json:"block_template"`
}

// Network type enum
type NetworkType int32

const (
	// Unspecified network type
	NetworkType_NETWORK_TYPE_UNSPECIFIED NetworkType = 0
	// Ethereum-like network type
	NetworkType_NETWORK_TYPE_ETH_LIKE NetworkType = 1
	// TRON network type
	NetworkType_NETWORK_TYPE_TRX NetworkType = 2
	// Solana network type
	NetworkType_NETWORK_TYPE_SOL NetworkType = 3
	// UTXO network type
	NetworkType_NETWORK_TYPE_UTXO NetworkType = 4
)

// Coin type reffering to SLIP-0044
type Coin_Type int32

const (
	// Bitcoin
	// We are intentionally using TYPE_BTC as the zero value for legacy compatibility.
	// buf:lint:ignore ENUM_ZERO_VALUE_SUFFIX
	Coin_TYPE_BTC Coin_Type = 0
	// Ethereum
	Coin_TYPE_ETH Coin_Type = 60
	// BNB
	Coin_TYPE_BNB Coin_Type = 714
	// Solana
	Coin_TYPE_SOL Coin_Type = 501
	// TRON
	Coin_TYPE_TRX Coin_Type = 195
)
