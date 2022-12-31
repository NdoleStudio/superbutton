<template>
  <v-container
    v-if="$store.getters.authUser && $store.getters.user"
    class="settings-page"
  >
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
                <v-alert dense text prominent color="info">
                  <v-row align="center">
                    <v-col cols="12">
                      <h1
                        class="subtitle-1 font-weight-bold text-uppercase mt-3"
                      >
                        {{ plan.name }}
                      </h1>
                      <p
                        v-if="plan.price && !subscriptionIsCancelled"
                        class="text--secondary"
                      >
                        Your next bill is for <b>${{ plan.price }}</b> on
                        <b>{{
                          new Date(
                            $store.getters.user.subscription_renews_at
                          ).toLocaleDateString()
                        }}</b>
                      </p>
                      <p
                        v-else-if="plan.price && subscriptionIsCancelled"
                        class="text--secondary"
                      >
                        You will be downgraded to the <b>FREE</b> plan on
                        <b>{{
                          new Date(
                            $store.getters.user.subscription_ends_at
                          ).toLocaleDateString()
                        }}</b>
                      </p>
                      <p v-else class="text--secondary">
                        {{ $store.getters.projects.length }}/3 projects
                      </p>
                    </v-col>
                    <v-col cols="12" class="d-flex mb-2 mt-n6">
                      <loading-button
                        v-if="!subscriptionIsCancelled"
                        color="primary"
                        :loading="loading"
                        @click="updateDetails"
                      >
                        Update Details
                      </loading-button>
                      <v-btn v-else color="primary" :href="checkoutURL"
                        >Upgrade Plan</v-btn
                      >
                      <v-spacer></v-spacer>
                      <v-dialog
                        v-if="!subscriptionIsCancelled"
                        v-model="dialog"
                        max-width="590"
                      >
                        <template #activator="{ on, attrs }">
                          <v-btn v-bind="attrs" color="error" text v-on="on">
                            Cancel Plan
                          </v-btn>
                        </template>
                        <v-card>
                          <v-card-text class="pt-4 mb-n6">
                            <h2 class="text--primary text-h5 mb-2">
                              Are you sure you want to cancel your subscription?
                            </h2>
                            <p>
                              You will be downgraded to the free plan at the end
                              of the current billing period on
                              <b>{{
                                new Date(
                                  $store.getters.user.subscription_renews_at
                                ).toLocaleDateString()
                              }}</b>
                            </p>
                          </v-card-text>
                          <v-card-actions>
                            <v-btn color="primary" @click="dialog = false">
                              Keep Subscription
                            </v-btn>
                            <v-spacer></v-spacer>
                            <loading-button
                              :text="true"
                              :loading="loading"
                              color="error"
                              @click="cancelPlan"
                            >
                              Cancel Plan
                            </loading-button>
                          </v-card-actions>
                        </v-card>
                      </v-dialog>
                    </v-col>
                  </v-row>
                </v-alert>
              </v-col>
            </v-row>
            <h2 v-if="plan.price === 0" class="text-h5 text--secondary mb-2">
              Upgrade Plan
            </h2>
            <v-row v-if="plan.price === 0">
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
      dialog: false,
      loading: false,
      mdiStarOutline,
      mdiCallMade,
      activeTab: 0,
      formWebsite: '',
      plans: [
        {
          name: 'Free',
          id: 'free',
          price: 0,
        },
        {
          name: 'PRO - Monthly',
          id: 'pro-monthly',
          price: 6,
        },
        {
          name: 'PRO - Yearly',
          id: 'pro-yearly',
          price: 60,
        },
      ],
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
    plan() {
      return this.plans.find(
        (x) => x.id === (this.$store.getters.user?.subscription_name || 'free')
      )
    },
    subscriptionIsCancelled() {
      return this.$store.getters.user?.subscription_status === 'cancelled'
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
    updateDetails() {
      this.loading = true
      this.$store
        .dispatch('getSubscriptionUpdateLink')
        .then((link) => {
          window.location.href = link
        })
        .catch(() => {
          this.loading = false
        })
    },
    cancelPlan() {
      this.loading = true
      this.$store
        .dispatch('cancelSubscription')
        .then(() => {
          this.$store.dispatch('addNotification', {
            message: 'Subscription cancelled successfully',
            type: 'success',
          })
          this.$router.push('/')
        })
        .catch(() => {
          this.loading = false
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
