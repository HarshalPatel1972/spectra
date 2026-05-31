"""
Sample Python file with deliberate cryptographic usage for Spectra testing.

This file contains various crypto patterns that Spectra's code scanner
should detect and flag for quantum risk assessment.
"""

import hashlib
import hmac

# PyCryptodome RSA key generation – quantum-vulnerable
from Crypto.PublicKey import RSA

# cryptography library RSA – also quantum-vulnerable
from cryptography.hazmat.primitives.asymmetric import rsa
from cryptography.hazmat.backends import default_backend


def generate_rsa_key_pycrypto():
    """Generate an RSA key using PyCryptodome (quantum-vulnerable)."""
    key = RSA.generate(2048)
    return key.export_key()


def generate_rsa_key_cryptography():
    """Generate an RSA key using the cryptography library (quantum-vulnerable)."""
    private_key = rsa.generate_private_key(
        public_exponent=65537,
        key_size=2048,
        backend=default_backend(),
    )
    return private_key


def hash_md5(data: bytes) -> str:
    """Compute an MD5 hash – cryptographically broken, quantum-weak."""
    return hashlib.md5(data).hexdigest()


def hash_sha1(data: bytes) -> str:
    """Compute a SHA-1 hash – cryptographically weakened, quantum-risky."""
    return hashlib.sha1(data).hexdigest()


def hash_sha256(data: bytes) -> str:
    """Compute a SHA-256 hash – considered quantum-safe."""
    return hashlib.sha256(data).hexdigest()


def hmac_sha1(key: bytes, msg: bytes) -> str:
    """Compute HMAC-SHA1 – uses SHA-1 internally."""
    return hmac.new(key, msg, hashlib.sha1).hexdigest()


if __name__ == "__main__":
    sample = b"quantum computing threat assessment data"
    print("MD5:", hash_md5(sample))
    print("SHA-1:", hash_sha1(sample))
    print("SHA-256:", hash_sha256(sample))
