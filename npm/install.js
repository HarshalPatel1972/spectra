const os = require('os');
const path = require('path');
const fs = require('fs');
const https = require('https');
const { execSync } = require('child_process');

const VERSION = require('./package.json').version;
const REPO = 'HarshalPatel1972/spectra';
const BIN_NAME = os.platform() === 'win32' ? 'spectra.exe' : 'spectra';

const PLATFORM_MAP = {
  win32: 'Windows',
  darwin: 'Darwin',
  linux: 'Linux'
};

const ARCH_MAP = {
  x64: 'x86_64',
  arm64: 'arm64',
  ia32: 'i386'
};

async function downloadAndExtract() {
  const platform = PLATFORM_MAP[os.platform()];
  const arch = ARCH_MAP[os.arch()];

  if (!platform || !arch) {
    console.error(`Unsupported platform/architecture: ${os.platform()}-${os.arch()}`);
    process.exit(1);
  }

  const ext = os.platform() === 'win32' ? 'zip' : 'tar.gz';
  const assetName = `spectra_${platform}_${arch}.${ext}`;
  // For testing/development, we point to latest or a specific release
  // Normally this would be: \`https://github.com/${REPO}/releases/download/v${VERSION}/${assetName}\`
  const downloadUrl = `https://github.com/${REPO}/releases/latest/download/${assetName}`;
  const binDir = path.join(__dirname, 'bin');
  
  if (!fs.existsSync(binDir)) {
    fs.mkdirSync(binDir, { recursive: true });
  }

  const archivePath = path.join(binDir, assetName);

  console.log(`Downloading Spectra from ${downloadUrl}...`);
  
  // Download file logic
  const file = fs.createWriteStream(archivePath);
  https.get(downloadUrl, (response) => {
    if (response.statusCode === 301 || response.statusCode === 302) {
      https.get(response.headers.location, (res) => {
        res.pipe(file);
        file.on('finish', () => {
          file.close(() => extractArchive(archivePath, binDir, ext));
        });
      });
    } else if (response.statusCode !== 200) {
      console.error(`Failed to download binary: HTTP ${response.statusCode}`);
      // Don't fail the install if the release isn't published yet
      console.log('Skipping binary download. Ensure you build it manually or release it on GitHub.');
      fs.unlinkSync(archivePath);
      process.exit(0);
    } else {
      response.pipe(file);
      file.on('finish', () => {
        file.close(() => extractArchive(archivePath, binDir, ext));
      });
    }
  }).on('error', (err) => {
    console.error(`Download error: ${err.message}`);
    process.exit(0);
  });
}

function extractArchive(archivePath, destDir, ext) {
  console.log('Extracting archive...');
  try {
    if (ext === 'zip') {
      const AdmZip = require('adm-zip');
      const zip = new AdmZip(archivePath);
      zip.extractAllTo(destDir, true);
    } else {
      const tar = require('tar');
      tar.x({
        file: archivePath,
        cwd: destDir,
        sync: true
      });
    }
    
    // Make binary executable
    const binPath = path.join(destDir, BIN_NAME);
    if (fs.existsSync(binPath) && os.platform() !== 'win32') {
      fs.chmodSync(binPath, 0o755);
    }
    
    console.log('Spectra installed successfully!');
    fs.unlinkSync(archivePath);
  } catch (err) {
    console.error('Error extracting archive:', err);
  }
}

// Only run if not in a CI environment to prevent failures during initial setup
if (!process.env.CI) {
  downloadAndExtract();
} else {
  console.log('Skipping download in CI environment.');
}
