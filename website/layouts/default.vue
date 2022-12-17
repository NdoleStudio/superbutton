<template>
  <v-app>
    <v-app-bar fixed app flat color="#ffffff">
      <v-container class="d-flex">
        <v-badge
          class="logo-badge"
          :class="{ 'logo-badge--mobile': $vuetify.breakpoint.mdAndDown }"
          color="#8338ec"
          content="Beta"
        >
          <nuxt-link to="/" class="text-decoration-none d-flex">
            <v-avatar tile size="33" class="mt-1">
              <v-img contain :src="require('@/static/logo.svg')"></v-img>
            </v-avatar>
            <h3
              v-if="$vuetify.breakpoint.lgAndUp"
              class="text-h4 font-weight-thin ml-1"
              :class="{
                'text--secondary': $vuetify.theme.dark,
                'text--primary': !$vuetify.theme.dark,
              }"
            >
              SuperButton
            </h3>
          </nuxt-link>
        </v-badge>
        <v-spacer></v-spacer>
        <client-only>
          <v-btn
            v-if="!isLoggedIn"
            :href="$config.dashoardURL"
            color="primary"
            large
          >
            <v-icon left>{{ mdiLoginVariant }}</v-icon>
            Get Started
          </v-btn>
          <v-btn v-else :href="$config.dashboardURL" color="primary" large>
            <v-icon left>{{ mdiArrowRight }}</v-icon>
            Dashboard
          </v-btn>
        </client-only>
      </v-container>
    </v-app-bar>
    <v-main>
      <Nuxt />
    </v-main>
  </v-app>
</template>

<script>
import { mdiLoginVariant, mdiArrowRight } from '@mdi/js'
import { hasAuthCookie } from '~/plugins/cookies'
export default {
  name: 'DefaultLayout',
  data() {
    return {
      mdiLoginVariant,
      mdiArrowRight,
    }
  },
  computed: {
    isLoggedIn() {
      return hasAuthCookie()
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
  &--mobile {
    .v-badge__wrapper {
      margin-left: 8px;
    }
  }
}
</style>
