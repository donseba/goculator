package llm

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var (
	// ChatGPTToken is the token for the ChatGPT function
	ChatGPTToken string
)

func ChatGPT(args []any) (any, error) {
	var question string

	for _, arg := range args {
		question = fmt.Sprintf("%s : %+v", question, arg)
	}

	url := "https://api.openai.com/v1/completions"
	var jsonData = []byte(fmt.Sprintf(`
        {
            "model": "gpt-3.5-turbo",
            "content": "%s"
        }`, question))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ChatGPTToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	_ = json.Unmarshal(body, &result)
	if choices, found := result["choices"].([]interface{}); found && len(choices) > 0 {
		firstChoice := choices[0].(map[string]interface{})
		if text, found := firstChoice["text"].(string); found {
			return text, nil
		}
	}

	if errs, found := result["error"].(map[string]interface{}); found {
		if message, found := errs["message"].(string); found {
			return "", errors.New(message)
		}
	}

	fmt.Printf("Unexpected response: %#v\n", result)

	return "", fmt.Errorf("unexpected response structure")
}
