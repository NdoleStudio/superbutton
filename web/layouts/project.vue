<template>
  <v-app>
    <v-app-bar app flat>
      <v-container class="py-0 fill-height">
        <v-badge class="logo-badge" color="#8338ec" content="Beta">
          <nuxt-link to="/" class="text-decoration-none d-flex">
            <v-avatar tile size="33" class="mt-1">
              <v-img contain :src="require('@/static/logo.svg')"></v-img>
            </v-avatar>
            <h3 class="text-h4 text--secondary font-weight-thin ml-1">
              Superbutton
            </h3>
          </nuxt-link>
        </v-badge>
        <div style="max-width: 250px">
          <v-select
            v-if="$store.getters.hasProjects"
            :value="$store.getters.activeProjectId"
            outlined
            dense
            label="Project"
            class="ml-16 mb-n4 mt-2"
            :items="projects"
            @change="onProjectChange"
          >
            <template #append-item>
              <div class="ml-3 mt-4">
                <v-btn
                  text
                  color="primary"
                  :to="{ name: 'projects-create' }"
                  small
                >
                  <v-icon small>{{ mdiPlus }}</v-icon>
                  Add Project
                </v-btn>
              </div>
            </template>
          </v-select>
        </div>
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
      <dashboard-loading
        v-if="!$store.getters.authStateChanged"
      ></dashboard-loading>
      <v-container v-else class="pt-8">
        <v-row>
          <v-col v-if="$vuetify.breakpoint.lgAndUp" cols="3" xl="2">
            <v-list color="transparent" rounded>
              <v-list-item
                v-for="item in items"
                :key="item.name"
                color="primary"
                link
                :to="item.route"
              >
                <v-list-item-icon>
                  <v-icon>{{ item.icon }}</v-icon>
                </v-list-item-icon>
                <v-list-item-content>
                  <v-list-item-title class="text-h5">{{
                    item.name
                  }}</v-list-item-title>
                </v-list-item-content>
              </v-list-item>
            </v-list>
          </v-col>
          <v-col cols="12" lg="9" xl="10">
            <Nuxt v-if="$store.getters.authStateChanged" />
          </v-col>
        </v-row>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import { mdiLogout, mdiPlus, mdiCogOutline, mdiLan, mdiXml } from '@mdi/js'

export default {
  name: 'ProjectLayout',
  data() {
    return {
      logoutIcon: mdiLogout,
      mdiPlus,
    }
  },
  computed: {
    projects() {
      return this.$store.getters.projects.map((project) => {
        return {
          text: project.name,
          value: project.id,
        }
      })
    },
    items() {
      return [
        {
          name: 'Settings',
          icon: mdiCogOutline,
          route: {
            name: 'projects-id-settings',
            params: {
              id: this.$store.getters.activeProjectId,
            },
          },
        },
        {
          name: 'Integrations',
          icon: mdiLan,
          route: {
            name: 'projects-id-integrations',
            params: {
              id: this.$store.getters.activeProjectId,
            },
          },
        },
        {
          name: 'Install Widget',
          icon: mdiXml,
          route: {
            name: 'projects-id-install',
            params: {
              id: this.$store.getters.activeProjectId,
            },
          },
        },
      ]
    },
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
    onProjectChange(item) {
      this.$store.dispatch('setActiveProjectId', item).then(() => {
        this.$router.push({
          name: 'projects-id-settings',
          params: {
            id: item,
          },
        })
      })
    },
  },
}
</script>

<style lang="scss">
.logo-badge {
  .v-badge__wrapper {
    span {
      margin-bottom: -12px;
    }
  }
}
</style>
