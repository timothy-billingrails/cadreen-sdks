#!/usr/bin/env node
const fs = require("fs");
const path = require("path");
const { execSync } = require("child_process");

const projectName = process.argv[2] || "my-cadreen-project";
const targetDir = path.resolve(projectName);

if (fs.existsSync(targetDir)) {
  console.error(`Directory "${projectName}" already exists.`);
  process.exit(1);
}

const templateDir = path.join(__dirname, "template");

function copyDir(src, dest) {
  fs.mkdirSync(dest, { recursive: true });
  for (const entry of fs.readdirSync(src, { withFileTypes: true })) {
    const srcPath = path.join(src, entry.name);
    const destPath = path.join(dest, entry.name);
    if (entry.isDirectory()) {
      copyDir(srcPath, destPath);
    } else {
      fs.copyFileSync(srcPath, destPath);
    }
  }
}

console.log(`Creating Cadreen project in ${targetDir}...`);
copyDir(templateDir, targetDir);

const pkgPath = path.join(targetDir, "package.json");
const pkg = JSON.parse(fs.readFileSync(pkgPath, "utf-8"));
pkg.name = projectName;
fs.writeFileSync(pkgPath, JSON.stringify(pkg, null, 2));

console.log("Installing dependencies...");
execSync("npm install", { cwd: targetDir, stdio: "inherit" });

console.log("Done! Start building:");
console.log(`  cd ${projectName}`);
console.log("Set CADREEN_API_KEY in .env, then:");
console.log("  npm run dev");
