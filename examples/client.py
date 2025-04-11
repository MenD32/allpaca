import openai

openai.base_url = "http://localhost:9090/v1/"
openai.api_key = "1234567890" # this is ignored anyhow

completion = openai.chat.completions.create(
  model="model",
  messages=[
    {"role": "developer", "content": "You are a helpful assistant. Your goal is to provide accurate, concise, and detailed answers to the user's questions. Always ensure your responses are clear and informative."},
    {"role": "user", "content": "Can you tell me what the capital of France is, along with some interesting facts about the city? "}
  ],
  stream=True
)

for chunk in completion:
  print(chunk.choices[0].delta.content, end=' ', flush=True) if chunk.choices[0].delta.content else print()
