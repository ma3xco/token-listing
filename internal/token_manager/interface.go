package tokenmanager

import (
	"context"
)

// ITokenManager is the interface for the token manager.
// it is used to manage the tokens in the memory.
// it is used to build the tokens into the assets.
// it is used to validate the tokens.
// it is used to create the token template.
type ITokenManager interface {
	// WalkThrough walks through the token list and load the tokens into the memory.
	// it returns the number of tokens loaded and an error if any.
	WalkThrough(ctx context.Context) (int, error)

	// CreateTokenTemplate creates a token template for the given token uid.
	// it returns an error if any.
	CreateTokenTemplate(ctx context.Context, uid string) error

	// ValidateTokens validates the tokens in the memory.
	// it returns an error if any.
	// the map key is the token uid, the value is the errors.
	ValidateTokens(ctx context.Context) map[string][]error

	// BuildTokens builds the tokens in the memory.into ./dist/***
	// the build assets contains
	// - tokens.json (all tokens list)
	// - :network_id/:tokenAddress.json (the token details with all token addresses Hashmap)
	// - :network_id/:tokenAddress/token_address.json (the token address details only)
	// - tokens/:tokenUid.json (the token Hashmap)
	// - :coin_marketcap_id.json (the coin marketcap Hashmap)
	// - tokens.featured.json (the featured tokens list)
	// it returns an error if any.
	BuildTokens(ctx context.Context) error
}
