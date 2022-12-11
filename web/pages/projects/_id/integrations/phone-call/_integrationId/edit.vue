<template>
  <v-container v-if="$store.getters.authUser">
    <project-page-title>Edit Phone Integration</project-page-title>
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
            placeholder="e.g Call us on +1(800) 555-0199"
            outlined
            required
          ></v-text-field>
          <no-ssr>
            <v-phone-input
              v-model="formPhoneNumber"
              :disabled="savingIntegration"
              outlined
              label="Phone Number"
              persistent-placeholder
              :error="$store.getters.errorMessages.has('phone_number')"
              :error-messages="$store.getters.errorMessages.get('phone_number')"
              placeholder="Whatsapp phone number"
            >
            </v-phone-input>
          </no-ssr>
          <div class="d-flex">
            <loading-button
              :loading="savingIntegration"
              :icon="mdiContentSave"
              :large="$vuetify.breakpoint.lgAndUp"
              @click="saveIntegration"
            >
              Update
              <span v-if="$vuetify.breakpoint.lgAndUp" class="px-1">Phone</span>
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
import { mdiMenuDown, mdiDelete, mdiContentSave } from '@mdi/js'

export default {
  name: 'ProjectsIntegrationsPhoneCallEdit',
  layout: 'project',
  data() {
    return {
      mdiMenuDown,
      mdiDelete,
      mdiContentSave,
      savingIntegration: false,
      formName: '',
      formText: '',
      formPhoneNumber: '',
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
        .dispatch('getPhoneCallIntegration', {
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
     * @param {EntitiesPhoneCallIntegration} integration
     */
    setDefaults(integration) {
      this.formName = integration.name
      this.formText = integration.text
      this.formPhoneNumber = integration.phone_number
    },

    deleteIntegration() {
      this.savingIntegration = true
      this.$store
        .dispatch('deletePhoneCallIntegration', {
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
        .dispatch('updatePhoneCallIntegration', {
          projectId: this.$store.getters.activeProjectId,
          integrationId: this.$route.params.integrationId,
          name: this.formName,
          text: this.formText,
          phone_number: this.formPhoneNumber,
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
