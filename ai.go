package main

import (
	"context"
	"sync"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

type AIEngine struct {
	modelName string
	llm       *ollama.LLM
	mu        sync.Mutex
}

func NewEngine(modelName string) *AIEngine {
	return &AIEngine{modelName: modelName}
}

// StreamPrompt connects to local Ollama and streams tokens to the provided channel.
func (e *AIEngine) StreamPrompt(ctx context.Context, prompt string, tokenChan chan<- string, doneChan chan<- bool) {
	defer func() { doneChan <- true }()

	e.mu.Lock()
	if e.llm == nil {
		llm, err := ollama.New(ollama.WithModel(e.modelName))
		if err != nil {
			tokenChan <- "\n\n[Error: Could not connect to Ollama. Is it running?]"
			return
		}
		e.llm = llm
	}
	e.mu.Unlock()

	systemPrompt := "You are Forge, an elite AI software architect running locally. Be concise, technical, and use markdown for code blocks."
	fullPrompt := systemPrompt + "\n\nUser Query: " + prompt

	_, err := e.llm.Call(ctx, fullPrompt,
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
				tokenChan <- string(chunk)
				return nil
			}
		}),
	)

	if err != nil && ctx.Err() == nil {
		tokenChan <- "\n\n[Error: " + err.Error() + "]"
	}
}