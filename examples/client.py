import openai

openai.base_url = "http://localhost:8080/v1/"
openai.api_key = "1234567890" # this is ignored anyhow

completion = openai.chat.completions.create(
  model="model",
  messages=[
    {"role": "developer", "content": "You are a helpful assistant."},
    {"role": "user", "content": "Hello!"}
  ],
  stream=True
)

for chunk in completion:
  print(chunk)
