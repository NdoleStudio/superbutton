<template>
  <v-btn
    color="secondary"
    :small="$vuetify.breakpoint.smAndDown"
    @click="goBack"
  >
    <v-icon>{{ mdiArrowLeft }}</v-icon>
    Go Back
  </v-btn>
</template>

<script>
import { mdiArrowLeft } from '@mdi/js'
export default {
  name: 'TermsAndConditions',
  data() {
    return {
      mdiArrowLeft,
    }
  },
  methods: {
    goBack() {
      if (window.history.length > 1) {
        this.$router.back()
        return
      }
      this.$router.push({ name: 'index' })
    },
  },
}
</script>
