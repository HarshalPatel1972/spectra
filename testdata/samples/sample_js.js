/**
 * Sample JavaScript file with deliberate cryptographic usage for Spectra testing.
 *
 * Contains various Node.js crypto API calls that the Spectra scanner should
 * detect and flag for quantum risk assessment.
 */

const crypto = require("crypto");

/**
 * Compute a SHA-1 hash – cryptographically weakened, quantum-risky.
 * @param {Buffer|string} data - Input data to hash.
 * @returns {string} Hex-encoded SHA-1 digest.
 */
function hashSHA1(data) {
  return crypto.createHash("sha1").update(data).digest("hex");
}

/**
 * Compute an MD5 hash – cryptographically broken.
 * @param {Buffer|string} data - Input data to hash.
 * @returns {string} Hex-encoded MD5 digest.
 */
function hashMD5(data) {
  return crypto.createHash("md5").update(data).digest("hex");
}

/**
 * Create an RSA-SHA256 signature – RSA is quantum-vulnerable via Shor's algorithm.
 * @param {Buffer} data - Data to sign.
 * @param {string} privateKey - PEM-encoded RSA private key.
 * @returns {Buffer} Digital signature.
 */
function signRSASHA256(data, privateKey) {
  const signer = crypto.createSign("RSA-SHA256");
  signer.update(data);
  return signer.sign(privateKey);
}

/**
 * Verify an ECDSA signature – ECDSA is quantum-vulnerable.
 * Uses the "ECDSA" algorithm family which is broken by Shor's algorithm.
 * @param {Buffer} data - Original data.
 * @param {Buffer} signature - Signature to verify.
 * @param {string} publicKey - PEM-encoded public key.
 * @returns {boolean} True if signature is valid.
 */
function verifyECDSA(data, signature, publicKey) {
  const verifier = crypto.createVerify("SHA256");
  verifier.update(data);
  // ECDSA verification – quantum-vulnerable
  return verifier.verify(publicKey, signature);
}

/**
 * Compute a SHA-256 hash – considered quantum-safe with 128-bit post-quantum security.
 * @param {Buffer|string} data - Input data to hash.
 * @returns {string} Hex-encoded SHA-256 digest.
 */
function hashSHA256(data) {
  return crypto.createHash("sha256").update(data).digest("hex");
}

// Example usage
if (require.main === module) {
  const sample = "quantum computing threat assessment data";
  console.log("SHA-1:", hashSHA1(sample));
  console.log("MD5:", hashMD5(sample));
  console.log("SHA-256:", hashSHA256(sample));
}

module.exports = { hashSHA1, hashMD5, signRSASHA256, verifyECDSA, hashSHA256 };
