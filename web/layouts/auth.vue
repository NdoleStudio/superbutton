<template>
  <v-app>
    <snackbar-notification></snackbar-notification>
    <Nuxt v-if="$store.getters.authStateChanged"/>
    <dashboard-loading v-else></dashboard-loading>
  </v-app>
</template>

<script lang="ts">
export default {
  name: "AuthLayout"
}
</script>
