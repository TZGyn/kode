package agent

import (
	"context"
	"errors"
	"fmt"
	"os"

	"charm.land/fantasy"
	"charm.land/fantasy/providers/openai"
	"charm.land/fantasy/providers/openrouter"
)

type Config struct {
	provider fantasy.Provider
	model    fantasy.LanguageModel
	ctx      context.Context
}

func NewProvider(provider string, apiKey string) (fantasy.Provider, error) {
	if provider == "openai" {
		provider, err := openai.New(openai.WithAPIKey(apiKey))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating OpenAI provider: %v\n", err)
			os.Exit(1)
		}

		return provider, nil
	}

	if provider == "openrouter" {
		provider, err := openrouter.New(openrouter.WithAPIKey(apiKey))

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating OpenRouter provider: %v\n", err)
			os.Exit(1)
		}
		return provider, nil
	}

	return nil, errors.New("Invalid Provider")
}

func NewModel(ctx context.Context, provider fantasy.Provider, modelID string) (fantasy.LanguageModel, error) {

	model, err := provider.LanguageModel(ctx, modelID)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return model, nil
}

func NewConfig(providerID string, modelID string, apiKey string) Config {
	ctx := context.Background()

	provider, err := NewProvider(providerID, apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating Config: %v\n", err)
		os.Exit(1)
	}

	model, err := NewModel(ctx, provider, modelID)

	return Config{
		provider: provider,
		model:    model,
		ctx:      ctx,
	}
}
