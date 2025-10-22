package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	tokenmanager "github.com/ma3xco/token-listing/internal/token_manager"
)

func main() {
	var isFork bool
	var hasScriptTag bool
	var changedFiles string

	flag.BoolVar(&isFork, "fork", false, "Whether the PR is from a fork")
	flag.BoolVar(&hasScriptTag, "script", false, "Whether the PR has a script tag")
	flag.StringVar(&changedFiles, "files", "", "Comma-separated list of changed files")
	flag.Parse()

	tm, err := tokenmanager.New(context.Background())
	if err != nil {
		log.Fatalf("failed to create token manager: %v", err)
	}
	count, err := tm.WalkThrough(context.Background())
	if err != nil {
		log.Fatalf("failed to walk through tokens: %v", err)
	}
	fmt.Printf("walked through %d tokens\n", count)

	// Apply fork-specific validation if needed
	if isFork && !hasScriptTag {
		fmt.Println("Fork PR detected - applying fork-specific validation rules")

		// Check if only /tokens/* files are changed
		if changedFiles != "" {
			files := strings.Split(changedFiles, ",")
			for _, file := range files {
				file = strings.TrimSpace(file)
				if file != "" && !strings.HasPrefix(file, "tokens/") {
					fmt.Printf("❌ Fork PRs can only modify files in /tokens/* directory. Found change in: %s\n", file)
					os.Exit(1)
				}
			}
		}

		// Apply fork-specific token validation
		forkValidationErrors := tm.ValidateTokensForFork(context.Background())
		if len(forkValidationErrors) > 0 {
			fmt.Println("❌ Fork-specific validation failed:")
			for tokenUid, errors := range forkValidationErrors {
				fmt.Printf("token %s has fork validation errors:\n", tokenUid)
				for _, error := range errors {
					fmt.Printf("  - %s\n", error)
				}
			}
			os.Exit(1)
		}
		fmt.Println("✅ Fork-specific validation passed")
	}

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
