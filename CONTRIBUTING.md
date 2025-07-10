# Contributing to Flow

We love your input! We want to make contributing to Flow as easy and transparent as possible, whether it's:

- Reporting a bug
- Discussing the current state of the code
- Submitting a fix
- Proposing new features
- Becoming a maintainer

## Philosophy

Flow follows these core principles:

- **Unix Philosophy**: Do one thing well
- **Minimalism**: Simple and focused
- **No Dependencies**: Pure Go, no external libraries
- **Privacy First**: No tracking, no cloud
- **Mindful Computing**: Technology that serves consciousness

## Development Process

We use GitHub to host code, to track issues and feature requests, as well as accept pull requests.

### Pull Requests

We actively welcome your pull requests.

1.  Fork the repo and create your branch from `main`.
2.  If you've added code that should be tested, add tests.
3.  Ensure the test suite passes with `make test`.
4.  Make sure your code is formatted with `make fmt`.
5.  Title your Pull Request using the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format. This helps us automate our release notes. For example:
    - `feat: Add a new command for exporting data`
    - `fix: Correctly handle paused session state`
    - `docs: Update the README with new instructions`
6.  Add a label to your Pull Request that matches the type of change. This is used to categorize the change in our release notes. The available labels are:
    - `feature` / `enhancement`
    - `bug` / `fix`
    - `documentation` / `docs`
    - `chore` / `refactor` / `ci`
7.  Issue that pull request!

### Release Process (for maintainers)

This project uses `release-drafter` to automate the creation of release notes.

### Development Setup

```bash
# Clone your fork
git clone https://github.com/yourusername/flow.git
cd flow

# Build and test
make build
make test

# Run development version
make dev
```

### Code Style

- Run `go fmt` on your code (or `make fmt`)
- Follow standard Go conventions
- Keep functions small and focused
- Write clear, descriptive variable names
- Add comments for complex logic

## Bug Reports

We use GitHub issues to track public bugs. Report a bug by [opening a new issue](https://github.com/e6a5/flow/issues).

**Great Bug Reports** tend to have:

- A quick summary and/or background
- Steps to reproduce
  - Be specific!
  - Give sample code if you can
- What you expected would happen
- What actually happens
- Notes (possibly including why you think this might be happening, or stuff you tried that didn't work)

## Feature Requests

We welcome feature requests, but remember Flow's philosophy of simplicity. Before proposing:

1. Check if it aligns with the Unix philosophy
2. Consider if it could be achieved through composition with other tools
3. Ensure it doesn't add complexity to the core use case

## License

By contributing, you agree that your contributions will be licensed under its MIT License.

## Questions?

Feel free to open an issue with the `question` label if you have any questions about contributing.
