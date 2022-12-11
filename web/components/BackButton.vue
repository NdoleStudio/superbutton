<template>
  <v-btn
    :color="icon ? 'primary' : 'default'"
    :small="$vuetify.breakpoint.smAndDown"
    :block="block"
    :icon="icon"
    :large="large"
    @click="goBack"
  >
    <v-icon left>{{ mdiArrowLeft }}</v-icon>
    <span v-if="!icon">Go Back</span>
  </v-btn>
</template>

<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator'
import { Location } from 'vue-router'
import { mdiKeyboardBackspace } from '@mdi/js'
@Component
export default class BackButton extends Vue {
  @Prop({ required: false }) route?: Location
  @Prop({ required: false, type: Boolean, default: false }) icon!: boolean
  @Prop({ required: false, type: Boolean, default: false }) block!: boolean
  @Prop({ required: false, type: Boolean, default: false }) large!: boolean
  mdiArrowLeft = mdiKeyboardBackspace
  goBack(): void {
    if (this.route) {
      this.$router.push(this.route)
      return
    }
    if (window.history.length > 1) {
      this.$router.back()
      return
    }
    this.$router.push({ name: 'index' })
  }
}
</script>
