# Customization Guide

Flow can be customized to better fit your workflow using automation hooks and environment variables.

## Automation Hooks

Flow can execute custom scripts based on session events. This allows you to integrate Flow with other tools, send notifications, start music, or automate any part of your deep work ritual.

Hooks are executed for the following events:

- `on_start`: Runs after a session is successfully started.
- `on_pause`: Runs after a session is paused.
- `on_resume`: Runs after a session is resumed.
- `on_end`: Runs after a session is successfully completed.

### How to Use Hooks

1.  **Create a hooks directory:**
    By default, Flow looks for hooks in `~/.config/flow/hooks/`. You may need to create this directory.

    ```bash
    mkdir -p ~/.config/flow/hooks
    ```

    _If you have `$XDG_CONFIG_HOME` set, the path will be `$XDG_CONFIG_HOME/flow/hooks/`._

2.  **Create your script:**
    Inside the `hooks` directory, create an executable file with the same name as the event. For example, to create a hook for the start event, you would create `~/.config/flow/hooks/on_start`.

    Here is an example `on_start` script that sends a desktop notification on macOS:

    ```bash
    #!/bin/bash

    # The session tag is passed as the first argument
    SESSION_TAG=$1

    # Send a notification
    osascript -e "display notification \"$SESSION_TAG\" with title \"Flow Session Started\""
    ```

3.  **Make the script executable:**
    ```bash
    chmod +x ~/.config/flow/hooks/on_start
    ```

Now, every time you run `flow start`, this script will be executed.

## Environment Variables

You can customize the file paths Flow uses for storing its data by setting the following environment variables. This is useful if you want to sync your Flow data using a service like Dropbox or keep it in a non-standard directory.

- **`FLOW_SESSION_PATH`**: Overrides the path for the active session file.

  - **Default:** `~/.local/share/flow/session` (or `$XDG_DATA_HOME/flow/session`)
  - **Example:** `export FLOW_SESSION_PATH=~/Dropbox/flow/session`

- **`FLOW_LOG_PATH`**: Overrides the base directory where the monthly log files are stored.
  - **Default:** `~/.local/share/flow/` (or `$XDG_DATA_HOME/flow/`)
  - **Example:** `export FLOW_LOG_PATH=~/Dropbox/flow/`

You can set these variables in your shell's configuration file (e.g., `~/.bashrc`, `~/.zshrc`) to make them permanent.

## Watcher Configuration

The `flow watch` command can be customized to adjust the timing of its reminders. This is done via a configuration file located at `~/.config/flow/config.yml`.

If the file does not exist, Flow will use the default timings. To customize them, create the `config.yml` file:

```bash
mkdir -p ~/.config/flow
touch ~/.config/flow/config.yml
```

_If you have `$XDG_CONFIG_HOME` set, the path will be `$XDG_CONFIG_HOME/flow/config.yml`._

### Available Options

You can specify the following durations in the YAML file. The values should be strings that can be parsed as a duration (e.g., "5m", "1h", "30s").

Here is a full example showing all available settings:

```yaml
# ~/.config/flow/config.yml
watch:
  # How often the watcher checks your session status.
  interval: "1m"

  # After 30 minutes of inactivity, suggest starting a session.
  remind_after_idle: "30m"

  # After a session has been paused for 10 minutes, suggest resuming.
  remind_after_pause: "10m"

  # After a session has been active for 90 minutes, suggest taking a break.
  remind_after_active: "1h30m"

# Set your daily focus goal (optional)
daily_goal: "4h"

# How long a session can run before being considered stale and auto-cleaned up
# Default: "8h" (8 hours)
stale_session_threshold: "6h"
```

### Stale Session Threshold

The `stale_session_threshold` setting controls how long a session can run before Flow considers it "stale" and automatically cleans it up when you start a new session. This is useful for preventing sessions that accumulate hundreds of hours when you forget to end them.

- **Default:** `"8h"` (8 hours)
- **Example values:** `"4h"`, `"6h"`, `"12h"`, `"24h"`
- **Format:** Any valid Go duration string (e.g., "30m", "2h30m", "1d")

When a session exceeds this threshold, Flow will:
1. Automatically detect it as stale when you run `flow start`
2. Log it as abandoned with an `[ABANDONED]` tag
3. Clean up the session file
4. Allow you to start a fresh session

This prevents the common problem of forgetting to end a session and ending up with inaccurate time tracking data.
