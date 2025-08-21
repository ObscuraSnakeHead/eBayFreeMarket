package apis

import (
	"strings"

	"github.com/xyproto/ollamaclient"
)

func AskLlama(prompt string) (*string, error) {
	oc := ollamaclient.NewWithModelAndAddr(APPLICATION_SETTINGS.OllamaModel, APPLICATION_SETTINGS.OllamaAPIUrl)
	oc.SetRandomOutput()

	if err := oc.PullIfNeeded(); err != nil {
		return nil, err
	}

	output, err := oc.GetOutput(prompt)
	if err != nil {
		return nil, err
	}
	out := strings.TrimSpace(output)
	return &out, nil
}

func TestSpam(title, description string) (*string, error) {
	prompt1 := strings.ReplaceAll(APPLICATION_SETTINGS.OllamaSystemPrompt, "<TITLE>", title)
	prompt2 := strings.ReplaceAll(prompt1, "<DESCRIPTION>", description)

	return AskLlama(prompt2)
}
