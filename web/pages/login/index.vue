<template>
  <v-main>
    <v-container fluid class="pa-0" fill-height>
      <v-row align="center" justify="center">
        <v-col
          cols="12"
          md="4"
          xl="3"
          :class="{ 'mt-n16': $vuetify.breakpoint.lgAndUp }"
        >
          <div class="text-center mb-5">
            <v-avatar tile size="45" class="mt-n8 mr-4">
              <v-img contain :src="require('@/static/logo.svg')"></v-img>
            </v-avatar>
            <span class="text-h3">Welcome</span>
          </div>
          <p class="text--secondary text-h6 text-center">
            Login to the {{ $store.getters.app.name }} dashboard
            <br>
            and manage your floating buttons
          </p>
          <v-card max-width="360" class="mx-auto">
            <v-card-text class="px-0">
              <no-ssr>
                <firebase-auth :to="to"></firebase-auth>
              </no-ssr>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </v-main>
</template>

<script lang="ts">
export default {
  name: "LoginIndex",
  layout: 'auth',
  middleware: ['guest'],
  data() {
    return {
      to: '/',
    }
  },
}
</script>
