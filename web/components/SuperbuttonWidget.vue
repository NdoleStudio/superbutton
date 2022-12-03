<template>
  <div class="sb-widget">
    <div v-if="windowOpen" class="sb-widget__window">
      <div
        class="sb-widget__window__header"
        :style="{ backgroundColor: settings.project.color }"
      >
        <div class="sb-row">
          <div class="sb-widget__window__header__project-name">
            {{ settings.project.name }}
          </div>
          <div
            class="sb-widget__window__header__close-button"
            @click="closeWidgetWindow"
          >
            <svg
              class="icon"
              style="width: 32px; height: 32px; fill: #f3f4f6"
              viewBox="0 0 32 32"
            >
              <path :d="closeIcon" />
            </svg>
          </div>
        </div>
      </div>
      <div class="sb-widget__window__body">
        <div
          v-if="activeIntegration"
          class="sb-widget__window__body__integration--active"
        >
          <div style="display: flex">
            <div
              class="sb-widget__window__body__integration--active__back-button"
              @click="closeActiveIntegration"
            >
              <svg
                class="icon"
                style="width: 32px; height: 32px; fill: #2196f3"
                viewBox="0 0 32 32"
              >
                <path :d="backIcon" />
              </svg>
            </div>
            <div class="sb-widget__window__body__integration--active__title">
              {{ activeIntegration.settings.title }}
            </div>
          </div>
        </div>
        <div
          v-if="activeIntegration"
          class="sb-widget__window__body__integration--active__text"
        >
          {{ activeIntegration.settings.text }}
        </div>
        <template v-for="integration in settings.integrations" v-else>
          <div
            v-if="integration.type === 'phone-call'"
            :key="integration.id"
            class="sb-widget__window__body__integration sb-widget__window__body__integration--phone-call"
            @click="openPhoneCall(integration.settings.phone_number)"
          >
            <div class="sb-widget__window__body__integration--phone-call__icon">
              <div
                class="sb-widget__window__body__integration--phone-call__image"
                :style="phoneCallIconStyle"
              ></div>
            </div>
            <div class="sb-widget__window__body__integration--phone-call__text">
              {{ integration.settings.text }}
            </div>
          </div>
          <div
            v-if="integration.type === 'link'"
            :key="integration.id"
            class="sb-widget__window__body__integration sb-widget__window__body__integration--link"
            @click="openLink(integration.settings.url)"
          >
            <div class="sb-widget__window__body__integration--link__icon">
              <div
                class="sb-widget__window__body__integration--link__image"
                :style="linkIconStyle"
              ></div>
            </div>
            <div class="sb-widget__window__body__integration--link__text">
              {{ integration.settings.text }}
            </div>
          </div>
          <div
            v-if="integration.type === 'whatsapp'"
            :key="integration.id"
            class="sb-widget__window__body__integration sb-widget__window__body__integration--whatsapp"
            @click="openWhatsappChat(integration.settings.phone_number)"
          >
            <div class="sb-widget__window__body__integration--whatsapp__icon">
              <div
                class="sb-widget__window__body__integration--whatsapp__image"
                :style="whatsappIconStyle"
              ></div>
            </div>
            <div class="sb-widget__window__body__integration--whatsapp__text">
              {{ integration.settings.text }}
            </div>
          </div>
          <div
            v-if="integration.type === 'content'"
            :key="integration.id"
            class="sb-widget__window__body__integration sb-widget__window__body__integration--content"
            @click="openContentIntegration(integration.id)"
          >
            <div style="flex-grow: 1">
              <div class="sb-widget__window__body__integration--content__title">
                {{ integration.settings.title }}
              </div>
              <div class="sb-widget__window__body__integration--content__text">
                {{ integration.settings.summary }}
              </div>
            </div>
            <div
              style="
                width: 24px;
                display: flex;
                justify-content: center;
                align-items: center;
              "
            >
              <svg
                class="icon"
                style="width: 24px; height: 24px; fill: #4b587c"
                viewBox="0 0 24 24"
              >
                <path :d="rightIcon" />
              </svg>
            </div>
          </div>
        </template>
        <div class="sb-widget__window__body__mention">
          Powered by
          <a target="_blank" href="https://superbutton.app?utm=widget"
            >Superbutton</a
          >
        </div>
      </div>
    </div>
    <div
      v-if="settingsLoaded && settings !== null"
      class="sb-widget__chat-head"
      :class="{
        'sb-widget--tooltip-active': tooltipActive,
        'sb-widget--tooltip-disabled': windowOpen,
      }"
      :style="widgetStyle"
      :aria-label="settings.project.greeting"
      data-microtip-position="left"
      role="sb-tooltip"
      @click="toggleWidgetWindow"
    >
      <div class="sb-widget__image" :style="widgetImageStyle"></div>
    </div>
  </div>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
