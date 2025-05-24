<template>
  <notifications />
  <div class="container px-4 py-5 my-5">
    <div class="logo-container text-center">
      <img class="d-block mx-auto mb-4" src="@/assets/logo.png" alt="Ephimail logo" width="150">
      <h1 class="display-5 fw-bold">ephimail</h1>
      <p class="lead mb-4">privacy-focused disposable email service.</p>
    </div>
    
    <div class="neo-card mb-4">
      <div class="d-grid gap-2 d-sm-flex justify-content-sm-center">
        <EmailForm @emailSubmitted="handleEmailSubmit"/>
      </div>
    </div>
    
    <!-- Always show both components, regardless of email selection -->
    <div class="row">
      <!-- Email Table Column -->
      <div class="col-md-8">
        <EmailTable :passedEmail="email"/>
      </div>
      
      <!-- Mailbox Reservation Component -->
      <div class="col-md-4">
        <MailboxReservation :email="email" />
      </div>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import EmailForm from '@/components/EmailForm.vue'
import EmailTable from '@/components/EmailTable.vue'
import MailboxReservation from '@/components/MailboxReservation.vue'

export default {
  name: 'HomeView',
  components: {
    EmailForm,
    EmailTable,
    MailboxReservation
  },
  data() {
    return {
      email: null,
    };
  },
  methods: {
    handleEmailSubmit(data) {
      this.email = data;
    }
  }
}
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Space+Grotesk:wght@400;500;700&display=swap');
@import '../assets/css/NeoBrutalism.css';

.logo-container {
  margin-bottom: 2rem;
}
</style>