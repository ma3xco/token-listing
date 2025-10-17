package models

// Token is the model for a token.
type Token struct {
	// Unique identifier for the token
	Uuid string `json:"uuid"`

	// Human-readable name of the token (e.g., "Ethereum", "USD Coin")
	Name string `json:"name"`

	// Symbol of the token (e.g., "ETH", "USDC")
	Symbol string `json:"symbol"`

	// Whether the token has gas sponsored.
	HasGasSponsored bool `json:"has_gas_sponsored"`

	// Whether the token is a stable token.
	IsStableToken bool `json:"is_stable_token"`

	// ID of the token that is wrapped by this token (if applicable)
	WrappedTokenUuid string `json:"wrapped_token_uuid"`

	// URL of the logo of the token in the form of png. must be 64x64.
	LogoPngUrl string `json:"logo_png_url"`

	// URL of the logo of the token in the form of svg.
	LogoSvgUrl string `json:"logo_svg_url"`

	// Description of the token.
	Description string `json:"description"`

	// ID of the token on CoinMarketCap. leave -1 if not on CoinMarketCap.
	CoinMarketCapId int64 `json:"coin_market_cap_id"`

	// Whether the token is featured. if the value is true, then the PR will be rejected.
	IsFeatured bool `json:"is_featured"`

	// The order index of the token. the lower the index, the higher the priority.
	// it should be greater than or equal to 100000.
	OrderIndex int64 `json:"order_index"`

	// The website of the token.
	WebsiteUrl string `json:"website_url"`

	// The X (Twitter) url of the token.
	XUrl string `json:"x_url"`

	// The discord url of the token.
	DiscordUrl string `json:"discord_url"`

	// The whitepaper url of the token.
	WhitepaperUrl string `json:"whitepaper_url"`

	// The live price url of the token. if the token is not listed on CoinMarketCap, then Required.
	// Historical price will be fetched from the CoinMarketCap only.
	// the url should return a simple json in the following format:
	// {
	// 	"price_usd": 1.23,
	//  "volume_24h": 100.00,
	// 	"volume_change_24h": 2.34, // 2.34%
	//  "percent_change_1h": 2.34, // 2.34%
	//  "percent_change_24h": 2.34, // 2.34%
	//  "percent_change_7d": 2.34, // 2.34%
	//  "percent_change_30d": 2.34, // 2.34%
	//  "percent_change_90d": 2.34, // 2.34%
	// }
	LivePriceUrl string `json:"live_price_url"`

	// the tags of the token.
	Tags []string `json:"tags"`

	// Whether the token is a scam.
	IsScam bool `json:"is_scam"`

	// Whether the token is disabled.
	// the token might be disabled due to security issues or other reasons.
	// the disabled token will not be displayed in the wallet.
	IsDisabled bool `json:"is_disabled"`

	// the addresses of the token on the networks.
	Addresses []TokenAddress `json:"addresses"`
}
