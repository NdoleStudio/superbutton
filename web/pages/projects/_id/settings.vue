<template>
  <v-container v-if="$store.getters.authUser">
    <v-row>
      <v-col>
        <h1 class="text-h4">Settings</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" lg="7">
        <v-row class="mb-1">
          <v-col>
            <v-btn
              rounded
              :text="activeTab !== 0"
              :color="activeTab === 0 ? 'primary' : 'default'"
              @click="activeTab = 0"
            >
              Project
            </v-btn>
            <v-btn
              rounded
              :text="activeTab !== 1"
              :color="activeTab === 1 ? 'primary' : 'default'"
              @click="activeTab = 1"
            >
              Integration
            </v-btn>
          </v-col>
        </v-row>
        <v-tabs-items v-model="activeTab">
          <v-tab-item>
            <v-card>
              <v-card-text>
                <v-form>
                  <v-text-field
                    v-model="formName"
                    :disabled="savingProject"
                    :counter="30"
                    label="Name"
                    persistent-placeholder
                    placeholder="Project name e.g Google"
                    outlined
                    class="mb-4 mt-2"
                    :error="$store.getters.errorMessages.has('name')"
                    :error-messages="$store.getters.errorMessages.get('name')"
                    required
                  ></v-text-field>
                  <v-text-field
                    v-model="formWebsite"
                    :disabled="savingProject"
                    :counter="255"
                    class="mb-4"
                    label="Website"
                    persistent-placeholder
                    :error="$store.getters.errorMessages.has('website')"
                    :error-messages="
                      $store.getters.errorMessages.get('website')
                    "
                    placeholder="Website URL e.g https://google.com"
                    outlined
                    required
                  ></v-text-field>
                  <v-select
                    v-model="formIcon"
                    class="mb-4"
                    :disabled="savingProject"
                    :items="projectIcons"
                    :error="$store.getters.errorMessages.has('icon')"
                    :error-messages="$store.getters.errorMessages.get('icon')"
                    outlined
                    aria-required="true"
                    label="Widget Icon"
                  >
                    <template #item="{ item, attrs, on }">
                      <v-list-item v-bind="attrs" v-on="on">
                        <v-list-item-title>{{ item.text }}</v-list-item-title>
                        <v-list-item-action>
                          <v-avatar size="40" :color="widgetColor">
                            <v-img
                              :src="getIconURL(item.value)"
                              width="20"
                              height="20"
                              contain
                            ></v-img>
                          </v-avatar>
                        </v-list-item-action>
                      </v-list-item>
                    </template>
                    <template #append>
                      <v-avatar
                        v-if="formIcon"
                        :color="widgetColor"
                        size="40"
                        class="mt-n2"
                      >
                        <v-img
                          :src="getIconURL(formIcon)"
                          width="20"
                          height="20"
                          contain
                        ></v-img>
                      </v-avatar>
                      <v-icon v-else>{{ mdiMenuDown }}</v-icon>
                    </template>
                  </v-select>
                  <v-text-field
                    v-model="formColor"
                    :disabled="savingProject"
                    :counter="7"
                    class="mb-4"
                    label="Color"
                    persistent-placeholder
                    :error="$store.getters.errorMessages.has('color')"
                    :error-messages="$store.getters.errorMessages.get('color')"
                    placeholder="#283593"
                    outlined
                    required
                  >
                    <template #append>
                      <v-icon :color="widgetColor">{{ mdiSquare }}</v-icon>
                    </template>
                  </v-text-field>
                  <v-text-field
                    v-model="formGreeting"
                    :disabled="savingProject"
                    :counter="30"
                    class="mb-4"
                    label="Greeting"
                    persistent-placeholder
                    :error="$store.getters.errorMessages.has('greeting')"
                    :error-messages="
                      $store.getters.errorMessages.get('greeting')
                    "
                    placeholder="e.g Need some help?"
                    outlined
                    required
                  ></v-text-field>
                  <v-text-field
                    v-model="formGreetingTimeoutSeconds"
                    :disabled="savingProject"
                    type="number"
                    class="mb-4"
                    label="Greeting Timeout (seconds)"
                    persistent-placeholder
                    :error="
                      $store.getters.errorMessages.has('greeting_timeout')
                    "
                    :error-messages="
                      $store.getters.errorMessages.get('greeting_timeout')
                    "
                    placeholder="10"
                    outlined
                    required
                  ></v-text-field>
                  <div class="d-flex">
                    <loading-button
                      :loading="savingProject"
                      :icon="mdiContentSave"
                      :large="true"
                      @click="updateProject"
                    >
                      Update Project
                    </loading-button>
                    <v-spacer></v-spacer>
                    <v-btn
                      large
                      :disabled="savingProject"
                      color="error"
                      text
                      @click="deleteProject"
                    >
                      <v-icon left>{{ mdiDelete }}</v-icon>
                      Delete
                    </v-btn>
                  </div>
                </v-form>
              </v-card-text>
            </v-card>
          </v-tab-item>
          <v-tab-item>
            <v-card>
              <v-card-title
                >Drag and drop integrations to change the order.</v-card-title
              >
              <v-card-text>
                <vue-draggable v-model="integrations" :disabled="savingProject">
                  <transition-group>
                    <v-card
                      v-for="integration in integrations"
                      :key="integration"
                      elevation="2"
                      class="mb-1 mt-1"
                      style="cursor: move"
                    >
                      <v-card-text class="pa-2">
                        <div class="d-flex align-center">
                          <v-avatar
                            size="32"
                            tile
                            :color="getIntegrationColor(integration)"
                          >
                            <v-img
                              :src="getIntegrationIcon(integration)"
                              width="20"
                              height="20"
                              contain
                            ></v-img>
                          </v-avatar>
                          <p class="subtitle-1 ml-4 mb-n1">
                            {{ getIntegrationText(integration) }}
                          </p>
                          <v-spacer></v-spacer>
                          <v-icon large>{{ mdiDrag }}</v-icon>
                        </div>
                      </v-card-text>
                    </v-card>
                  </transition-group>
                </vue-draggable>
                <loading-button
                  class="mt-4"
                  :loading="savingProject"
                  :icon="mdiContentSave"
                  :large="true"
                  @click="updateProjectIntegrations"
                >
                  Update Integrations
                </loading-button>
              </v-card-text>
            </v-card>
          </v-tab-item>
        </v-tabs-items>
      </v-col>
      <v-col cols="12" lg="5" class="mt-12">
        <install-widget-html color="default"></install-widget-html>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import {
  mdiPlus,
  mdiContentSave,
  mdiSquare,
  mdiMenuDown,
  mdiDelete,
  mdiDrag,
} from '@mdi/js'

