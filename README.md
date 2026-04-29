# ts-koans

Inspired by [Elixir Koans](https://github.com/elixirkoans/elixir-koans), ts-koans is an interactive way to learn about TypeScript's type system, from the basic to the complex. 

## Prerequisites

 - [Node.js](https://nodejs.org/en/download/) >= 20
 - [TypeScript](https://www.typescriptlang.org/docs/handbook/typescript-tooling-in-5-minutes.html)

Note that if you use a version manager such as [nvm](https://github.com/nvm-sh/nvm), [n](https://github.com/tj/n), [asdf](https://asdf-vm.com/) or [mise](https://mise.jdx.dev/), your `tsc` installation might not be globally available. Please make sure they're in your `$PATH` before running ts-koans.

Installing from npm (either globally or locally) will download `tsc` and make it available, so running `ts-koans` will use the one installed by npm, if available.

## Running

### NPM

```bash
npm i -g ts-koans
```

and then

```bash
ts-koans
```

### GitHub Release

Download the latest release for your architecture, then extract and run `tskoans` from your favorite terminal emulator.

### From Source

Alternatively, you may clone this repo and run `go run .` from the root. This requires golang to be available in your `$PATH`.

## Problems?

Please open an issue if you encounter any errors! This is still very early in development. It is not "battle-tested" or "hardened." In fact it is quite soft and pleasantly squishy.

## Contributing

What, you think you can make it better? You probably can! I'm not a golang expert, or a typescript expert, or an expert in many other things. If you want to open a PR with some enhancements, I will do my best to review it (but I make no promises).