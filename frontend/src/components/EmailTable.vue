<template>
  <!-- Remove the v-if entirely so the component always renders -->
  <div class="neo-card mb-4">
    <div v-if="!passedEmail" class="p-4 text-center">
      <h5>No mailbox selected</h5>
      <p class="text-muted">Enter an email address above to view your inbox</p>
    </div>
    
    <div v-else class="row">
      <div class="col-lg-5">
        <div class="email-list-wrapper p-2 pt-4">
          <div v-if="encryptionError" class="neo-alert neo-alert-danger">
            <strong>Decryption Error:</strong> {{ encryptionError }}
            <p>Make sure you have the correct private key in your URL.</p>
          </div>
          
          <ul class="list-group emails">
            <li class="list-group-item email" v-if="emails.length == 0">
              <div class="row">
                <div class="col-sm-9">
                  <div>
                    <p class="mb-1 text-singleline placeholder">Example (example@example.com)</p>
                    <p class="placeholder mb-0 text-singleline">Waiting for incoming emails...</p>
                  </div>
                </div>
              </div>
            </li>
            <li class="list-group-item email" v-for="email in emails" :key="email.message_id" 
                @click="showEmail(email)" 
                v-bind:class="{ 'email-active shadow-sm rounded': email.message_id == this.emailSelectedId }">
              <div class="row">
                <div class="col-sm-9">
                  <div>
                    <p class="mb-1 text-singleline fw-bold">
                      {{ email.from && email.from[0] ? email.from[0].name || 'Unknown Sender' : 'Unknown Sender' }} 
                      ({{ email.from && email.from[0] ? email.from[0].email || 'unknown@example.com' : 'unknown@example.com' }})
                      <i v-if="email.encrypted" class="fas fa-lock text-success ml-1" title="Encrypted"></i>
                    </p>
                    <p class="text-muted mb-0 text-singleline">{{ email.body ? email.body.substring(0, 25) : 'No content' }}</p>
                  </div>
                </div>
                <div class="col-sm-3">
                  <div>
                    <p class="mb-1 text-muted">
                      {{ email.orig_date && email.orig_date.hour ? email.orig_date.hour : '--' }}:{{ email.orig_date && email.orig_date.minute ? email.orig_date.minute : '--' }}
                    </p>
                    <p class="text-muted mb-0 text-singleline d-flex justify-content-end">
                      <font-awesome-icon icon="fa-solid fa-download" />
                    </p>
                  </div>
                </div>
              </div>
            </li>
          </ul>
        </div>
      </div>
      
      <div class="col-lg-7">
        <div class="email-content p-2 pt-4">
          <p class="mb-1" ref="emailDate" v-bind:class="{ 'placeholder': !emailSelected }">
            01/01/1970 13:37
          </p>
          <h5 class="mb-1" v-bind:class="{ 'placeholder': !emailSelected }" ref="emailSubject">
            Subject
          </h5>
          <div class="mt-0 mb-3">
            <span class="d-inline" v-bind:class="{ 'placeholder': !emailSelected }">From: </span>
            <p class="d-inline" v-bind:class="{ 'placeholder': !emailSelected }" ref="emailFrom">
              example@example.com
            </p>
            <span v-if="currentEmailEncrypted" class="neo-badge neo-badge-success ml-2">
              <i class="fas fa-lock"></i> Encrypted
            </span>
          </div>

          <div v-if="decryptionInProgress" class="text-center p-4">
            <div class="neo-spinner mx-auto"></div>
            <p class="mt-2">Decrypting email...</p>
          </div>
          
          <div v-else style="font-family: 'Space Grotesk', sans-serif;" 
               class="mb-5 email-body" 
               v-bind:class="{ 'placeholder': !emailSelected }" 
               ref="emailBody">
            &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
  
<script>
import { Email } from 'typescript-email-parser';
import EncryptionService from '../services/encryption';
import { webSocketService } from '../services/websocket';

