<template>
  <v-container v-if="$store.getters.authUser">
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Settings</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" lg="7">
        <v-form>
          <v-text-field
            v-model="formName"
            :disabled="savingProject"
            :counter="30"
            label="Name"
            persistent-placeholder
            placeholder="Project name e.g Google"
            outlined
            class="mb-4"
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
            :error-messages="$store.getters.errorMessages.get('website')"
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
            :error-messages="$store.getters.errorMessages.get('greeting')"
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
            :error="$store.getters.errorMessages.has('greeting_timeout')"
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
      </v-col>
      <v-col cols="12" lg="5">
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
      savingProject: false,
      formName: '',
      formWebsite: '',
      formColor: '',
      formIcon: '',
      copyButtonActive: true,
      formGreeting: '',
      formGreetingTimeoutSeconds: '',
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

    if (this.$store.getters.activeProjectId !== this.$route.params.id) {
      await this.$router.push('/')
    }
  },
  methods: {
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
