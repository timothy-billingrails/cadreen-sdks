#!/usr/bin/env node

const https = require("https");
const fs = require("fs");
const path = require("path");

const REPO = "timothy-billingrails/cadreen-sdks";
const BINARY = "cadreen";
const VERSION = require("./package.json").version;
const MAX_REDIRECTS = 5;

function getPlatform() {
  const platform = process.platform;
  switch (platform) {
    case "darwin":
      return "darwin";
    case "linux":
      return "linux";
    case "win32":
      return "windows";
    default:
      throw new Error(`Unsupported platform: ${platform}`);
  }
}

function getArch() {
  const arch = process.arch;
  switch (arch) {
    case "x64":
      return "amd64";
    case "arm64":
      return "arm64";
    default:
      throw new Error(`Unsupported architecture: ${arch}`);
  }
}

function getDownloadURL() {
  const platform = getPlatform();
  const arch = getArch();
  const ext = platform === "windows" ? ".exe" : "";
  return `https://github.com/${REPO}/releases/download/cli-v${VERSION}/cadreen_${platform}_${arch}${ext}`;
}

function download(url, redirects) {
  redirects = redirects || 0;
  return new Promise((resolve, reject) => {
    if (redirects > MAX_REDIRECTS) {
      reject(new Error("Too many redirects"));
      return;
    }

    const request = https.get(url, { timeout: 30000 }, (response) => {
      if (
        response.statusCode === 302 ||
        response.statusCode === 301 ||
        response.statusCode === 307 ||
        response.statusCode === 308
      ) {
        return download(response.headers.location, redirects + 1).then(
          resolve,
          reject
        );
      }
      if (response.statusCode !== 200) {
        reject(new Error(`Download failed: HTTP ${response.statusCode}`));
        return;
      }

      const chunks = [];
      response.on("data", (chunk) => chunks.push(chunk));
      response.on("end", () => resolve(Buffer.concat(chunks)));
      response.on("error", reject);
    });
    request.on("error", reject);
    request.on("timeout", () => {
      request.destroy();
      reject(new Error("Download timed out"));
    });
  });
}

async function install() {
  const binDir = path.join(__dirname, "bin");
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const ext = process.platform === "win32" ? ".exe" : "";
  const binPath = path.join(binDir, `${BINARY}${ext}`);

  console.log(`Downloading Cadreen CLI v${VERSION}...`);

  const url = getDownloadURL();
  console.log(`  URL: ${url}`);

  try {
    const data = await download(url);
    fs.writeFileSync(binPath, data);
    fs.chmodSync(binPath, 0o755);
    console.log(`  Installed to: ${binPath}`);
    console.log("");
    console.log("Next steps:");
    console.log("  cadreen init    — Set up your account");
    console.log("  cadreen --help  — See all commands");
  } catch (error) {
    console.error(`  Download failed: ${error.message}`);
    console.error("");
    console.error("Install manually:");
    console.error(`  ${url}`);
    process.exit(1);
  }
}

install();
