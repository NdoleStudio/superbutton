<template>
  <v-app>
    <v-app-bar app flat color="#ffffff" :elevate-on-scroll="true" elevation="4">
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
        <v-btn
          v-if="$vuetify.breakpoint.mdAndUp"
          text
          class="mr-4"
          :large="$vuetify.breakpoint.mdAndUp"
          @click.stop="goToPricing"
        >
          Pricing
        </v-btn>
        <client-only>
          <v-btn
            v-if="!isLoggedIn"
            :href="$config.dashboardURL"
            color="primary"
            :large="$vuetify.breakpoint.mdAndUp"
          >
            <v-icon left>{{ mdiLoginVariant }}</v-icon>
            Get Started
          </v-btn>
          <v-btn
            v-else
            :href="$config.dashboardURL"
            color="primary"
            :large="$vuetify.breakpoint.mdAndUp"
          >
            <v-icon left>{{ mdiArrowRight }}</v-icon>
            Dashboard
          </v-btn>
        </client-only>
      </v-container>
    </v-app-bar>
    <v-main>
      <Nuxt />
    </v-main>
    <v-footer class="mt-8">
      <v-container>
        <v-row>
          <v-col cols="12" md="4">
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
            <p class="subtitle-2 text--secondary">
              Made With <v-icon color="#cf1112">{{ mdiHeart }}</v-icon> in
              Tallinn
              <v-img
                class="d-inline-block"
                width="20"
                src="https://upload.wikimedia.org/wikipedia/commons/8/8f/Flag_of_Estonia.svg"
              ></v-img>
            </p>
            <p>
              <v-btn
                large
                href="https://twitter.com/superbuttonHQ"
                icon
                color="#1DA1F2"
              >
                <v-icon large>{{ mdiTwitter }}</v-icon>
              </v-btn>
              <v-btn
                large
                href="https://github.com/NdoleStudio/superbutton"
                icon
                color="000000"
              >
                <v-icon large>{{ mdiGithub }}</v-icon>
              </v-btn>
            </p>
          </v-col>
          <v-col cols="12" md="4">
            <h2 class="text-h6 mb-2">Resources</h2>
            <ul style="list-style: none" class="pa-0">
              <li class="mb-2">
                <v-hover v-slot="{ hover }">
                  <a
                    class="text--primary text-decoration-none"
                    :class="{ 'text-decoration-underline': hover }"
                    @click.stop="goToPricing"
                    >Pricing</a
                  >
                </v-hover>
              </li>
              <li>
                <v-hover v-slot="{ hover }">
                  <a
                    href="https://status.superbutton.app"
                    class="text--primary text-decoration-none"
                    :class="{ 'text-decoration-underline': hover }"
                    >Site status</a
                  >
                </v-hover>
              </li>
            </ul>
          </v-col>
          <v-col cols="12" md="4">
            <h2 class="text-h6 mb-2">Legal</h2>
            <ul style="list-style: none" class="pa-0">
              <li class="mb-2">Terms & Conditions</li>
              <li>Privacy Policy</li>
            </ul>
          </v-col>
        </v-row>
      </v-container>
    </v-footer>
  </v-app>
</template>

<script>
import {
  mdiLoginVariant,
  mdiArrowRight,
  mdiGithub,
  mdiTwitter,
  mdiHeart,
} from '@mdi/js'
import { hasAuthCookie } from '~/plugins/cookies'
export default {
  name: 'DefaultLayout',
  data() {
    return {
      mdiLoginVariant,
      mdiArrowRight,
      mdiGithub,
      mdiHeart,
      mdiTwitter,
    }
  },
  computed: {
    isLoggedIn() {
      return hasAuthCookie()
    },
  },
  methods: {
    goToPricing() {
      this.$vuetify.goTo('#pricing')
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
