package main

import (
	"chatbot/invoke"
	"context"
	"flag"
	"fmt"
	"log"
)

func main() {

	region := flag.String("region", "us-east-1", "The AWS region")
	flag.Parse()

	fmt.Printf("Using AWS region: %s\n", *region)

	ctx := context.Background()

	client := invoke.CreateClient(ctx, region)

	prompt := "Hello, there!"

	resp, respError := invoke.InvokeClaude(ctx, prompt, client)
	if respError != nil {
		log.Fatal(respError)
	}

	fmt.Println(resp)

}
