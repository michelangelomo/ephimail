<template>
  <notifications />
  <div class="app-container">
    <!-- Enhanced Left Sidebar -->
    <div class="sidebar">
      <div class="logo-section">
        <div class="logo-wrapper">
          <img class="logo" src="@/assets/logo.png" alt="Ephimail logo">
          <div class="logo-glow"></div>
        </div>
        <h1 class="brand-title">ephimail</h1>
        <p class="brand-subtitle">privacy-focused disposable email service</p>
      </div>
    </div>

    <!-- Enhanced Main Content Area -->
    <div class="main-content">
      <!-- Dynamic Email Form Container -->
      <div class="email-form-container" 
           :class="{ 
             'minimized': isMinimized, 
             'centered': !isMinimized,
             'animated': isAnimating 
           }">
        <transition name="form-transition" mode="out-in">
          <!-- Large Form (Initial State) -->
          <div v-if="!isMinimized" class="form-wrapper-large">
            <div class="hero-section">
              <h2 class="form-title">
                Create Your Temporary Mailbox
                <span class="title-accent">✨</span>
              </h2>
              <p class="form-subtitle">
                Generate a disposable email address instantly. No logs and encryption available!
              </p>
            </div>
            <div class="form-card">
              <EmailForm @emailSubmitted="handleEmailSubmit" />
            </div>
            
            <!-- How it works section -->
            <div class="how-it-works">
              <h3>How it works</h3>
              <div class="steps">
                <div class="step" v-for="(step, index) in howItWorksSteps" :key="index">
                  <div class="step-number">{{ index + 1 }}</div>
                  <div class="step-content">
                    <h4>{{ step.title }}</h4>
                    <p>{{ step.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Minimized Form (Top Bar) -->
          <div v-else class="form-wrapper-small">
            <div class="form-header">
              <!-- Email preview when inbox is selected -->
              <div class="current-mailbox-preview" v-if="email">
                <div class="preview-label">Your temporary mailbox:</div>
                <div class="preview-email">
                  <span class="preview-username">{{ email.split('@')[0] }}</span>
                  <span class="preview-at">@</span>
                  <span class="preview-domain">{{ email.split('@')[1] }}</span>
                </div>
              </div>
              
              <div class="form-actions" v-if="!email">
                <EmailForm @emailSubmitted="handleEmailSubmit" />
              </div>
              
              <div class="header-actions">
                <button @click="toggleSettings" class="neo-btn neo-btn-secondary settings-btn">
                  <i class="fas fa-cog"></i>
                  Settings
                </button>
                <button @click="createNewMailbox" class="neo-btn neo-btn-primary new-btn">
                  <i class="fas fa-plus"></i>
                  New Mailbox
                </button>
              </div>
            </div>
          </div>
        </transition>
      </div>

      <!-- Enhanced Email Content Area -->
      <transition name="content-slide">
        <div v-if="email && isMinimized" class="email-content-area">
          <!-- Main Email Interface -->
          <div class="email-interface">
            <EmailTable :passedEmail="email" />
          </div>
        </div>
      </transition>
    </div>
    
    <!-- Settings Modal -->
    <transition name="modal-fade">
      <div v-if="showSettings" class="settings-overlay" @click="closeSettings">
        <div class="settings-modal" @click.stop>
          <div class="settings-header">
            <h3>
              <i class="fas fa-cog"></i>
              Mailbox Settings
            </h3>
            <button @click="closeSettings" class="close-btn">
              <i class="fas fa-times"></i>
            </button>
          </div>
          <div class="settings-content">
            <MailboxReservation :email="email" @reservationChanged="handleReservationChange" />
          </div>
        </div>
      </div>
    </transition>
    
    <!-- Loading overlay -->
    <transition name="fade">
      <div v-if="isAnimating" class="loading-overlay">
        <div class="loading-content">
          <div class="neo-spinner-large"></div>
          <p>Setting up your mailbox...</p>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
import EmailForm from '@/components/EmailForm.vue'
import EmailTable from '@/components/EmailTable.vue'
import MailboxReservation from '@/components/MailboxReservation.vue'

export default {
  name: 'EphimailApp',
  components: {
    EmailForm,
    EmailTable,
    MailboxReservation
  },
  data() {
    return {
      email: null,
      isMinimized: false,
      isAnimating: false,
      showSettings: false,
      isReserved: false,
      howItWorksSteps: [
        {
          title: 'Choose an address',
          description: 'Pick any username and domain combination'
        },
        {
          title: 'Start receiving',
          description: 'Your inbox is ready immediately'
        },
        {
          title: 'Stay secure',
          description: 'Optionally reserve and encrypt your mailbox'
        }
      ]
    };
  },
  methods: {
    async handleEmailSubmit(data) {
      this.isAnimating = true;
      this.email = data;
      
      // Smooth transition with proper timing
      setTimeout(() => {
        this.isMinimized = true;
        setTimeout(() => {
          this.isAnimating = false;
        }, 800);
      }, 500);
    },
    
    toggleSettings() {
      if (!this.email) return;
      this.showSettings = !this.showSettings;
    },
    
    closeSettings() {
      this.showSettings = false;
    },
    
    createNewMailbox() {
      // Reset to initial state
      this.email = null;
      this.isMinimized = false;
      this.showSettings = false;
      this.isReserved = false;
      
      // Clear URL hash properly without causing page reload
      history.replaceState(null, '', window.location.pathname + window.location.search);
    },
    
    handleReservationChange(reservationData) {
      this.isReserved = reservationData.reserved;
      this.showSettings = false; // Auto-close settings after action
    }
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@300;400;500;600;700&display=swap');
@import '@/assets/css/NeoBrutalism.css';
@import '@/assets/css/app.css';
</style>