#!/usr/bin/env node
// Thin shim that execs the platform-specific tskoans binary that
// install.js placed alongside this file.

const path = require("path");
const fs = require("fs");
const { spawnSync } = require("child_process");

const binary =
  process.platform === "win32"
    ? path.join(__dirname, "tskoans.exe")
    : path.join(__dirname, "tskoans");

if (!fs.existsSync(binary)) {
  console.error(
    "ts-koans: binary not found at " +
      binary +
      "\nTry reinstalling: npm i -g ts-koans",
  );
  process.exit(1);
}

const result = spawnSync(binary, process.argv.slice(2), {
  stdio: "inherit",
});

if (result.error) {
  console.error("ts-koans: failed to launch binary:", result.error.message);
  process.exit(1);
}

process.exit(result.status ?? 0);
