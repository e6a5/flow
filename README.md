# 🌊 Flow: A Terminal-Based Tool for Deep Work

[![CI](https://github.com/e6a5/flow/actions/workflows/ci.yml/badge.svg)](https://github.com/e6a5/flow/actions/workflows/ci.yml)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/e6a5/flow)](https://github.com/e6a5/flow/releases/latest)
[![Go Version](https://img.shields.io/github/go-mod/go-version/e6a5/flow)](https://go.dev/)
[![GitHub license](https://img.shields.io/github/license/e6a5/flow)](https://github.com/e6a5/flow/blob/main/LICENSE)

**Flow is a minimalist command-line tool for focused, single-tasking work sessions. It protects your attention, helps you build a deep work habit, and provides powerful insights into your focus patterns—all without leaving your terminal.**

![Flow Dashboard](assets/dashboard.png)

It's designed for developers, writers, and anyone who wants to build a more mindful and effective relationship with their work.

---

## Table of Contents

- [The Philosophy](#the-philosophy-your-attention-is-sacred)
- [Features](#features)
- [Installation](#installation)
- [Getting Started](#getting-started-a-typical-workflow)
- [Commands](#full-command-reference)
- [Customization](#customization)
- [Contributing](#contributing)

---

## The Philosophy: Your Attention is Sacred

> In a world of constant distraction, your ability to focus is a superpower. Flow is built on a simple idea: **one thing at a time**. It's not about complex productivity metrics or chasing a never-ending task list. It's about creating a clear, intentional boundary around your work, allowing you to engage deeply and mindfully.
>
> Flow helps you answer a simple question: "What am I working on right now?" And by logging your completed sessions, it helps you reflect on a more important one: "How am I investing my attention?"

## Features

- **Mindful Focus**: A single active session at a time to encourage deep, single-tasking work.
- **Rich Dashboard**: A beautiful, GitHub-style contribution graph to visualize your focus history over the last year.
- **Powerful Exports**: Export your work sessions to CSV or JSON for invoicing, analysis, or personal records.
- **Privacy First**: Your data is yours. Everything is stored locally in plain text files. No cloud, no tracking.
- **Shell Integration**: Seamlessly display your current focus session in your shell prompt (`bash` and `zsh` supported).
- **Automation Hooks**: Trigger custom scripts on session events (`on_start`, `on_pause`, `on_end`).

## Installation

```bash
curl -sSL https://raw.githubusercontent.com/e6a5/flow/main/install.sh | bash
```
*The installer will add the `flow` binary to `/usr/local/bin` and check for necessary dependencies.*

For other installation methods (Go, manual), see the [Installation Guide](docs/INSTALL.md).

## Getting Started: A Typical Workflow

Flow is designed to be intuitive. Here's how a typical session works:

1.  **Start a session** when you're ready to focus. Give it a tag to describe your task.
    ```bash
    flow start --tag "Writing the first draft of the new feature spec"
    ```

2.  **Check your status** at any time.
    ```bash
    flow status
    > 🌊 Deep work: Writing the first draft of the new feature spec (Active for 1h 15m)
    ```

3.  **Take a break** when you need one.
    ```bash
    flow pause
    ```
    Then, **resume** when you're ready to get back to it.
    ```bash
    flow resume
    ```

4.  **End the session** when the work is complete. Your focus time is automatically logged.
    ```bash
    flow end
    > ✨ Session complete: Writing the first draft of the new feature spec
    > Total focus time: 2h 30m
    ```

## Gain Insights from Your Work

Once you've logged a few sessions, you can use Flow's data tools to understand your work patterns.

- **Visualize your consistency** with the dashboard. The color of each day represents your total focus time:
  - **Lightest Blue**: 1 minute - 2 hours
  - **Light Blue**: 2 - 4 hours
  - **Medium Blue**: 4 - 6 hours
  - **Darkest Blue**: More than 6 hours
  ```bash
  flow dashboard
  ```

- **Review your session history**:
  ```bash
  flow log --week --stats
  ```
- **Export your data for invoicing or analysis**:
  ```bash
  flow export --month 2023-10 --format csv --output "october-invoice.csv"
  ```

## Full Command Reference

### Core Session Commands
| Command | Description |
| ------- | ----------- |
| `start [--tag "name"]` | Begin a deep work session. |
| `status [--raw]` | Check the current session status. |
| `pause` | Pause the active session. |
| `resume` | Resume a paused session. |
| `end` | Complete the session and log it. |

### Data & Analysis Commands
| Command | Description |
| ------- | ----------- |
| `log [flags]` | View completed session history. See `flow log --help` for flags. |
| `dashboard` | Show a yearly contribution graph of your focus sessions. |
| `export [flags]` | Export session data to CSV or JSON. See `flow export --help` for flags. |

### Utility Commands
| Command | Description |
| ------- | ----------- |
| `completion [bash\|zsh]` | Generate shell completion scripts. |

## Customization

You can extend Flow to fit your unique workflow using hooks and environment variables.

- **Automation Hooks**: Trigger custom scripts on session events.
- **Configuration**: Customize storage paths using environment variables.

For detailed information, see the [Customization Guide](docs/CUSTOMIZATION.md).

## Contributing

Flow is built for the community, and we welcome contributions! Whether it's a bug report, a feature request, or a pull request, we'd love to hear from you. Please see our [Contributing Guidelines](CONTRIBUTING.md) to get started.

## License

Flow is open-source software licensed under the [MIT License](LICENSE).

---

*One thing at a time. Runs offline. Powered by presence.*