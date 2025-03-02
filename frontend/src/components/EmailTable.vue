<template>
  <div v-if="passedEmail != null">
    <div class="container p-3 email-container">
      <div class="border shadow rounded bg-white">
        <div class="row justify-content-center">
          <div class="col-lg-5 d-flex">
            <div class="email-list-wrapper p-2 pt-4">
              <div v-if="encryptionError" class="alert alert-danger">
                <strong>Decryption Error:</strong> {{ encryptionError }}
                <p>Make sure you have the correct private key in your URL.</p>
              </div>
              
              <ul class="list-group emails">
                <li class="list-group-item email" v-if="emails.length == 0">
                  <div class="row">
                    <div class="col-sm-9">
                      <div>
                        <p class="mb-1 text-singleline placeholder">Example (example@example.com)</p>
                        <p class="placeholder mb-0 text-singleline">aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa</p>
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
                        <p class="mb-1 text-singleline">
                          {{ email.from[0].name }} ({{ email.from[0].email }})
                          <i v-if="email.encrypted" class="fas fa-lock text-success ml-1" title="Encrypted"></i>
                        </p>
                        <p class="text-muted mb-0 text-singleline">{{ email.body.substring(0, 25) }}</p>
                      </div>
                    </div>
                    <div class="col-sm-3">
                      <div>
                        <p class="mb-1 text-muted">{{ email.orig_date.hour }}:{{ email.orig_date.minute }}</p>
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
          <div class="col-lg-7 col-7 d-flex">
            <div class="p-2 pt-4">
              <p aria-hidden="true" class="mb-1" ref="emailDate" v-bind:class="{ 'placeholder': !emailSelected }">
                01/01/1970 13:37
              </p>
              <h5 class="mb-1" v-bind:class="{ 'placeholder': !emailSelected }" aria-hidden="true" ref="emailSubject">
                Subject
              </h5>
              <div class="mt-0 mb-3">
                <span class="d-inline" v-bind:class="{ 'placeholder': !emailSelected }">From: </span>
                <p class="d-inline" v-bind:class="{ 'placeholder': !emailSelected }" aria-hidden="true" ref="emailFrom">
                  example@example.com
                </p>
                <span v-if="currentEmailEncrypted" class="badge bg-success ml-2">
                  <i class="fas fa-lock"></i> Encrypted
                </span>
              </div>

              <div v-if="decryptionInProgress" class="text-center p-4">
                <div class="spinner-border text-primary" role="status">
                  <span class="visually-hidden">Decrypting...</span>
                </div>
                <p class="mt-2">Decrypting email...</p>
              </div>
              
              <div v-else style="font-family: Arial, sans-serif;" 
                   class="mb-5" 
                   v-bind:class="{ 'placeholder': !emailSelected }" 
                   aria-hidden="true" 
                   ref="emailBody">
                &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;
              </div>
            </div>
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
      this.emails = [];
      this.isLoading = true;
      this.emailSelected = false;
      this.emailSelectedId = '';
      this.encryptionError = null;
      
      try {
        const response = await fetch(`${process.env.VUE_APP_BACKEND_URL}/inbox/${this.passedEmail}`);
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

      // Format the date
      let d = `${email.orig_date.day}/${email.orig_date.month}/${email.orig_date.year} ${email.orig_date.hour}:${email.orig_date.minute}`;
      this.$refs.emailDate.innerText = d;
      
      // Set the from and subject
      this.$refs.emailFrom.innerText = email.from[0].name + " (" + email.from[0].email + ")";
      this.$refs.emailFrom.class = "";
      this.$refs.emailSubject.innerText = email.subject;
      this.$refs.emailSubject.class = "";
      
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
          const decryptedEmail = Email.parse(decryptedContent);
          
          // Update the displayed email content
          this.$refs.emailBody.innerText = decryptedEmail.body || "No content";
          this.$refs.emailBody.class = "";
          
          // Update subject if it was encrypted
          if (decryptedEmail.subject) {
            this.$refs.emailSubject.innerText = decryptedEmail.subject;
          }
          
          this.encryptionError = null;
        } catch (error) {
          console.error("Decryption error:", error);
          this.encryptionError = error.message;
          this.$refs.emailBody.innerText = "⚠️ This email is encrypted and could not be decrypted. Please check your private key.";
          this.$refs.emailBody.class = "";
        } finally {
          this.decryptionInProgress = false;
        }
      } else if (email.encrypted) {
        // Email is encrypted but no private key
        this.$refs.emailBody.innerText = "🔒 This email is encrypted. You need the private key to view it.";
        this.$refs.emailBody.class = "";
      } else {
        // Regular unencrypted email
        this.$refs.emailBody.innerText = email.body;
        this.$refs.emailBody.class = "";
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

<style>
.email-list-wrapper {
  background-color: #f2f4f6;
}
.emails {
  border: none !important;
}
.email {
  background-color: transparent !important;
  border: none !important;
}
.email-active {
  background-color: #fff !important;
}
.text-singleline {
  text-overflow: ellipsis; 
  overflow: hidden; 
  white-space: nowrap;
}
.email-container {
  min-width: 800px !important;
}
</style>