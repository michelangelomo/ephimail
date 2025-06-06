<template>
    <form @submit.prevent="submitEmail" class="email-form" :class="{ 'compact': isCompact }">
        <div class="neo-input-group" :class="{ 'loading': isLoading }">
            <div class="input-wrapper">
                <input 
                    v-model="email" 
                    :placeholder="emailGenerated" 
                    type="text" 
                    class="neo-input email-input" 
                    :class="{ 'has-error': hasError }"
                    aria-label="Email username"
                    @input="handleInput"
                    @focus="handleFocus"
                    @blur="handleBlur">
            </div>
            
            <div class="separator">
                <span class="at-symbol">@</span>
            </div>
            
            <div class="select-wrapper">
                <select v-model="selectedOption" class="neo-select domain-select" @change="handleDomainChange">
                    <option v-for="option in options" :key="option" :value="option">
                        {{ option }}
                    </option>
                </select>
            </div>
            
            <button type="submit" 
                    class="neo-btn neo-btn-primary submit-btn" 
                    :disabled="!canSubmit"
                    :class="{ 'loading': isLoading }">
                <span v-if="isLoading" class="btn-loading">
                    <div class="neo-spinner"></div>
                </span>
                <span v-else class="btn-text">
                    <span v-if="isCompact">→</span>
                    <span v-else>Go! 👀</span>
                </span>
            </button>
        </div>
        
        <!-- Error message -->
        <transition name="error-slide">
            <div v-if="errorMessage" class="error-message">
                <span>⚠️</span>
                {{ errorMessage }}
            </div>
        </transition>
        
        <!-- Preview (only in large mode) -->
        <div v-if="!isCompact && (email || emailGenerated)" class="email-preview">
            <div class="preview-label">Your temporary email:</div>
            <div class="preview-container">
                <div class="preview-email">
                    <span class="preview-username">{{ email || emailGenerated }}</span>
                    <span class="preview-at">@</span>
                    <span class="preview-domain">{{ selectedOption }}</span>
                </div>
                <button type="button" @click="copyPreviewEmail" class="preview-copy-btn" title="Copy email address">
                    <span v-if="!previewCopied">📋</span>
                    <span v-else>✅</span>
                </button>
            </div>
        </div>
    </form>
</template>

<script>
import EncryptionService from '../services/encryption';
import configService from '../services/config';

