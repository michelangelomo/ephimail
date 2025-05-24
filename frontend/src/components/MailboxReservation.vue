<template>
  <div class="reservation-container">
    <!-- Show this when no email is selected -->
    <div v-if="!email" class="neo-card info-card">
      <div class="info-header">
        <i class="fas fa-info-circle"></i>
        <h5>Mailbox Features</h5>
      </div>
      <div class="info-content">
        <p>Select an email address to access these features:</p>
        <div class="feature-list">
          <div class="feature-point">
            <i class="fas fa-lock"></i>
            <span>Reserve a mailbox for private use</span>
          </div>
          <div class="feature-point">
            <i class="fas fa-shield-alt"></i>
            <span>Enable end-to-end encryption</span>
          </div>
          <div class="feature-point">
            <i class="fas fa-clock"></i>
            <span>Automatically delete old emails</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Reservation form when email is selected but not reserved -->
    <div v-else-if="!reserved" class="reservation-form">
      <div class="form-section">
        <h6 class="section-title">
          <i class="fas fa-calendar-alt"></i>
          Reservation Duration
        </h6>
        <div class="duration-options">
          <label v-for="option in durationOptions" :key="option.value" 
                 class="duration-option" :class="{ 'selected': duration === option.value }">
            <input v-model="duration" :value="option.value" type="radio" name="duration">
            <div class="option-content">
              <div class="option-icon">
                <i :class="option.icon"></i>
              </div>
              <div class="option-text">
                <span class="option-label">{{ option.label }}</span>
                <span class="option-desc">{{ option.description }}</span>
              </div>
            </div>
          </label>
        </div>
      </div>
      
      <div class="form-section">
        <div class="encryption-toggle">
          <label class="neo-checkbox enhanced">
            <input v-model="useEncryption" type="checkbox">
            <span class="neo-checkbox-mark"></span>
            <div class="checkbox-content">
              <span class="checkbox-title">
                <i class="fas fa-lock"></i>
                Enable End-to-End Encryption
              </span>
              <span class="checkbox-desc">Secure your emails with client-side encryption</span>
            </div>
          </label>
        </div>
        
        <transition name="info-slide">
          <div v-if="useEncryption" class="encryption-info">
            <div class="info-box">
              <div class="info-icon">
                <i class="fas fa-exclamation-triangle"></i>
              </div>
              <div class="info-text">
                <h6>About End-to-End Encryption</h6>
                <ul>
                  <li>All emails are encrypted with a unique key only you have</li>
                  <li>Your private key is stored in the URL - bookmark it!</li>
                  <li><strong>If you lose the URL, you cannot recover your emails</strong></li>
                </ul>
              </div>
            </div>
          </div>
        </transition>
      </div>
      
      <div class="form-actions">
        <button @click="reserveMailbox" 
                class="neo-btn neo-btn-primary reserve-btn" 
                :disabled="isLoading">
          <span v-if="isLoading" class="btn-spinner">
            <div class="neo-spinner d-inline-block"></div>
            Processing...
          </span>
          <span v-else class="btn-content">
            <i class="fas fa-shield-alt"></i>
            Reserve Mailbox
          </span>
        </button>
      </div>
    </div>
    
    <!-- Reserved mailbox status -->
    <div v-else class="reservation-status">
      <div class="status-header">
        <div class="status-icon success">
          <i class="fas fa-check-circle"></i>
        </div>
        <div class="status-text">
          <h6>Mailbox Reserved</h6>
          <p>Active until <strong>{{ formatDate(expiresAt) }}</strong></p>
        </div>
      </div>
      
      <div class="status-details">
        <div class="detail-item">
          <span class="detail-label">Status</span>
          <span class="detail-value success">
            <i class="fas fa-shield-alt"></i>
            Reserved
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">Encryption</span>
          <span class="detail-value" :class="encrypted ? 'success' : 'neutral'">
            <i :class="encrypted ? 'fas fa-lock' : 'fas fa-unlock'"></i>
            {{ encrypted ? 'Enabled' : 'Disabled' }}
          </span>
        </div>
        <div class="detail-item">
          <span class="detail-label">Expires</span>
          <span class="detail-value">{{ timeUntilExpiry }}</span>
        </div>
      </div>
      
      <div v-if="encrypted" class="encryption-warning">
        <div class="warning-header">
          <i class="fas fa-key"></i>
          <span>Important Security Information</span>
        </div>
        <p>Your unique access URL contains your private key:</p>
        <div class="url-display">
          <input type="text" class="url-input" :value="reservationUrl" readonly>
          <button class="neo-btn copy-btn" @click="copyUrl">
            <i class="fas fa-copy"></i>
          </button>
        </div>
        <div class="security-reminder">
          <i class="fas fa-exclamation-triangle"></i>
          <span>Bookmark this URL immediately! It's the only way to decrypt your emails.</span>
        </div>
      </div>
      
      <div class="status-actions">
        <button @click="extendReservation" class="neo-btn neo-btn-secondary">
          <i class="fas fa-clock"></i>
          Extend Time
        </button>
        <button @click="releaseReservation" class="neo-btn neo-btn-danger">
          <i class="fas fa-trash"></i>
          Release
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
      required: false
    }
  },
  
  emits: ['reservationChanged'],
  
  data() {
    return {
      duration: '24h',
      useEncryption: false,
      isLoading: false,
      reserved: false,
      encrypted: false,
      expiresAt: null,
      reservationUrl: '',
      privateKey: null,
      durationOptions: [
        {
          value: '1h',
          label: '1 Hour',
          description: 'Quick temporary use',
          icon: 'fas fa-clock'
        },
        {
          value: '24h',
          label: '1 Day',
          description: 'Most popular choice',
          icon: 'fas fa-calendar-day'
        },
        {
          value: '168h',
          label: '1 Week',
          description: 'Extended protection',
          icon: 'fas fa-calendar-week'
        }
      ]
    };
  },
  
  computed: {
    timeUntilExpiry() {
      if (!this.expiresAt) return '';
      
      const now = new Date();
      const expires = new Date(this.expiresAt);
      const diff = expires - now;
      
      if (diff <= 0) return 'Expired';
      
      const hours = Math.floor(diff / (1000 * 60 * 60));
      const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
      
      if (hours > 24) {
        const days = Math.floor(hours / 24);
        return `${days} day${days > 1 ? 's' : ''} remaining`;
      } else if (hours > 0) {
        return `${hours}h ${minutes}m remaining`;
      } else {
        return `${minutes}m remaining`;
      }
    }
  },
  
  watch: {
    email: {
      immediate: true,
      handler() {
        if (this.email) {
          this.checkReservation();
        }
      }
    }
  },
  
  mounted() {
    this.privateKey = EncryptionService.getPrivateKeyFromUrl();
    
    // Update expiry time every minute
    this.expiryInterval = setInterval(() => {
      this.$forceUpdate();
    }, 60000);
  },
  
  beforeUnmount() {
    if (this.expiryInterval) {
      clearInterval(this.expiryInterval);
    }
  },
  
  methods: {
    async reserveMailbox() {
      this.isLoading = true;
      
      try {
        let publicKey = '';
        
        if (this.useEncryption) {
          const keyPair = await EncryptionService.generateKeyPair();
          publicKey = keyPair.publicKey;
          this.privateKey = keyPair.privateKey;
        }
        
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
        
        this.reserved = true;
        this.encrypted = data.encrypted;
        this.expiresAt = new Date(data.expires_at);
        
        if (this.encrypted && this.privateKey) {
          EncryptionService.setPrivateKeyInUrl(this.privateKey);
          this.reservationUrl = window.location.href;
          
          this.$notify({
            title: 'Encryption Enabled',
            text: 'Bookmark this URL to keep access to your encrypted emails!',
            type: 'warning',
            duration: 10000
          });
        }
        
        this.$emit('reservationChanged', {
          reserved: true,
          encrypted: this.encrypted,
          expiresAt: this.expiresAt
        });
        
        this.$notify({
          title: 'Mailbox Reserved',
          text: `Your mailbox is now reserved for ${this.getDurationLabel(this.duration)}`,
          type: 'success'
        });
        
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
      if (!this.email) return;
      
      try {
        const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/inbox/${this.email}/reservation`);
        
        if (response.ok) {
          const data = await response.json();
          
          this.reserved = true;
          this.encrypted = data.encrypted;
          this.expiresAt = new Date(data.expires_at);
          
          this.$emit('reservationChanged', {
            reserved: true,
            encrypted: this.encrypted,
            expiresAt: this.expiresAt
          });
          
          if (this.encrypted && !this.privateKey) {
            this.$notify({
              title: 'Encryption Key Missing',
              text: 'This mailbox is encrypted but your access key is missing.',
              type: 'error',
              duration: 0
            });
          }
        }
      } catch (error) {
        console.error('Error checking reservation:', error);
      }
    },
    
    async extendReservation() {
      // Implementation for extending reservation
      this.$notify({
        title: 'Feature Coming Soon',
        text: 'Reservation extension will be available in a future update.',
        type: 'info'
      });
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
        const url = new URL(window.location);
        url.hash = '';
        window.history.replaceState({}, '', url);
        
        this.$emit('reservationChanged', {
          reserved: false,
          encrypted: false,
          expiresAt: null
        });
        
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
    
    getDurationLabel(duration) {
      const option = this.durationOptions.find(opt => opt.value === duration);
      return option ? option.label.toLowerCase() : duration;
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
          this.$notify({
            title: 'Copy Failed',
            text: 'Could not copy URL to clipboard',
            type: 'error'
          });
        });
    }
  }
};
</script>

<style scoped>
.reservation-container {
  width: 100%;
}

/* Info Card */
.info-card {
  border: none;
  background: linear-gradient(135deg, #f8fafc 0%, #e2e8f0 100%);
  box-shadow: none;
}

.info-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  color: var(--text-light);
}

.info-header i {
  font-size: 1.2rem;
}

.info-content p {
  margin-bottom: 1rem;
  color: var(--text-light);
}

.feature-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.feature-point {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: var(--text);
  font-size: 0.9rem;
}

.feature-point i {
  width: 20px;
  color: var(--primary);
}

/* Reservation Form */
.reservation-form {
  padding: 1.5rem;
}

.form-section {
  margin-bottom: 2rem;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  color: var(--text);
  font-weight: 600;
  font-size: 0.95rem;
}

/* Duration Options */
.duration-options {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.duration-option {
  position: relative;
  cursor: pointer;
  display: block;
}

.duration-option input {
  position: absolute;
  opacity: 0;
  pointer-events: none;
}

.option-content {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border: 2px solid var(--neutral);
  border-radius: 8px;
  background: white;
  transition: all 0.3s ease;
}

.duration-option.selected .option-content {
  border-color: var(--primary);
  background: rgba(248, 113, 113, 0.05);
  box-shadow: 0 4px 12px rgba(248, 113, 113, 0.15);
}

.duration-option:hover .option-content {
  border-color: var(--primary);
}

.option-icon {
  width: 40px;
  height: 40px;
  background: var(--neutral);
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text);
  transition: all 0.3s ease;
}

.duration-option.selected .option-icon {
  background: var(--primary);
  color: white;
}

.option-text {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.option-label {
  font-weight: 600;
  color: var(--text);
  margin-bottom: 0.25rem;
}

.option-desc {
  font-size: 0.85rem;
  color: var(--text-light);
}

/* Enhanced Checkbox */
.neo-checkbox.enhanced {
  padding-left: 0;
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1rem;
  border: 2px solid var(--neutral);
  border-radius: 8px;
  background: white;
  transition: all 0.3s ease;
  margin-bottom: 0;
}

.neo-checkbox.enhanced:hover {
  border-color: var(--primary);
}

.neo-checkbox.enhanced input:checked + .neo-checkbox-mark + .checkbox-content {
  opacity: 1;
}

.neo-checkbox.enhanced .neo-checkbox-mark {
  position: static;
  margin: 0;
  flex-shrink: 0;
}

.checkbox-content {
  display: flex;
  flex-direction: column;
  opacity: 0.7;
  transition: opacity 0.3s ease;
}

.checkbox-title {
  font-weight: 600;
  color: var(--text);
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.25rem;
}

.checkbox-desc {
  font-size: 0.85rem;
  color: var(--text-light);
  line-height: 1.3;
}

/* Encryption Info */
.encryption-info {
  margin-top: 1rem;
}

.info-box {
  display: flex;
  gap: 1rem;
  padding: 1rem;
  background: rgba(251, 191, 36, 0.1);
  border: 2px solid var(--accent);
  border-radius: 8px;
}

.info-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--accent-dark);
  flex-shrink: 0;
}

.info-text h6 {
  margin: 0 0 0.5rem 0;
  color: var(--text);
  font-size: 0.9rem;
}

.info-text ul {
  margin: 0;
  padding-left: 1.2rem;
  font-size: 0.85rem;
  color: var(--text-light);
  line-height: 1.4;
}

.info-text li {
  margin-bottom: 0.25rem;
}

/* Form Actions */
.form-actions {
  padding-top: 1rem;
  border-top: 1px solid var(--neutral);
}

.reserve-btn {
  width: 100%;
  padding: 0.875rem 1.5rem;
  font-weight: 600;
  position: relative;
}

.btn-spinner,
.btn-content {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

/* Reservation Status */
.reservation-status {
  padding: 1.5rem;
}

.status-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: rgba(52, 211, 153, 0.1);
  border-radius: 8px;
  border: 2px solid var(--success);
}

.status-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
}

.status-icon.success {
  background: var(--success);
  color: white;
}

.status-text h6 {
  margin: 0 0 0.25rem 0;
  color: var(--text);
  font-weight: 600;
}

.status-text p {
  margin: 0;
  color: var(--text-light);
  font-size: 0.9rem;
}

.status-details {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1.5rem;
  padding: 1rem;
  background: var(--background);
  border-radius: 8px;
}

.detail-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.detail-label {
  font-size: 0.9rem;
  color: var(--text-light);
}

.detail-value {
  font-weight: 600;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
}

.detail-value.success {
  color: var(--success-dark);
}

.detail-value.neutral {
  color: var(--text-light);
}

/* Encryption Warning */
.encryption-warning {
  background: rgba(251, 191, 36, 0.1);
  border: 2px solid var(--accent);
  border-radius: 8px;
  padding: 1rem;
  margin-bottom: 1.5rem;
}

.warning-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
  font-weight: 600;
  color: var(--accent-dark);
}

.encryption-warning p {
  margin: 0 0 0.75rem 0;
  font-size: 0.9rem;
  color: var(--text);
}

.url-display {
  display: flex;
  gap: 0.5rem;
  margin-bottom: 0.75rem;
}

.url-input {
  flex: 1;
  padding: 0.5rem;
  border: 2px solid var(--neutral);
  border-radius: 6px;
  font-size: 0.8rem;
  font-family: monospace;
  background: white;
}

.copy-btn {
  padding: 0.5rem 0.75rem;
  font-size: 0.8rem;
}

.security-reminder {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.85rem;
  color: var(--accent-dark);
  font-weight: 500;
}

/* Status Actions */
.status-actions {
  display: flex;
  gap: 0.75rem;
}

.status-actions .neo-btn {
  flex: 1;
  padding: 0.75rem 1rem;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

/* Transitions */
.info-slide-enter-active,
.info-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.info-slide-enter-from,
.info-slide-leave-to {
  opacity: 0;
  transform: translateY(-10px);
  max-height: 0;
}

.info-slide-enter-to,
.info-slide-leave-from {
  opacity: 1;
  transform: translateY(0);
  max-height: 200px;
}

/* Responsive Design */
@media (max-width: 768px) {
  .reservation-form,
  .reservation-status {
    padding: 1rem;
  }
  
  .status-actions {
    flex-direction: column;
  }
  
  .duration-options {
    gap: 0.5rem;
  }
  
  .option-content {
    padding: 0.75rem;
  }
  
  .option-icon {
    width: 32px;
    height: 32px;
  }
}

@media (max-width: 480px) {
  .status-header {
    flex-direction: column;
    text-align: center;
    gap: 0.75rem;
  }
  
  .detail-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.25rem;
  }
  
  .url-display {
    flex-direction: column;
  }
}
</style>