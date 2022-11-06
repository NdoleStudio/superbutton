import { Auth } from 'firebase/auth'
import { Framework } from 'vuetify'

interface Firebase {
  auth: Auth
}

export interface SelectItem {
  text: string
  value: string | number
}

declare module 'vue/types/vue' {
  interface Vue {
    $vuetify: Framework
    $fire: Firebase
  }
}
