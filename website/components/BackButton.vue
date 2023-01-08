<template>
  <v-btn
    color="#212121"
    class="white--text"
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
  name: 'BackButton',
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
