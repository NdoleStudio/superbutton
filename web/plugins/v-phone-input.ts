import 'flag-icons/css/flag-icons.min.css'
import 'v-phone-input/dist/v-phone-input.css'
import Vue from 'vue'
import { createVPhoneInput } from 'v-phone-input'

const vPhoneInput = createVPhoneInput({
  countryIconMode: 'svg',
})

Vue.use(vPhoneInput)
