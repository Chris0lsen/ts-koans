#!/usr/bin/env node
// Downloads the platform-specific tskoans binary from the matching GitHub
// release and extracts it into ./bin so the `ts-koans` shim can exec it.
//
// All release archives are .tar.gz (including Windows). Windows 10 1803+
// and Windows 11 ship `tar.exe` in System32, so we can rely on it.

const fs = require("fs");
const path = require("path");
const os = require("os");
const https = require("https");
const { spawnSync } = require("child_process");

const REPO = "Chris0lsen/ts-koans";
const pkg = require("./package.json");
const VERSION = process.env.TSKOANS_VERSION || pkg.version;

const PLATFORM_MAP = {
  darwin: "darwin",
  linux: "linux",
  win32: "windows",
};

const ARCH_MAP = {
  x64: "amd64",
  arm64: "arm64",
};

// Match the build matrix in .github/workflows/release.yml
const SUPPORTED = new Set([
  "darwin-amd64",
  "darwin-arm64",
  "linux-amd64",
  "windows-amd64",
]);

function detectTarget() {
  const goos = PLATFORM_MAP[process.platform];
  const goarch = ARCH_MAP[process.arch];
  if (!goos || !goarch) {
    throw new Error(
      `Unsupported platform: ${process.platform}/${process.arch}.`,
    );
  }
  const key = `${goos}-${goarch}`;
  if (!SUPPORTED.has(key)) {
    throw new Error(
      `No prebuilt binary for ${key}. ` +
        `Please open an issue at https://github.com/${REPO}/issues.`,
    );
  }
  return {
    archive: `tskoans-${goos}-${goarch}.tar.gz`,
    binary: goos === "windows" ? "tskoans.exe" : "tskoans",
  };
}

function get(url) {
  return new Promise((resolve, reject) => {
    https
      .get(
        url,
        { headers: { "User-Agent": "ts-koans-npm-installer" } },
        (res) => {
          if (
            res.statusCode &&
            res.statusCode >= 300 &&
            res.statusCode < 400 &&
            res.headers.location
          ) {
            res.resume();
            resolve(get(res.headers.location));
            return;
          }
          if (res.statusCode !== 200) {
            reject(new Error(`Download failed: ${url} -> ${res.statusCode}`));
            return;
          }
          const chunks = [];
          res.on("data", (c) => chunks.push(c));
          res.on("end", () => resolve(Buffer.concat(chunks)));
          res.on("error", reject);
        },
      )
      .on("error", reject);
  });
}

function extractTarGz(buf, destDir, binaryName) {
  const tmp = path.join(os.tmpdir(), `tskoans-${Date.now()}.tar.gz`);
  fs.writeFileSync(tmp, buf);
  try {
    const r = spawnSync("tar", ["-xzf", tmp, "-C", destDir], {
      stdio: "inherit",
    });
    if (r.error) throw r.error;
    if (r.status !== 0) {
      throw new Error(
        `tar exited with status ${r.status}. ` +
          `If you're on Windows, ensure tar.exe is available.`,
      );
    }
  } finally {
    try {
      fs.unlinkSync(tmp);
    } catch {}
  }
  const out = path.join(destDir, binaryName);
  fs.chmodSync(out, 0o755);
  return out;
}

async function main() {
  // Skip in dev contexts where the package isn't actually being installed.
  if (process.env.TSKOANS_SKIP_DOWNLOAD === "1") {
    console.log("TSKOANS_SKIP_DOWNLOAD set; skipping binary download.");
    return;
  }

  const { archive, binary } = detectTarget();
  const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${archive}`;
  const binDir = path.join(__dirname, "bin");
  fs.mkdirSync(binDir, { recursive: true });

  console.log(`Downloading ${url}`);
  const buf = await get(url);
  extractTarGz(buf, binDir, binary);
  console.log(`Installed ${binary} into ${binDir}`);
}

main().catch((err) => {
  console.error("ts-koans install failed:", err.message);
  process.exit(1);
});