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
        
        <!-- Enhanced Features list -->
        <div class="features-list">
          <div class="feature-item" v-for="(feature, index) in features" :key="index" 
               :style="{ 'animation-delay': `${index * 0.1}s` }">
            <div class="feature-icon">
              <i :class="feature.icon"></i>
            </div>
            <div class="feature-content">
              <span class="feature-title">{{ feature.title }}</span>
              <span class="feature-desc">{{ feature.description }}</span>
            </div>
          </div>
        </div>
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
                Generate a disposable email address instantly. No registration required.
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
      features: [
        {
          icon: 'fas fa-shield-alt',
          title: 'End-to-end encryption',
          description: 'Your emails are encrypted client-side'
        },
        {
          icon: 'fas fa-clock',
          title: 'Auto-expire emails',
          description: 'Messages automatically delete after 24h'
        },
        {
          icon: 'fas fa-user-secret',
          title: 'Anonymous & private',
          description: 'No registration or personal data required'
        },
        {
          icon: 'fas fa-zap',
          title: 'Instant delivery',
          description: 'Real-time email notifications'
        }
      ],
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
      
      // Clear URL hash
      window.location.hash = '';
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
@import '../assets/css/NeoBrutalism.css';

.app-container {
  display: flex;
  min-height: 100vh;
  background: linear-gradient(135deg, var(--background) 0%, #f8fafc 100%);
  overflow-x: hidden;
}

/* Enhanced Sidebar */
.sidebar {
  width: 340px;
  background: linear-gradient(180deg, var(--primary) 0%, var(--primary-dark) 100%);
  padding: 2rem;
  display: flex;
  flex-direction: column;
  box-shadow: 8px 0 32px rgba(0, 0, 0, 0.15);
  position: relative;
  overflow: hidden;
  z-index: 10;
}

.sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23ffffff' fill-opacity='0.08'%3E%3Ccircle cx='30' cy='30' r='4'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E") repeat;
  pointer-events: none;
}

.logo-section {
  position: relative;
  z-index: 1;
}

.logo-wrapper {
  position: relative;
  text-align: center;
  margin-bottom: 1.5rem;
  animation: float 6s ease-in-out infinite;
}

.logo {
  width: 90px;
  height: 90px;
  border-radius: 24px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.25);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 2;
}

.logo-glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 110px;
  height: 110px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.3) 0%, transparent 70%);
  border-radius: 50%;
  animation: pulse 3s ease-in-out infinite;
  z-index: 1;
}

.logo:hover {
  transform: scale(1.05) rotate(5deg);
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.3);
}

.brand-title {
  color: white;
  font-size: 2.8rem;
  font-weight: 700;
  text-align: center;
  margin-bottom: 0.5rem;
  text-shadow: 2px 2px 8px rgba(0, 0, 0, 0.3);
  letter-spacing: -0.02em;
}

.brand-subtitle {
  color: rgba(255, 255, 255, 0.9);
  text-align: center;
  font-size: 1rem;
  margin-bottom: 3rem;
  font-weight: 400;
  line-height: 1.5;
  text-shadow: 1px 1px 3px rgba(0, 0, 0, 0.2);
}

.features-list {
  margin-top: auto;
}

.feature-item {
  display: flex;
  align-items: flex-start;
  color: rgba(255, 255, 255, 0.95);
  margin-bottom: 1.5rem;
  padding: 1rem;
  border-radius: 12px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  opacity: 0;
  transform: translateX(-20px);
  animation: slideInLeft 0.6s ease-out forwards;
  backdrop-filter: blur(10px);
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.15);
  transform: translateX(8px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}

.feature-icon {
  margin-right: 1rem;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  flex-shrink: 0;
}

.feature-content {
  display: flex;
  flex-direction: column;
}

.feature-title {
  font-weight: 600;
  font-size: 1rem;
  margin-bottom: 0.25rem;
}

.feature-desc {
  font-size: 0.85rem;
  opacity: 0.8;
  line-height: 1.3;
}

