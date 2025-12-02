# Contributing to Vulcan

Thank you for your interest in contributing to **Vulcan**! 🎉

We welcome all contributions — whether it's bug fixes, new features, documentation updates, or suggestions.

## How to Contribute

### 1. Fork the repository
Click the **Fork** button on GitHub and clone your fork locally.

```bash
git clone https://github.com/OhMyDitzzy/vulcan.git
cd vulcan
````

### 2. Create a feature branch

```bash
git checkout -b feature/my-feature
```

### 3. Make changes

Follow the existing project structure and coding standards.

* Go code → `go fmt` + idiomatic Go
* React/TS → ESLint + Prettier

### 4. Run tests

```bash
make test
cd vulcan-web
npm test
```

### 5. Commit with clear messages

```bash
git commit -m "feat: add new miner dashboard widget"
```

### 6. Push and create a Pull Request

```bash
git push origin feature/my-feature
```

Open a PR on GitHub and describe your changes.

## Code Style

### Go

* Use tabs for indentation
* Follow idiomatic Go conventions
* Document exported functions
* Keep functions small and focused

### TypeScript / React

* Use functional components + hooks
* Use TypeScript types whenever possible
* Keep components pure
* Follow the existing folder structure

## Reporting Issues

Open an issue on GitHub with:

* Steps to reproduce
* Expected behavior
* Actual behavior
* Logs or screenshots

## Security Issues

**Do not** open public issues for security vulnerabilities.
Email the maintainer directly instead.

## Thank You

Your contributions help make Vulcan better for everyone! 🚀