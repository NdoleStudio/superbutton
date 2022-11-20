<template>
  <v-container v-if="$store.getters.authUser">
    <v-row class="pt-16">
      <v-col lg="4" md="8" offset-md="2" offset-lg="4">
        <v-card class="mt-16 pa-8">
          <v-card-text class="pt-3 text-center">
            <h4 class="text-h4 text--primary">Welcome</h4>
            <h5 class="pt-1 pb-6 text--secondary subtitle-1">
              Start by setting up your new project
            </h5>
            <v-form>
              <v-text-field
                v-model="formName"
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
              <loading-button
                :loading="$store.getters.creatingProject"
                :icon="mdiPlus"
                :block="true"
                :large="true"
                @click="createProject"
              >
                Create Project
              </loading-button>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiPlus } from '@mdi/js'

export default {
  name: 'ProjectsCreate',
  layout: 'default',
  data() {
    return {
      to: '/',
      formName: '',
      mdiPlus,
      formWebsite: '',
    }
  },
  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }
    await Promise.all([this.$store.dispatch('loadProjects')])

    if (!this.$store.getters.hasProjects) {
      await this.$router.push('/projects/create')
    }
  },
  methods: {
    createProject() {
      this.$store
        .dispatch('createProject', {
          name: this.formName,
          website: this.formWebsite,
        })
        .then((payload) => {
          this.$router.push({
            name: 'projects-id-settings',
            params: {
              id: payload.id,
            },
          })
        })
    },
  },
}
</script>
