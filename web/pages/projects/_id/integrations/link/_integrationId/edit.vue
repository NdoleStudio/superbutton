<template>
  <v-container v-if="$store.getters.authUser">
    <v-row>
      <v-col class="d-flex">
        <back-button :icon="true" :large="true"></back-button>
        <h1 class="text-h4 ml-2 mb-4">Edit Link Integration</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" lg="8">
        <v-form>
          <v-text-field
            v-model="formName"
            :disabled="savingIntegration"
            :counter="30"
            label="Name"
            persistent-placeholder
            placeholder="e.g Customer Service"
            outlined
            class="mb-4"
            :error="$store.getters.errorMessages.has('name')"
            :error-messages="$store.getters.errorMessages.get('name')"
            required
          ></v-text-field>
          <v-text-field
            v-model="formText"
            :disabled="savingIntegration"
            :counter="30"
            class="mb-4"
            label="Text"
            persistent-placeholder
            :error="$store.getters.errorMessages.has('text')"
            :error-messages="$store.getters.errorMessages.get('text')"
            placeholder="e.g Visit our FAQ page"
            outlined
            required
          ></v-text-field>
          <v-text-field
            v-model="formWebsite"
            :disabled="savingIntegration"
            :counter="30"
            class="mb-4"
            label="Website"
            persistent-placeholder
            :error="$store.getters.errorMessages.has('website')"
            :error-messages="$store.getters.errorMessages.get('website')"
            placeholder="e.g https://example.com"
            outlined
            required
          ></v-text-field>
          <div class="d-flex">
            <loading-button
              :loading="savingIntegration"
              :icon="mdiPlus"
              :large="true"
              @click="saveIntegration"
            >
              Update Phone Integration
            </loading-button>
            <v-spacer></v-spacer>
            <v-btn
              large
              :disabled="savingIntegration"
              color="error"
              text
              @click="deleteIntegration"
            >
              <v-icon left>{{ mdiDelete }}</v-icon>
              Delete
            </v-btn>
          </div>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiPlus, mdiMenuDown, mdiDelete } from '@mdi/js'

export default {
  name: 'ProjectsIntegrationsLinkEdit',
  layout: 'project',
  data() {
    return {
      mdiPlus,
      mdiMenuDown,
      mdiDelete,
      savingIntegration: false,
      formName: '',
      formText: '',
      formWebsite: '',
    }
  },
  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }
    await Promise.all([this.$store.dispatch('loadProjects')])
    this.loadIntegration()
  },
  methods: {
    loadIntegration() {
      this.savingIntegration = true
      this.$store
        .dispatch('getLinkIntegration', {
          projectId: this.$store.getters.activeProjectId,
          integrationId: this.$route.params.integrationId,
        })
        .then((payload) => {
          this.setDefaults(payload)
        })
        .finally(() => {
          this.savingIntegration = false
        })
    },
    /**
     *
     * @param {EntitiesLinkIntegration} integration
     */
    setDefaults(integration) {
      this.formName = integration.name
      this.formText = integration.text
      this.formWebsite = integration.url
    },

    deleteIntegration() {
      this.savingIntegration = true
      this.$store
        .dispatch('deleteLinkIntegration', {
          projectId: this.$store.getters.activeProjectId,
          integrationId: this.$route.params.integrationId,
        })
        .then(() => {
          this.$router.push({
            name: 'projects-id-integrations',
            params: {
              id: this.$store.getters.activeProjectId,
            },
          })
        })
        .finally(() => {
          this.savingIntegration = false
        })
    },

    saveIntegration() {
      this.savingIntegration = true
      this.$store
        .dispatch('updateLinkIntegration', {
          projectId: this.$store.getters.activeProjectId,
          integrationId: this.$route.params.integrationId,
          name: this.formName,
          text: this.formText,
          website: this.formWebsite,
        })
        .then(() => {
          this.$router.push({
            name: 'projects-id-integrations',
            params: {
              id: this.$store.getters.activeProjectId,
            },
          })
        })
        .finally(() => {
          this.savingIntegration = false
        })
    },
  },
}
</script>
