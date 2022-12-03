<template>
  <v-container v-if="$store.getters.authUser">
    <v-row>
      <v-col class="d-flex">
        <back-button :icon="true" :large="true"></back-button>
        <h1 class="ml-2 text-h4 mb-4">Create Link Integration</h1>
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
            placeholder="e.g Visit our FAQ Page"
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
          <loading-button
            :loading="savingIntegration"
            :icon="mdiPlus"
            :large="true"
            @click="saveIntegration"
          >
            Add Link Integration
          </loading-button>
        </v-form>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiPlus, mdiMenuDown } from '@mdi/js'

export default {
  name: 'ProjectsIntegrationsLinkCreate',
  layout: 'project',
  data() {
    return {
      mdiPlus,
      mdiMenuDown,
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

    if (this.$store.getters.activeProjectId !== this.$route.params.id) {
      await this.$router.push('/')
    }
  },
  methods: {
    saveIntegration() {
      this.savingIntegration = true
      this.$store
        .dispatch('addLinkIntegration', {
          projectId: this.$store.getters.activeProjectId,
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
