<template>
  <v-container v-if="$store.getters.authUser">
    <project-page-title>Edit Link Integration</project-page-title>
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
            :counter="255"
            class="mb-4"
            label="Website"
            persistent-placeholder
            :error="$store.getters.errorMessages.has('website')"
            :error-messages="$store.getters.errorMessages.get('website')"
            placeholder="e.g https://example.com"
            outlined
            required
          ></v-text-field>
          <v-select
            v-model="formIcon"
            class="mb-4"
            :disabled="savingIntegration"
            :items="linkIcons"
            :error="$store.getters.errorMessages.has('icon')"
            :error-messages="$store.getters.errorMessages.get('icon')"
            outlined
            aria-required="true"
            label="Icon"
          >
            <template #item="{ item, attrs, on }">
              <v-list-item v-bind="attrs" v-on="on">
                <v-list-item-title>{{ item.text }}</v-list-item-title>
                <v-list-item-action>
                  <v-avatar size="40" tile :color="widgetColor">
                    <v-img
                      :src="getIconURL(item.value)"
                      width="25"
                      height="25"
                      contain
                    ></v-img>
                  </v-avatar>
                </v-list-item-action>
              </v-list-item>
            </template>
            <template #append>
              <v-avatar
                v-if="formIcon"
                tile
                color="#1E88E5"
                size="40"
                class="mt-n2"
              >
                <v-img
                  :src="getIconURL(formIcon)"
                  width="25"
                  height="25"
                  contain
                ></v-img>
              </v-avatar>
              <v-icon v-else>{{ mdiMenuDown }}</v-icon>
            </template>
          </v-select>
          <v-text-field
            v-model="formColor"
            :disabled="savingIntegration"
            :counter="7"
            class="mb-4"
            label="Color"
            persistent-placeholder
            :error="$store.getters.errorMessages.has('color')"
            :error-messages="$store.getters.errorMessages.get('color')"
            placeholder="e.g #1E88E5"
            outlined
            required
          >
            <template #append>
              <v-icon :color="widgetColor">{{ mdiSquare }}</v-icon>
            </template>
          </v-text-field>
          <div class="d-flex">
            <loading-button
              :loading="savingIntegration"
              :icon="mdiContentSave"
              :large="$vuetify.breakpoint.lgAndUp"
              @click="saveIntegration"
            >
              Update
              <span v-if="$vuetify.breakpoint.lgAndUp" class="px-1">Link</span>
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
import { mdiMenuDown, mdiDelete, mdiContentSave, mdiSquare } from '@mdi/js'

export default {
  name: 'ProjectsIntegrationsLinkEdit',
  layout: 'project',
  data() {
    return {
      mdiContentSave,
      mdiMenuDown,
      mdiDelete,
      mdiSquare,
      savingIntegration: false,
      formName: '',
      formIcon: 'link',
      formText: '',
      formWebsite: '',
      formColor: '#1E88E5',
      linkIcons: [
        {
          text: 'Link Icon',
          value: 'link',
        },
        {
          text: 'Documentation Icon',
          value: 'documentation',
        },
        {
          text: 'Mail Icon',
          value: 'mail',
        },
        {
          text: 'Github Icon',
          value: 'github',
        },
      ],
    }
  },
  computed: {
    widgetColor() {
      return /^#[0-9A-F]{6}$/i.test(this.formColor) ? this.formColor : '#1E88E5'
    },
  },
  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }
    await Promise.all([this.$store.dispatch('loadProjects')])
    this.loadIntegration()
  },
  methods: {
    getIconURL(value) {
      return window.location.origin + '/icons/' + value + '.svg'
    },
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
      this.formColor = integration.color ? integration.color : this.formColor
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
          icon: this.formIcon,
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
