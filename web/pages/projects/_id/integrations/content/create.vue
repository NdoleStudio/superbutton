<template>
  <v-container v-if="$store.getters.authUser">
    <v-row>
      <v-col class="d-flex">
        <back-button :icon="true" :large="true"></back-button>
        <h1 class="text-h4 ml-2 mb-4">Add Content</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" lg="8" xl="6">
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
          <loading-button
            :loading="savingIntegration"
            :icon="mdiPlus"
            :large="true"
            @click="saveIntegration"
          >
            Add Content
          </loading-button>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiPlus, mdiMenuDown } from '@mdi/js'

export default {
  name: 'ProjectsIntegrationsContentCreate',
  layout: 'project',
  data() {
    return {
      mdiPlus,
      mdiMenuDown,
      savingIntegration: false,
      formName: '',
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
  },
  methods: {
    saveIntegration() {
      this.savingIntegration = true
      this.$store
        .dispatch('addContentIntegration', {
          projectId: this.$store.getters.activeProjectId,
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
