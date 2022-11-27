<template>
  <v-container v-if="$store.getters.authUser">
    <v-row>
      <v-col>
        <h1 class="text-h4 mb-4">Integrations</h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-expansion-panels :value="0" class="mb-4" readonly>
          <v-expansion-panel>
            <v-expansion-panel-header>
              <div class="d-flex">
                <v-icon color="#25D366" class="mr-4">{{ mdiWhatsapp }}</v-icon>
                <h3 class="text-h6 font-weight-bold">Whatsapp</h3>
                <v-progress-circular
                  v-if="loadingIntegrations"
                  class="ml-2 mt-2"
                  size="16"
                  width="1"
                  indeterminate
                  color="primary"
                ></v-progress-circular>
              </div>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-divider></v-divider>
              <v-simple-table v-if="whatsappIntegrations.length">
                <template #default>
                  <thead class="text-uppercase">
                    <tr>
                      <th class="text-left" style="width: 30%">Name</th>
                      <th class="text-left" style="width: 50%">Identifier</th>
                      <th class="">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="item in whatsappIntegrations" :key="item.id">
                      <td>{{ item.name }}</td>
                      <td>{{ item.integration_id }}</td>
                      <td>
                        <v-btn
                          small
                          class="secondary"
                          :to="{
                            name: 'projects-id-integrations-whatsapp-integrationId-edit',
                            params: {
                              id: $store.getters.activeProjectId,
                              integrationId: item.integration_id,
                            },
                          }"
                        >
                          <v-icon left>{{ mdiSquareEditOutline }}</v-icon>
                          Edit
                        </v-btn>
                      </td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
              <v-btn
                :to="{
                  name: 'projects-id-integrations-whatsapp-create',
                  params: { id: $store.getters.activeProjectId },
                }"
                class="primary mt-4"
                small
              >
                <v-icon left>{{ mdiPlus }}</v-icon>
                Add Whatsapp
              </v-btn>
            </v-expansion-panel-content>
          </v-expansion-panel>
        </v-expansion-panels>
        <v-expansion-panels :value="0" class="mb-4" readonly>
          <v-expansion-panel>
            <v-expansion-panel-header>
              <div class="d-flex">
                <v-icon color="primary" class="mr-4">{{
                  mdiStickerText
                }}</v-icon>
                <h3 class="text-h6 font-weight-bold">Content</h3>
                <v-progress-circular
                  v-if="loadingIntegrations"
                  class="ml-2 mt-2"
                  size="16"
                  width="1"
                  indeterminate
                  color="primary"
                ></v-progress-circular>
              </div>
            </v-expansion-panel-header>
            <v-expansion-panel-content>
              <v-divider></v-divider>
              <v-simple-table v-if="whatsappIntegrations.length">
                <template #default>
                  <thead class="text-uppercase">
                    <tr>
                      <th class="text-left" style="width: 30%">Name</th>
                      <th class="text-left" style="width: 50%">Identifier</th>
                      <th class="">Action</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="item in contentIntegrations" :key="item.id">
                      <td>{{ item.name }}</td>
                      <td>{{ item.integration_id }}</td>
                      <td>
                        <v-btn
                          small
                          class="secondary"
                          :to="{
                            name: 'projects-id-integrations-content-integrationId-edit',
                            params: {
                              id: $store.getters.activeProjectId,
                              integrationId: item.integration_id,
                            },
                          }"
                        >
                          <v-icon left>{{ mdiSquareEditOutline }}</v-icon>
                          Edit
                        </v-btn>
                      </td>
                    </tr>
                  </tbody>
                </template>
              </v-simple-table>
              <v-btn
                :to="{
                  name: 'projects-id-integrations-content-create',
                  params: { id: $store.getters.activeProjectId },
                }"
                class="primary mt-4"
                small
              >
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
import {
  mdiPlus,
  mdiStickerText,
  mdiWhatsapp,
  mdiSquareEditOutline,
} from '@mdi/js'

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
      mdiSquareEditOutline,
      projectIntegrations: [],
      formWebsite: '',
      loadingIntegrations: false,
    }
  },
  computed: {
    /**
     * @returns {EntitiesProjectIntegration[]}
     */
    whatsappIntegrations() {
      return this.projectIntegrations.filter((integration) => {
        return integration.type === 'whatsapp'
      })
    },
    /**
     * @returns {EntitiesProjectIntegration[]}
     */
    contentIntegrations() {
      return this.projectIntegrations.filter((integration) => {
        return integration.type === 'content'
      })
    },
  },
  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }
    await Promise.all([this.$store.dispatch('loadProjects')])
    this.loadIntegrations()
  },
  methods: {
    loadIntegrations() {
      this.loadingIntegrations = true
      this.$store
        .dispatch('getProjectIntegrations', this.$store.getters.activeProjectId)
        .then((payload) => {
          this.projectIntegrations = payload
        })
        .finally(() => {
          this.loadingIntegrations = false
        })
    },
  },
}
</script>
