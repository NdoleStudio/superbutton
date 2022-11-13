<template>
  <div>
    <v-app-bar app flat>
      <v-container class="py-0 fill-height">
        <v-avatar tile size="32">
          <v-img contain :src="require('@/static/logo.svg')"></v-img>
        </v-avatar>
        <h3 class="text-h4 text--secondary font-weight-thin ml-1">
          Superbutton
        </h3>
        <v-spacer></v-spacer>
        <v-menu v-if="$store.getters.authUser" left bottom>
          <template #activator="{ on }">
            <v-btn icon x-large v-on="on">
              <v-avatar size="36" color="black">
                <img
                  v-if="$store.getters.authUser.photoURL"
                  :src="$store.getters.authUser.photoURL"
                  :alt="$store.getters.authUser.displayName"
                />
                <span v-else class="white--text text-h5">A</span>
              </v-avatar>
            </v-btn>
          </template>
          <v-list class="px-2" nav :dense="$vuetify.breakpoint.mdAndDown">
            <v-list-item-group>
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
      <v-container>
        <v-row>
          <v-col cols="2">
            <v-sheet rounded="lg">
              <v-list color="transparent">
                <v-list-item v-for="n in 5" :key="n" link>
                  <v-list-item-content>
                    <v-list-item-title> List Item {{ n }} </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>

                <v-divider class="my-2"></v-divider>

                <v-list-item link color="grey lighten-4">
                  <v-list-item-content>
                    <v-list-item-title> Refresh </v-list-item-title>
                  </v-list-item-content>
                </v-list-item>
              </v-list>
            </v-sheet>
          </v-col>

          <v-col>
            <v-sheet min-height="70vh" rounded="lg">
              <!--  -->
            </v-sheet>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </div>
</template>

<script>
import { mdiLogout } from '@mdi/js'

export default {
  name: 'ProjectsCreate',
  layout: 'auth',
  data() {
    return {
      to: '/',
      logoutIcon: mdiLogout,
    }
  },
  mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }
    this.$store.dispatch('getUser')
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
