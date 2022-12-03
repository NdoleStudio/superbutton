<template>
  <div
    v-if="settingsLoaded && settings !== null"
    class="sb-widget"
    :class="{ 'sb-widget--tooltip-active': tooltipActive }"
    :style="widgetStyle"
    :aria-label="settings.project.greeting"
    data-microtip-position="left"
    role="sb-tooltip"
  >
    <div class="sb-widget__image" :style="widgetImageStyle"></div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'

interface WhatsappIntegration {
  type: 'whatsapp'
  settings: {
    enabled: boolean
    text: string
    phone_number: string
    icon: string
  }
}

interface ContentIntegration {
  type: 'content'
  settings: {
    enabled: boolean
    title: string
    summary: string
    text: string
    icon: string
  }
}

interface Project {
  name: string
  icon: string
  greeting: string
  greeting_timeout_seconds: number
  color: string
}

interface Settings {
  project: Project | null
  integrations: Array<WhatsappIntegration | ContentIntegration>
}

@Component
export default class SuperbuttonWidget extends Vue {
  settingsLoaded = false
  tooltipActive = false
  settings: Settings | null = null

  get backgroundImage() {
    return (
      window.location.origin + '/icons/' + this.settings?.project?.icon + '.svg'
    )
  }

  get widgetStyle() {
    return {
      backgroundColor: this.settings?.project?.color,
    }
  }

  get widgetImageStyle() {
    return {
      background: `url(${this.backgroundImage})`,
      backgroundRepeat: 'no-repeat',
      backgroundSize: 'cover',
    }
  }

  get isMobile(): boolean {
    return /iPhone|iPod|Android/i.test(navigator.userAgent)
  }

  mounted() {
    this.loadSettings(
      '9DMHezLb9NV7Had2PY003K8KRVn2',
      '0f097a15-3a7b-4602-9c6d-ed2b00683a47'
    )
  }

  loadSettings(userId: string, projectId: string) {
    fetch(`http://localhost:8000/v1/settings/${userId}/projects/${projectId}`)
      .then((response) => response.json())
      .then((response) => {
        this.settings = response.data
        this.settingsLoaded = true
        if (this.settings?.project?.greeting && !this.isMobile) {
          setTimeout(() => {
            this.tooltipActive = true
          }, this.settings.project.greeting_timeout_seconds * 1000)
        }
      })
  }
}
</script>

<style scoped lang="scss">
@import 'assets/microtip.scss';
.sb-widget {
  height: 72px;
  width: 72px;
  cursor: pointer;
  display: flex;
  justify-content: center;
  align-items: center;
  position: fixed;
  right: 48px;
  border-radius: 50%;
  z-index: 10000;
  bottom: 48px;

  &__image {
    background-repeat: no-repeat;
    background-size: cover;
    height: 44px;
    width: 44px;
  }
}
</style>
