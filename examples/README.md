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

## Notes

- Customize the client script in `examples/client.py` to test different prompts and parameters.
- Refer to the [Allpaca documentation](https://github.com/MenD32/allpaca) for advanced usage.
