# Neh

[![Go Report Card](https://goreportcard.com/badge/github.com/neh-cli/neh)](https://goreportcard.com/report/github.com/neh-cli/neh)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/neh-cli/neh?sort=semver)

Neh is a sleek CLI application designed for high-speed, real-time interaction with AI specializing in Large Language Models.

<p align="center">
  <img src="https://raw.githubusercontent.com/neh-cli/neh/refs/heads/main/screencast/screencast.gif" alt="Screencast">
</p>

## Installation

### For macOS

```bash
brew install neh-cli/tap/neh
```

## Subcommands Completion

### Bash Shell Completion

To enable bash shell completion for the `neh` command, you can use the following command:

```bash
source <(neh completion bash)
```

This command will generate the necessary completion script and source it into your current shell session. To make this change permanent, you can add the command to your shell's startup file (e.g., `~/.bashrc` or `~/.bash_profile`).

### Zsh Shell Completion

Similarly, for `zsh` shell completion, you can use the following command:

```bash
source <(neh completion zsh)
```

And add it to your `~/.zshrc` file to ensure the completion is available in every new terminal session.

### Fish Shell Completion

For `fish` shell users, you can enable completion with:

```bash
neh completion fish | source
```

To make this permanent, you can write the output to a file in your `~/.config/fish/completions` directory.


## FAQ

### Q: Getting "zsh: no matches found" error

When running `neh` commands in zsh with special characters like `?` or `*`, you might encounter this error:

```
% neh o 1 + 1 = ?
zsh: no matches found: ?
```

This occurs because zsh interprets certain characters as glob patterns. To resolve this, follow these steps:

1. Check if alias already exists:
   ```bash
   grep "alias neh" ~/.zshrc
   ```

2. If not found, add the alias:
   ```bash
   echo "alias neh='noglob neh'" >> ~/.zshrc
   ```

3. Reload your shell configuration:
   ```bash
   source ~/.zshrc
   ```

This disables glob expansion for the `neh` command, allowing special characters to be passed as arguments.