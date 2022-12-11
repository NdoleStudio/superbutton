<template>
  <v-container v-if="$store.getters.authUser">
    <project-page-title>Edit Content Integration</project-page-title>
    <v-row>
      <v-col cols="12" lg="8">
        <v-form>
          <v-text-field
            v-model="formName"
            :disabled="savingIntegration"
            :counter="30"
            label="Name"
            persistent-placeholder
            placeholder="e.g FAQ"
            outlined
            class="mb-4"
            :error="$store.getters.errorMessages.has('name')"
            :error-messages="$store.getters.errorMessages.get('name')"
            required
          ></v-text-field>
          <v-text-field
            v-model="formTitle"
            :disabled="savingIntegration"
            :counter="50"
            class="mb-4"
            label="Title"
            persistent-placeholder
            :error="$store.getters.errorMessages.has('title')"
            :error-messages="$store.getters.errorMessages.get('title')"
            placeholder="e.g How to create smart buttons."
            outlined
            required
          ></v-text-field>
          <v-textarea
            v-model="formSummary"
            :disabled="savingIntegration"
            :counter="300"
            class="mb-4"
            label="Summary"
            :rows="3"
            persistent-placeholder
            :error="$store.getters.errorMessages.has('summary')"
            :error-messages="$store.getters.errorMessages.get('summary')"
            placeholder="e.g This is a summary that appears under the title"
            outlined
            required
          ></v-textarea>
          <v-textarea
            v-model="formText"
            :disabled="savingIntegration"
            :counter="1000"
            class="mb-4"
            label="Text"
            :rows="6"
            persistent-placeholder
            :error="$store.getters.errorMessages.has('text')"
            :error-messages="$store.getters.errorMessages.get('text')"
            placeholder="This is a complete text for the content you want to add. It can be up to 1000 characters"
            outlined
            required
          ></v-textarea>
          <div class="d-flex">
            <loading-button
              :loading="savingIntegration"
              :icon="mdiContentSave"
              :large="$vuetify.breakpoint.lgAndUp"
              @click="saveIntegration"
            >
              Update
              <span v-if="$vuetify.breakpoint.lgAndUp" class="px-1"
                >Content</span
              >
              Integration
            </loading-button>
            <v-spacer></v-spacer>
            <v-btn
              :large="$vuetify.breakpoint.lgAndUp"
              :disabled="savingIntegration"
              color="error"
              text
              @click="deleteIntegration"
            >
              <v-icon v-if="$vuetify.breakpoint.lgAndUp" left>{{
                mdiDelete
              }}</v-icon>
              Delete
            </v-btn>
          </div>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiPlus, mdiMenuDown, mdiDelete, mdiContentSave } from '@mdi/js'

export default {
  name: 'ProjectsIntegrationsWhatsappEdit',
  layout: 'project',
  data() {
    return {
      mdiPlus,
      mdiMenuDown,
      mdiDelete,
      mdiContentSave,
      savingIntegration: false,
      formName: '',
      content: '',
      formTitle: '',
      formSummary: '',
      formText: '',
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
        .dispatch('getContentIntegration', {
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
     * @param {EntitiesContentIntegration} integration
     */
    setDefaults(integration) {
      this.formName = integration.name
      this.formText = integration.text
      this.formSummary = integration.summary
      this.formTitle = integration.title
    },

    deleteIntegration() {
      this.savingIntegration = true
      this.$store
        .dispatch('deleteContentIntegration', {
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
        .dispatch('updateContentIntegration', {
          projectId: this.$store.getters.activeProjectId,
          integrationId: this.$route.params.integrationId,
          name: this.formName,
          text: this.formText,
          title: this.formTitle,
          summary: this.formSummary,
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
