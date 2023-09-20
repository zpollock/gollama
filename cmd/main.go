package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	chatPrompt = flag.String("chat-prompt", "A chat between a curious user and an artificial intelligence assistant. The assistant follows the given rules no matter what.\n", "the top prompt in chat completions")
	userName   = flag.String("user-name", "\nUSER: ", "USER name in chat completions")
	aiName     = flag.String("ai-name", "\nASSISTANT: ", "ASSISTANT name in chat completions")
	systemName = flag.String("system-name", "\nASSISTANT's RULE: ", "SYSTEM name in chat completions")
	stop       = flag.String("stop", "</s>", "the end of the response in chat completions")
	llamaAPI   = flag.String("llama-api", "http://127.0.0.1:8080", "Set the address of server.cpp in llama.cpp")
	apiKey     = flag.String("api-key", "", "Set the API key to allow only a few users")
	host       = flag.String("host", "0.0.0.0", "Set the IP address to listen")
	port       = flag.Int("port", 8081, "Set the port to listen")
)

func main() {
	flag.Parse()
	r := gin.Default()

	r.POST("/chat/completions", chatCompletionsHandler)
	r.POST("/completions", completionsHandler)

	addr := *host + ":" + strconv.Itoa(*port)
	r.Run(addr)
}

func chatCompletionsHandler(c *gin.Context) {

	if *apiKey != "" {
		authHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer "+*apiKey) {
			c.JSON(http.StatusForbidden, nil)
			return
		}
	}

	var body map[string]interface{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stream := false
	tokenize := false

	if s, ok := body["stream"].(bool); ok {
		stream = s
	}

	if t, ok := body["tokenize"].(bool); ok {
		tokenize = t
	}

	postData := makePostData(body, true, stream)
	promptToken := []string{}

	if tokenize {
		tokenData := makeTokenizeRequest(postData["prompt"].(string))
		if tokens, ok := tokenData["tokens"].([]string); ok {
			promptToken = tokens
		}
	}

	var resData map[string]interface{}
	if !stream {
		data := makeCompletionRequest(postData)
		log.Println(data)
		resData = makeResData(data, true, promptToken)
		c.JSON(http.StatusOK, resData)
	} else {
		c.Stream(func(w io.Writer) bool {
			data := makeStreamCompletionRequest(postData)
			timeNow := int(time.Now().Unix())
			for _, line := range data {
				resData := makeResDataStream(line, true, timeNow)
				w.Write([]byte("data: " + resData + "\n"))
				w.(http.Flusher).Flush()
			}
			return false
		})
	}
}

func completionsHandler(c *gin.Context) {
	if *apiKey != "" {
		authHeader := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer "+*apiKey) {
			c.JSON(http.StatusForbidden, nil)
			return
		}
	}

	var body map[string]interface{}
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stream := false
	tokenize := false

	if s, ok := body["stream"].(bool); ok {
		stream = s
	}

	if t, ok := body["tokenize"].(bool); ok {
		tokenize = t
	}

	postData := makePostData(body, false, stream)
	promptToken := []string{}

	if tokenize {
		tokenData := makeTokenizeRequest(postData["prompt"].(string))
		if tokens, ok := tokenData["tokens"].([]string); ok {
			promptToken = tokens
		}
	}

	var resData map[string]interface{}
	if !stream {
		data := makeCompletionRequest(postData)
		log.Println(data)
		resData = makeResData(data, false, promptToken)
		c.JSON(http.StatusOK, resData)
	} else {
		c.Stream(func(w io.Writer) bool {
			data := makeStreamCompletionRequest(postData)
			timeNow := int(time.Now().Unix())
			for _, line := range data {
				resData := makeResDataStream(line, false, timeNow)
				w.Write([]byte("data: " + resData + "\n"))
				w.(http.Flusher).Flush()
			}
			return false
		})
	}
}

