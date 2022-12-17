<template>
  <v-container>
    <v-row>
      <v-col cols="12" class="text-center mt-16" lg="8" offset-lg="2">
        <h1
          :class="{
            'text-h2': $vuetify.breakpoint.lgAndUp,
            'text-h3': !$vuetify.breakpoint.lgAndUp,
          }"
        >
          Communicate with your customers using their preferred channels
        </h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col
        cols="12"
        class="text-center"
        md="8"
        offset-md="2"
        lg="6"
        offset-lg="3"
      >
        <h2
          class="text--secondary font-weight-light"
          :class="{
            'text-h5': $vuetify.breakpoint.lgAndUp,
            'text-h5': !$vuetify.breakpoint.lgAndUp,
          }"
        >
          Capture more leads on your website using a configurable
          <a
            class="text-decoration-none"
            href="https://github.com/NdoleStudio/superbutton"
            >open source</a
          >
          floating widget. You don't need to know how to code and
          <span style="color: #8338ec">it is free.</span>
        </h2>
      </v-col>
    </v-row>
    <v-row class="mt-n4">
      <v-col
        cols="12"
        class="text-center"
        md="8"
        offset-md="2"
        lg="6"
        offset-lg="3"
      >
        <v-tooltip top>
          <template #activator="{ on, attrs }">
            <v-avatar
              class="mx-2"
              color="#25d366"
              rounded
              size="32"
              v-bind="attrs"
              v-on="on"
            >
              <v-img max-width="24" contain :src="iconLink('whatsapp')"></v-img>
            </v-avatar>
          </template>
          <span>WhatsApp</span>
        </v-tooltip>
        <v-tooltip top>
          <template #activator="{ on, attrs }">
            <v-avatar
              class="mx-2"
              size="32"
              rounded
              color="#1e88e5"
              v-bind="attrs"
              v-on="on"
            >
              <v-img max-width="24" contain :src="iconLink('link')"></v-img>
            </v-avatar>
          </template>
          <span>External Links</span>
        </v-tooltip>
        <v-tooltip top>
          <template #activator="{ on, attrs }">
            <v-avatar
              class="mx-2"
              size="32"
              rounded
              color="#25d366"
              v-bind="attrs"
              v-on="on"
            >
              <v-img
                max-width="24"
                contain
                :src="iconLink('phone-call')"
              ></v-img>
            </v-avatar>
          </template>
          <span>Phone Call</span>
        </v-tooltip>
        <v-tooltip top>
          <template #activator="{ on, attrs }">
            <v-avatar
              class="mx-2"
              size="32"
              rounded
              color="#1e88e5"
              v-bind="attrs"
              v-on="on"
            >
              <v-img
                max-width="24"
                contain
                :src="iconLink('documentation')"
              ></v-img>
            </v-avatar>
          </template>
          <span>Documentation</span>
        </v-tooltip>
      </v-col>
    </v-row>
    <v-row>
      <v-col
        cols="12"
        class="text-center"
        md="8"
        offset-md="2"
        lg="6"
        offset-lg="3"
      >
        <v-btn :href="$config.dashoardURL" color="primary" x-large>
          Get Started for Free
        </v-btn>
      </v-col>
      <v-col
        cols="12"
        class="text-center mt-n4"
        md="8"
        offset-md="2"
        lg="6"
        offset-lg="3"
      >
        <v-icon class="mt-n1">{{ mdiCheckCircle }}</v-icon>
        <span class="text--secondary">No credit card required</span>
      </v-col>
    </v-row>
    <v-row class="mt-8">
      <v-col>
        <vue-glow
          color="#8338ec"
          mode="hex"
          elevation="24"
          intensity="2.07"
          intense
        >
          <v-img contain :src="require('@/static/dashboard.png')"></v-img>
        </vue-glow>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiCheckCircle } from '@mdi/js'
export default {
  name: 'IndexPage',
  data() {
    return {
      mdiCheckCircle,
    }
  },
  methods: {
    /**
     * @param {string} icon
     */
    iconLink(icon) {
      return `${this.$config.cdnURL}/icons/${icon}.svg`
    },
  },
}
</script>
