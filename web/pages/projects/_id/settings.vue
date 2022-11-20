<template>
  <v-container v-if="$store.getters.authUser">
    <v-row class="pt-16">
      <v-col class="text-center">
        <h1 class="text-h2">Project Settings</h1>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiPlus } from '@mdi/js'

export default {
  name: 'ProjectsSettings',
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
      // await this.$router.push('/projects/create')
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
