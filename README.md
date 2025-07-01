# 🌊 Flow

> Protect Your Attention

A mindful boundary for deep work, designed in the spirit of [Zenta](https://github.com/e6a5/zenta).

## Installation

### Quick Install (Recommended)

```bash
curl -sSL https://raw.githubusercontent.com/e6a5/flow/main/install.sh | bash
```

*The install script automatically detects your platform and downloads the appropriate binary to `/usr/local/bin`.*

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

* Creates a mindful boundary around your work
* Maintains one focused session at a time
* Runs quietly in background while you work
* Provides gentle awareness of current focus
* Respects natural work rhythms with pause/resume

## Philosophy

Flow is not a productivity app.  
It's a ritual for protecting attention.  
It's a boundary you place around your mind.

Following Unix philosophy, Flow does one thing well: session boundaries.  
Following Zenta spirit, it serves consciousness, not metrics.

## Core Commands

| Command | Purpose | Example |
|---------|---------|---------|
| `start` | Begin deep work session | `flow start --tag "refactoring"` |
| `status` | Check current session | `flow status` |
| `pause` | Mindful break | `flow pause` |
| `resume` | Continue working | `flow resume` |
| `end` | Complete session | `flow end` |

## Examples

### Starting Work
```bash
# Begin focused work
flow start --tag "writing documentation"

🌊 Starting deep work: writing documentation
   Clear your mind
   Focus on what matters  
   Let distractions pass
Session active in background.
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

## Composing with Zenta

```bash
flow end && zenta           # Breathe after completion
flow start --tag "deep work" # Enter focus
flow pause && zenta now     # Mindful break
```

## Design Principles

- 🧘 **Mindful**: Ritual approach, not productivity hack
- 🎯 **Focused**: One session at a time, no multitasking
- 🔒 **Private**: No tracking, no cloud, purely local
- 🌿 **Gentle**: Respects natural work rhythms
- 🔄 **Simple**: Clean Unix-like commands
- 💫 **Present**: Awareness without optimization

## Why Flow?

**Real vs Fake Focus Tools:**

✅ **Real (Flow's way):**
- Creates sacred space for work
- One thing at a time enforcement  
- Mindful boundaries, not metrics
- Serves attention, not productivity

❌ **Fake:**
- Tracking productivity scores
- Optimizing work performance
- Gamifying focus time
- Measuring output metrics

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

Flow follows the same principles as Zenta: mindful computing that serves consciousness.

## License

[MIT License](LICENSE) - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Zenta](https://github.com/e6a5/zenta) - mindfulness for terminal users
- Built with intention for conscious developers

---

*One thing at a time. Runs offline. Powered by presence.* 