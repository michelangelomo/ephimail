<template>
  <div class="reservation-container mt-4">
    <div v-if="!reserved" class="card">
      <div class="card-header bg-primary text-white">
        <h5 class="mb-0">Reserve this mailbox</h5>
      </div>
      <div class="card-body">
        <p class="card-text">Reserving this mailbox prevents others from accessing emails sent to it.</p>
        
        <div class="mb-3">
          <label for="duration" class="form-label">Reservation Duration</label>
          <select v-model="duration" id="duration" class="form-select">
            <option value="1h">1 Hour</option>
            <option value="24h">1 Day</option>
            <option value="168h">1 Week</option>
          </select>
        </div>
        
        <div class="form-check mb-3">
          <input v-model="useEncryption" class="form-check-input" type="checkbox" id="encryption">
          <label class="form-check-label" for="encryption">
            Enable End-to-End Encryption
          </label>
        </div>
        
        <div v-if="useEncryption" class="alert alert-info" role="alert">
          <h6><i class="fas fa-lock"></i> About End-to-End Encryption</h6>
          <p class="mb-0">
            When enabled, all emails will be encrypted with a key that only you have.
            Your unique private key will be stored in the URL - bookmark it to maintain access.
            <strong>If you lose this URL, you cannot recover your emails!</strong>
          </p>
        </div>
        
        <button @click="reserveMailbox" class="btn btn-primary" :disabled="isLoading">
          <span v-if="isLoading" class="spinner-border spinner-border-sm" role="status" aria-hidden="true"></span>
          {{ isLoading ? 'Processing...' : 'Reserve Mailbox' }}
        </button>
      </div>
    </div>
    
    <div v-else class="card">
      <div class="card-header bg-success text-white">
        <h5 class="mb-0">Mailbox Reserved</h5>
      </div>
      <div class="card-body">
        <p class="card-text">
          This mailbox is reserved until <strong>{{ formatDate(expiresAt) }}</strong>
        </p>
        
        <div v-if="encrypted" class="alert alert-warning">
          <h6><i class="fas fa-exclamation-triangle"></i> Important Security Information</h6>
          <p>This mailbox uses end-to-end encryption. Your unique access URL is:</p>
          <div class="input-group mb-3">
            <input type="text" class="form-control" :value="reservationUrl" readonly>
            <button class="btn btn-outline-secondary" type="button" @click="copyUrl">
              <i class="fas fa-copy"></i> Copy
            </button>
          </div>
          <p class="mb-0">
            <strong>Bookmark this URL immediately!</strong> It contains your private key.
            If you lose this URL, you will not be able to decrypt your emails.
          </p>
        </div>
        
        <button @click="releaseReservation" class="btn btn-danger">
          Release Reservation
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import EncryptionService from '../services/encryption';

export default {
  props: {
    email: {
      type: String,
      required: true
    }
  },
  
  data() {
    return {
      duration: '24h',
      useEncryption: false,
      isLoading: false,
      reserved: false,
      encrypted: false,
      expiresAt: null,
      reservationUrl: '',
      privateKey: null
    };
  },
  
  mounted() {
    // Check if there's a private key in the URL
    this.privateKey = EncryptionService.getPrivateKeyFromUrl();
    
    // Check if mailbox is already reserved
    this.checkReservation();
  },
  
  methods: {
    async reserveMailbox() {
      this.isLoading = true;
      
      try {
        let publicKey = '';
        
        // Generate key pair if encryption is enabled
        if (this.useEncryption) {
          const keyPair = await EncryptionService.generateKeyPair();
          publicKey = keyPair.publicKey;
          this.privateKey = keyPair.privateKey;
        }
        
        // Send reservation request
        const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/inbox/reserve`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            email: this.email,
            duration: this.duration,
            public_key: publicKey
          })
        });
        
        if (!response.ok) {
          throw new Error(`Failed to reserve mailbox: ${response.statusText}`);
        }
        
        const data = await response.json();
        
        // Update component state
        this.reserved = true;
        this.encrypted = data.encrypted;
        this.expiresAt = new Date(data.expires_at);
        
        // If encrypted, set the private key in URL and create access URL
        if (this.encrypted && this.privateKey) {
          EncryptionService.setPrivateKeyInUrl(this.privateKey);
          this.reservationUrl = window.location.href;
          
          // Notify user to bookmark the URL
          this.$notify({
            title: 'Encryption Enabled',
            text: 'Bookmark this URL immediately to keep access to your emails!',
            type: 'warning',
            duration: 10000
          });
        }
      } catch (error) {
        console.error('Reservation error:', error);
        this.$notify({
          title: 'Reservation Failed',
          text: error.message,
          type: 'error'
        });
      } finally {
        this.isLoading = false;
      }
    },
    
    async checkReservation() {
      try {
        const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/inbox/${this.email}/reservation`);
        
        if (response.ok) {
          const data = await response.json();
          
          this.reserved = true;
          this.encrypted = data.encrypted;
          this.expiresAt = new Date(data.expires_at);
          
          // If encrypted but no private key in URL, show warning
          if (this.encrypted && !this.privateKey) {
            this.$notify({
              title: 'Encryption Key Missing',
              text: 'This mailbox is encrypted but your access key is missing. You will not be able to read emails.',
              type: 'error',
              duration: 0 // No auto-dismiss
            });
          }
        }
      } catch (error) {
        console.error('Error checking reservation:', error);
      }
    },
    
    async releaseReservation() {
      if (!confirm('Are you sure you want to release this reservation? This action cannot be undone.')) {
        return;
      }
      
      try {
        const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/inbox/${this.email}/reservation`, {
          method: 'DELETE'
        });
        
        if (!response.ok) {
          throw new Error(`Failed to release reservation: ${response.statusText}`);
        }
        
        this.reserved = false;
        this.encrypted = false;
        this.expiresAt = null;
        this.privateKey = null;
        this.reservationUrl = '';
        
        // Remove private key from URL
        window.location.hash = '';
        
        this.$notify({
          title: 'Reservation Released',
          text: 'Your mailbox reservation has been released.',
          type: 'success'
        });
      } catch (error) {
        console.error('Error releasing reservation:', error);
        this.$notify({
          title: 'Error',
          text: error.message,
          type: 'error'
        });
      }
    },
    
    formatDate(date) {
      if (!date) return '';
      return new Date(date).toLocaleString();
    },
    
    copyUrl() {
      navigator.clipboard.writeText(this.reservationUrl)
        .then(() => {
          this.$notify({
            title: 'Copied',
            text: 'URL copied to clipboard',
            type: 'success'
          });
        })
        .catch(err => {
          console.error('Failed to copy URL:', err);
        });
    }
  }
};
</script>

<style scoped>
.reservation-container {
  max-width: 800px;
  margin: 0 auto;
}
</style>