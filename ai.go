package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/tools"
)

type AIEngine struct {
	modelName string
	llm       *ollama.LLM
	mu        sync.Mutex
}

func NewEngine(modelName string) *AIEngine {
	return &AIEngine{modelName: modelName}
}

// StreamPrompt now uses an Agent loop to handle Tool Calling
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

	// 1. Setup the Agent with our custom Tools
	toolsList := GetTools()
	
	// We use LangChainGo's built-in Agent executor
	agent := agents.NewConversationalAgent(e.llm, toolsList)
	executor := agents.NewExecutor(agent, toolsList, 
		agents.WithMaxIterations(3), // Prevent infinite loops
	)

	// 2. We need to intercept the chain to stream the final answer
	// Note: LangChainGo's agent streaming can be complex. 
	// For this TUI, we will run the agent and stream the final output.
	
	// To keep the TUI responsive, we run the agent in the background
	result, err := chains.Run(ctx, executor, prompt)
	
	if err != nil {
		tokenChan <- fmt.Sprintf("\n\n[Agent Error: %s]", err.Error())
		return
	}

	// Stream the final result to the UI
	for _, char := range result {
		select {
		case <-ctx.Done():
			return
		default:
			tokenChan <- string(char)
		}
	}
}