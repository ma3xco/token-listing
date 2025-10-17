package tokenmanager

import (
	"context"
	"regexp"

	"github.com/ma3xco/token-listing/internal/models"
	"github.com/sirupsen/logrus"
)

type tokenManager struct {
	logger logrus.FieldLogger

	// State --------------------------------------------------------------

	// the key is the network id, the value is the network.
	networks map[int64]models.Network

	// the key is the network id, the value is the address regex.
	networkAddressRegex map[int64]*regexp.Regexp

	// list of tokens, the key is the token uid, the value is the token.
	tokens map[string]*models.Token

	// list of featured tokens, the key is the token uid, the value is the token.
	featuredTokens map[string]struct{}

	// the key is the coin marketcap id, the value is the token uid.
	// if the value is empty, then the token is not on CoinMarketCap.
	coinMarketcapIdToTokenUid map[int64]string

	// the key is the network id, the value is the token address map to the token uid.
	networkTokenAddresses map[int64]map[string]string
}

var _ ITokenManager = (*tokenManager)(nil)

func New(ctx context.Context, ops ...Option) (ITokenManager, error) {
	tm := new(tokenManager)
	tm.setDefaults()
	for _, op := range ops {
		if err := op(tm); err != nil {
			return nil, err
		}
	}
	err := tm.init(ctx)
	if err != nil {
		return nil, err
	}
	return tm, nil
}