func isPresent(data map[string]interface{}, key string) bool {
	_, ok := data[key]
	return ok
}

func convertChat(messages []interface{}) string {
	var prompt strings.Builder
	prompt.WriteString(*chatPrompt)

	for _, msg := range messages {
		messageMap, ok := msg.(map[string]interface{})
		if !ok {
			fmt.Println("Message is not in the expected format")
			continue
		}

		role, roleOk := messageMap["role"].(string)
		content, contentOk := messageMap["content"].(string)

		if !roleOk || !contentOk {
			fmt.Println("Message fields are not in the expected format")
			continue
		}

		prompt.WriteString(getPrompt(role, content))
	}

	prompt.WriteString(strings.TrimRight(*aiName, "\n"))
	prompt.WriteString(*stop)
	return prompt.String()
}

func getPrompt(role, content string) string {
	switch role {
	case "system":
		return *systemName + content
	case "user":
		return *userName + content
	case "assistant":
		return *aiName + content
	}
	return ""
}

func makePostData(body map[string]interface{}, chat bool, stream bool) map[string]interface{} {
	postData := make(map[string]interface{})

	if chat {
		postData["prompt"] = convertChat(body["messages"].([]interface{}))
	} else {
		postData["prompt"] = body["prompt"].(string)

		if isPresent(body, "stop") {
			stops := body["stop"].([]interface{})
			for _, v := range stops {
				postData["stop"] = append(postData["stop"].([]string), v.(string))
			}
		}
	}

	postData["stop"] = []string{*stop}
	if isPresent(body, "stop") {
		stops := body["stop"].([]interface{})
		for _, v := range stops {
			postData["stop"] = append(postData["stop"].([]string), v.(string))
		}
	}

	if isPresent(body, "temperature") {
		postData["temperature"] = body["temperature"].(float64)
	}
	if isPresent(body, "top_k") {
		postData["top_k"] = int(body["top_k"].(float64))
	}
	if isPresent(body, "top_p") {
		postData["top_p"] = body["top_p"].(float64)
	}
	if isPresent(body, "max_tokens") {
		postData["n_predict"] = int(body["max_tokens"].(float64))
	}
	if isPresent(body, "presence_penalty") {
		postData["presence_penalty"] = body["presence_penalty"].(float64)
	}
	if isPresent(body, "frequency_penalty") {
		postData["frequency_penalty"] = body["frequency_penalty"].(float64)
	}
	if isPresent(body, "repeat_penalty") {
		postData["repeat_penalty"] = body["repeat_penalty"].(float64)
	}
	if isPresent(body, "mirostat") {
		postData["mirostat"] = body["mirostat"].(string)
	}
	if isPresent(body, "mirostat_tau") {
		postData["mirostat_tau"] = body["mirostat_tau"].(float64)
	}
	if isPresent(body, "mirostat_eta") {
		postData["mirostat_eta"] = body["mirostat_eta"].(float64)
	}
	if isPresent(body, "seed") {
		postData["seed"] = int(body["seed"].(float64))
	}

	if isPresent(body, "logit_bias") {
		logitBias := make(map[int]float64)
		for k, v := range body["logit_bias"].(map[string]interface{}) {
			key, _ := strconv.Atoi(k)
			logitBias[key] = v.(float64)
		}
		postData["logit_bias"] = logitBias
	}

	postData["n_keep"] = -1
	postData["stream"] = stream

	fmt.Printf("Request: %v\n", postData)
	return postData
}

