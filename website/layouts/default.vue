<template>
  <v-app>
    <v-app-bar fixed app flat>
      <v-container class="d-flex">
        <v-badge
          v-if="$vuetify.breakpoint.lgAndUp"
          class="logo-badge"
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
        <v-spacer></v-spacer>
        <v-btn color="primary" large>
          <v-icon left>{{ mdiLoginVariant }}</v-icon>
          Get Started
        </v-btn>
      </v-container>
    </v-app-bar>
    <v-main style="background-color: #f5f5f5">
      <v-container>
        <Nuxt />
      </v-container>
    </v-main>
  </v-app>
</template>

<script>
import { mdiLoginVariant } from '@mdi/js'
export default {
  name: 'DefaultLayout',
  data() {
    return {
      mdiLoginVariant,
    }
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
