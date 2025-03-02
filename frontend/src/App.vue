<template>
  <notifications />
  <div class="px-4 py-5 my-5">
    <img class="d-block mx-auto mb-4" src="@/assets/logo.png" alt="" width="256">
    <div class="col-lg-6 mx-auto text-center">
      <h1 class="display-5 fw-bold">ephimail</h1>
      <p class="lead mb-4">privacy-focus disposable email service.</p>
      <div class="d-grid gap-2 d-sm-flex justify-content-sm-center">
        <EmailForm @emailSubmitted="handleEmailSubmit"/>
      </div>
    </div>
    
    <!-- Mailbox Reservation Component (only shown when an email is selected) -->
    <div v-if="email" class="d-grid gap-2 d-sm-flex justify-content-sm-center">
      <MailboxReservation :email="email" />
    </div>
    
    <!-- Email Table -->
    <div class="d-grid gap-2 d-sm-flex justify-content-sm-center">
      <EmailTable :passedEmail="email"/>
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
body {
  background-color: #98c48e !important;
}
</style>