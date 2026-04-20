# ts-koans

Interactive TypeScript type-system koans. This npm package downloads a
prebuilt Go binary from the matching [GitHub release][releases] on install.

## Install

```sh
npm i -g ts-koans
```

Then run:

```sh
ts-koans
```

## Prerequisites

- Node.js >= 18 (only required for installation; the koans themselves
  shell out to `tsc`, so you'll also want TypeScript on your `PATH`).

## Supported platforms

- macOS (x64, arm64)
- Linux (x64)
- Windows (x64)

For other platforms, build from source — see the
[project README][repo].

[releases]: https://github.com/chris0lsen/ts-koans/releases
[repo]: https://github.com/chris0lsen/ts-koans
