import colors from 'vuetify/es5/util/colors'

export default {
  // Target: https://go.nuxtjs.dev/config-target
  target: 'static',

  // Global page headers: https://go.nuxtjs.dev/config-head
  head: {
    titleTemplate: '%s - Super Button',
    title: 'Awesome floating button for your website',
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
      { hid: 'description', name: 'description', content: '' },
      { name: 'format-detection', content: 'telephone=no' },
    ],
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
  },

  // Global CSS: https://go.nuxtjs.dev/config-css
  css: [],

  // Plugins to run before rendering page: https://go.nuxtjs.dev/config-plugins
  plugins: [
    { src: '~/plugins/v-phone-input', ssr: false },
    { src: '~/plugins/vue-draggable', ssr: false },
  ],

  // Auto import components: https://go.nuxtjs.dev/config-components
  components: true,

  // Modules for dev and build (recommended): https://go.nuxtjs.dev/config-modules
  buildModules: [
    '@nuxtjs/dotenv',
    // https://go.nuxtjs.dev/typescript
    '@nuxt/typescript-build',
    // https://go.nuxtjs.dev/vuetify
    '@nuxtjs/vuetify',
  ],

  // Modules: https://go.nuxtjs.dev/config-modules
  modules: [
    [
      '@nuxtjs/firebase',
      {
        config: {
          apiKey: 'AIzaSyDsX2aZleeS8yc5w5TwruJKhrdpyixkiw0',
          authDomain: 'superbutton-8a6ad.firebaseapp.com',
          projectId: 'superbutton-8a6ad',
          storageBucket: 'superbutton-8a6ad.appspot.com',
          messagingSenderId: '1038134704954',
          appId: '1:1038134704954:web:a7046561e1b24fa4ed2e03',
          measurementId: 'G-GTL1VHGNKJ',
        },
        services: {
          auth: {
            persistence: 'local', // default
            initialize: {
              onAuthStateChangedAction: 'onAuthStateChanged',
              onIdTokenChangedAction: 'onAuthStateChanged',
              subscribeManually: false,
            },
            ssr: false, // default
          },
        },
      },
    ],
    '@nuxtjs/sentry',
  ],

  // Vuetify module configuration: https://go.nuxtjs.dev/config-vuetify
  vuetify: {
    customVariables: ['~/assets/variables.scss'],
    theme: {
      dark: true,
      themes: {
        dark: {
          primary: colors.blue.darken2,
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
  sentry: {
    dsn: process.env.SENTRY_DSN,
    config: {
      environment: process.env.SENTRY_ENVIRONMENT,
      release:
        process.env.SENTRY_PROJECT_NAME +
        '@' +
        process.env.NUXT_ENV_CURRENT_GIT_SHA,
    },
    tracing: {
      tracesSampleRate: 1.0,
      vueOptions: {
        tracing: true,
        tracingOptions: {
          hooks: ['mount', 'update'],
          timeout: 2000,
          trackComponents: true,
        },
      },
    },
  },

  publicRuntimeConfig: {
    checkoutURL: process.env.CHECKOUT_URL,
  },

  // Build Configuration: https://go.nuxtjs.dev/config-build
  build: {},
}
