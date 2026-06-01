const crypto = require('crypto');

function generateVulnerableHash(data) {
    // Spectra should flag this MD5 usage
    const hash = crypto.createHash('md5');
    hash.update(data);
    return hash.digest('hex');
}
