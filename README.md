# 🌊 Flow

> A Command-Line Tool for Deep Work

Flow is a minimalist command-line tool designed to help you focus on single-tasking by creating a mindful boundary around your work. It emphasizes mindfulness and presence, not productivity metrics.

## ⚡ Quick Start

### Install (one-liner)

```bash
curl -sSL https://raw.githubusercontent.com/e6a5/flow/main/install.sh | bash
```

### Start a Session

```bash
flow start --tag "writing docs"
```

### Enable Shell Integration (Optional)

For ambient awareness in your shell prompt:

```bash
# For bash users
echo 'eval "$(flow completion bash)"' >> ~/.bashrc

# For zsh users  
echo 'eval "$(flow completion zsh)"' >> ~/.zshrc
```

## 🌟 Key Features

- **Mindful Focus**: Helps you maintain a single-tasking mindset with one active session at a time
- **Session Management**: Start, pause, resume, and end sessions easily
- **Shell Integration**: Seamlessly integrates with your shell for ambient awareness
- **Automation Hooks**: Customize workflows with session event hooks (`on_start`, `on_pause`, `on_resume`, `on_end`)
- **Privacy First**: No tracking, no cloud, purely local with XDG-compliant storage
- **Script Friendly**: Raw output mode for integration with other tools

## 🌿 Philosophy

Flow believes in creating a mindful boundary around your work. It's about protecting your attention and engaging deeply with one task at a time. The tool serves your consciousness, not metrics.

## 💡 Commands

| Command | Description | Examples |
| ------- | ----------- | -------- |
| `start [--tag "name"]` | Begin a deep work session | `flow start --tag "code review"` |
| `status [--raw]` | Check the current session | `flow status` or `flow status --raw` |
| `pause` | Take a mindful break | `flow pause` |
| `resume` | Continue working | `flow resume` |
| `end` | Complete the session | `flow end` |
| `completion [bash\|zsh]` | Generate shell completions | `flow completion bash` |

## 🔧 Installation and Usage

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/e6a5/flow/main/install.sh | bash
```

### Go Install

```bash
go install github.com/e6a5/flow@latest
```

### Download Binaries

Download pre-built binaries from the [releases page](https://github.com/e6a5/flow/releases).

### Build from Source

```bash
git clone https://github.com/e6a5/flow.git
cd flow
make build
```

## 🎣 Automation Hooks

Flow supports automation hooks that trigger on session events. Create executable scripts in `~/.config/flow/hooks/`:

- `on_start` - Triggered when starting a session
- `on_pause` - Triggered when pausing a session  
- `on_resume` - Triggered when resuming a session
- `on_end` - Triggered when ending a session

Each hook receives the session tag as its first argument.

Example hook (`~/.config/flow/hooks/on_start`):
```bash
#!/bin/bash
echo "Starting work on: $1" | notify-send "Flow" 
```

## 🛠️ Configuration

Flow follows XDG Base Directory standards:

- **Session data**: `$XDG_DATA_HOME/flow/session` (default: `~/.local/share/flow/session`)
- **Hooks**: `$XDG_CONFIG_HOME/flow/hooks/` (default: `~/.config/flow/hooks/`)
- **Custom session path**: Set `FLOW_SESSION_PATH` environment variable



## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines. This is a tool for mindful computing that serves consciousness.

## License

[MIT License](LICENSE) - see the [LICENSE](LICENSE) file for details.

---

*One thing at a time. Runs offline. Powered by presence.*