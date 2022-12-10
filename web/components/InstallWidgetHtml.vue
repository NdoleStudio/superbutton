<template>
  <v-card>
    <v-card-subtitle class="subtitle-1"
      >Paste the widget code below in the <code>&lt;head&gt;</code> section of
      your HTML page to install the widget on your website.</v-card-subtitle
    >
    <v-card-text>
      <v-alert text color="primary" style="word-break: break-all">{{
        widgetCode
      }}</v-alert>
      <v-btn
        :color="color"
        :block="block"
        :disabled="!copyButtonActive"
        @click="copyWidgetCode"
      >
        <v-icon left>{{ mdiContentCopy }}</v-icon>
        Copy Code
      </v-btn>
    </v-card-text>
  </v-card>
</template>

<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator'
import { mdiContentCopy } from '@mdi/js'
@Component
export default class InstallWidgetHtml extends Vue {
  @Prop({ required: false, type: Boolean, default: false }) block!: boolean
  @Prop({ required: false, type: String, default: 'primary' }) color!: string
  mdiContentCopy = mdiContentCopy
  copyButtonActive = true

  get widgetCode(): string {
    return (
      `<script type="text/javascript">window.SB_USER_ID="${this.$store.getters.authUser.uid}";window.SB_PROJECT_ID="${this.$store.getters.activeProjectId}";(function(){const d=document;const s=d.createElement("script");s.src="https://cdn.superbutton.app/widget.js";s["async"]=true;d.getElementsByTagName("head")[0].appendChild(s)})();` +
      '</' +
      'script>'
    )
  }

  copyWidgetCode() {
    this.copyButtonActive = false
    navigator.clipboard
      .writeText(this.widgetCode)
      .then(() => {
        this.$store.dispatch('addNotification', {
          message: 'Code copied to clipboard',
          type: 'success',
        })
      })
      .catch(this.$sentry.captureException)
      .finally(() => {
        setTimeout(() => {
          this.copyButtonActive = true
        }, 3000)
      })
  }
}
</script>
