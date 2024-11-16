package invoke

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// InvokeModelWrapper encapsulates Amazon Bedrock actions used in the examples.
// It contains a Bedrock Runtime client that is used to invoke foundation models.
type InvokeModelWrapper struct {
	BedrockRuntimeClient *bedrockruntime.Client
}

// Each model provider has their own individual request and response formats.
// For the format, ranges, and default values for Anthropic Claude, refer to:
// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters-claude.html

type ClaudeRequest struct {
	Prompt            string   `json:"prompt"`
	MaxTokensToSample int      `json:"max_tokens_to_sample"`
	Temperature       float64  `json:"temperature,omitempty"`
	StopSequences     []string `json:"stop_sequences,omitempty"`
}

type ClaudeResponse struct {
	Completion string `json:"completion"`
}

// Invokes Anthropic Claude on Amazon Bedrock to run an inference using the input
// provided in the request body.
func CreateClient(ctx context.Context, region *string) *bedrockruntime.Client {
	cfg, err := createSDKConfig(ctx, region)
	if err != nil {
		fmt.Println("Could not load config. Please check config and verify your AWS account is set up correctly")
		fmt.Println(err)
		return nil
	}
	client := bedrockruntime.NewFromConfig(cfg)

	return client
}

func createSDKConfig(ctx context.Context, region *string) (aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(*region))
	if err != nil {
		return aws.Config{}, err
	}
	return cfg, nil
}

func InvokeClaude(ctx context.Context, prompt string, client *bedrockruntime.Client) (string, error) {
	modelId := "anthropic.claude-3-5-sonnet-20240620-v1:0"

	// Anthropic Claude requires enclosing the prompt as follows:
	enclosedPrompt := "Human: " + prompt + "\n\nAssistant:"

	body, err := json.Marshal(ClaudeRequest{
		Prompt:            enclosedPrompt,
		MaxTokensToSample: 200,
		Temperature:       0.5,
		StopSequences:     []string{"\n\nHuman:"},
	})

	if err != nil {
		return "", fmt.Errorf("failed to marshal: %s", err)
	}

	output, err := client.InvokeModel(ctx, &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(modelId),
		ContentType: aws.String("application/json"),
		Body:        body,
	})

	if err != nil {
		log.Fatal(err, modelId)
	}

	var response ClaudeResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		return "", fmt.Errorf("failed to unmarshal: %s", err)
	}

	return response.Completion, nil
}
