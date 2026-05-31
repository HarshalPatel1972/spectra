#!/usr/bin/env node

const { spawnSync } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const BIN_NAME = os.platform() === 'win32' ? 'spectra.exe' : 'spectra';
const binPath = path.join(__dirname, 'bin', BIN_NAME);

if (!fs.existsSync(binPath)) {
  console.error(`Error: Spectra binary not found at ${binPath}`);
  console.error('Please run "npm run postinstall" or ensure the binary is installed.');
  process.exit(1);
}

const args = process.argv.slice(2);

const result = spawnSync(binPath, args, { stdio: 'inherit' });

if (result.error) {
  console.error(`Error executing Spectra: ${result.error.message}`);
  process.exit(1);
}

process.exit(result.status || 0);
