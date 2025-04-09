generate:
    curl http://localhost:8080/v1/chat/completions \
        -H "Content-Type: application/json" \
        -d '{ "stream": true, "stream_options": {"include_usage": true}, "model": "model", "messages": [ { "role": "user", "content": "Who are you?" } ] }'
