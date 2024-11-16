package main

import (
	"bufio"
	"chatbot/invoke"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	region := flag.String("region", "us-east-2", "The AWS region")
	flag.Parse()

	fmt.Printf("Using AWS region: %s\n", *region)

	ctx := context.Background()

	client := invoke.CreateClient(ctx, region)

	input := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("User: ")

		input.Scan()
		prompt := input.Text()

		if len(prompt) == 0 {
			continue
		}

		resp, respError := invoke.InvokeClaude(ctx, prompt, client)
		if respError != nil {
			log.Fatal(respError)
		}

		fmt.Printf("Claude:%s\n", resp)
	}

}
