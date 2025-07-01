# 🌊 Flow

> A Command-Line Tool for Deep Work

Flow is a minimalist command-line tool for creating a mindful boundary around your work. It helps you protect your attention and engage in single-tasking.

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/e6a5/flow/main/install.sh | bash
```

*The install script automatically detects your platform, downloads the appropriate binary to `/usr/local/bin`, and provides instructions for enabling shell completion.*

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

## Usage

```bash
flow start --tag "writing docs"   # Enter deep work
flow status                       # Check current session
flow pause                        # Take a mindful break
flow resume                       # Continue working
flow end                          # Complete session
```

## What It Does

*   Creates a mindful boundary around your work
*   Maintains one focused session at a time
*   Runs quietly while you work
*   Provides gentle awareness of your current focus
*   Respects natural work rhythms with pause/resume

## Philosophy

Flow is not a productivity app.
It's a ritual for protecting your attention.
It's a boundary you place around your mind.

Following the Unix philosophy, Flow does one thing well: creating and managing focus sessions. It serves your consciousness, not metrics.

## Core Commands

| Command | Purpose | Example |
|---|---|---|
| `start` | Begin deep work session | `flow start --tag "refactoring"` |
| `status` | Check current session | `flow status` |
| `pause` | Take a mindful break | `flow pause` |
| `resume` | Continue working | `flow resume` |
| `end` | Complete the session | `flow end` |
| `completion`| Generate shell completion script | `flow completion zsh` |


## Examples

### Starting Work

```bash
# Begin focused work
flow start --tag "writing documentation"

🌊 Starting deep work: writing documentation

   Clear your mind
   Focus on what matters
   Let distractions pass

Session active. Use 'flow status' to check.
```

### Checking Progress

```bash
flow status
🌊 Deep work: writing documentation
Active for 1h 23m.
```

### Taking Breaks

```bash
flow pause
⏸️  Paused: writing documentation
Worked for 1h 23m. Use 'flow resume' when ready.

# Later...
flow resume
🌊 Resumed: writing documentation
Continue your deep work.
```

### Completing Work

```bash
flow end
✨ Session complete: writing documentation
Total focus time: 2h 15m
Carry this focus forward.
```

## Shell Completion

To make using `flow` even easier, you can enable shell completion. The `install.sh` script will provide instructions, but you can also generate the script manually.

For example, to enable completion for Zsh, add this to your `.zshrc`:

```bash
source <(flow completion zsh)
```

Supported shells: `bash`, `zsh`.

## Configuration

The session file is stored locally on your machine. The path is determined in the following order:

1.  `$FLOW_SESSION_PATH`: An explicit file path.
2.  `$XDG_DATA_HOME/flow/session`: The standard path on Linux.
3.  `~/.local/share/flow/session`: The fallback XDG path.
4.  `~/.flow-session`: For backward compatibility.

## Design Principles

-   🧘 **Mindful**: A ritualistic approach, not a productivity hack.
-   🎯 **Focused**: One session at a time, no multitasking.
-   🔒 **Private**: No tracking, no cloud, purely local.
-   🌿 **Gentle**: Respects natural work rhythms.
-   🔄 **Simple**: Clean, Unix-like commands.
-   💫 **Present**: Promotes awareness without optimization.


## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines. This is a tool for mindful computing that serves consciousness.

## License

[MIT License](LICENSE) - see the [LICENSE](LICENSE) file for details.

---

*One thing at a time. Runs offline. Powered by presence.* 