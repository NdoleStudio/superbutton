<template>
  <v-app>
    <v-app-bar app>
      <v-container class="py-0 fill-height">
        <v-badge class="logo-badge" color="#8338ec" content="Beta">
          <nuxt-link to="/" class="text-decoration-none d-flex">
            <v-avatar tile size="33" class="mt-1">
              <v-img contain :src="require('@/static/logo.svg')"></v-img>
            </v-avatar>
            <h3
              class="text--secondary font-weight-thin ml-1"
              :class="{
                'text-h4': $vuetify.breakpoint.lgAndUp,
                'text-h5': !$vuetify.breakpoint.lgAndUp,
              }"
            >
              SuperButton
            </h3>
          </nuxt-link>
        </v-badge>
        <v-spacer></v-spacer>

        <v-menu offset-y left bottom>
          <template #activator="{ on, attrs }">
            <v-btn icon x-large v-bind="attrs" v-on="on">
              <v-avatar size="36" color="black">
                <img
                  v-if="
                    $store.getters.authUser && $store.getters.authUser.photoURL
                  "
                  :src="$store.getters.authUser.photoURL"
                  :alt="$store.getters.authUser.displayName"
                />
                <span
                  v-else-if="$store.getters.authUser"
                  class="white--text text-h5"
                >
                  {{ $store.getters.authUser.email.charAt(0).toUpperCase() }}
                </span>
              </v-avatar>
            </v-btn>
          </template>
          <v-list class="px-2" nav :dense="$vuetify.breakpoint.mdAndDown">
            <v-list-item-group>
              <v-list-item to="/projects/create">
                <v-list-item-icon class="pl-2">
                  <v-icon dense>{{ mdiPlus }}</v-icon>
                </v-list-item-icon>
                <v-list-item-content class="ml-n3">
                  <v-list-item-title class="pr-16">
                    Create Project</v-list-item-title
                  >
                </v-list-item-content>
              </v-list-item>
              <v-list-item @click="logout">
                <v-list-item-icon class="pl-2">
                  <v-icon dense>{{ logoutIcon }}</v-icon>
                </v-list-item-icon>
                <v-list-item-content class="ml-n3">
                  <v-list-item-title class="pr-16"> Logout </v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list-item-group>
          </v-list>
        </v-menu>
      </v-container>
    </v-app-bar>
    <v-main>
      <snackbar-notification></snackbar-notification>
      <Nuxt v-if="$store.getters.authStateChanged" />
      <dashboard-loading v-else></dashboard-loading>
    </v-main>
  </v-app>
</template>

<script>
import { mdiLogout, mdiPlus } from '@mdi/js'

export default {
  name: 'DefaultLayout',
  data() {
    return {
      logoutIcon: mdiLogout,
      mdiPlus,
      items: [],
    }
  },
  methods: {
    logout() {
      this.$fire.auth.signOut().then(() => {
        this.$store.dispatch('addNotification', {
          type: 'info',
          message: 'You have successfully logged out',
        })
        this.$router.push('/')
      })
    },
  },
}
</script>
