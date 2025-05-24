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
                <div class="input-icon" v-if="!isCompact">
                    <i class="fas fa-user"></i>
                </div>
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
                <div class="select-icon">
                    <i class="fas fa-chevron-down"></i>
                </div>
            </div>
            
            <button type="submit" 
                    class="neo-btn neo-btn-primary submit-btn" 
                    :disabled="!canSubmit"
                    :class="{ 'loading': isLoading }">
                <span v-if="isLoading" class="btn-loading">
                    <div class="neo-spinner"></div>
                </span>
                <span v-else class="btn-text">
                    <i class="fas fa-arrow-right" v-if="isCompact"></i>
                    <span v-else>Go! 👀</span>
                </span>
            </button>
        </div>
        
        <!-- Error message -->
        <transition name="error-slide">
            <div v-if="errorMessage" class="error-message">
                <i class="fas fa-exclamation-triangle"></i>
                {{ errorMessage }}
            </div>
        </transition>
        
        <!-- Preview (only in large mode) -->
        <div v-if="!isCompact && (email || emailGenerated)" class="email-preview">
            <div class="preview-label">Your temporary email:</div>
            <div class="preview-email">
                <span class="preview-username">{{ email || emailGenerated }}</span>
                <span class="preview-at">@</span>
                <span class="preview-domain">{{ selectedOption }}</span>
            </div>
        </div>
    </form>
</template>

