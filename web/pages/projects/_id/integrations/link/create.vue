<template>
  <v-container v-if="$store.getters.authUser">
    <project-page-title>Create Link Integration</project-page-title>
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
          <v-text-field
            v-model="formColor"
            :disabled="savingIntegration"
            :counter="7"
            class="mb-4"
            label="Color"
            persistent-placeholder
            :error="$store.getters.errorMessages.has('color')"
            :error-messages="$store.getters.errorMessages.get('color')"
            placeholder="#1E88E5"
            outlined
            required
          >
            <template #append>
              <v-icon :color="widgetColor">{{ mdiSquare }}</v-icon>
            </template>
          </v-text-field>
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
                :color="widgetColor"
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
          <loading-button
            :loading="savingIntegration"
            :icon="mdiPlus"
            :large="true"
            :block="$vuetify.breakpoint.mdAndDown"
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
import { mdiPlus, mdiMenuDown, mdiSquare } from '@mdi/js'

export default {
  name: 'ProjectsIntegrationsLinkCreate',
  layout: 'project',
  data() {
    return {
      mdiPlus,
      mdiMenuDown,
      mdiSquare,
      savingIntegration: false,
      formName: '',
      formText: '',
      formIcon: 'link',
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
        {
          text: 'Map Icon',
          value: 'map',
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

    if (this.$store.getters.activeProjectId !== this.$route.params.id) {
      await this.$router.push('/')
    }
  },
  methods: {
    getIconURL(value) {
      return window.location.origin + '/icons/' + value + '.svg'
    },
    saveIntegration() {
      this.savingIntegration = true
      this.$store
        .dispatch('addLinkIntegration', {
          projectId: this.$store.getters.activeProjectId,
          name: this.formName,
          text: this.formText,
          color: this.formColor,
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
