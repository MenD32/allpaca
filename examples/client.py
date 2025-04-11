import openai
import datetime

openai.base_url = "http://localhost:9090/v1/"
openai.api_key = "1234567890" # this is ignored anyhow

start_time = datetime.datetime.now()
completion = openai.chat.completions.create(
  model="model",
  messages=[
    {"role": "developer", "content": "You are a helpful assistant. Your goal is to provide accurate, concise, and detailed answers to the user's questions. Always ensure your responses are clear and informative."},
    {"role": "user", "content": "Can you tell me what the capital of France is, along with some interesting facts about the city? "}
  ],
  stream=True
)

token_times = []
for chunk in completion:
  token_times.append(datetime.datetime.now())
  print(chunk.choices[0].delta.content, end=' ', flush=True) if chunk.choices[0].delta.content else print()

ttft = token_times[0].timestamp() - start_time.timestamp()
itl = (token_times[-1].timestamp() - token_times[0].timestamp()) / (len(token_times) - 1) if len(token_times) > 1 else 0

print(f"Time to first token: {ttft} seconds")
print(f"Average time between tokens: {itl} seconds")
