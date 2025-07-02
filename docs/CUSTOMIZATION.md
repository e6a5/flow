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
    *If you have `$XDG_CONFIG_HOME` set, the path will be `$XDG_CONFIG_HOME/flow/hooks/`.*

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