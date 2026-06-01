const crypto = require('crypto');

function insecureLegacyFunction() {
    // Spectra should catch this SHA-1 usage!
    const shasum = crypto.createHash('sha1');
    shasum.update('hunter2');
    return shasum.digest('hex');
}

// Trigger webhook!
