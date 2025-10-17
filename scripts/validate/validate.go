package main

import (
	"context"
	"fmt"
	"log"
	"os"

	tokenmanager "github.com/ma3xco/token-listing/internal/token_manager"
)

func main() {

	tm, err := tokenmanager.New(context.Background())
	if err != nil {
		log.Fatalf("failed to create token manager: %v", err)
	}
	count, err := tm.WalkThrough(context.Background())
	if err != nil {
		log.Fatalf("failed to walk through tokens: %v", err)
	}
	fmt.Printf("walked through %d tokens\n", count)
	validationErrors := tm.ValidateTokens(context.Background())
	if len(validationErrors) > 0 {
		for tokenUid, errors := range validationErrors {
			fmt.Printf("token %s has errors:\n", tokenUid)
			for _, error := range errors {
				fmt.Printf("  - %s\n", error)
			}
		}
		os.Exit(1)
	} else {
		fmt.Printf("all tokens are valid\n")
		fmt.Println("validation completed")
	}
}
