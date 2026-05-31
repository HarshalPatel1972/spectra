package samples;

import java.security.KeyPairGenerator;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import javax.crypto.Cipher;
import javax.crypto.NoSuchPaddingException;

/**
 * Sample Java file with deliberate cryptographic usage for Spectra testing.
 *
 * Contains various crypto API calls that the Spectra scanner should detect
 * and flag for quantum risk assessment.
 */
public class sample_java {

    /**
     * Generate an RSA key pair – quantum-vulnerable via Shor's algorithm.
     */
    public static void generateRSAKeyPair() throws NoSuchAlgorithmException {
        KeyPairGenerator kpg = KeyPairGenerator.getInstance("RSA");
        kpg.initialize(2048);
        kpg.generateKeyPair();
    }

    /**
     * Create a DES cipher – broken even classically, worse under quantum.
     */
    public static Cipher createDESCipher()
            throws NoSuchAlgorithmException, NoSuchPaddingException {
        return Cipher.getInstance("DES");
    }

    /**
     * Compute an MD5 digest – cryptographically broken.
     */
    public static byte[] hashMD5(byte[] data) throws NoSuchAlgorithmException {
        MessageDigest md = MessageDigest.getInstance("MD5");
        return md.digest(data);
    }

    /**
     * Compute a SHA-1 digest – weakened, not quantum-safe.
     */
    public static byte[] hashSHA1(byte[] data) throws NoSuchAlgorithmException {
        MessageDigest md = MessageDigest.getInstance("SHA-1");
        return md.digest(data);
    }

    /**
     * Compute a SHA-256 digest – considered quantum-safe.
     */
    public static byte[] hashSHA256(byte[] data) throws NoSuchAlgorithmException {
        MessageDigest md = MessageDigest.getInstance("SHA-256");
        return md.digest(data);
    }

    public static void main(String[] args) throws Exception {
        byte[] sample = "quantum computing threat assessment data".getBytes();
        generateRSAKeyPair();
        createDESCipher();
        byte[] md5 = hashMD5(sample);
        byte[] sha1 = hashSHA1(sample);
        byte[] sha256 = hashSHA256(sample);
    }
}
