<template>
  <div class="email-interface">
    <!-- No mailbox selected state -->
    <div v-if="!passedEmail" class="no-mailbox-state">
      <div class="empty-state-content">
        <div class="empty-state-icon">
          <span>📬</span>
        </div>
        <h3>No mailbox selected</h3>
        <p>Enter an email address above to view your inbox</p>
        <div class="quick-tips">
          <div class="tip-item">
            <span>💡</span>
            <span>Tip: Any email address works instantly</span>
          </div>
          <div class="tip-item">
            <span>🛡️</span>
            <span>Reserve your mailbox for privacy</span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Main email interface -->
    <div v-else class="email-container">
      <!-- Email list sidebar -->
      <div class="email-list-section">
        <div class="list-header">
          <div class="header-content">
            <h4>
              <span>📬</span>
              Inbox
            </h4>
            <div class="email-count" v-if="emails.length > 0">
              {{ emails.length }} email{{ emails.length !== 1 ? 's' : '' }}
            </div>
          </div>
          <div class="header-actions">
            <button @click="refreshEmails" class="neo-btn neo-btn-sm refresh-btn" :disabled="isLoading">
              <span :class="{ 'spinning': isLoading }">🔄</span>
            </button>
          </div>
        </div>
        
        <!-- Encryption error -->
        <div v-if="encryptionError" class="neo-alert neo-alert-danger">
          <span>⚠️</span>
          <div>
            <strong>Decryption Error:</strong>
            <p>{{ encryptionError }}</p>
            <small>Make sure you have the correct private key in your URL.</small>
          </div>
        </div>
        
        <!-- Email list -->
        <div class="email-list-wrapper">
          <div v-if="isLoading && emails.length === 0" class="loading-state">
            <div class="neo-spinner"></div>
            <p>Loading emails...</p>
          </div>
          
          <div v-else-if="emails.length === 0" class="empty-inbox">
            <div class="empty-inbox-icon">
              <span>📭</span>
            </div>
            <h5>No emails yet</h5>
            <p>Emails sent to <strong>{{ passedEmail }}</strong> will appear here</p>
            <div class="waiting-indicator">
              <div class="pulse-dot"></div>
              <span>Waiting for incoming emails...</span>
            </div>
          </div>
          
          <div v-else class="emails-list">
            <div 
              v-for="email in emails" 
              :key="email.message_id" 
              class="email-item"
              :class="{ 
                'active': email.message_id === emailSelectedId,
                'encrypted': email.encrypted,
                'unread': email.unread 
              }"
              @click="showEmail(email)">
              
              <div class="email-header">
                <div class="sender-info">
                  <div class="sender-avatar">
                    {{ getSenderInitials(email.from) }}
                  </div>
                  <div class="sender-details">
                    <div class="sender-name">
                      {{ getSenderName(email.from) }}
                    </div>
                    <div class="sender-email">
                      {{ getSenderEmail(email.from) }}
                    </div>
                  </div>
                </div>
                <div class="email-meta">
                  <div class="email-time">
                    {{ formatEmailTime(email.orig_date) }}
                  </div>
                  <div class="email-status">
                    <span v-if="email.encrypted" title="Encrypted">🔒</span>
                    <span v-if="email.unread" class="unread-indicator" title="Unread">●</span>
                  </div>
                </div>
              </div>
              
              <div class="email-preview">
                <div class="email-subject">
                  {{ email.subject || 'No Subject' }}
                </div>
                <div class="email-snippet">
                  {{ getEmailSnippet(email) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Email content area -->
      <div class="email-content-section">
        <div v-if="!emailSelected" class="no-email-selected">
          <div class="no-selection-content">
            <div class="no-selection-icon">
              <span>👆</span>
            </div>
            <h4>Select an email to view</h4>
            <p>Choose an email from the list to read its contents</p>
          </div>
        </div>
        
        <div v-else class="email-viewer">
          <!-- Email header -->
          <div class="email-content-header">
            <div class="email-title-section">
              <h3 class="email-subject" ref="emailSubject">
                {{ selectedEmailData.subject || 'No Subject' }}
              </h3>
              <div class="email-badges">
                <span v-if="currentEmailEncrypted" class="neo-badge neo-badge-success">
                  🔒
                  Encrypted
                </span>
                <span v-if="selectedEmailData.unread" class="neo-badge neo-badge-primary">
                  New
                </span>
              </div>
            </div>
            
            <div class="email-metadata">
              <div class="metadata-row">
                <span class="metadata-label">From:</span>
                <span class="metadata-value" ref="emailFrom">
                  {{ getSenderName(selectedEmailData.from) }} &lt;{{ getSenderEmail(selectedEmailData.from) }}&gt;
                </span>
              </div>
              <div class="metadata-row">
                <span class="metadata-label">Date:</span>
                <span class="metadata-value" ref="emailDate">
                  {{ formatFullDate(selectedEmailData.orig_date) }}
                </span>
              </div>
              <div class="metadata-row">
                <span class="metadata-label">To:</span>
                <span class="metadata-value">{{ passedEmail }}</span>
              </div>
            </div>
            
            <div class="email-toolbar">
              <!--<button @click="deleteEmail(selectedEmailData)" class="neo-btn neo-btn-sm neo-btn-danger">
                <span>🗑️</span>
                Delete
              </button>-->
              <button @click="downloadEmail" class="neo-btn neo-btn-sm">
                <span>📥</span>
                Download
              </button>
            </div>
          </div>
          
          <!-- Email content -->
          <div class="email-content-body">
            <div v-if="decryptionInProgress" class="decryption-state">
              <div class="neo-spinner"></div>
              <p>Decrypting email...</p>
            </div>
            
            <div v-else class="email-body-content" ref="emailBody">
              <div v-if="emailContentHtml" class="html-content" v-html="emailContentText"></div>
              <div v-else class="text-content">{{ emailContentText }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
  
<script>
import PostalMime from 'postal-mime';
import EncryptionService from '../services/encryption';
import { webSocketService } from '../services/websocket';

export default {
  props: {
    'passedEmail': String,
  },
  
  data() {
    return {
      isLoading: false,
      emails: [],
      emailSelected: false,
      emailSelectedId: '',
      selectedEmailData: {},
      privateKey: null,
      encryptionError: null,
      decryptionInProgress: false,
      currentEmailEncrypted: false,
      emailContentHtml: false,
      emailContentText: '',
      rawEmails: {},
    }
  },
  
  watch: {
    passedEmail: {
      deep: true,
      handler: function(){
        this.loadEmails()
      }
    }
  },
  
  mounted() {
    this.privateKey = EncryptionService.getPrivateKeyFromUrl();
    
    if (this.passedEmail) {
      this.loadEmails();
      this.setupWebSocket();
    }
  },
  
  beforeUnmount() {
    this.cleanupWebSocket();
  },
  
  methods: {
    /**
     * Parse email address from header value
     * Handles formats like: "Name <email@domain.com>" or "email@domain.com"
     * @param {string} headerValue - The raw from header value
     * @returns {object} Object with name and email properties
     */
    parseFromHeader(headerValue) {
      if (!headerValue || typeof headerValue !== 'string') {
        return { name: "Unknown Sender", email: "unknown@example.com" };
      }
      
      // Remove any extra whitespace
      headerValue = headerValue.trim();
      
      // Check for format: "Name <email@domain.com>"
      const nameEmailMatch = headerValue.match(/^(.+?)\s*<(.+?)>$/);
      if (nameEmailMatch) {
        const name = nameEmailMatch[1].trim().replace(/^["']|["']$/g, ''); // Remove quotes if present
        const email = nameEmailMatch[2].trim();
        return { name, email };
      }
      
      // Check if it's just an email address
      const emailOnlyMatch = headerValue.match(/^[^\s@]+@[^\s@]+\.[^\s@]+$/);
      if (emailOnlyMatch) {
        return { name: headerValue, email: headerValue };
      }
      
      // Fallback
      return { name: headerValue, email: "unknown@example.com" };
    },

    /**
     * Extract from information from email headers
     * @param {array} headers - Array of header objects
     * @returns {array} Array with from information in PostalMime format
     */
    extractFromHeaders(headers) {
      if (!headers || !Array.isArray(headers)) {
        return [{ name: "Unknown Sender", email: "unknown@example.com" }];
      }
      
      const fromHeader = headers.find(header => 
        header.key && header.key.toLowerCase() === 'from'
      );
      
      if (!fromHeader || !fromHeader.value) {
        return [{ name: "Unknown Sender", email: "unknown@example.com" }];
      }
      
      const parsed = this.parseFromHeader(fromHeader.value);
      return [parsed];
    },

    async loadEmails(silent = false) {
      if (!this.passedEmail) {
        this.emails = [];
        this.isLoading = false;
        return;
      }
      
      if (!silent) {
        this.emails = [];
        this.isLoading = true;
        this.emailSelected = false;
        this.emailSelectedId = '';
        this.encryptionError = null;
      }
      
      try {
        const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/inbox/${this.passedEmail}`);
        if (!response.ok) {
          throw new Error(`Failed to fetch emails: ${response.statusText}`);
        }
        
        const data = await response.json();
        this.rawEmails = data;
        
        const parsedEmails = [];
        
        for(let d in data){
          let emailContent = data[d];
          let isEncrypted = EncryptionService.isEncrypted(emailContent);
          
          try {
            let e = await PostalMime.parse(emailContent);
            console.log("Parsed email:", e);
            
            // Use message_id from PostalMime if available, otherwise extract from key
            if(!e.message_id) {
              e.message_id = d.split(":")[1];
            }
            
            // Fix the from field parsing issue
            if (!e.from || !Array.isArray(e.from) || e.from.length === 0 || 
                (e.from[0].name === "Unknown Sender" && e.from[0].email === "unknown@example.com")) {
              // Extract from information from headers instead
              e.from = this.extractFromHeaders(e.headers);
            }
            
            e.orig_date = e.date || new Date();
            
            // PostalMime uses 'text' for the main content, 'body' as fallback
            // Priority: text > body > default message
            e.content = e.text || e.body || "No content";
            
            e.subject = e.subject || "No Subject";
            e.encrypted = isEncrypted;
            e.unread = true; // Mark as unread initially
            
            parsedEmails.push(e);
          } catch (parseError) {
            console.error("Failed to parse email:", parseError);
            
            // Even for unparseable emails, try to extract from header if available
            let fromInfo = [{ name: "Unknown Sender", email: "unknown@example.com" }];
            try {
              // Try to extract headers from the raw email content
              const headerMatch = emailContent.match(/^From:\s*(.+?)$/m);
              if (headerMatch) {
                fromInfo = [this.parseFromHeader(headerMatch[1])];
              }
            } catch (headerError) {
              console.warn("Could not extract from header from unparseable email");
            }
            
            parsedEmails.push({
              message_id: d.split(":")[1],
              from: fromInfo,
              subject: "Unparseable Email",
              content: "This email couldn't be parsed properly.",
              text: "This email couldn't be parsed properly.",
              orig_date: {
                day: new Date().getDate(),
                month: new Date().getMonth() + 1,
                year: new Date().getFullYear(),
                hour: new Date().getHours(),
                minute: new Date().getMinutes()
              },
              encrypted: isEncrypted,
              unread: true
            });
          }
        }
        
        // Sort emails by date (newest first)
        parsedEmails.sort((a, b) => {
          const getDateValue = (origDate) => {
            if (origDate instanceof Date) return origDate;
            if (typeof origDate === 'object' && origDate.year) {
              return new Date(origDate.year, (origDate.month || 1) - 1, origDate.day || 1, origDate.hour || 0, origDate.minute || 0);
            }
            return new Date(origDate);
          };
          
          const dateA = getDateValue(a.orig_date);
          const dateB = getDateValue(b.orig_date);
          return dateB - dateA;
        });
              
        this.emails = parsedEmails;
        this.isLoading = false;
        
        if(this.emails.length > 0 && !this.emailSelected) {
          this.emailSelectedId = this.emails[0].message_id;
          this.showEmail(this.emails[0]);
        }
      } catch (error) {
        console.error("Failed to fetch emails:", error);
        this.isLoading = false;
        this.$notify({
          title: "Error",
          text: "Failed to load emails: " + error.message,
          type: "error"
        });
      }
    },
    
    async showEmail(email) {
      this.emailSelected = true;
      this.emailSelectedId = email.message_id;
      this.selectedEmailData = email;
      this.currentEmailEncrypted = email.encrypted;
      this.emailContentHtml = false;
      this.emailContentText = '';
      
      // Mark email as read
      email.unread = false;

      if (email.encrypted && this.privateKey) {
        try {
          this.decryptionInProgress = true;
          
          const emailKey = Object.keys(this.rawEmails).find(key => key.includes(email.message_id));
          if (!emailKey) {
            throw new Error("Email content not found");
          }
          
          const encryptedContent = this.rawEmails[emailKey];
          const decryptedContent = await EncryptionService.decryptEmail(encryptedContent, this.privateKey);
          
          let decryptedEmail;
          try {
            decryptedEmail = await PostalMime.parse(decryptedContent);
          } catch (parseError) {
            console.error("Error parsing decrypted email:", parseError);
            decryptedEmail = {
              text: decryptedContent,
              subject: email.subject || "Decrypted Email"
            };
          }
          
          // Use text field first, then body, then fallback
          this.emailContentText = decryptedEmail.text || decryptedEmail.body || "No content";
          this.emailContentHtml = this.isHtmlContent(this.emailContentText);
          
          if (decryptedEmail.subject) {
            this.selectedEmailData.subject = decryptedEmail.subject;
          }
          
          this.encryptionError = null;
        } catch (error) {
          console.error("Decryption error:", error);
          this.encryptionError = error.message;
          this.emailContentText = "⚠️ This email is encrypted and could not be decrypted. Please check your private key.";
        } finally {
          this.decryptionInProgress = false;
        }
      } else if (email.encrypted) {
        this.emailContentText = "🔒 This email is encrypted. You need the private key to view it.";
      } else {
        // Use the content field we set during parsing (text > body > default)
        this.emailContentText = email.content || "No content";
        this.emailContentHtml = this.isHtmlContent(this.emailContentText);
      }
    },
    
    isHtmlContent(content) {
      return /<[a-z][\s\S]*>/i.test(content);
    },
    
    getSenderName(fromArray) {
      if (!fromArray || !Array.isArray(fromArray) || fromArray.length === 0) {
        return 'Unknown Sender';
      }
      return fromArray[0].name || fromArray[0].email || 'Unknown Sender';
    },
    
    getSenderEmail(fromArray) {
      if (!fromArray || !Array.isArray(fromArray) || fromArray.length === 0) {
        return 'unknown@example.com';
      }
      return fromArray[0].email || 'unknown@example.com';
    },
    
    getSenderInitials(fromArray) {
      const name = this.getSenderName(fromArray);
      const parts = name.split(' ');
      if (parts.length >= 2) {
        return (parts[0][0] + parts[1][0]).toUpperCase();
      }
      return name.substring(0, 2).toUpperCase();
    },
    
    getEmailSnippet(email) {
      // Use the same priority as in parsing: text > body > default
      const content = email.text || email.body || email.content || 'No content';
      
      // Remove HTML tags and get first 60 characters
      const textContent = content.replace(/<[^>]*>/g, '').trim();
      return textContent.length > 60 ? textContent.substring(0, 60) + '...' : textContent;
    },
    
    formatEmailTime(origDate) {
      if (!origDate) return '--:--';
      
      // Handle both Date objects and custom object format
      let emailDate;
      if (origDate instanceof Date) {
        emailDate = origDate;
      } else if (typeof origDate === 'object' && origDate.year) {
        // Handle custom object format
        emailDate = new Date(
          origDate.year,
          (origDate.month || 1) - 1, // Convert 1-based month to 0-based
          origDate.day || 1,
          origDate.hour || 0,
          origDate.minute || 0
        );
      } else {
        // Try to parse as string
        emailDate = new Date(origDate);
      }
      
      // Validate the date
      if (isNaN(emailDate.getTime())) {
        return '--:--';
      }
      
      const now = new Date();
      const diffMs = now - emailDate;
      const diffMins = Math.floor(diffMs / (1000 * 60));
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
      
      if (diffMins < 1) return 'Just now';
      if (diffMins < 60) return `${diffMins}m ago`;
      if (diffHours < 24) return `${diffHours}h ago`;
      if (diffDays < 7) return `${diffDays}d ago`;
      
      return `${emailDate.getDate()}/${emailDate.getMonth() + 1}`;
    },
    
    formatFullDate(origDate) {
      if (!origDate) return 'Date unavailable';
      
      // Handle both Date objects and custom object format
      let emailDate;
      if (origDate instanceof Date) {
        emailDate = origDate;
      } else if (typeof origDate === 'object' && origDate.year) {
        // Handle custom object format
        emailDate = new Date(
          origDate.year,
          (origDate.month || 1) - 1, // Convert 1-based month to 0-based
          origDate.day || 1,
          origDate.hour || 0,
          origDate.minute || 0
        );
      } else {
        // Try to parse as string
        emailDate = new Date(origDate);
      }
      
      // Validate the date
      if (isNaN(emailDate.getTime())) {
        return 'Date unavailable';
      }
      
      const day = String(emailDate.getDate()).padStart(2, '0');
      const month = String(emailDate.getMonth() + 1).padStart(2, '0');
      const year = emailDate.getFullYear();
      const hour = String(emailDate.getHours()).padStart(2, '0');
      const minute = String(emailDate.getMinutes()).padStart(2, '0');
      
      return `${day}/${month}/${year} ${hour}:${minute}`;
    },
    
    async refreshEmails() {
      await this.loadEmails();
      this.$notify({
        title: "Refreshed",
        text: "Email list has been updated",
        type: "success"
      });
    },
    
    async deleteEmail(email) {
      if (!confirm('Are you sure you want to delete this email?')) {
        return;
      }
      
      try {
        const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/api/emails/${email.message_id}`, {
          method: 'DELETE'
        });
        
        if (!response.ok) {
          throw new Error('Failed to delete email');
        }
        
        // Remove from local list
        this.emails = this.emails.filter(e => e.message_id !== email.message_id);
        
        // If this was the selected email, clear selection
        if (this.emailSelectedId === email.message_id) {
          this.emailSelected = false;
          this.emailSelectedId = '';
          this.selectedEmailData = {};
        }
        
        this.$notify({
          title: "Deleted",
          text: "Email has been deleted",
          type: "success"
        });
        
      } catch (error) {
        console.error('Error deleting email:', error);
        this.$notify({
          title: "Error",
          text: "Failed to delete email",
          type: "error"
        });
      }
    },
    
    downloadEmail() {
      const emailText = `
From: ${this.getSenderName(this.selectedEmailData.from)} <${this.getSenderEmail(this.selectedEmailData.from)}>
To: ${this.passedEmail}
Date: ${this.formatFullDate(this.selectedEmailData.orig_date)}
Subject: ${this.selectedEmailData.subject}

${this.emailContentText}
      `.trim();
      
      const blob = new Blob([emailText], { type: 'text/plain' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `email-${this.selectedEmailData.message_id}.txt`;
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    },

    setupWebSocket() {
      webSocketService.connect();
      webSocketService.subscribeToInbox(this.passedEmail);
      webSocketService.addEventListener('new_email', this.handleNewEmail);
      webSocketService.addEventListener('email_deleted', this.handleEmailDeleted);
    },

    cleanupWebSocket() {
      if (this.passedEmail) {
        webSocketService.unsubscribeFromInbox(this.passedEmail);
      }
      webSocketService.removeEventListener('new_email', this.handleNewEmail);
      webSocketService.removeEventListener('email_deleted', this.handleEmailDeleted);
    },

    handleNewEmail(data) {
      if (data.email === this.passedEmail) {
        this.loadEmails(true);
        this.$notify({
          title: "New Email",
          text: "You have received a new email.",
          type: "info"
        });
      }
    },

    handleEmailDeleted(data) {
      if (data.email === this.passedEmail) {
        this.loadEmails(true);
      }
    }
  }
};
</script>

<style scoped>
@import '@/assets/css/emailtable.css';
</style>