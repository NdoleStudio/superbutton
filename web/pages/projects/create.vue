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
        <v-row class="pt-16">
          <v-col lg="4" md="6" offset-md="3" offset-lg="4">
            <v-card class="mt-16 pa-8">
              <v-card-text class="pt-3 text-center">
                <h4 class="text-h4 text--primary">Welcome</h4>
                <h5 class="pt-1 pb-6 text--secondary subtitle-1">
                  Start by setting up your new project
                </h5>
                <v-form>
                  <v-text-field
                    v-model="formName"
                    :counter="30"
                    label="Name"
                    persistent-placeholder
                    placeholder="Project name e.g Google"
                    outlined
                    class="mb-4"
                    :error="$store.getters.errorMessages.has('name')"
                    :error-messages="$store.getters.errorMessages.get('name')"
                    required
                  ></v-text-field>
                  <v-text-field
                    v-model="formWebsite"
                    :counter="255"
                    class="mb-4"
                    label="Website"
                    persistent-placeholder
                    :error="$store.getters.errorMessages.has('website')"
                    :error-messages="
                      $store.getters.errorMessages.get('website')
                    "
                    placeholder="Website URL e.g https://google.com"
                    outlined
                    required
                  ></v-text-field>
                  <loading-button
                    :loading="$store.getters.creatingProject"
                    :icon="mdiPlus"
                    :block="true"
                    :large="true"
                    @click="createProject"
                  >
                    Create Project
                  </loading-button>
                </v-form>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </div>
</template>

<script>
import { mdiLogout, mdiPlus } from '@mdi/js'

export default {
  name: 'ProjectsCreate',
  layout: 'auth',
  data() {
    return {
      to: '/',
      formName: '',
      mdiPlus,
      formWebsite: '',
      logoutIcon: mdiLogout,
    }
  },
  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }
    await Promise.all([
      this.$store.dispatch('loadProjects'),
      this.$store.dispatch('loadUser'),
    ])

    if (!this.$store.getters.hasProjects) {
      await this.$router.push('/projects/create')
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

    createProject() {
      this.$store
        .dispatch('createProject', {
          name: this.formName,
          website: this.formWebsite,
        })
        .then((payload) => {
          this.$router.push('/projects/' + payload.id)
        })
    },
  },
}
</script>