import { mdiArrowLeft, mdiChevronRight, mdiClose } from '@mdi/js'

interface WhatsappIntegration {
  id: string
  type: 'whatsapp'
  settings: {
    enabled: boolean
    text: string
    phone_number: string
    icon: string
  }
}

interface PhoneCallConfiguration {
  id: string
  type: 'phone-call'
  settings: {
    enabled: boolean
    text: string
    phone_number: string
    icon: string
  }
}

export interface LinkIntegration {
  type: 'link'
  id: string
  settings: {
    enabled: boolean
    title: string
    url: string
    text: string
    icon: string
  }
}

export interface ContentIntegration {
  type: 'content'
  id: string
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
  integrations: Array<
    | WhatsappIntegration
    | ContentIntegration
    | PhoneCallConfiguration
    | LinkIntegration
  >
}

@Component
export default class SuperbuttonWidget extends Vue {
  settingsLoaded = false
  tooltipActive = false
  windowOpen = false
  showGreeting = false
  rightIcon: string = mdiChevronRight
  closeIcon: string = mdiClose
  backIcon: string = mdiArrowLeft
  settings: Settings | null = null
  activeIntegrationId: string | null = null

  get activeIntegration(): ContentIntegration | null {
    const integration = this.settings?.integrations.find(
      (x) => x.id === this.activeIntegrationId
    )
    if (!integration) {
      return null
    }
    return integration as ContentIntegration
  }

  get backgroundImage() {
    if (this.windowOpen) {
      return this.iconUrl('close')
    }
    return this.iconUrl(this.settings?.project?.icon as string)
  }

  get whatsappIconStyle() {
    return {
      backgroundImage: `url(${this.iconUrl('whatsapp')}`,
      backgroundRepeat: 'no-repeat',
      height: '16px',
      width: '16px',
      backgroundSize: 'cover',
    }
  }

  get phoneCallIconStyle() {
    return {
      backgroundImage: `url(${this.iconUrl('phone-call')}`,
      backgroundRepeat: 'no-repeat',
      height: '16px',
      width: '16px',
      backgroundSize: 'cover',
    }
  }

  get linkIconStyle() {
    return {
      backgroundImage: `url(${this.iconUrl('link')}`,
      backgroundRepeat: 'no-repeat',
      height: '16px',
      width: '16px',
      backgroundSize: 'cover',
    }
  }

  get widgetStyle() {
    return {
      backgroundColor: this.settings?.project?.color,
      float: 'right',
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
      this.$store.getters.authUser.uid,
      this.$store.getters.activeProjectId
    )
  }

  toggleWidgetWindow() {
    if (this.windowOpen) {
      this.closeWidgetWindow()
      return
    }
    this.openWidgetWindow()
  }

  openWidgetWindow() {
    this.windowOpen = true
    this.tooltipActive = false
  }

  openWhatsappChat(phoneNumber: string) {
    window
      .open(`https://wa.me/${phoneNumber.replace('+', '')}`, '_blank')
      ?.focus()
  }

  openPhoneCall(phoneNumber: string) {
    window.open(`tel:${phoneNumber}`)?.focus()
  }

  openLink(url: string) {
    window.open(url, '_blank')?.focus()
  }

  openContentIntegration(integrationId: string) {
    this.activeIntegrationId = integrationId
  }

  closeWidgetWindow() {
    this.activeIntegrationId = null
    this.windowOpen = false
    if (this.showGreeting) {
      this.tooltipActive = true
    }
  }

  displayGreeting() {
    if (this.settings?.project?.greeting && !this.isMobile) {
      setTimeout(() => {
        if (!this.windowOpen) {
          this.tooltipActive = true
        }
        this.showGreeting = true
      }, this.settings.project.greeting_timeout_seconds * 1000)
    }
  }

