import { Framework } from 'vuetify'

declare module 'vue/types/vue' {
  interface Vue {
    $vuetify: Framework
  }
}
