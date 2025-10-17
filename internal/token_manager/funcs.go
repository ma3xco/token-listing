package tokenmanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/ma3xco/token-listing/internal/models"
	"github.com/sirupsen/logrus"
)

func (tm *tokenManager) setDefaults() {
	logger := logrus.New()
	logger.ReportCaller = true
	tm.logger = logger

	tm.networks = make(map[int64]models.Network)
	tm.networkAddressRegex = make(map[int64]*regexp.Regexp)
	tm.tokens = make(map[string]*models.Token)
	tm.featuredTokens = make(map[string]struct{})
	tm.coinMarketcapIdToTokenUid = make(map[int64]string)
	tm.networkTokenAddresses = make(map[int64]map[string]string)
}

func (tm *tokenManager) init(ctx context.Context) error {
	err := tm.loadNetworks(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (tm *tokenManager) WalkThrough(ctx context.Context) (int, error) {
	entries, err := os.ReadDir("tokens")
	if err != nil {
		return 0, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == "_example" {
			continue
		}
		tknUid := entry.Name()

		metaFile, err := os.ReadFile(fmt.Sprintf("tokens/%s/meta.json", tknUid))
		if err != nil {
			return 0, err
		}
		var token models.Token
		err = json.Unmarshal(metaFile, &token)
		if err != nil {
			return 0, err
		}
		_, err = os.ReadFile(fmt.Sprintf("tokens/%s/logo.png", tknUid))
		if err != nil {
			return 0, err
		}
		// Token exists and logo exists, load the token into the memory.
		tm.tokens[tknUid] = &token
		if token.IsFeatured {
			tm.featuredTokens[tknUid] = struct{}{}
		}
		for _, address := range token.Addresses {
			reg, ok := tm.networkAddressRegex[int64(address.NetworkId)]
			if !ok {
				return 0, errors.New("network address regex not found")
			}
			if !reg.MatchString(address.Address) && address.NetworkId != 1 {
				return 0, fmt.Errorf("network address regex not match for network %d and address %s", address.NetworkId, address.Address)
			}
			if _, ok := tm.networkTokenAddresses[int64(address.NetworkId)][address.Address]; ok {
				return 0, fmt.Errorf("token address already exists for network %d and address %s", address.NetworkId, address.Address)
			}
			if tm.networkTokenAddresses[int64(address.NetworkId)] == nil {
				tm.networkTokenAddresses[int64(address.NetworkId)] = make(map[string]string)
			}
			tm.networkTokenAddresses[int64(address.NetworkId)][address.Address] = tknUid
			tm.coinMarketcapIdToTokenUid[int64(address.NetworkId)] = tknUid
		}

	}

	return len(tm.tokens), nil
}

func (tm *tokenManager) CreateTokenTemplate(ctx context.Context, uid string) error {
	return errors.New("not implemented")
}

// ValidateTokens validates the tokens in the memory.
// Validation rules:
// - The token must have a name.
// - The token must have a symbol.
// - The token must have a logo.
// - The token must have a description.
// - The token must have a coin marketcap id or price url.
// - logo png is not too large and have 64x64 size.
func (tm *tokenManager) ValidateTokens(ctx context.Context) map[string][]error {
	validationErrors := make(map[string][]error)

	for tokenUid, token := range tm.tokens {
		var errors []error

		// Validate token name
		if strings.TrimSpace(token.Name) == "" {
			errors = append(errors, fmt.Errorf("token name is required"))
		}

		// Validate token symbol
		if strings.TrimSpace(token.Symbol) == "" {
			errors = append(errors, fmt.Errorf("token symbol is required"))
		}

		// Validate token description
		if strings.TrimSpace(token.Description) == "" {
			errors = append(errors, fmt.Errorf("token description is required"))
		}

		// Validate logo PNG URL
		if strings.TrimSpace(token.LogoPngUrl) == "" {
			errors = append(errors, fmt.Errorf("logo PNG URL is required"))
		} else {
			// Validate logo URL format
			if _, err := url.Parse(token.LogoPngUrl); err != nil {
				errors = append(errors, fmt.Errorf("invalid logo PNG URL format: %v", err))
			}
		}

		// Validate logo file exists and is 64x64 PNG
		logoPath := fmt.Sprintf("tokens/%s/logo.png", tokenUid)
		if err := tm.validateLogoFile(logoPath); err != nil {
			errors = append(errors, fmt.Errorf("logo file validation failed: %v", err))
		}

		// Validate either CoinMarketCap ID or LivePriceUrl is provided
		if token.CoinMarketCapId == -1 && strings.TrimSpace(token.LivePriceUrl) == "" {
			errors = append(errors, fmt.Errorf("either CoinMarketCap ID or LivePriceUrl must be provided"))
		}

		// Validate CoinMarketCap ID if provided
		if token.CoinMarketCapId != -1 && token.CoinMarketCapId <= 0 {
			errors = append(errors, fmt.Errorf("CoinMarketCap ID must be positive"))
		}

		// Validate LivePriceUrl format if provided
		if strings.TrimSpace(token.LivePriceUrl) != "" {
			if _, err := url.Parse(token.LivePriceUrl); err != nil {
				errors = append(errors, fmt.Errorf("invalid LivePriceUrl format: %v", err))
			}
		}

		// Validate URLs format
		if strings.TrimSpace(token.WebsiteUrl) != "" {
			if _, err := url.Parse(token.WebsiteUrl); err != nil {
				errors = append(errors, fmt.Errorf("invalid website URL format: %v", err))
			}
		}

		if strings.TrimSpace(token.XUrl) != "" {
			if _, err := url.Parse(token.XUrl); err != nil {
				errors = append(errors, fmt.Errorf("invalid X (Twitter) URL format: %v", err))
			}
		}

		if strings.TrimSpace(token.DiscordUrl) != "" {
			if _, err := url.Parse(token.DiscordUrl); err != nil {
				errors = append(errors, fmt.Errorf("invalid Discord URL format: %v", err))
			}
		}

		if strings.TrimSpace(token.WhitepaperUrl) != "" {
			if _, err := url.Parse(token.WhitepaperUrl); err != nil {
				errors = append(errors, fmt.Errorf("invalid whitepaper URL format: %v", err))
			}
		}

		// Validate SVG logo URL format if provided
		if strings.TrimSpace(token.LogoSvgUrl) != "" {
			if _, err := url.Parse(token.LogoSvgUrl); err != nil {
				errors = append(errors, fmt.Errorf("invalid logo SVG URL format: %v", err))
			}
		}

		// Validate token addresses
		if len(token.Addresses) == 0 {
			errors = append(errors, fmt.Errorf("at least one token address is required"))
		}

		for i, address := range token.Addresses {
			addressErrors := tm.validateTokenAddress(address, i)
			errors = append(errors, addressErrors...)
		}

		// Validate wrapped token UUID if provided
		if strings.TrimSpace(token.WrappedTokenUuid) != "" {
			if _, exists := tm.tokens[token.WrappedTokenUuid]; !exists {
				errors = append(errors, fmt.Errorf("wrapped token UUID '%s' does not exist", token.WrappedTokenUuid))
			}
		}

		// Validate scam flag
		if token.IsScam {
			errors = append(errors, fmt.Errorf("scam tokens are not allowed"))
		}

		if len(errors) > 0 {
			validationErrors[tokenUid] = errors
		}
	}

	return validationErrors
}

// validateLogoFile validates that the logo file exists and is a 64x64 PNG
func (tm *tokenManager) validateLogoFile(logoPath string) error {
	// Check if file exists
	fileInfo, err := os.Stat(logoPath)
	if err != nil {
		return fmt.Errorf("logo file does not exist: %v", err)
	}

	// Check file size (should not be too large, e.g., max 1MB)
	const maxFileSize = 1024 * 1024 // 1MB
	if fileInfo.Size() > maxFileSize {
		return fmt.Errorf("logo file is too large: %d bytes (max: %d bytes)", fileInfo.Size(), maxFileSize)
	}

	// Check if file is too small
	const minFileSize = 200 // 200 bytes
	if fileInfo.Size() < minFileSize {
		return fmt.Errorf("logo file is too small: %d bytes (min: %d bytes)", fileInfo.Size(), minFileSize)
	}

	// Open and decode the image to check dimensions
	file, err := os.Open(logoPath)
	if err != nil {
		return fmt.Errorf("cannot open logo file: %v", err)
	}
	defer file.Close()

	// First, decode PNG config to get actual dimensions without loading full image
	config, err := png.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("cannot decode PNG config: %v", err)
	}

	// Check dimensions using PNG config decoder
	width := config.Width
	height := config.Height

	if width != 64 || height != 64 {
		return fmt.Errorf("logo must be exactly 64x64 pixels, got: %dx%d", width, height)
	}

	return nil
}

// validateTokenAddress validates a single token address
func (tm *tokenManager) validateTokenAddress(address models.TokenAddress, index int) []error {
	var errors []error

	// Validate address is not empty
	if strings.TrimSpace(address.Address) == "" {
		errors = append(errors, fmt.Errorf("address[%d]: address is required", index))
	}

	// Validate token UID matches
	if strings.TrimSpace(address.TokenUid) == "" {
		errors = append(errors, fmt.Errorf("address[%d]: token UID is required", index))
	}

	// Validate network ID exists
	network, exists := tm.networks[int64(address.NetworkId)]
	if !exists {
		errors = append(errors, fmt.Errorf("address[%d]: network ID %d does not exist", index, address.NetworkId))
	} else {
		// Validate address format using network regex
		if address.NetworkId != 1 { // Skip validation for Bitcoin (network ID 1) as it has special handling
			regex, ok := tm.networkAddressRegex[int64(address.NetworkId)]
			if ok && !regex.MatchString(address.Address) {
				errors = append(errors, fmt.Errorf("address[%d]: address '%s' does not match network %s regex", index, address.Address, network.Name))
			}
		}
	}

	// Validate decimals
	if address.Decimals > 18 {
		errors = append(errors, fmt.Errorf("address[%d]: decimals cannot exceed 18", index))
	}

	// Validate token type
	validTokenTypes := []string{"ERC20", "ERC721", "ERC1155", "SPL", "SPL2022", "COIN"}
	isValidType := false
	for _, validType := range validTokenTypes {
		if address.TokenType == validType {
			isValidType = true
			break
		}
	}
	if !isValidType {
		errors = append(errors, fmt.Errorf("address[%d]: invalid token type '%s', must be one of: %v", index, address.TokenType, validTokenTypes))
	}

	// Validate gas sponsored strategy
	if address.GasSponsoredStrategy < 0 || address.GasSponsoredStrategy > 5 {
		errors = append(errors, fmt.Errorf("address[%d]: gas sponsored strategy must be between 0 and 5", index))
	}

	// Validate name matches token name
	if strings.TrimSpace(address.Name) == "" {
		errors = append(errors, fmt.Errorf("address[%d]: name is required", index))
	}

	// Validate symbol matches token symbol
	if strings.TrimSpace(address.Symbol) == "" {
		errors = append(errors, fmt.Errorf("address[%d]: symbol is required", index))
	}

	// Validate logo PNG URL format if provided
	if strings.TrimSpace(address.LogoPngUrl) != "" {
		if _, err := url.Parse(address.LogoPngUrl); err != nil {
			errors = append(errors, fmt.Errorf("address[%d]: invalid logo PNG URL format: %v", index, err))
		}
	}

	// Validate logo SVG URL format if provided
	if strings.TrimSpace(address.LogoSvgUrl) != "" {
		if _, err := url.Parse(address.LogoSvgUrl); err != nil {
			errors = append(errors, fmt.Errorf("address[%d]: invalid logo SVG URL format: %v", index, err))
		}
	}

	return errors
}

// BuildTokens builds the tokens in the memory.into ./dist/***
// the build assets contains
// - tokens.json (all tokens list) - done
// - :network_id/:tokenAddress.json (the token details with all token addresses Hashmap) - done
// - :network_id/:tokenAddress/token_address.json (the token address details only) - done
// - tokens/:tokenUid.json (the token Hashmap) - done
// - :coin_marketcap_id.json (the coin marketcap Hashmap) - SKIPPED
// - tokens.featured.json (the featured tokens list) - done
// it returns an error if any.
func (tm *tokenManager) BuildTokens(ctx context.Context) error {
	// clean the dist directory
	err := os.RemoveAll("./dist")
	if err != nil {
		return err
	}
	err = os.Mkdir("./dist", 0755)
	if err != nil {
		return err
	}
	// build tokens.json
	{
		var tokens []models.Token
		for i := range tm.tokens {
			tokens = append(tokens, *tm.tokens[i])
		}
		bytes, err := json.Marshal(tokens)
		if err != nil {
			return err
		}
		err = os.WriteFile("./dist/tokens.json", bytes, 0644)
		if err != nil {
			return err
		}
	}
	// build tokens.featured.json
	{
		var tokens []models.Token
		for i := range tm.featuredTokens {
			tokens = append(tokens, *tm.tokens[i])
		}
		bytes, err := json.Marshal(tokens)
		if err != nil {
			return err
		}
		err = os.WriteFile("./dist/tokens.featured.json", bytes, 0644)
		if err != nil {
			return err
		}
	}

	// build :network_id/:tokenAddress.json
	{
		for networkId, tokenAddresses := range tm.networkTokenAddresses {
			os.Mkdir(fmt.Sprintf("./dist/%d", networkId), 0755)
			for tokenAddress, tokenUid := range tokenAddresses {
				token, ok := tm.tokens[tokenUid]
				if !ok {
					return fmt.Errorf("token %s not found", tokenUid)
				}
				bytes, err := json.Marshal(token)
				if err != nil {
					return err
				}
				err = os.WriteFile(fmt.Sprintf("./dist/%d/%s.json", networkId, tokenAddress), bytes, 0644)
				if err != nil {
					return err
				}

			}
		}

	}
	// build tokens/:tokenUid.json & tokens/:tokenUid.png
	{
		os.Mkdir("./dist/tokens", 0755)
		for tokenUid := range tm.tokens {
			token, ok := tm.tokens[tokenUid]
			if !ok {
				return fmt.Errorf("token %s not found", tokenUid)
			}
			bytes, err := json.Marshal(token)
			if err != nil {
				return err
			}
			err = os.WriteFile(fmt.Sprintf("./dist/tokens/%s.json", tokenUid), bytes, 0644)
			if err != nil {
				return err
			}
			// load the logo.png
			logoPng, err := os.ReadFile(fmt.Sprintf("tokens/%s/logo.png", tokenUid))
			if err != nil {
				return err
			}
			err = os.WriteFile(fmt.Sprintf("./dist/tokens/%s.png", tokenUid), logoPng, 0644)
			if err != nil {
				return err
			}
		}
	}
	// :network_id/:tokenAddress/token_address.json (the token address details only)
	{
		for networkId := range tm.networks {
			err := os.Mkdir(fmt.Sprintf("./dist/%d", networkId), 0755)
			if err != nil && !os.IsExist(err) {
				return err
			}
		}
		for index := range tm.tokens {
			token := tm.tokens[index]
			for _, address := range token.Addresses {
				os.Mkdir(fmt.Sprintf("./dist/%d/%s", address.NetworkId, address.Address), 0755)
				bytes, err := json.Marshal(address)
				if err != nil {
					return err
				}
				err = os.WriteFile(fmt.Sprintf("./dist/%d/%s/token_address.json", address.NetworkId, address.Address), bytes, 0644)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (tm *tokenManager) loadNetworks(ctx context.Context) error {
	bytes, err := os.ReadFile("networks/networks.json")
	if err != nil {
		return err
	}
	var networks []models.Network
	err = json.Unmarshal(bytes, &networks)
	if err != nil {
		return err
	}
	for _, network := range networks {
		tm.networks[int64(network.Id)] = network
		tm.networkAddressRegex[int64(network.Id)] = regexp.MustCompile(network.AddressRegex)
	}
	return nil
}
