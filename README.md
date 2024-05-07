# Go API Wrapper for llama.cpp

This project provides a simple and easy-to-deploy API wrapper around the llama.cpp server using Go and Docker. It allows you to quickly set up and run the llama.cpp server alongside a Gin server in a Docker container.

## Prerequisites

Before you begin, ensure you have the following prerequisites:

- Docker installed on your machine.
- Docker Compose installed on your machine.
- Obtain a local llama model in bin or gguf format.

## Getting Started

**Clone the Repository:**

Clone this Git repository to your local machine:

git clone https://github.com/zpollock/gollama.git

cd gollama


**Configuration:**

Configure the project by setting the necessary environment variables in the `.env` file. You will need to specify the following fields:

- `LLAMA_MODEL_PATH`: The path to mount the llama model in the container.
- `LLAMA_MODEL`: The name of the llama model file.
- `LLAMA_CPP_SERVER`: The URL of the llama.cpp server.
- `PORT`: The port for the API

**Get llama.cpp*:*

Clone the llama.cpp into the gollama/chat folder:

cd chat

git clone https://github.com/ggerganov/llama.cpp.git

**Build the Docker Container:**

Build the Docker container using Docker Compose:

docker-compose build


**Deploy the Container:**

Run the Docker container using Docker Compose:

docker-compose up


The Gin server and llama.cpp server will be started in the same container, and your API will be available at `http://127.0.0.1:<PORT>`.

## Usage

You can now send requests to your API to interact with the llama.cpp server. For example:

```bash
#Request
curl --request POST \
     --url http://127.0.0.1:8081/chat/completions \
     --header "Content-Type: application/json" \
     --data '{
       "messages": [
         {
           "role": "USER",
           "content": "What is the secret of your success?"
         }
       ],
       "n_predict": 128
     }'
#Response
{
    "choices":[
    {
        "finish_reason":"stop",
        "index":0,"message":{
            "content":" Hello! *adjusts glasses* I'm glad you asked! My success is due to my programming, which includes a vast knowledge base and the ability to process natural language. However, I must inform you that I am not capable of revealing any secrets or confidential information. It is important to respect privacy and security, both online and offline. Is there anything else I can help you with? *smiles*\n","role":"assistant"
        }
    }
    ],
    "created":1695205085,
    "id":"chatcmpl",
    "model":"LLaMA_CPP",
    "object":"chat.completion",
    "truncated":false,
    "usage":{
        "completion_tokens":39,
        "prompt_tokens":0,
        "total_tokens":39
    }
}  
```

## Command-Line Flags

The following command-line flags are available when running the API:

1. **-chat-prompt**: Specifies the top prompt in chat completions. Default value is "A chat between a curious user and an artificial intelligence assistant. The assistant follows the given rules no matter what."

2. **-user-name**: Sets the USER name in chat completions. Default value is "\nUSER: ".

3. **-ai-name**: Sets the ASSISTANT name in chat completions. Default value is "\nASSISTANT: ".

4. **-system-name**: Sets the SYSTEM name in chat completions. Default value is "\nASSISTANT's RULE: ".

5. **-stop**: Defines the end of the response in chat completions. Default value is "\</s>".

6. **-llama-api**: Sets the address of the server.cpp in llama.cpp. Default value is "http://127.0.0.1:8080".

7. **-api-key**: Sets the API key to allow only a few users. Default value is an empty string.

8. **-host**: Sets the IP address to listen on. Default value is "0.0.0.0".

9. **-port**: Sets the API port to listen on. Default value is 8081.

## Post Body Variables

The API accepts POST requests with a JSON payload in the request body. Below are the available variables that can be included in the JSON payload:

1. **stream**: A boolean indicating whether to stream the response. If set to true, the API will stream the response data. If false, the API will return a single response. Default value is false.

2. **tokenize**: A boolean indicating whether to tokenize the input prompt. If set to true, the API will tokenize the prompt before processing. Default value is false.

3. **messages** (Only for chat completions): An array of message objects representing a conversation. Each message object should have the following properties:
   - **role**: A string representing the role of the sender (e.g., "system", "user", "assistant").
   - **content**: The content of the message as a string.

4. **prompt** (Only for text completions): A string containing the input prompt for the completion.

5. **stop**: An array of strings defining the stopping conditions for completion. This can be used to specify when the completion should stop based on certain keywords or conditions.

6. **temperature**: A float64 value controlling the randomness of the output. Higher values make the output more random. Default value is not specified in the code.

7. **top_k**: An integer value controlling the number of top tokens to consider. Default value is not specified in the code.

8. **top_p**: A float64 value controlling the nucleus sampling parameter. Default value is not specified in the code.

9. **max_tokens**: An integer value specifying the maximum number of tokens in the output. Default value is not specified in the code.

10. **presence_penalty**: A float64 value controlling the presence penalty. Default value is not specified in the code.

11. **frequency_penalty**: A float64 value controlling the frequency penalty. Default value is not specified in the code.

12. **repeat_penalty**: A float64 value controlling the repeat penalty. Default value is not specified in the code.

13. **mirostat**: A string value not specified in the code.

14. **mirostat_tau**: A float64 value not specified in the code.

15. **mirostat_eta**: A float64 value not specified in the code.

16. **seed**: An integer value specifying the random seed. Default value is not specified in the code.

17. **logit_bias**: A map of integer keys to float64 values. Not specified in the code.


## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please create a GitHub issue or submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
