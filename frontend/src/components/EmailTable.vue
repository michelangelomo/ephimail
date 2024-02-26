<template>
  <div v-if="passedEmail != null">
    <div class="container p-3 email-container">
          <div class="border shadow rounded bg-white">
              <div class="row justify-content-center">
                  <div class="col-lg-5 d-flex">
                      <div class="email-list-wrapper p-2 pt-4">
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
                              <li class="list-group-item email" v-for="email in emails" :key="email.message_id" @click="this.showEmail(email)" v-bind:class="{ 'email-active shadow-sm rounded': email.message_id == this.emailSelectedId }">
                                  <div class="row">
                                      <div class="col-sm-9">
                                          <div>
                                              <p class="mb-1 text-singleline">{{ email.from[0].name }} ({{ email.from[0].email }})</p>
                                              <p class="text-muted mb-0 text-singleline">{{ email.body.substring(0, 25) }}</p>
                                          </div>
                                      </div>
                                      <div class="col-sm-3">
                                          <div>
                                              <p class="mb-1 text-muted">{{ email.orig_date.hour }}:{{ email.orig_date.minute }}</p>
                                              <p class="text-muted mb-0 text-singleline d-flex justify-content-end"><font-awesome-icon icon="fa-solid fa-download" /></p>
                                          </div>
                                      </div>
                                  </div>
                              </li>
                          </ul>
                      </div>
                  </div>
                  <div class="col-lg-7 col-7 d-flex">
                      <div class="p-2 pt-4">
                          <p aria-hidden="true" class="mb-1" ref="emailDate" v-bind:class="{ 'placeholder': !emailSelected }">01/01/1970 13:37</p>
                          <h5 class="mb-1" v-bind:class="{ 'placeholder': !emailSelected }" aria-hidden="true" ref="emailSubject">Subject</h5>
                          <div class="mt-0 mb-3">
                            <span class="d-inline" v-bind:class="{ 'placeholder': !emailSelected }">From: </span><p class="d-inline" v-bind:class="{ 'placeholder': !emailSelected }" aria-hidden="true" ref="emailFrom">example@example.com</p>
                          </div>

                          <div style="font-family: Arial, sans-serif;" class="mb-5" v-bind:class="{ 'placeholder': !emailSelected }" aria-hidden="true" ref="emailBody">
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
// import { ContentLoader } from "vue-content-loader"
import { Email } from 'typescript-email-parser';

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
    }
  },
  // components: { ContentLoader },
  watch: {
    passedEmail: {
      deep: true,
      handler: function(){
        this.loadEmails()
      }
    }
  },
  methods: {
    async loadEmails() {
      this.emails = [];
      this.isLoading = true;
      this.emailSelected = false;
      this.emailSelectedId = '';
      try {
          const response = await fetch("http://localhost:8000/inbox/{email}".replace('{email}', this.passedEmail));
          const data = await response.json();
          for(let d in data){
            let e = Email.parse(data[d])
            if(!e.message_id) {
              e.message_id = d.split(":")[1]
            }
            this.emails.push(e)            
          }
          this.isLoading = false
          if(this.emails.length > 0) {
            this.emailSelectedId = this.emails[0].message_id
            this.showEmail(this.emails[0])
          }
      } catch (error) {
          console.error("Failed to fetch select options:", error);
          this.isLoading = false
      }
    },
    showEmail(email) {
      this.emailSelected = true;
      this.emailSelectedId = email.message_id;

      let d = `${email.orig_date.day}/${email.orig_date.month}/${email.orig_date.year} ${email.orig_date.hour}:${email.orig_date.minute}`
      this.$refs.emailDate.innerText = d;
      this.$refs.emailFrom.innerText = email.from[0].name + " (" + email.from[0].email + ")"
      this.$refs.emailFrom.class = ""
      this.$refs.emailSubject.innerText = email.subject
      this.$refs.emailSubject.class = ""
      this.$refs.emailBody.innerText = email.body
      this.$refs.emailBody.class = ""
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
        .email-active{
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
