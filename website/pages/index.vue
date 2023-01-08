<template>
  <div>
    <v-container>
      <v-row>
        <v-col
          cols="12"
          class="text-center"
          lg="8"
          offset-lg="2"
          :class="{
            'mt-16': $vuetify.breakpoint.lgAndUp,
            'mt-4': !$vuetify.breakpoint.lgAndUp,
          }"
        >
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
                <v-img
                  max-width="24"
                  contain
                  :src="iconLink('whatsapp')"
                ></v-img>
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
          <v-tooltip top>
            <template #activator="{ on, attrs }">
              <v-avatar
                class="mx-2"
                size="32"
                rounded
                color="#000000"
                v-bind="attrs"
                v-on="on"
              >
                <v-img max-width="24" contain :src="iconLink('github')"></v-img>
              </v-avatar>
            </template>
            <span>Github</span>
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
          <v-btn :href="$config.dashboardURL" color="primary" x-large>
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
            :intensity="2.07"
            intense
          >
            <v-img contain :src="require('@/static/dashboard.png')"></v-img>
          </vue-glow>
        </v-col>
      </v-row>
      <v-row class="mt-16">
        <v-col cols="12" class="d-flex justify-center">
          <v-divider
            v-if="$vuetify.breakpoint.lgAndUp"
            class="mt-4 primary"
          ></v-divider>
          <p class="px-4 text-h5 text-center">
            Trusted by companies who care about their customers
          </p>
          <v-divider
            v-if="$vuetify.breakpoint.lgAndUp"
            class="mt-4 primary"
          ></v-divider>
        </v-col>
        <v-col cols="12" class="d-flex justify-center flex-wrap">
          <a
            href="http://hepilo.com"
            target="_blank"
            class="text-decoration-none"
          >
            <v-img
              contain
              class="mx-4 my-2 logo-image"
              max-width="220px"
              :src="require('@/static/hepilo.png')"
            ></v-img>
          </a>
          <a
            href="http://nyangapay.com"
            target="_blank"
            class="text-decoration-none"
          >
            <v-img
              contain
              class="mx-4 my-2 logo-image"
              max-width="220px"
              :src="require('@/static/nyangapay.png')"
            ></v-img>
          </a>
          <a
            href="http://httpsms.com"
            target="_blank"
            class="text-decoration-none mt-1"
          >
            <v-img
              contain
              class="mx-4 my-2 logo-image"
              max-width="220px"
              :src="require('@/static/httpsms.png')"
            ></v-img>
          </a>
        </v-col>
        <v-col v-if="!$vuetify.breakpoint.mdAndUp" class="px-16 mb-n16">
          <v-divider class="primary"></v-divider>
        </v-col>
      </v-row>
      <v-row class="mt-16 mb-16">
        <v-col cols="12" md="6" class="d-flex align-center">
          <div
            class="text-center"
            :class="{ 'mt-8': $vuetify.breakpoint.lgAndUp }"
          >
            <h3 class="text-h2 mb-4">Customizable</h3>
            <h5 class="text-h5 font-weight-light">
              You can configure the <b>color</b>
              <v-text-field
                v-model="formColor"
                class="d-inline-block pb-0 mb-1 inline-field"
                dense
                placeholder="HEX color"
                style="max-width: 100px"
                solo-inverted
              ></v-text-field>
              <b>icon</b>
              <v-select
                v-model="formIcon"
                class="d-inline-block pb-0 mb-1 inline-field"
                dense
                :items="icons"
                placeholder="HEX color"
                style="max-width: 130px"
                solo-inverted
              >
              </v-select>
              <b>greeting text</b>
              <v-text-field
                v-model="formGreetingText"
                class="d-inline-block pb-0 mb-1 inline-field"
                dense
                placeholder="👋 We're here to help."
                style="max-width: 180px"
                solo-inverted
              ></v-text-field>
              of the widget so it has the same look and feel as your website.
            </h5>
          </div>
        </v-col>
        <v-col
          cols="12"
          md="6"
          class="text-center d-flex align-center"
          :class="{
            'justify-center': $vuetify.breakpoint.mdAndUp,
            'justify-end px-4': !$vuetify.breakpoint.mdAndUp,
          }"
        >
          <div
            class="mt-16"
            :class="{ 'tooltip-active': formGreetingText }"
            :aria-label="formGreetingText"
            data-microtip-position="left"
            :role="formGreetingText ? 'tooltip' : ''"
          >
            <v-avatar size="74" :color="formColor">
              <v-img contain max-width="44" :src="iconLink(formIcon)"></v-img>
            </v-avatar>
          </div>
        </v-col>
      </v-row>
      <v-row class="mt-16">
        <v-col cols="12" md="6" class="d-flex align-center" :order-md="2">
          <div
            class="text-center"
            :class="{ 'mt-16': $vuetify.breakpoint.mdAndUp }"
          >
            <h3
              class="text-h2 mb-4"
              :class="{ 'mt-n12': $vuetify.breakpoint.mdAndUp }"
            >
              Delight Your Visitors
            </h3>
            <h5 class="text-h5 font-weight-light">
              Some of your website visitors like contact forms while others
              prefer to chat directly using whatsapp or email, some prefer phone
              calls, and some prefer to read your documentation. Support every
              customer using their favorite channel.
            </h5>
          </div>
        </v-col>
        <v-col
          :order-md="1"
          cols="12"
          md="6"
          class="justify-center d-flex align-center flex-wrap"
          :class="{ 'mt-16': $vuetify.breakpoint.mdAndUp }"
        >
          <v-tooltip top>
            <template #activator="{ on, attrs }">
              <v-avatar
                class="mx-2 mb-4"
                color="#25d366"
                rounded
                size="72"
                v-bind="attrs"
                v-on="on"
              >
                <v-img
                  max-width="44"
                  contain
                  :src="iconLink('whatsapp')"
                ></v-img>
              </v-avatar>
            </template>
            <span>WhatsApp</span>
          </v-tooltip>
          <v-tooltip top>
            <template #activator="{ on, attrs }">
              <v-avatar
                class="mx-2 mb-4"
                size="72"
                rounded
                color="#1e88e5"
                v-bind="attrs"
                v-on="on"
              >
                <v-img max-width="44" contain :src="iconLink('mail')"></v-img>
              </v-avatar>
            </template>
            <span>Email</span>
          </v-tooltip>
          <v-tooltip top>
            <template #activator="{ on, attrs }">
              <v-avatar
                class="mx-2 mb-4"
                size="72"
                rounded
                color="#25d366"
                v-bind="attrs"
                v-on="on"
              >
                <v-img
                  max-width="44"
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
                class="mx-2 mb-4"
                size="72"
                rounded
                color="#1e88e5"
                v-bind="attrs"
                v-on="on"
              >
                <v-img
                  max-width="44"
                  contain
                  :src="iconLink('documentation')"
                ></v-img>
              </v-avatar>
            </template>
            <span>Documentation</span>
          </v-tooltip>
          <v-tooltip top>
            <template #activator="{ on, attrs }">
              <v-avatar
                class="mx-2 mb-4"
                size="72"
                rounded
                color="#000000"
                v-bind="attrs"
                v-on="on"
              >
                <v-img max-width="44" contain :src="iconLink('github')"></v-img>
              </v-avatar>
            </template>
            <span>GitHub</span>
          </v-tooltip>
          <v-tooltip top>
            <template #activator="{ on, attrs }">
              <v-avatar
                class="mx-2 mb-4"
                size="72"
                rounded
                color="#1e88e5"
                v-bind="attrs"
                v-on="on"
              >
                <v-img max-width="44" contain :src="iconLink('link')"></v-img>
              </v-avatar>
            </template>
            <span>External Links</span>
          </v-tooltip>
        </v-col>
      </v-row>
    </v-container>
    <v-sheet class="mt-16">
      <v-container>
        <v-row>
          <v-col class="mt-16" md="6" offset-md="3">
            <h2
              id="pricing"
              class="text-center text-h2 mb-4 text-decoration-underline"
              style="color: #8338ec"
            >
              <span class="text--primary">Pricing</span>
            </h2>
            <h4 class="text-center text-h5 text--secondary">
              Most of the SuperButton features are completely
              <span style="color: #8338ec">free</span> but if you want a little
              extra, you can go pro
            </h4>
            <div class="d-flex justify-center mt-4 align-center">
              <p class="text-h6 mr-3 mt-3">Monthly</p>
              <v-switch v-model="yearlyPricing"></v-switch>
              <p class="text-h6 ml-3 mt-3">
                Yearly <v-chip color="primary" small>-20%</v-chip>
              </p>
            </div>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" lg="4" offset-lg="2">
            <v-card elevation="4">
              <v-card-text>
                <h1 class="text-center text-h3 text--primary">Free</h1>
                <p class="subtitle-1 text-center">
                  Great for hobby websites and personal sites and experiments.
                </p>
                <p class="text-center text--primary">
                  <span class="text-h3">$0</span>
                </p>
                <p class="text-center mt-n3">No credit card required</p>
                <v-btn
                  block
                  color="black"
                  :href="$config.dashboardURL"
                  class="white--text"
                  large
                  >Get Started</v-btn
                >
                <p class="mt-8 subtitle-1">
                  <v-icon color="primary" class="mt-n1" left>{{
                    mdiCheckCircle
                  }}</v-icon
                  >Enable multiple channels
                </p>
                <p class="subtitle-1 mt-n4">
                  <v-icon color="primary" class="mt-n1" left>{{
                    mdiCheckCircle
                  }}</v-icon
                  >Use the widget on multiple websites
                </p>
                <p class="subtitle-1 mt-n4">
                  <v-icon color="primary" class="mt-n1" left>{{
                    mdiCheckCircle
                  }}</v-icon
                  >Customize the color, icon, etc
                </p>
                <p class="subtitle-1 mt-n4">
                  <v-icon color="primary" class="mt-n1" left>{{
                    mdiCheckCircle
                  }}</v-icon
                  >Basic Support
                </p>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" lg="4">
            <v-card elevation="4">
              <v-card-text>
                <h1 class="text-center text-h3" style="color: #8338ec">Pro</h1>
                <p class="subtitle-1 text-center">
                  Advanced features to delight your website visitors.
                </p>
                <p v-if="!yearlyPricing" class="text-center text--primary">
                  <span class="text-h3">$6</span>/month
                </p>
                <p v-else class="text-center text--primary">
                  <span class="text-h3">$60</span>/year
                </p>
                <p v-if="!yearlyPricing" class="text-center mt-n3">
                  or <b>$72</b> per year
                </p>
                <p v-else class="text-center mt-n3">or <b>$5</b> per month</p>
                <v-btn block color="primary" :href="$config.dashboardURL" large
                  >Try For Free</v-btn
                >
                <p class="mt-8 subtitle-1">
                  <v-icon color="primary" class="mt-n1" left>{{
                    mdiCheckCircle
                  }}</v-icon
                  >Everything in Free plan
                </p>
                <p class="subtitle-1 mt-n4">
                  <v-icon color="primary" class="mt-n1" left>{{
                    mdiCheckCircle
                  }}</v-icon
                  >20MB file uploads in contact form
                </p>
                <p class="subtitle-1 mt-n4">
                  <v-icon color="primary" class="mt-n1" left>{{
                    mdiCheckCircle
                  }}</v-icon
                  >Priority Support
                </p>
                <p class="subtitle-1 mt-n4">
                  <v-icon color="primary" class="mt-n1" left>{{
                    mdiCheckCircle
                  }}</v-icon
                  >Remove superbutton branding
                </p>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </v-sheet>
  </div>
</template>

<script>
import { mdiCheckCircle } from '@mdi/js'
export default {
  name: 'IndexPage',
  data() {
    return {
      mdiCheckCircle,
      formColor: '#8338EC',
      formIcon: 'chat',
      yearlyPricing: false,
      icons: [
        {
          text: 'Chat Icon',
          value: 'chat',
        },
        {
          text: 'Whatsapp Icon',
          value: 'whatsapp',
        },
        {
          text: 'Help Icon',
          value: 'help-chat',
        },
        {
          text: 'Documentation Icon',
          value: 'documentation',
        },
      ],
      formGreetingText: "👋 We're here to help.",
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

<style lang="scss">
.v-application {
  .logo-image {
    &:not(:hover) {
      filter: grayscale(100%);
    }
  }
  .inline-field {
    .v-text-field__details {
      display: none;
    }
  }
}
</style>
