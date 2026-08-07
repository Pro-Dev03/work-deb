# Claude CLI / Anthropic SDK Setup Guide

## ✅ Installation Status

**Anthropic Python SDK:** Installed successfully

## 🔑 Getting API Key

1. Go to: https://console.anthropic.com
2. Sign up or log in
3. Go to API Keys section
4. Create new API key
5. Copy the key

## ⚙️ Configuration

### Option 1: Environment Variable (Recommended)

```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

Add to your `~/.bashrc` or `~/.zshrc` for persistence:

```bash
echo 'export ANTHROPIC_API_KEY="your-api-key-here"' >> ~/.bashrc
source ~/.bashrc
```

### Option 2: Direct in Python Script

```python
import anthropic

client = anthropic.Anthropic(api_key="your-api-key-here")
```

## 📝 Usage Examples

### Simple Text Generation

```python
import anthropic

client = anthropic.Anthropic()

message = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=1024,
    messages=[
        {
            "role": "user",
            "content": "Hello, Claude!"
        }
    ]
)

print(message.content[0].text)
```

### Streaming Response

```python
import anthropic

client = anthropic.Anthropic()

with client.messages.stream(
    model="claude-3-5-sonnet-20241022",
    max_tokens=1024,
    messages=[{"role": "user", "content": "Hello!"}]
) as stream:
    for text in stream.text_stream:
        print(text, end="", flush=True)
```

### With System Prompt

```python
import anthropic

client = anthropic.Anthropic()

message = client.messages.create(
    model="claude-3-5-sonnet-20241022",
    max_tokens=1024,
    system="You are a helpful assistant.",
    messages=[
        {
            "role": "user",
            "content": "Explain quantum computing"
        }
    ]
)

print(message.content[0].text)
```

## 📊 Available Models

- **claude-3-5-sonnet-20241022** - Balanced, fast, affordable
- **claude-3-5-haiku-20241022** - Fastest, cheapest
- **claude-3-opus-20240229** - Most capable, most expensive

## 💰 Pricing (as of 2024)

| Model | Input | Output |
|-------|-------|--------|
| Claude 3.5 Sonnet | $3 / 1M tokens | $15 / 1M tokens |
| Claude 3.5 Haiku | $0.25 / 1M tokens | $1.25 / 1M tokens |
| Claude 3 Opus | $15 / 1M tokens | $75 / 1M tokens |

## 🧪 Test Installation

Create a test file `test_claude.py`:

```python
import anthropic

client = anthropic.Anthropic()

try:
    message = client.messages.create(
        model="claude-3-5-haiku-20241022",
        max_tokens=100,
        messages=[{"role": "user", "content": "Say 'Hello World'"}]
    )
    print("✅ Success!")
    print(message.content[0].text)
except Exception as e:
    print(f"❌ Error: {e}")
```

Run it:
```bash
python3 test_claude.py
```

## 📚 Documentation

Full documentation: https://docs.anthropic.com

## 🔧 Troubleshooting

### Error: "No API key provided"
- Set `ANTHROPIC_API_KEY` environment variable
- Or pass `api_key` parameter when creating client

### Error: "Rate limit exceeded"
- You've hit your rate limit
- Wait or upgrade your plan

### Error: "Insufficient credits"
- Add credits to your Anthropic account
- Check billing at console.anthropic.com

## 🎯 Quick Start Script

```bash
# Set API key
export ANTHROPIC_API_KEY="your-api-key-here"

# Test
python3 -c "
import anthropic
client = anthropic.Anthropic()
message = client.messages.create(
    model='claude-3-5-haiku-20241022',
    max_tokens=100,
    messages=[{'role': 'user', 'content': 'Hello!'}]
)
print(message.content[0].text)
"
```

---

*Last Updated: 2024*
