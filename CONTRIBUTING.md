# Contributing to My App

Thank you for your interest in contributing to this project! This document
outlines our coding style guidelines and expectations for contributions.

## Go Style Guide

All Go code in this repository must follow
[Google's Go Style Guide](https://google.github.io/styleguide/go/).

Key points:

- **Formatting**: All Go source files must be formatted with `gofmt`.
  Run `gofmt -w .` before committing.
- **Naming**: Follow the conventions described in
  [Google Go Style Decisions — Naming](https://google.github.io/styleguide/go/decisions#naming).
  Use MixedCaps (exported) or mixedCaps (unexported). Avoid underscores
  in Go names.
- **Error handling**: Return errors rather than using `panic`. Wrap errors
  with `fmt.Errorf("context: %w", err)` to preserve the error chain.
- **Comments**: Exported functions, types, and package declarations must
  have doc comments that begin with the name of the element they describe.
- **Imports**: Group imports into standard library, third-party, and local
  packages, separated by blank lines. Use `goimports` to manage import
  ordering automatically.

For the full guide, see:

- [Google Go Style Guide](https://google.github.io/styleguide/go/)
- [Google Go Style Best Practices](https://google.github.io/styleguide/go/best-practices)
- [Google Go Style Decisions](https://google.github.io/styleguide/go/decisions)

## HTML / CSS Style Guide

All HTML and CSS in this repository must follow
[Google's HTML/CSS Style Guide](https://google.github.io/styleguide/htmlcssguide.html).

Key points:

- **Document type**: Use `<!DOCTYPE html>` (HTML5).
- **Encoding**: Specify `UTF-8` as the character encoding via
  `<meta charset="UTF-8">`.
- **Indentation**: Use 2 spaces for indentation (no tabs).
- **Lowercase**: Use lowercase for HTML element names, attributes, and
  CSS selectors, properties, and values (except strings).
- **Trailing whitespace**: Remove trailing whitespace.
- **Semantic markup**: Use HTML elements for their intended purpose
  (e.g., `<a>` for links, `<button>` for actions).
- **Separation of concerns**: Keep structure (HTML), presentation (CSS),
  and behavior (JavaScript) separate where practical.
- **Protocol**: Omit the protocol (`http:`, `https:`) from URLs pointing
  to images, stylesheets, scripts, and other media files unless the
  resource is not available over both protocols.

For the full guide, see:

- [Google HTML/CSS Style Guide](https://google.github.io/styleguide/htmlcssguide.html)

## General Guidelines

- Keep pull requests focused on a single change.
- Write meaningful commit messages that explain *why* a change was made.
- Add or update tests for any behavioral changes.
- Ensure all tests pass before submitting a pull request.