  iconUrl(icon: string) {
    return window.location.origin + '/icons/' + icon + '.svg'
  }

  closeActiveIntegration() {
    this.activeIntegrationId = null
  }

  loadSettings(userId: string, projectId: string) {
    fetch(
      `${process.env.BASE_URL_BACKEND}/v1/settings/${userId}/projects/${projectId}`
    )
      .then((response) => response.json())
      .then((response) => {
        this.settings = response.data
        this.settingsLoaded = true
        this.displayGreeting()
      })
  }
}
</script>

<style scoped lang="scss">
@import 'assets/microtip.scss';
.sb-widget {
  position: fixed;
  right: 48px;
  z-index: 10000;
  bottom: 48px;

  &__chat-head {
    height: 72px;
    width: 72px;
    display: flex;
    cursor: pointer;
    border-radius: 50%;
    justify-content: center;
    align-items: center;
  }

  &__window {
    width: 400px;
    height: 650px;
    background-color: #f3f4f6;
    border-radius: 16px;
    margin-bottom: 16px;

    &__header {
      height: 90px;
      border-top-left-radius: 12px;
      border-top-right-radius: 12px;
      padding: 16px;

      &__project-name {
        font-size: 2.25rem;
        font-weight: 400;
      }
      &__close-button {
        margin-left: auto;
        margin-top: -8px;
        margin-right: -16px;
        cursor: pointer;
      }
    }

    &__body {
      height: 560px;
      width: 100%;
      color: #21293c;
      position: relative;

      &__integration--active {
        padding-top: 8px;
        padding-left: 12px;
        border-bottom: 1px solid #4b587c;
        &__title {
          font-weight: bold;
          font-size: 16px;
        }
        &__text {
          padding-left: 16px;
          padding-right: 16px;
          margin-top: -8px;
          white-space: pre-line;
        }
        &__back-button {
          cursor: pointer;
          display: flex;
          justify-content: center;
          align-items: center;
          &:hover {
            path {
              fill: #32388d;
            }
          }
        }
      }

      &__mention {
        width: 100%;
        color: #4b587c;
        text-align: center;
        font-size: 13px;
        position: absolute;
        bottom: 16px;
        a {
          text-decoration: none !important;
          color: #4b587c;
          font-weight: bold;
        }
      }

      &__integration {
        width: 90%;
        padding: 12px;
        border-radius: 4px;
        margin: 12px auto;
        font-size: 18px;
        background-color: white;
        cursor: pointer;
        display: flex;
        box-shadow: 0 5px 7px -3.5px rgba(0, 0, 0, 0.2),
          0 6.5px 6px 1.5px rgba(0, 0, 0, 0.14),
          0 4.5px 11px 4px rgba(0, 0, 0, 0.12) !important;
      }

      &__integration--content {
        &:hover {
          background-color: #2196f3;
          color: white;
          path {
            fill: white;
          }
        }

        &__title {
          font-size: 16px;
          font-weight: bold;
        }
        &__text {
          font-size: 13px;
        }
      }

      &__integration--phone-call {
        display: flex;
        &:hover {
          background-color: #7ed766;
          color: white;
        }

        &__icon {
          background-color: #7ed766;
          border-radius: 2px;
          margin-right: 8px;
          height: 24px;
          width: 24px;
          padding: 4px;
        }
      }

      &__integration--link {
        display: flex;
        &:hover {
          background-color: #1e88e5;
          color: white;
        }

        &__icon {
          background-color: #1e88e5;
          border-radius: 2px;
          margin-right: 8px;
          height: 24px;
          width: 24px;
          padding: 4px;
        }
      }

      &__integration--whatsapp {
        display: flex;
        &:hover {
          background-color: #25d366;
          color: white;
        }

        &__icon {
          background-color: #25d366;
          border-radius: 2px;
          margin-right: 8px;
          height: 24px;
          width: 24px;
          padding: 4px;
        }
      }
    }
  }

  .sb-row {
    width: 100%;
    display: flex;
  }

  &__image {
    background-repeat: no-repeat;
    background-size: cover;
    transition: background-image 0.3s ease-in-out;
    height: 44px;
    width: 44px;
  }
}
</style>
