export default class EncryptionService {
    /**
     * Generate a new RSA key pair for email encryption
     * @returns {Promise<{publicKey: string, privateKey: string}>} The key pair as base64 strings
     */
    static async generateKeyPair() {
      try {
        // Generate RSA key pair with Web Crypto API
        const keyPair = await window.crypto.subtle.generateKey(
          {
            name: "RSA-OAEP",
            modulusLength: 2048,
            publicExponent: new Uint8Array([1, 0, 1]), // 65537
            hash: "SHA-256",
          },
          true, // extractable
          ["encrypt", "decrypt"] // key usages
        );
  
        // Export the public key to SPKI format
        const publicKeyBuffer = await window.crypto.subtle.exportKey(
          "spki",
          keyPair.publicKey
        );
  
        // Export the private key to PKCS8 format
        const privateKeyBuffer = await window.crypto.subtle.exportKey(
          "pkcs8",
          keyPair.privateKey
        );
  
        // Convert ArrayBuffer to Base64
        const publicKeyBase64 = this._arrayBufferToBase64(publicKeyBuffer);
        const privateKeyBase64 = this._arrayBufferToBase64(privateKeyBuffer);
  
        return {
          publicKey: publicKeyBase64,
          privateKey: privateKeyBase64
        };
      } catch (error) {
        console.error("Error generating key pair:", error);
        throw new Error("Failed to generate encryption keys");
      }
    }
  
    /**
     * Import a private key from base64 string
     * @param {string} privateKeyBase64 - The base64 encoded private key
     * @returns {Promise<CryptoKey>} The imported private key
     */
    static async importPrivateKey(privateKeyBase64) {
      try {
        const privateKeyBuffer = this._base64ToArrayBuffer(privateKeyBase64);
        
        return await window.crypto.subtle.importKey(
          "pkcs8",
          privateKeyBuffer,
          {
            name: "RSA-OAEP",
            hash: "SHA-256",
          },
          false, // not extractable
          ["decrypt"] // key usage
        );
      } catch (error) {
        console.error("Error importing private key:", error);
        throw new Error("Invalid private key");
      }
    }
  
    /**
     * Decrypt an email using the private key
     * @param {string} encryptedData - Base64 encoded encrypted data
     * @param {string} privateKeyBase64 - Base64 encoded private key
     * @returns {Promise<string>} The decrypted email
     */
    static async decryptEmail(encryptedData, privateKeyBase64) {
      try {
        // Import the private key
        const privateKey = await this.importPrivateKey(privateKeyBase64);
        
        // Convert base64 encrypted data to ArrayBuffer
        const encryptedBuffer = this._base64ToArrayBuffer(encryptedData);
        
        // Decrypt the data
        const decryptedBuffer = await window.crypto.subtle.decrypt(
          {
            name: "RSA-OAEP",
          },
          privateKey,
          encryptedBuffer
        );
        
        // Convert ArrayBuffer to string
        const decoder = new TextDecoder();
        return decoder.decode(decryptedBuffer);
      } catch (error) {
        console.error("Error decrypting email:", error);
        throw new Error("Failed to decrypt email. Invalid key or corrupted data.");
      }
    }
  
    /**
     * Check if a string is base64 encoded (likely encrypted)
     * @param {string} data - The data to check
     * @returns {boolean} True if the data appears to be encrypted
     */
    static isEncrypted(data) {
      // Basic check for base64 encoding pattern
      const base64Regex = /^[A-Za-z0-9+/]+={0,2}$/;
      return base64Regex.test(data);
    }
  
    /**
     * Extracts the private key from URL
     * @returns {string|null} The private key or null if not found
     */
    static getPrivateKeyFromUrl() {
      const url = new URL(window.location);
      const urlParams = new URLSearchParams(url.hash.substring(1).split('?')[1] || '');
      return urlParams.get('key') || null;
    }
  
    /**
     * Sets the private key in the URL while preserving the inbox path
     * @param {string} privateKey - The private key to set
     * @param {string} email - The email address for the inbox
     */
    static setPrivateKeyInUrl(privateKey, email) {
      const url = new URL(window.location);
      const newHash = `/${email}?key=${encodeURIComponent(privateKey)}`;
      
      // Update the URL without triggering a page reload
      history.replaceState(null, '', `${url.origin}${url.pathname}${url.search}#${newHash}`);
    }

    /**
     * Gets the current inbox email from URL
     * @returns {string|null} The email address or null if not found
     */
    static getEmailFromUrl() {
      const hashLocation = window.location.hash;
      if (!hashLocation) return null;
      
      const path = hashLocation.substring(1);
      const emailPart = path.split('?')[0]; // Remove query parameters
      const extractedPath = emailPart.replace("/", "");
      
      // Basic email validation
      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
      if (emailRegex.test(extractedPath)) {
        return extractedPath;
      }
      
      return null;
    }

    /**
     * Updates URL with email and optionally private key
     * @param {string} email - The email address
     * @param {string|null} privateKey - Optional private key
     */
    static updateUrlWithEmail(email, privateKey = null) {
      const url = new URL(window.location);
      let newHash = `/${email}`;
      
      if (privateKey) {
        newHash += `?key=${encodeURIComponent(privateKey)}`;
      }
      
      // Update the URL without triggering a page reload
      history.replaceState(null, '', `${url.origin}${url.pathname}${url.search}#${newHash}`);
    }
  
    /**
     * Convert ArrayBuffer to Base64 string
     * @private
     */
    static _arrayBufferToBase64(buffer) {
      const bytes = new Uint8Array(buffer);
      let binary = '';
      for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
      }
      return btoa(binary);
    }
  
    /**
     * Convert Base64 string to ArrayBuffer
     * @private
     */
    static _base64ToArrayBuffer(base64) {
      const binaryString = atob(base64);
      const bytes = new Uint8Array(binaryString.length);
      for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
      }
      return bytes.buffer;
    }
  }