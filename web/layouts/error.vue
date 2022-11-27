<template>
  <v-main>
    <v-container fluid class="pa-0" fill-height>
      <v-row align="center" justify="center">
        <v-col cols="12" md="4" xl="3" class="mt-16 text-center">
          <div class="text-center mb-5">
            <span class="text-h3 error--text">An Error Occurred</span>
          </div>
          <p class="text--secondary text-h6 text-center">
            And error occurred while handling your last request. This
            application is currently in active development and some errors are
            unavoidable. We have already been automatically notified about this
            error and we are working hard to fix it.
          </p>
          <v-btn class="primary mx-auto" rounded to="/">Take Me Home</v-btn>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script>
export default {
  name: 'ErrorLayout',
  layout: 'auth',
  props: {
    error: {
      type: Object,
      default: null,
    },
  },
  data() {
    return {
      pageNotFound: '404 Not Found',
      otherError: 'An error occurred',
    }
  },
  head() {
    const title =
      this.error.statusCode === 404 ? this.pageNotFound : this.otherError
    return {
      title,
    }
  },
}
</script>

<style scoped>
h1 {
  font-size: 20px;
}
</style>
