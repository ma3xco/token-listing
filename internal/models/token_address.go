package models

// TokenAddress is the model for a token address.
type TokenAddress struct {
	// The address of the token on the network.
	Address string `json:"address"`

	// The UID of the token on the network. must be the same as the token uid.
	TokenUid string `json:"token_uid"`

	// The ID of the network. refer to networks directory for other networks.
	NetworkId int32 `json:"network_id"`

	// Whether the token is verified on the network.
	// the token verification must satisfy the verification criteria.
	IsVerified bool `json:"is_verified"`

	// The number of decimals used to get its user representation.
	Decimals uint32 `json:"decimals"`

	// Whether the token is the native token of the network.
	IsNative bool `json:"is_native"`

	// the type of the token, ERC20, ERC721, ERC1155, SPL, SPL2022,
	// for the Tron, BNB and other ethereum-like networks should use ERC prefix.
	TokenType string `json:"token_type"`

	// Whether the token has a proxy.
	Upgradeable bool `json:"upgradeable"`

	// Whether the token has a blue checkmark.
	HasBlueCheckmark bool `json:"has_blue_checkmark"`

	// The gas sponsored strategy of the token.
	// 0: no gas sponsored, 1: full matrix strategy, 2: Authorized transfer, 3: permit, 4: gas-transfer, 5: co-signer(SOL only)
	GasSponsoredStrategy int32 `json:"gas_sponsored_strategy"`

	// The name of the token. must be the same as the token.
	Name string `json:"name"`

	// The symbol of the token. Shall be the same as the token.
	Symbol string `json:"symbol"`

	// The URL of the logo of the token. PNG is required
	LogoPngUrl string `json:"logo_png_url"`

	// The URL of the logo of the token in the form of svg.
	LogoSvgUrl string `json:"logo_svg_url"`
}
