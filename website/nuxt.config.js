import colors from 'vuetify/es5/util/colors'

export default {
  // Target: https://go.nuxtjs.dev/config-target
  target: 'static',

  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    titleTemplate: '%s',
    title:
      'Communicate with your customers using their preferred channels - SuperButton',
    htmlAttrs: {
      lang: 'en',
    },
    script: [
      {
        hid: 'sb-widget',
        src: '/superbutton.js',
        async: true,
        defer: true,
      },
    ],
    meta: [
      { charset: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      {
        hid: 'description',
        name: 'description',
        content:
          "Capture more leads on your website using a configurable open source floating widget. You don't need to know how to code and it is free.",
      },
      {
        hid: 'og-title',
        property: 'og:title',
        content:
          'Communicate with your customers using their preferred channels',
      },
      {
        hid: 'og-desc',
        property: 'og:description',
        content:
          "Capture more leads on your website using a configurable open source floating widget. You don't need to know how to code and it is free.",
      },
      {
        hid: 'og-image',
        property: 'og:image',
        content: 'https://superbutton.app/header.png',
      },
      {
        hid: 'twitter-card',
        name: 'twitter:card',
        content: 'summary_large_image',
      },
      { name: 'format-detection', content: 'telephone=no' },
    ],
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
  },

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: [],

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [{ src: '~/plugins/vue-glow', ssr: false }],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: true,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
    // https://go.nuxtjs.dev/typescript
    '@nuxt/typescript-build',
    // https://go.nuxtjs.dev/vuetify
    '@nuxtjs/vuetify',
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: [
    '@nuxtjs/sitemap',
    // https://go.nuxtjs.dev/axios
    '@nuxtjs/axios',
  ],

  // Axios module configuration: https://go.nuxtjs.dev/config-axios
  axios: {
    // Workaround to avoid enforcing hard-coded localhost:3000: https://github.com/nuxt-community/axios-module/issues/308
    baseURL: '/',
  },

  // Vuetify module configuration: https://go.nuxtjs.dev/config-vuetify
  vuetify: {
    customVariables: ['~/assets/variables.scss'],
    theme: {
      dark: false,
      themes: {
        light: {
          primary: '#8338ec',
          accent: colors.grey.darken3,
          secondary: colors.amber.darken3,
          info: colors.teal.lighten1,
          warning: colors.amber.base,
          error: colors.deepOrange.accent4,
          success: colors.green.accent3,
        },
      },
    },
  },

  sitemap: {
    hostname: process.env.APP_URL,
  },

  // Build Configuration: https://go.nuxtjs.dev/config-build
  build: {},

  publicRuntimeConfig: {
    appURL: process.env.APP_URL,
    cdnURL: process.env.CDN_URL,
    dashboardURL: process.env.DASHBOARD_URL,
  },

  server: {
    port: 3333,
  },
}
