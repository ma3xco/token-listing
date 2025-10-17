package main

import (
	"context"
	"fmt"
	"log"

	tokenmanager "github.com/ma3xco/token-listing/internal/token_manager"
)

func main() {
	// build the tokens
	tm, err := tokenmanager.New(context.Background())
	if err != nil {
		log.Fatalf("failed to create token manager: %v", err)
	}
	count, err := tm.WalkThrough(context.Background())
	if err != nil {
		log.Fatalf("failed to walk through tokens: %v", err)
	}
	fmt.Printf("walked through %d tokens\n", count)
	err = tm.BuildTokens(context.Background())
	if err != nil {
		log.Fatalf("failed to build tokens: %v", err)
	}
	fmt.Println("build assets completed")
}
