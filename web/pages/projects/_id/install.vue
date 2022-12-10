<template>
  <v-container v-if="$store.getters.authUser">
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Install Widget</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-expansion-panels :value="0" class="mb-4">
          <v-expansion-panel>
            <v-expansion-panel-header>
              <div class="d-flex">
                <v-icon color="primary" large left>{{ mdiXml }}</v-icon>
                <h3 class="text-h4 font-weight-bold">HTML</h3>
              </div>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-divider></v-divider>
              <install-widget-html></install-widget-html>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiXml } from '@mdi/js'

export default {
  name: 'ProjectsInstall',
  layout: 'project',
  data() {
    return {
      mdiXml,
    }
  },
  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }
    await Promise.all([this.$store.dispatch('loadProjects')])
  },
}
</script>