.stats-section {
  background: rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 1rem;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.stat-item {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 0.8rem;
  color: rgba(255, 255, 255, 0.7);
  margin-bottom: 0.25rem;
}

/* Enhanced Main Content */
.main-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  position: relative;
  background: var(--background);
}

.email-form-container {
  transition: all 1s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  z-index: 5;
}

.email-form-container.centered {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 2rem;
}

.email-form-container.minimized {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(20px);
  border-bottom: 3px solid var(--text);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
  padding: 0;
}

.email-form-container.animated {
  transform: translateY(-20px);
}

/* Email Content Area */
.email-content-area {
  flex: 1;
  padding: 2rem;
  background: var(--background);
}

.email-interface {
  height: 100%;
  margin: 0 auto;
}

/* Settings Modal */
.settings-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 2rem;
  backdrop-filter: blur(4px);
}

.settings-modal {
  background: white;
  border: 3px solid var(--text);
  border-radius: 12px;
  box-shadow: 8px 8px 0 var(--shadow-color);
  max-width: 600px;
  width: 100%;
  max-height: 80vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.settings-header {
  background: var(--secondary);
  color: white;
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 3px solid var(--text);
}

.settings-header h3 {
  margin: 0;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.close-btn {
  background: none;
  border: none;
  color: white;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 4px;
  transition: all 0.3s ease;
  font-size: 1rem;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.2);
}

.settings-content {
  flex: 1;
  overflow-y: auto;
}

/* Hero Section */
.form-wrapper-large {
  max-width: 700px;
  width: 100%;
  text-align: center;
}

.hero-section {
  margin-bottom: 3rem;
}

.form-title {
  font-size: 3rem;
  font-weight: 700;
  color: var(--text);
  margin-bottom: 1rem;
  letter-spacing: -0.03em;
  line-height: 1.1;
}

.title-accent {
  display: inline-block;
  animation: bounce 2s infinite;
}

.form-subtitle {
  font-size: 1.2rem;
  color: var(--text-light);
  margin-bottom: 0;
  font-weight: 400;
  line-height: 1.6;
}

.form-card {
  background: white;
  border-radius: 16px;
  padding: 2rem;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.1);
  border: 3px solid var(--text);
  margin-bottom: 3rem;
  position: relative;
  overflow: hidden;
}

.form-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, var(--primary), var(--secondary), var(--success));
}

/* How it works section */
.how-it-works {
  text-align: left;
  margin-top: 4rem;
}

.how-it-works h3 {
  text-align: center;
  font-size: 1.5rem;
  margin-bottom: 2rem;
  color: var(--text);
}

.steps {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.step {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.5rem;
  background: white;
  border-radius: 12px;
  border: 2px solid var(--neutral);
  transition: all 0.3s ease;
}

.step:hover {
  border-color: var(--primary);
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

.step-number {
  width: 32px;
  height: 32px;
  background: var(--primary);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 700;
  flex-shrink: 0;
}

.step-content h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1rem;
  font-weight: 600;
}

.step-content p {
  margin: 0;
  font-size: 0.9rem;
  color: var(--text-light);
  line-height: 1.4;
}

/* Minimized Form */
.form-wrapper-small {
  padding: 1.5rem 2rem;
}

.form-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
}

.current-mailbox-preview {
  flex: 1;
}

