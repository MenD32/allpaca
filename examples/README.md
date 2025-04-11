# Allpaca Tool Example

This example demonstrates how to use the Allpaca tool to set up an OpenAI mock server and test it with a Python client script.

## Prerequisites

- Install Allpaca by following the [installation guide](https://github.com/MenD32/allpaca).
- Ensure Python is installed on your system.
- Install the OpenAI API Python package by running:

    ```bash
    pip install openai
    ```

## Steps

1. **Start the Mock Server**

    Run the following command to start the OpenAI mock server:

    ```bash
    allpaca mock-server start
    ```

    The server will start on `http://localhost:8080` by default.

2. **Use the Client Script**

    A sample client script is provided in the `examples/client.py` file. You can use it to send requests to the mock server.

3. **Run the Client Script**

    Execute the client script to send a request to the mock server:

    ```bash
    python examples/client.py
    ```

    You should see a mock response printed in the terminal.

## Example with Custom Config

This example demonstrates how to use a custom configuration file with Allpaca.

1. **Create a Custom Config File (Optional)**

    Create a JSON configuration file, for example `example_config.json`, with the desired mock server settings:

    ```json
    {
        "port": 9090,
        "chat_endpoint": "/v1/chat/completions",
        "model": "model",
        "address": "127.0.0.1",
        "itl_val": 0.1,
        "ttft_val": 0.2
    }
    ```

2. **Start the Mock Server with the Custom Config**

    Use the following command to start the mock server with the configuration file:

    ```bash
    allpaca run -c examples/example_config.json
    ```

    The server will now start on `http://localhost:9090` and use the custom responses defined in the configuration file.

3. **Run the Client Script**

    Update the client script in `examples/client.py` to point to the new server URL (`http://localhost:9090`) and execute it:

    ```bash
    python examples/client.py
    ```

    You should see the custom mock response printed in the terminal.

## Notes

- Customize the client script in `examples/client.py` to test different prompts and parameters.
- Refer to the [Allpaca documentation](https://github.com/MenD32/allpaca) for advanced usage.
- Use the custom configuration file to simulate various scenarios and responses.
