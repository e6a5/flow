# Installation Guide

There are several ways to install Flow. The recommended method for most users is the installer script.

## Installer Script (macOS & Linux)

This is the simplest way to install Flow. The script will detect your operating system and architecture, download the latest release binary, and move it to `/usr/local/bin`.

```bash
curl -sSL https://raw.githubusercontent.com/e6a5/flow/main/install.sh | bash
```

You may be prompted for your password if the script needs `sudo` permissions to write to `/usr/local/bin`.

## Using `go install`

If you have a Go environment set up, you can install Flow directly using `go install`. This will download the source code and build the binary for your system.

```bash
go install github.com/e6a5/flow@latest
```

This will install the `flow` binary into your `$GOPATH/bin` directory. Make sure this directory is in your system's `PATH`.

## Building from Source

You can also build Flow from the source code for full control.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/e6a5/flow.git
    cd flow
    ```

2.  **Build the binary:**
    The `Makefile` provides a convenient way to build the project. This will create a `flow` binary in the root of the project directory.
    ```bash
    make build
    ```

3.  **Move the binary to your PATH:**
    You can then move this binary to any directory in your system's `PATH`. For example:
    ```bash
    mv ./flow /usr/local/bin/
    ```

## Shell Completion

After installing, it's highly recommended to set up shell completion for a better user experience. Add the appropriate line for your shell to your shell's configuration file (e.g., `~/.bashrc`, `~/.zshrc`).

**For Bash:**
```bash
echo 'eval "$(flow completion bash)"' >> ~/.bashrc
```

**For Zsh:**
```bash
echo 'eval "$(flow completion zsh)"' >> ~/.zshrc
```

You will need to restart your shell for the changes to take effect. 