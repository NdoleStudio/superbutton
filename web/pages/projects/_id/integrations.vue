<template>
  <v-container v-if="$store.getters.authUser">
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Integrations</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-expansion-panels :value="0">
          <v-expansion-panels :value="0" class="mb-4">
            <v-expansion-panel>
              <v-expansion-panel-header>
                <div class="d-flex">
                  <v-icon color="#25D366" class="mr-4">{{
                    mdiWhatsapp
                  }}</v-icon>
                  <h3 class="text-h6 font-weight-bold">Whatsapp</h3>
                </div>
              </v-expansion-panel-header>
              <v-expansion-panel-content>
                <v-btn class="primary" small>
                  <v-icon left>{{ mdiPlus }}</v-icon>
                  Add Whatsapp
                </v-btn>
              </v-expansion-panel-content>
            </v-expansion-panel>
          </v-expansion-panels>
          <v-expansion-panel>
            <v-expansion-panel-header>
              <div class="d-flex">
                <v-icon color="primary" class="mr-4">{{
                  mdiStickerText
                }}</v-icon>
                <h3 class="text-h6 font-weight-bold">Content</h3>
              </div>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-btn class="primary" small>
                <v-icon left>{{ mdiPlus }}</v-icon>
                Add Content
              </v-btn>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiPlus, mdiStickerText, mdiWhatsapp } from '@mdi/js'

export default {
  name: 'ProjectsIntegrations',
  layout: 'project',
  data() {
    return {
      to: '/',
      formName: '',
      mdiPlus,
      mdiStickerText,
      mdiWhatsapp,
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