<script>
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
    
    mounted() {
        this.fetchSelectOptions();
        this.generateName();
        this.extractPathFromURL();
    },
    
    methods: {
        extractPathFromURL() {
            const hashLocation = window.location.hash;
            let extractedPath = hashLocation.substring(1).replace("/", "");
            if(this.isValidEmail(extractedPath)) {
                let splitted = extractedPath.split("@");
                if(this.options.includes(splitted[1])) {
                    return
                }
                this.selectedOption = splitted[1];
                this.email = splitted[0];
                this.$emit('emailSubmitted', extractedPath);
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
                
                // Update URL
                window.location.hash = `/${completeEmail}`;
                
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
                const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/domains`);
                
                if (!response.ok) {
                    throw new Error('Failed to fetch domains');
                }
                
                const data = await response.json();
                this.options = data;
                
                if(data.length > 0) {
                    this.selectedOption = data[0];
                }
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
.email-form {
    width: 100%;
}

.neo-input-group {
    position: relative;
    margin-bottom: 1rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.neo-input-group.loading {
    opacity: 0.8;
}

.neo-input-group:not(.compact) {
    box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
    border-radius: 12px;
    overflow: hidden;
    background: white;
    border: 3px solid var(--text);
}

.email-form.compact .neo-input-group {
    border-radius: 8px;
    border: 2px solid var(--text);
    overflow: hidden;
    background: white;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

/* Input styling */
.input-wrapper {
    position: relative;
    flex: 1;
}

.email-input {
    border: none !important;
    box-shadow: none !important;
    border-radius: 0 !important;
    font-size: 1rem;
    font-weight: 500;
    transition: all 0.3s ease;
    padding-right: 3rem;
}

.email-form:not(.compact) .email-input {
    padding: 1rem 3rem 1rem 1.25rem;
    font-size: 1.1rem;
}

.email-form.compact .email-input {
    padding: 0.75rem 2.5rem 0.75rem 1rem;
    font-size: 0.95rem;
}

.email-input:focus {
    outline: none;
    background: rgba(248, 113, 113, 0.02);
}

.email-input.has-error {
    background: rgba(239, 68, 68, 0.05);
}

.email-input::placeholder {
    color: var(--text-light);
    font-style: italic;
}

.input-icon {
    position: absolute;
    right: 1rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-light);
    pointer-events: none;
    transition: all 0.3s ease;
}

.email-input:focus + .input-icon {
    color: var(--primary);
}

/* Separator */
.separator {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--primary);
    color: white;
    font-weight: 700;
    font-size: 1.1rem;
    padding: 0 1rem;
    min-width: 40px;
}

.email-form.compact .separator {
    padding: 0 0.75rem;
    font-size: 1rem;
    min-width: 32px;
}

/* Select styling */
.select-wrapper {
    position: relative;
    flex: 1;
    max-width: 180px;
}

.domain-select {
    border: none !important;
    box-shadow: none !important;
    border-radius: 0 !important;
    appearance: none;
    background: white;
    cursor: pointer;
    font-weight: 500;
    padding-right: 2.5rem !important;
}

.email-form:not(.compact) .domain-select {
    padding: 1rem 2.5rem 1rem 1.25rem !important;
    font-size: 1.1rem;
}

.email-form.compact .domain-select {
    padding: 0.75rem 2rem 0.75rem 1rem !important;
    font-size: 0.95rem;
}

.domain-select:focus {
    outline: none;
    background: rgba(248, 113, 113, 0.02);
}

.select-icon {
    position: absolute;
    right: 1rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-light);
    pointer-events: none;
    transition: all 0.3s ease;
    font-size: 0.8rem;
}

.domain-select:focus + .select-icon {
    color: var(--primary);
}

/* Submit button */
.submit-btn {
    border-radius: 0 !important;
    box-shadow: none !important;
    border: none !important;
    font-weight: 700;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    position: relative;
    overflow: hidden;
}

.email-form:not(.compact) .submit-btn {
    padding: 1rem 2rem;
    font-size: 1rem;
    min-width: 120px;
}

.email-form.compact .submit-btn {
    padding: 0.75rem 1.5rem;
    font-size: 0.9rem;
    min-width: 80px;
}

.submit-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.submit-btn:not(:disabled):hover {
    background: var(--primary-dark) !important;
    transform: translateY(-1px);
}

.submit-btn:not(:disabled):active {
    transform: translateY(0);
}

.btn-loading,
.btn-text {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    transition: all 0.3s ease;
}

.btn-loading .neo-spinner {
    width: 1rem;
    height: 1rem;
    border-width: 2px;
}

/* Error message */
.error-message {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--danger);
    font-size: 0.9rem;
    font-weight: 500;
    padding: 0.75rem 1rem;
    background: rgba(239, 68, 68, 0.1);
    border: 2px solid var(--danger);
    border-radius: 8px;
    margin-top: 1rem;
}

/* Suggestions */
.suggestions {
    margin-top: 2rem;
    padding: 1.5rem;
    background: white;
    border: 2px solid var(--neutral);
    border-radius: 12px;
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.05);
}

.suggestion-header {
    font-size: 0.9rem;
    color: var(--text-light);
    margin-bottom: 1rem;
    font-weight: 500;
}

.suggestion-list {
    display: grid;
    gap: 0.75rem;
}

.suggestion-item {
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    padding: 1rem;
    border: 2px solid var(--neutral);
    border-radius: 8px;
    background: white;
    cursor: pointer;
    transition: all 0.3s ease;
    text-align: left;
}

.suggestion-item:hover {
    border-color: var(--primary);
    background: rgba(248, 113, 113, 0.02);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.suggestion-email {
    font-weight: 600;
    color: var(--text);
    font-family: monospace;
    margin-bottom: 0.25rem;
}

.suggestion-desc {
    font-size: 0.85rem;
    color: var(--text-light);
}

/* Email preview */
.email-preview {
    margin-top: 1.5rem;
    padding: 1.25rem;
    background: linear-gradient(135deg, rgba(248, 113, 113, 0.05) 0%, rgba(59, 130, 246, 0.05) 100%);
    border: 2px solid var(--primary);
    border-radius: 12px;
    text-align: center;
}

.preview-label {
    font-size: 0.9rem;
    color: var(--text-light);
    margin-bottom: 0.75rem;
    font-weight: 500;
}

.preview-email {
    font-size: 1.2rem;
    font-weight: 600;
    font-family: monospace;
    color: var(--text);
    word-break: break-all;
}

.preview-username {
    color: var(--primary);
}

.preview-at {
    color: var(--text-light);
    margin: 0 0.25rem;
}

.preview-domain {
    color: var(--secondary);
}

/* Transitions */
.error-slide-enter-active,
.error-slide-leave-active {
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.error-slide-enter-from,
.error-slide-leave-to {
    opacity: 0;
    transform: translateY(-10px);
    max-height: 0;
    margin-top: 0;
    margin-bottom: 0;
    padding-top: 0;
    padding-bottom: 0;
}

/* Responsive Design */
@media (max-width: 768px) {
    .neo-input-group {
        flex-direction: column;
    }
    
    .separator {
        order: 1;
        padding: 0.5rem;
        justify-content: flex-start;
        min-width: auto;
    }
    
    .input-wrapper {
        order: 0;
    }
    
    .select-wrapper {
        order: 2;
        max-width: none;
    }
    
    .submit-btn {
        order: 3;
        border-radius: 0 0 8px 8px !important;
    }
    
    .suggestions {
        margin-top: 1rem;
        padding: 1rem;
    }
    
    .email-preview {
        margin-top: 1rem;
        padding: 1rem;
    }
    
    .preview-email {
        font-size: 1rem;
    }
}

@media (max-width: 480px) {
    .email-form:not(.compact) .email-input,
    .email-form:not(.compact) .domain-select {
        padding: 0.875rem 1rem;
        font-size: 1rem;
    }
    
    .email-form:not(.compact) .submit-btn {
        padding: 0.875rem 1.5rem;
        font-size: 0.95rem;
    }
    
    .suggestion-item {
        padding: 0.75rem;
    }
}

/* Focus states */
@media (prefers-reduced-motion: no-preference) {
    .email-input:focus,
    .domain-select:focus {
        animation: gentle-pulse 2s infinite;
    }
}

@keyframes gentle-pulse {
    0%, 100% { background: rgba(248, 113, 113, 0.02); }
    50% { background: rgba(248, 113, 113, 0.05); }
}
</style>