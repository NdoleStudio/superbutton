<template>
  <v-container v-if="$store.getters.authUser" class="settings-page">
    <v-row>
      <v-col>
        <h1 class="text-h4" :class="{ 'mt-n4': $vuetify.breakpoint.mdAndDown }">
          Settings
        </h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12">
        <v-row class="mb-1">
          <v-col>
            <v-btn
              rounded
              :text="activeTab !== 0"
              :color="activeTab === 0 ? 'primary' : 'default'"
              @click="activeTab = 0"
            >
              Billing & Usage
            </v-btn>
          </v-col>
        </v-row>
        <v-tabs-items v-model="activeTab">
          <v-tab-item>
            <h2 class="text-h5 text--secondary mb-2">Current Plan</h2>
            <v-row>
              <v-col>
                <v-alert dense text :icon="undefined" prominent color="info">
                  <v-row align="center">
                    <v-col class="grow">
                      <h1
                        class="subtitle-1 font-weight-bold text-uppercase mt-3"
                      >
                        Free
                      </h1>
                      <p class="text--secondary">1/3 projects</p>
                    </v-col>
                    <v-col class="shrink">
                      <v-btn color="primary" :href="checkoutURL">
                        <v-icon left>{{ mdiStarOutline }}</v-icon>
                        Upgrade
                      </v-btn>
                    </v-col>
                  </v-row>
                </v-alert>
              </v-col>
            </v-row>
            <h2 class="text-h5 text--secondary mb-2">Upgrade Plan</h2>
            <v-row>
              <v-col cols="12" md="6">
                <v-card :href="checkoutURL" outlined>
                  <v-card-text>
                    <v-row align="center">
                      <v-col class="grow">
                        <h1
                          class="subtitle-1 font-weight-bold text-uppercase mt-3"
                        >
                          Pro - Monthly
                        </h1>
                        <p class="text--secondary">Unlimited projects</p>
                      </v-col>
                      <v-col class="shrink">
                        <span class="text-h5 text--primary">$6</span>/month
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
              <v-col cols="12" md="6">
                <v-card :href="checkoutURL" outlined>
                  <v-card-text>
                    <v-row align="center">
                      <v-col class="grow">
                        <h1
                          class="subtitle-1 font-weight-bold text-uppercase mt-3"
                        >
                          Pro - Yearly
                          <v-chip small color="primary" class="mt-n1"
                            >Save 20%</v-chip
                          >
                        </h1>
                        <p class="text--secondary">Unlimited projects</p>
                      </v-col>
                      <v-col class="shrink">
                        <span class="text-h5 text--primary">$5</span>/month
                      </v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
            <h2 class="text-h5 text--secondary mb-2 mt-4">Usage Statistics</h2>
            <v-row>
              <v-col cols="12" md="4">
                <v-alert
                  dark
                  dense
                  :icon="mdiCallMade"
                  prominent
                  color="primary"
                  text
                >
                  <h2 class="text-h4 font-weight-bold mt-4">
                    {{ $store.getters.projects.length }}
                  </h2>
                  <p class="text--secondary mt-n1">Active Projects</p>
                </v-alert>
              </v-col>
            </v-row>
          </v-tab-item>
        </v-tabs-items>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiPlus, mdiCallMade, mdiStarOutline } from '@mdi/js'

export default {
  name: 'SettingsIndex',
  layout: 'user',
  data() {
    return {
      to: '/',
      formName: '',
      mdiPlus,
      mdiStarOutline,
      mdiCallMade,
      activeTab: 0,
      formWebsite: '',
    }
  },
  computed: {
    checkoutURL() {
      const url = new URL(this.$config.checkoutURL)
      const user = this.$store.getters.authUser
      url.searchParams.append('checkout[custom][user_id]', user.uid)
      url.searchParams.append('checkout[email]', user.email)
      url.searchParams.append('checkout[name]', user.displayName)
      return url.toString()
    },
  },
  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }
    await Promise.all([
      this.$store.dispatch('loadUser'),
      this.$store.dispatch('loadProjects'),
    ])
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
<style lang="scss" scoped>
.settings-page {
  .v-tabs-items {
    background-color: transparent !important;
  }
}
</style>
