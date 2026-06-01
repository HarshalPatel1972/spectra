const crypto = require('crypto');

function generateWeakHash(data) {
    // Spectra CI will detect this MD5 usage
    const hash = crypto.createHash('md5');
    hash.update(data);
    return hash.digest('hex');
}

// Trigger webhook #2
// Trigger webhook #3