export default {
  name: 'ProjectsSettings',
  layout: 'project',
  data() {
    return {
      to: '/',
      mdiPlus,
      mdiContentSave,
      mdiSquare,
      mdiMenuDown,
      mdiDelete,
      mdiDrag,
      savingProject: false,
      formName: '',
      formWebsite: '',
      formColor: '',
      formIcon: '',
      activeTab: 0,
      copyButtonActive: true,
      formGreeting: '',
      formGreetingTimeoutSeconds: '',
      projectSettings: null,
      projectIcons: [
        {
          text: 'Chat Icon',
          value: 'chat',
        },
        {
          text: 'Whatsapp Icon',
          value: 'whatsapp',
        },
        {
          text: 'Help Chat Icon',
          value: 'help-chat',
        },
      ],
      integrations: [],
    }
  },
  computed: {
    widgetColor() {
      return /^#[0-9A-F]{6}$/i.test(this.formColor) ? this.formColor : '#283593'
    },
  },

  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }

    await Promise.all([this.$store.dispatch('loadProjects')])

    this.setFormDefaults(this.$store.getters.activeProject)
    this.loadIntegrations()

    if (this.$store.getters.activeProjectId !== this.$route.params.id) {
      await this.$router.push('/')
    }
  },
  methods: {
    loadIntegrations() {
      this.$store
        .dispatch('getProjectSettings', this.$store.getters.activeProjectId)
        .then((settings) => {
          this.projectSettings = settings
          this.integrations = settings.integrations.map(
            (integration) => integration.id
          )
        })
    },

    /**
     * @param {string} integrationId
     * @returns {string}
     */
    getIntegrationIcon(integrationId) {
      const integration = this.getIntegration(integrationId)
      return this.getIconURL(integration.settings.icon ?? integration.type)
    },

    /**
     * @param {string} integrationId
     * @returns {string}
     */
    getIntegrationColor(integrationId) {
      switch (this.getIntegration(integrationId).type) {
        case 'whatsapp':
          return '#25d366'
        case 'phone-call':
          return '#7ed766'
        case 'content':
          return '#2196f3'
        case 'link':
          return '#1e88e5'
      }
    },

    /**
     * @param {string} integrationId
     * @returns {string}
     */
    getIntegrationText(integrationId) {
      const integration = this.getIntegration(integrationId)
      switch (integration.type) {
        case 'content':
          return integration.settings.title
        default:
          return integration.settings.text
      }
    },

    /**
     *
     * @param {string} integrationId
     * @returns {EntitiesProjectSettingsIntegration}
     */
    getIntegration(integrationId) {
      return this.projectSettings.integrations.find(
        (x) => x.id === integrationId
      )
    },

    getIconURL(value) {
      return window.location.origin + '/icons/' + value + '.svg'
    },
    /**
     * @param {EntitiesProject} project
     */
    setFormDefaults(project) {
      if (project === null) {
        return
      }
      this.formIcon = project.icon
      this.formName = project.name
      this.formWebsite = project.url
      this.formColor = project.color
      this.formGreetingTimeoutSeconds = project.greeting_timeout_seconds
      this.formGreeting = project.greeting
    },
    updateProject() {
      this.savingProject = true
      this.$store
        .dispatch('updateProject', {
          projectId: this.$store.getters.activeProjectId,
          name: this.formName,
          website: this.formWebsite,
          icon: this.formIcon,
          greeting: this.formGreeting,
          greeting_timeout: parseInt(this.formGreetingTimeoutSeconds) ?? 0,
          color: this.formColor,
        })
        .finally(() => {
          this.savingProject = false
        })
    },
    updateProjectIntegrations() {
      this.savingProject = true
      this.$store
        .dispatch('updateProjectIntegrations', {
          projectId: this.$store.getters.activeProjectId,
          order: this.integrations,
        })
        .finally(() => {
          this.savingProject = false
        })
    },
    deleteProject() {
      this.savingProject = true
      this.$store
        .dispatch('deleteProject', this.$store.getters.activeProjectId)
        .then(() => {
          if (this.$store.getters.hasProjects) {
            this.$router.push({
              name: 'projects-id-integrations',
              params: {
                id: this.$store.getters.activeProjectId,
              },
            })
            return
          }
          this.$router.push({
            name: 'projects-create',
          })
        })
        .finally(() => {
          this.savingProject = false
        })
    },
  },
}
</script>