export default {
  props: {
    'passedEmail': String,
  },
  data() {
    return {
      isLoading: true,
      emails: [],
      emailSelected: false,
      emailSelectedId: '',
      privateKey: null,
      encryptionError: null,
      decryptionInProgress: false,
      currentEmailEncrypted: false,
      rawEmails: {}, // Store the raw emails for decryption
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
    // Check for private key in URL
    this.privateKey = EncryptionService.getPrivateKeyFromUrl();
    
    if (this.passedEmail) {
      this.loadEmails();
      this.setupWebSocket();
    }
  },
  methods: {
    async loadEmails() {
      if (!this.passedEmail) {
        this.emails = [];
        this.isLoading = false;
        return;
      }
      
      this.emails = [];
      this.isLoading = true;
      this.emailSelected = false;
      this.emailSelectedId = '';
      this.encryptionError = null;
      
      try {
        const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/inbox/${this.passedEmail}`);
        if (!response.ok) {
          throw new Error(`Failed to fetch emails: ${response.statusText}`);
        }
        
        const data = await response.json();
        this.rawEmails = data; // Store raw emails
        
        for(let d in data){
          let emailContent = data[d];
          let isEncrypted = EncryptionService.isEncrypted(emailContent);
          
          // Try to parse as is first
          try {
            let e = Email.parse(emailContent);
            if(!e.message_id) {
              e.message_id = d.split(":")[1];
            }
            
            // Ensure required fields exist
            if (!e.from || !Array.isArray(e.from) || e.from.length === 0) {
              e.from = [{ name: "Unknown Sender", email: "unknown@example.com" }];
            }
            
            if (!e.orig_date) {
              e.orig_date = {
                day: "--",
                month: "--",
                year: "----",
                hour: "--",
                minute: "--"
              };
            }
            
            if (!e.body) {
              e.body = "No content";
            }
            
            if (!e.subject) {
              e.subject = "No Subject";
            }
            
            // Flag if it's encrypted
            e.encrypted = isEncrypted;
            
            this.emails.push(e);
          } catch (parseError) {
            console.error("Failed to parse email:", parseError);
            
            // Add a placeholder for unparseable emails
            this.emails.push({
              message_id: d.split(":")[1],
              from: [{ name: "Unknown Sender", email: "unknown@example.com" }],
              subject: "Unparseable Email",
              body: "This email couldn't be parsed properly.",
              orig_date: {
                day: "--",
                month: "--",
                year: "----",
                hour: "--",
                minute: "--"
              },
              encrypted: isEncrypted
            });
          }
        }
        
        this.isLoading = false;
        if(this.emails.length > 0) {
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
      this.currentEmailEncrypted = email.encrypted;

      // Format the date with proper fallbacks
      let d = 'Date unavailable';
      if (email.orig_date) {
        const day = email.orig_date.day || '--';
        const month = email.orig_date.month || '--';
        const year = email.orig_date.year || '----';
        const hour = email.orig_date.hour || '--';
        const minute = email.orig_date.minute || '--';
        d = `${day}/${month}/${year} ${hour}:${minute}`;
      }
      
      if (this.$refs.emailDate) {
        this.$refs.emailDate.innerText = d;
      }
      
      // Set the from and subject with fallbacks
      if (this.$refs.emailFrom) {
        const sender = email.from && email.from[0] 
          ? `${email.from[0].name || 'Unknown'} (${email.from[0].email || 'unknown@example.com'})`
          : 'Unknown Sender';
        this.$refs.emailFrom.innerText = sender;
        this.$refs.emailFrom.class = "";
      }
      
      if (this.$refs.emailSubject) {
        this.$refs.emailSubject.innerText = email.subject || 'No Subject';
        this.$refs.emailSubject.class = "";
      }
      
      // Handle encrypted emails
      if (email.encrypted && this.privateKey) {
        try {
          this.decryptionInProgress = true;
          
          // Find the raw email content
          const emailKey = Object.keys(this.rawEmails).find(key => key.includes(email.message_id));
          if (!emailKey) {
            throw new Error("Email content not found");
          }
          
          const encryptedContent = this.rawEmails[emailKey];
          
          // Decrypt the email
          const decryptedContent = await EncryptionService.decryptEmail(encryptedContent, this.privateKey);
          
          // Parse the decrypted email
          let decryptedEmail;
          try {
            decryptedEmail = Email.parse(decryptedContent);
          } catch (parseError) {
            console.error("Error parsing decrypted email:", parseError);
            // Create a fallback email object
            decryptedEmail = {
              body: decryptedContent,
              subject: email.subject || "Decrypted Email"
            };
          }
          
          // Update the displayed email content safely
          if (this.$refs.emailBody) {
            this.$refs.emailBody.innerText = decryptedEmail.body || "No content";
            this.$refs.emailBody.class = "";
          }
          
          // Update subject if it was encrypted
          if (decryptedEmail.subject && this.$refs.emailSubject) {
            this.$refs.emailSubject.innerText = decryptedEmail.subject;
          }
          
          this.encryptionError = null;
        } catch (error) {
          console.error("Decryption error:", error);
          this.encryptionError = error.message;
          if (this.$refs.emailBody) {
            this.$refs.emailBody.innerText = "⚠️ This email is encrypted and could not be decrypted. Please check your private key.";
            this.$refs.emailBody.class = "";
          }
        } finally {
          this.decryptionInProgress = false;
        }
      } else if (email.encrypted) {
        // Email is encrypted but no private key
        if (this.$refs.emailBody) {
          this.$refs.emailBody.innerText = "🔒 This email is encrypted. You need the private key to view it.";
          this.$refs.emailBody.class = "";
        }
      } else {
        // Regular unencrypted email
        if (this.$refs.emailBody) {
          this.$refs.emailBody.innerText = email.body || "No content";
          this.$refs.emailBody.class = "";
        }
      }
    },

    setupWebSocket() {
      // Connect to WebSocket
      webSocketService.connect();
      
      // Subscribe to the inbox
      webSocketService.subscribeToInbox(this.passedEmail);
      
      // Listen for new emails
      webSocketService.addEventListener('new_email', this.handleNewEmail);
      
      // Listen for deleted emails
      webSocketService.addEventListener('email_deleted', this.handleEmailDeleted);
    },

    handleNewEmail(data) {
      if (data.email === this.passedEmail) {
        // Reload emails when a new one arrives
        this.loadEmails();
        
        // Notify the user
        this.$notify({
          title: "New Email",
          text: "You have received a new email.",
          type: "info"
        });
      }
    },

    handleEmailDeleted(data) {
      if (data.email === this.passedEmail) {
        // Reload emails when one is deleted
        this.loadEmails();
      }
    },

    // Add to the beforeDestroy lifecycle hook to clean up WebSocket
    beforeDestroy() {
      // Unsubscribe from the inbox
      if (this.passedEmail) {
        webSocketService.unsubscribeFromInbox(this.passedEmail);
      }
      
      // Remove event listeners
      webSocketService.removeEventListener('new_email', this.handleNewEmail);
      webSocketService.removeEventListener('email_deleted', this.handleEmailDeleted);
    }
  }
};
</script>

<style scoped>
.email-list-wrapper {
  background-color: var(--card-bg) !important;
  border-radius: 8px;
  height: 100%;
  padding: 1rem;
}

.emails {
  border: none !important;
  max-height: 400px;
  overflow-y: auto;
  padding-left: 0;
  list-style-type: none;
}

.email {
  padding: 1rem;
  margin-bottom: 0.5rem;
  transition: all 0.2s ease;
  border-radius: 6px !important;
  border: 2px solid transparent !important;
}

.email-active {
  background-color: var(--neutral) !important;
  border: 2px solid var(--text) !important;
  box-shadow: 4px 4px 0 var(--shadow-color);
}

.email:hover:not(.email-active) {
  background-color: rgba(0, 0, 0, 0.05) !important;
  cursor: pointer;
}

.email-content {
  background-color: var(--card-bg);
  min-height: 400px;
  border: 3px solid var(--text);
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 6px 6px 0 var(--shadow-color);
}

.email-body {
  line-height: 1.6;
  min-height: 300px;
}

.text-singleline {
  text-overflow: ellipsis; 
  overflow: hidden; 
  white-space: nowrap;
}

.neo-spinner {
  width: 40px;
  height: 40px;
  border-width: 4px;
}

/* Remove padding from card when containing emails */
:deep(.neo-card) {
  padding: 0;
}
</style>