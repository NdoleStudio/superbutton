<template>
  <v-main>
    <v-container fluid class="pa-0" fill-height>
      <v-row align="center" justify="center">
        <v-col cols="12" md="4" xl="3" class="text-center px-4">
          <h2 class="text-h4 text--secondary mt-4 mb-4">
            <v-img
              contain
              class="mx-auto d-inline-block mt-n16"
              :src="require('@/static/logo.svg')"
              max-width="32"
            ></v-img>
            Loading the superbutton dashboard.
          </h2>
          <v-progress-linear
            indeterminate
            color="secondary"
          ></v-progress-linear>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script>
export default {
  name: 'DashboardLoading',
}
</script>

<style lang="scss">
.fade-enter-active {
  transition: opacity 0.5s;
}

.fade-enter,
.fade-leave-active {
  opacity: 0;
}
</style>