.current-mailbox-preview .preview-label {
  font-size: 0.85rem;
  color: var(--text-light);
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.current-mailbox-preview .preview-email {
  font-size: 1.1rem;
  font-weight: 600;
  font-family: monospace;
  color: var(--text);
}

.current-mailbox-preview .preview-username {
  color: var(--primary);
}

.current-mailbox-preview .preview-at {
  color: var(--text-light);
  margin: 0 0.25rem;
}

.current-mailbox-preview .preview-domain {
  color: var(--secondary);
}

.form-actions {
  display: flex;
  align-items: center;
  flex: 1;
  justify-content: center;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
  flex-shrink: 0;
}

.settings-btn,
.new-btn {
  white-space: nowrap;
  padding: 0.5rem 1rem;
  font-size: 0.9rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

/* Email Content Area */
.email-content-area {
  flex: 1;
  padding: 2rem;
  width: 100%;
}

/* Enhanced Reservation Panel */
.reservation-panel {
  background: white;
  border: 3px solid var(--text);
  border-radius: 12px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: all 0.3s ease;
}

.reservation-panel.expanded {
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.15);
}

.panel-header {
  background: var(--secondary);
  color: white;
  padding: 1rem 1.5rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: all 0.3s ease;
}

.panel-header:hover {
  background: var(--secondary-dark);
}

.panel-header h3 {
  margin: 0;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.toggle-icon {
  transition: transform 0.3s ease;
}

.toggle-icon.rotated {
  transform: rotate(180deg);
}

.panel-content {
  padding: 0;
}

.email-interface {
  flex: 1;
  min-height: 600px;
}

/* Loading Overlay */
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(5px);
}

.loading-content {
  text-align: center;
  color: white;
}

.loading-content p {
  margin-top: 1rem;
  font-size: 1.1rem;
  font-weight: 500;
}

.neo-spinner-large {
  width: 60px;
  height: 60px;
  border: 4px solid rgba(255, 255, 255, 0.3);
  border-right-color: white;
  border-radius: 50%;
  animation: spinner 1s linear infinite;
}

/* Animations */
@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-12px); }
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.8; transform: translate(-50%, -50%) scale(1.1); }
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-10px); }
  60% { transform: translateY(-5px); }
}

@keyframes slideInLeft {
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes spinner {
  to { transform: rotate(360deg); }
}

/* Animations */
@keyframes float {
  0%, 100% { transform: translateY(0px); }
  50% { transform: translateY(-12px); }
}

@keyframes pulse {
  0%, 100% { opacity: 0.5; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.8; transform: translate(-50%, -50%) scale(1.1); }
}

@keyframes bounce {
  0%, 20%, 50%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-10px); }
  60% { transform: translateY(-5px); }
}

@keyframes slideInLeft {
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes spinner {
  to { transform: rotate(360deg); }
}

/* Transitions */
.form-transition-enter-active,
.form-transition-leave-active {
  transition: all 0.6s cubic-bezier(0.4, 0, 0.2, 1);
}

.form-transition-enter-from {
  opacity: 0;
  transform: translateY(30px) scale(0.95);
}

.form-transition-leave-to {
  opacity: 0;
  transform: translateY(-30px) scale(1.05);
}

.content-slide-enter-active {
  transition: all 0.8s cubic-bezier(0.4, 0, 0.2, 1);
  transition-delay: 0.4s;
}

.content-slide-enter-from {
  opacity: 0;
  transform: translateY(60px);
}

.panel-slide-enter-active,
.panel-slide-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.panel-slide-enter-from,
.panel-slide-leave-to {
  opacity: 0;
  transform: translateY(-20px);
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Responsive Design */
@media (max-width: 1200px) {
  .sidebar {
    width: 300px;
  }
  
  .form-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .form-actions {
    width: 100%;
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .app-container {
    flex-direction: column;
  }
  
  .sidebar {
    width: 100%;
    height: auto;
    padding: 1.5rem;
  }
  
  .features-list {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
    margin: 2rem 0;
  }
  
  .stats-section {
    margin-top: 1rem;
  }
  
  .form-title {
    font-size: 2.2rem;
  }
  
  .steps {
    grid-template-columns: 1fr;
  }
  
  .email-content-area {
    padding: 1rem;
  }
  
  .form-header {
    padding: 1rem;
  }
  
  .current-mailbox {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.5rem;
  }
}

@media (max-width: 480px) {
  .sidebar {
    padding: 1rem;
  }
  
  .brand-title {
    font-size: 2.2rem;
  }
  
  .form-title {
    font-size: 1.8rem;
  }
  
  .form-card {
    padding: 1.5rem;
  }
  
  .features-list {
    grid-template-columns: 1fr;
  }
  
  .feature-item {
    padding: 0.75rem;
  }
}
</style>