func makeResData(data map[string]interface{}, chat bool, promptToken []string) map[string]interface{} {
	resData := make(map[string]interface{})

	if chat {
		resData["id"] = "chatcmpl"
		resData["object"] = "chat.completion"
	} else {
		resData["id"] = "cmpl"
		resData["object"] = "text_completion"
	}

	resData["created"] = int(time.Now().Unix())
	resData["truncated"] = data["truncated"]
	resData["model"] = "LLaMA_CPP"
	resData["usage"] = map[string]interface{}{
		"prompt_tokens":     len(promptToken),
		"completion_tokens": data["tokens_predicted"],
		"total_tokens":      len(promptToken) + int(data["tokens_predicted"].(float64)),
	}

	if len(promptToken) != 0 {
		resData["promptToken"] = promptToken
	}

	if chat {
		choices := []map[string]interface{}{
			{
				"index": 0,
				"message": map[string]interface{}{
					"role":    "assistant",
					"content": data["content"],
				},
				"finish_reason": "stop",
			},
		}
		if !data["stopped_eos"].(bool) && !data["stopped_word"].(bool) {
			choices[0]["finish_reason"] = "length"
		}
		resData["choices"] = choices
	} else {
		choices := []map[string]interface{}{
			{
				"text":          data["content"],
				"index":         0,
				"logprobs":      nil,
				"finish_reason": "stop",
			},
		}
		if !data["stopped_eos"].(bool) && !data["stopped_word"].(bool) {
			choices[0]["finish_reason"] = "length"
		}
		resData["choices"] = choices
	}

	return resData
}

func makeTokenizeRequest(prompt string) map[string]interface{} {
	data := map[string]interface{}{
		"content": prompt,
	}
	jsonData, _ := json.Marshal(data)
	response, err := http.Post(*llamaAPI+"/tokenize", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)
	var tokenData map[string]interface{}
	json.Unmarshal(responseBody, &tokenData)

	return tokenData
}

func makeCompletionRequest(postData map[string]interface{}) map[string]interface{} {
	jsonData, _ := json.Marshal(postData)
	response, err := http.Post(*llamaAPI+"/completion", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseBody, _ := ioutil.ReadAll(response.Body)
	var data map[string]interface{}
	json.Unmarshal(responseBody, &data)

	return data
}

func makeStreamCompletionRequest(postData map[string]interface{}) []map[string]interface{} {
	jsonData, _ := json.Marshal(postData)
	response, err := http.Post(*llamaAPI+"/completion", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var data []map[string]interface{}
	decoder := json.NewDecoder(response.Body)

	for decoder.More() {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			log.Fatal(err)
		}
		data = append(data, chunk)
	}

	return data
}

func makeResDataStream(data map[string]interface{}, chat bool, timeNow int) string {
	type Choice struct {
		FinishReason interface{} `json:"finish_reason"`
		Index        int         `json:"index"`
		Delta        interface{} `json:"delta,omitempty"`
		Text         string      `json:"text,omitempty"`
	}

	type ResData struct {
		ID      string   `json:"id"`
		Object  string   `json:"object"`
		Created int      `json:"created"`
		Model   string   `json:"model"`
		Choices []Choice `json:"choices"`
	}

	resData := ResData{
		ID:      "chatcmpl",
		Object:  "chat.completion.chunk",
		Created: timeNow,
		Model:   "LLaMA_CPP",
		Choices: []Choice{
			{
				FinishReason: nil,
				Index:        0,
			},
		},
	}

	if chat {
		if data["start"].(bool) {
			resData.Choices[0].Delta = map[string]interface{}{
				"role": "assistant",
			}
		} else {
			resData.Choices[0].Delta = map[string]interface{}{
				"content": data["content"],
			}
			if data["stop"].(bool) {
				if data["stopped_eos"].(bool) || data["stopped_word"].(bool) {
					resData.Choices[0].FinishReason = "stop"
				} else {
					resData.Choices[0].FinishReason = "length"
				}
			}
		}
	} else {
		resData.Choices[0].Text = data["content"].(string)
		if data["stop"].(bool) {
			if data["stopped_eos"].(bool) || data["stopped_word"].(bool) {
				resData.Choices[0].FinishReason = "stop"
			} else {
				resData.Choices[0].FinishReason = "length"
			}
		}
	}

	jsonData, _ := json.Marshal(resData)
	return string(jsonData)
}