export default {
    props: {
        isCompact: {
            type: Boolean,
            default: false
        }
    },
    
    data() {
        return {
            email: '',
            emailGenerated: '',
            options: [],
            selectedOption: '',
            isLoading: false,
            hasError: false,
            errorMessage: '',
            isFocused: false,
            backendUrl: '',
            previewCopied: false,
            reg: /^(([^<>()\]\\.,;:\s@"]+(\.[^<>()\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,24}))$/
        };
    },
    
    computed: {
        canSubmit() {
            return (this.options.length > 0) && !this.isLoading && !this.hasError;
        },
        
        fullEmail() {
            const username = this.email || this.emailGenerated;
            return `${username}@${this.selectedOption}`;
        }
    },
    
    async mounted() {
        // Get backend URL from config service
        this.backendUrl = configService.getBackendUrl();
        
        this.fetchSelectOptions();
        this.generateName();
        this.extractPathFromURL();
    },
    
    methods: {
        async copyPreviewEmail() {
            const emailToCopy = this.fullEmail;
            
            try {
                await navigator.clipboard.writeText(emailToCopy);
                this.previewCopied = true;
                
                this.$notify({
                    title: 'Copied!',
                    text: `Email address ${emailToCopy} copied to clipboard`,
                    type: 'success',
                    duration: 3000
                });
                
                // Reset the checkmark after 2 seconds
                setTimeout(() => {
                    this.previewCopied = false;
                }, 2000);
                
            } catch (error) {
                console.error('Failed to copy email address:', error);
                
                // Fallback for older browsers
                this.fallbackCopyTextToClipboard(emailToCopy);
            }
        },
        
        fallbackCopyTextToClipboard(text) {
            const textArea = document.createElement('textarea');
            textArea.value = text;
            textArea.style.top = '0';
            textArea.style.left = '0';
            textArea.style.position = 'fixed';
            
            document.body.appendChild(textArea);
            textArea.focus();
            textArea.select();
            
            try {
                const successful = document.execCommand('copy');
                if (successful) {
                    this.previewCopied = true;
                    this.$notify({
                        title: 'Copied!',
                        text: `Email address ${text} copied to clipboard`,
                        type: 'success',
                        duration: 3000
                    });
                    
                    setTimeout(() => {
                        this.previewCopied = false;
                    }, 2000);
                } else {
                    throw new Error('execCommand failed');
                }
            } catch (err) {
                this.$notify({
                    title: 'Copy Failed',
                    text: 'Could not copy email address to clipboard',
                    type: 'error',
                    duration: 5000
                });
            }
            
            document.body.removeChild(textArea);
        },
        
        extractPathFromURL() {
            // Use the new encryption service method
            const extractedEmail = EncryptionService.getEmailFromUrl();
            
            if (extractedEmail && this.isValidEmail(extractedEmail)) {
                const splitted = extractedEmail.split("@");
                
                // Check if this domain is in our options list
                if (this.options.length === 0) {
                    // Options haven't loaded yet, wait for them
                    this.$nextTick(() => {
                        if (this.options.includes(splitted[1])) {
                            this.selectedOption = splitted[1];
                            this.email = splitted[0];
                            this.$emit('emailSubmitted', extractedEmail);
                        }
                    });
                } else if (this.options.includes(splitted[1])) {
                    this.selectedOption = splitted[1];
                    this.email = splitted[0];
                    this.$emit('emailSubmitted', extractedEmail);
                }
            }
        },
        
        async submitEmail() {
            if (this.isLoading) return;
            
            this.clearError();
            this.isLoading = true;
            
            try {
                // Use generated name if no email provided
                if(this.email === ''){
                    this.email = this.emailGenerated
                }
                
                const completeEmail = this.email + "@" + this.selectedOption;
                
                // Validate email
                if(!this.isValidEmail(completeEmail)){
                    throw new Error('Please enter a valid email address');
                }
                
                // Check if there's already a private key in the URL (for encrypted mailboxes)
                const existingPrivateKey = EncryptionService.getPrivateKeyFromUrl();
                
                // Update URL - preserve existing private key if it exists
                EncryptionService.updateUrlWithEmail(completeEmail, existingPrivateKey);
                
                // Add a small delay for better UX
                await new Promise(resolve => setTimeout(resolve, 500));
                
                this.$emit('emailSubmitted', completeEmail);
                
            } catch (error) {
                this.showError(error.message);
            } finally {
                this.isLoading = false;
            }
        },
        
        async fetchSelectOptions() {
            try {
                const response = await fetch(`${this.backendUrl}/domains`);
                
                if (!response.ok) {
                    throw new Error('Failed to fetch domains');
                }
                
                const data = await response.json();
                this.options = data;
                
                if(data.length > 0) {
                    this.selectedOption = data[0];
                }
                
                // Re-check URL after options are loaded
                this.extractPathFromURL();
                
            } catch (error) {
                console.error("Failed to fetch select options:", error);
                this.showError("Could not load available domains. Please try again.");
                
                // Fallback domains
                this.options = ['example.com', 'example.org'];
                this.selectedOption = this.options[0];
            }
        },
        
        randomInt(min, max) {
            return Math.floor(Math.random() * (max - min)) + min;
        },
        
        generateName() {
            const name1 = ["swift","bright","cool","clever","smart","quick","zen","calm","bold","epic","neat","pure","safe","wise","free","happy","magic","super","fast","easy"];
            const name2 = ["mail","box","post","temp","test","demo","user","guest","inbox","secure","private","quick","instant","random","anon","safe","temp","trial","beta","dev"];
            
            const name = name1[this.randomInt(0, name1.length)] + '-' + name2[this.randomInt(0, name2.length)] + this.randomInt(10, 99);
            this.emailGenerated = name;
        },
        
        isValidEmail(email) {
            return this.reg.test(email);
        },
        
        handleInput() {
            this.clearError();
            
            // Real-time validation
            if (this.email && this.selectedOption) {
                const fullEmail = `${this.email}@${this.selectedOption}`;
                if (!this.isValidEmail(fullEmail)) {
                    this.hasError = true;
                } else {
                    this.hasError = false;
                }
            }
        },
        
        handleFocus() {
            this.isFocused = true;
            this.clearError();
        },
        
        handleBlur() {
            this.isFocused = false;
        },
        
        handleDomainChange() {
            this.clearError();
            if (this.email && this.selectedOption) {
                this.handleInput();
            }
        },
        
        showError(message) {
            this.errorMessage = message;
            this.hasError = true;
        },
        
        clearError() {
            this.errorMessage = '';
            this.hasError = false;
        }
    },
};
</script>

<style scoped>
@import '@/assets/css/emailform.css';

/* Preview container styles */
.preview-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 0.5rem;
}

.preview-copy-btn {
  background: rgba(255, 255, 255, 0.9);
  border: 2px solid var(--text);
  border-radius: 6px;
  padding: 0.5rem;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 40px;
  min-height: 40px;
  backdrop-filter: blur(10px);
  box-shadow: 2px 2px 0 var(--shadow-color);
  flex-shrink: 0;
}

.preview-copy-btn:hover {
  background: rgba(255, 255, 255, 1);
  transform: translateY(-1px);
  box-shadow: 3px 3px 0 var(--shadow-color);
}

.preview-copy-btn:active {
  transform: translateY(0);
  box-shadow: 1px 1px 0 var(--shadow-color);
}

.preview-copy-btn:focus {
  outline: 2px solid var(--primary);
  outline-offset: 2px;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .preview-container {
    flex-direction: column;
    gap: 0.5rem;
    align-items: flex-start;
  }
  
  .preview-copy-btn {
    align-self: center;
    min-width: 44px;
    min-height: 44px;
  }
}

@media (max-width: 480px) {
  .preview-copy-btn {
    padding: 0.75rem;
    font-size: 0.85rem;
  }
}
</style>