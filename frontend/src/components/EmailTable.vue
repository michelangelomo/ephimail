<template>
  <div class="email-interface">
    <!-- No mailbox selected state -->
    <div v-if="!passedEmail" class="no-mailbox-state">
      <div class="empty-state-content">
        <div class="empty-state-icon">
          <i class="fas fa-inbox"></i>
        </div>
        <h3>No mailbox selected</h3>
        <p>Enter an email address above to view your inbox</p>
        <div class="quick-tips">
          <div class="tip-item">
            <i class="fas fa-lightbulb"></i>
            <span>Tip: Any email address works instantly</span>
          </div>
          <div class="tip-item">
            <i class="fas fa-shield-alt"></i>
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
              <i class="fas fa-inbox"></i>
              Inbox
            </h4>
            <div class="email-count" v-if="emails.length > 0">
              {{ emails.length }} email{{ emails.length !== 1 ? 's' : '' }}
            </div>
          </div>
          <div class="header-actions">
            <button @click="refreshEmails" class="neo-btn neo-btn-sm refresh-btn" :disabled="isLoading">
              <i class="fas fa-sync-alt" :class="{ 'spinning': isLoading }"></i>
            </button>
          </div>
        </div>
        
        <!-- Encryption error -->
        <div v-if="encryptionError" class="neo-alert neo-alert-danger">
          <i class="fas fa-exclamation-triangle"></i>
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
              <i class="fas fa-envelope-open"></i>
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
                    <i v-if="email.encrypted" class="fas fa-lock" title="Encrypted"></i>
                    <i v-if="email.unread" class="fas fa-circle unread-indicator" title="Unread"></i>
                  </div>
                </div>
              </div>
              
              <div class="email-preview">
                <div class="email-subject">
                  {{ email.subject || 'No Subject' }}
                </div>
                <div class="email-snippet">
                  {{ getEmailSnippet(email.body) }}
                </div>
              </div>
              
              <div class="email-actions">
                <button @click.stop="deleteEmail(email)" class="action-btn delete-btn" title="Delete">
                  <i class="fas fa-trash"></i>
                </button>
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
              <i class="fas fa-mouse-pointer"></i>
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
                  <i class="fas fa-lock"></i>
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
              <button @click="deleteEmail(selectedEmailData)" class="neo-btn neo-btn-sm neo-btn-danger">
                <i class="fas fa-trash"></i>
                Delete
              </button>
              <button @click="downloadEmail" class="neo-btn neo-btn-sm">
                <i class="fas fa-download"></i>
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
              <div v-if="emailContentHtml" class="html-content"></div>
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
            if(!e.message_id) {
              e.message_id = d.split(":")[1];
            }
            
            // Ensure required fields exist
            e.from = e.from && Array.isArray(e.from) && e.from.length > 0 
              ? e.from 
              : [{ name: "Unknown Sender", email: "unknown@example.com" }];
            
            e.orig_date = e.orig_date || {
              day: new Date().getDate(),
              month: new Date().getMonth() + 1,
              year: new Date().getFullYear(),
              hour: new Date().getHours(),
              minute: new Date().getMinutes()
            };
            
            e.body = e.body || "No content";
            e.subject = e.subject || "No Subject";
            e.encrypted = isEncrypted;
            e.unread = true; // Mark as unread initially
            
            parsedEmails.push(e);
          } catch (parseError) {
            console.error("Failed to parse email:", parseError);
            
            parsedEmails.push({
              message_id: d.split(":")[1],
              from: [{ name: "Unknown Sender", email: "unknown@example.com" }],
              subject: "Unparseable Email",
              body: "This email couldn't be parsed properly.",
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
          const dateA = new Date(a.orig_date.year, a.orig_date.month - 1, a.orig_date.day, a.orig_date.hour, a.orig_date.minute);
          const dateB = new Date(b.orig_date.year, b.orig_date.month - 1, b.orig_date.day, b.orig_date.hour, b.orig_date.minute);
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
            decryptedEmail = PostalMime.parse(decryptedContent);
          } catch (parseError) {
            console.error("Error parsing decrypted email:", parseError);
            decryptedEmail = {
              body: decryptedContent,
              subject: email.subject || "Decrypted Email"
            };
          }
          
          this.emailContentText = decryptedEmail.body || "No content";
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
        this.emailContentText = email.body || "No content";
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
    
    getEmailSnippet(body) {
      if (!body) return 'No content';
      
      // Remove HTML tags and get first 60 characters
      const textContent = body.replace(/<[^>]*>/g, '').trim();
      return textContent.length > 60 ? textContent.substring(0, 60) + '...' : textContent;
    },
    
    formatEmailTime(origDate) {
      if (!origDate) return '--:--';
      
      const now = new Date();
      const emailDate = new Date(
        origDate.year || now.getFullYear(),
        (origDate.month || now.getMonth() + 1) - 1,
        origDate.day || now.getDate(),
        origDate.hour || 0,
        origDate.minute || 0
      );
      
      const diffMs = now - emailDate;
      const diffMins = Math.floor(diffMs / (1000 * 60));
      const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
      const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
      
      if (diffMins < 1) return 'Just now';
      if (diffMins < 60) return `${diffMins}m ago`;
      if (diffHours < 24) return `${diffHours}h ago`;
      if (diffDays < 7) return `${diffDays}d ago`;
      
      return `${origDate.day}/${origDate.month}`;
    },
    
    formatFullDate(origDate) {
      if (!origDate) return 'Date unavailable';
      
      const day = String(origDate.day || '--').padStart(2, '0');
      const month = String(origDate.month || '--').padStart(2, '0');
      const year = origDate.year || '----';
      const hour = String(origDate.hour || '--').padStart(2, '0');
      const minute = String(origDate.minute || '--').padStart(2, '0');
      
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
.email-interface {
  height: 100%;
  display: flex;
  flex-direction: column;
}

/* No mailbox state */
.no-mailbox-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 500px;
  padding: 2rem;
}

.empty-state-content {
  text-align: center;
  max-width: 400px;
}

.empty-state-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 1.5rem;
  background: var(--neutral-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  color: var(--text-light);
  border: 3px solid var(--text);
}

.empty-state-content h3 {
  color: var(--text);
  margin-bottom: 0.5rem;
}

.empty-state-content p {
  color: var(--text-light);
  margin-bottom: 2rem;
}

.quick-tips {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.tip-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 1rem;
  background: white;
  border: 2px solid var(--neutral);
  border-radius: 8px;
  font-size: 0.9rem;
  color: var(--text);
}

.tip-item i {
  color: var(--primary);
  width: 20px;
  text-align: center;
}

/* Main email container */
.email-container {
  display: grid;
  grid-template-columns: 400px 1fr;
  height: 100%;
  gap: 1.5rem;
  min-height: 600px;
}

/* Email list section */
.email-list-section {
  display: flex;
  flex-direction: column;
  background: white;
  border: 3px solid var(--text);
  border-radius: 12px;
  box-shadow: 6px 6px 0 var(--shadow-color);
  overflow: hidden;
}

.list-header {
  background: var(--primary);
  color: white;
  padding: 1rem 1.25rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 3px solid var(--text);
}

.header-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.header-content h4 {
  margin: 0;
  font-size: 1.1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.email-count {
  background: rgba(255, 255, 255, 0.2);
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8rem;
  font-weight: 600;
}

.refresh-btn {
  background: rgba(255, 255, 255, 0.2);
  border: 2px solid rgba(255, 255, 255, 0.3);
  color: white;
  padding: 0.5rem;
  width: 36px;
  height: 36px;
}

.refresh-btn:hover {
  background: rgba(255, 255, 255, 0.3);
}

.refresh-btn i.spinning {
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

.email-list-wrapper {
  flex: 1;
  overflow-y: auto;
  padding: 0.5rem;
}

/* Loading and empty states */
.loading-state,
.empty-inbox {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 2rem;
  text-align: center;
}

.loading-state p,
.empty-inbox p {
  color: var(--text-light);
  margin-top: 1rem;
}

.empty-inbox-icon {
  width: 60px;
  height: 60px;
  background: var(--neutral-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: var(--text-light);
  margin-bottom: 1rem;
}

.waiting-indicator {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-top: 1.5rem;
  color: var(--text-light);
  font-size: 0.9rem;
}

.pulse-dot {
  width: 8px;
  height: 8px;
  background: var(--primary);
  border-radius: 50%;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.4; transform: scale(1); }
  50% { opacity: 1; transform: scale(1.2); }
}

/* Email list items */
.emails-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.email-item {
  background: white;
  border: 2px solid var(--neutral);
  border-radius: 8px;
  padding: 1rem;
  cursor: pointer;
  transition: all 0.3s ease;
  position: relative;
}

.email-item:hover {
  border-color: var(--primary);
  background: rgba(248, 113, 113, 0.02);
  transform: translateX(4px);
}

.email-item.active {
  border-color: var(--primary);
  background: rgba(248, 113, 113, 0.1);
  box-shadow: 2px 2px 0 var(--shadow-light);
}

.email-item.encrypted::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  bottom: 0;
  width: 4px;
  background: var(--success);
  border-radius: 8px 0 0 8px;
}

.email-item.unread {
  border-width: 3px;
  border-color: var(--primary);
}

.email-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 0.75rem;
}

.sender-info {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  flex: 1;
}

.sender-avatar {
  width: 36px;
  height: 36px;
  background: var(--primary);
  color: white;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.8rem;
  font-weight: 700;
  flex-shrink: 0;
}

.sender-details {
  flex: 1;
  min-width: 0;
}

.sender-name {
  font-weight: 600;
  color: var(--text);
  font-size: 0.9rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sender-email {
  color: var(--text-light);
  font-size: 0.8rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.25rem;
}

.email-time {
  font-size: 0.8rem;
  color: var(--text-light);
  white-space: nowrap;
}

.email-status {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.email-status i {
  font-size: 0.7rem;
}

.email-status .fa-lock {
  color: var(--success);
}

.unread-indicator {
  color: var(--primary);
  font-size: 0.5rem;
}

.email-preview {
  margin-bottom: 0.5rem;
}

.email-subject {
  font-weight: 600;
  color: var(--text);
  font-size: 0.9rem;
  margin-bottom: 0.25rem;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.email-snippet {
  color: var(--text-light);
  font-size: 0.8rem;
  line-height: 1.3;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.email-actions {
  display: flex;
  justify-content: flex-end;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.email-item:hover .email-actions {
  opacity: 1;
}

.action-btn {
  background: none;
  border: none;
  padding: 0.5rem;
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.3s ease;
  color: var(--text-light);
}

.delete-btn:hover {
  background: var(--danger-light);
  color: var(--danger);
}

/* Email content section */
.email-content-section {
  background: white;
  border: 3px solid var(--text);
  border-radius: 12px;
  box-shadow: 6px 6px 0 var(--shadow-color);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.no-email-selected {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 400px;
}

.no-selection-content {
  text-align: center;
}

.no-selection-icon {
  width: 60px;
  height: 60px;
  background: var(--neutral-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.5rem;
  color: var(--text-light);
  margin: 0 auto 1rem;
}

.no-selection-content h4 {
  color: var(--text);
  margin-bottom: 0.5rem;
}

.no-selection-content p {
  color: var(--text-light);
}

.email-viewer {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.email-content-header {
  padding: 1.5rem;
  border-bottom: 3px solid var(--text);
  background: var(--background);
}

.email-title-section {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.email-subject {
  color: var(--text);
  font-size: 1.25rem;
  margin: 0;
  flex: 1;
  margin-right: 1rem;
}

.email-badges {
  display: flex;
  gap: 0.5rem;
  flex-shrink: 0;
}

.email-metadata {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.metadata-row {
  display: grid;
  grid-template-columns: 60px 1fr;
  gap: 1rem;
  font-size: 0.9rem;
}

.metadata-label {
  color: var(--text-light);
  font-weight: 600;
}

.metadata-value {
  color: var(--text);
  word-break: break-word;
}

.email-toolbar {
  display: flex;
  gap: 0.75rem;
}

.email-content-body {
  flex: 1;
  padding: 1.5rem;
  overflow-y: auto;
}

.decryption-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  text-align: center;
}

.decryption-state p {
  margin-top: 1rem;
  color: var(--text-light);
}

.email-body-content {
  line-height: 1.6;
  color: var(--text);
}

.html-content {
  font-family: inherit;
}

.html-content img {
  max-width: 100%;
  height: auto;
}

.text-content {
  white-space: pre-wrap;
  font-family: inherit;
}

/* Responsive design */
@media (max-width: 1024px) {
  .email-container {
    grid-template-columns: 1fr;
    grid-template-rows: 300px 1fr;
  }
  
  .email-list-section {
    order: 2;
  }
  
  .email-content-section {
    order: 1;
  }
}

@media (max-width: 768px) {
  .email-container {
    grid-template-rows: 250px 1fr;
    gap: 1rem;
  }
  
  .list-header {
    padding: 0.75rem 1rem;
  }
  
  .email-item {
    padding: 0.75rem;
  }
  
  .sender-avatar {
    width: 32px;
    height: 32px;
    font-size: 0.7rem;
  }
  
  .email-content-header {
    padding: 1rem;
  }
  
  .email-subject {
    font-size: 1.1rem;
  }
  
  .email-toolbar {
    flex-wrap: wrap;
  }
  
  .email-content-body {
    padding: 1rem;
  }
}

@media (max-width: 480px) {
  .tip-item {
    padding: 0.75rem;
    font-size: 0.85rem;
  }
  
  .email-item {
    padding: 0.5rem;
  }
  
  .sender-info {
    gap: 0.5rem;
  }
  
  .sender-name,
  .email-subject {
    font-size: 0.85rem;
  }
  
  .sender-email,
  .email-snippet {
    font-size: 0.75rem;
  }
  
  .email-content-header {
    padding: 0.75rem;
  }
  
  .email-toolbar .neo-btn {
    font-size: 0.8rem;
    padding: 0.5rem 0.75rem;
  }
}
</style>