/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "LLM chatbot",
	Long:  `LLM chatbot`,
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		// Set up a channel to listen for interrupt signals
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			<-sigChan
			fmt.Println("\nInterrupt signal received. Exiting...")
			os.Exit(0)
		}()

		llm, err := openai.New()
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()

		// Initial LLM prompt phase
		fmt.Print("Enter initial prompt for LLM: ")
		initialPrompt, _ := reader.ReadString('\n')
		initialPrompt = strings.TrimSpace(initialPrompt)
		content := []llms.MessageContent{
			llms.TextParts(llms.ChatMessageTypeSystem, initialPrompt),
		}
		fmt.Println("Initial prompt received. Entering chat mode...")

		for {
			fmt.Print("> ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			switch input {
			case "quit", "exit":
				fmt.Println("Exiting...")
				os.Exit(0)
			default:
				// Process user input with the LLM here
				response := ""
				content = append(content, llms.TextParts(llms.ChatMessageTypeHuman, input))
				llm.GenerateContent(ctx, content,
					llms.WithMaxTokens(1024),
					llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
						fmt.Print(string(chunk))
						response = response + string(chunk)
						return nil
					}),
				)
				content = append(content, llms.TextParts(llms.ChatMessageTypeSystem, response))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chatCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chatCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
