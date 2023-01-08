<template>
  <div>
    <v-badge
      class="logo-badge"
      :class="{ 'logo-badge--mobile': $vuetify.breakpoint.mdAndDown }"
      color="#8338ec"
      content="Beta"
    >
      <nuxt-link to="/" class="text-decoration-none d-flex">
        <v-avatar tile size="33" class="mt-1">
          <v-img contain :src="require('@/static/logo.svg')"></v-img>
        </v-avatar>
        <h3
          class="text-h4 font-weight-thin ml-1"
          :class="{
            'text--secondary': $vuetify.theme.dark,
            'text--primary': !$vuetify.theme.dark,
          }"
        >
          SuperButton
        </h3>
      </nuxt-link>
    </v-badge>
    <p>
      SuperButton is an
      <a
        class="text-decoration-none"
        href="https://github.com/NdoleStudio/superbutton"
        >open source</a
      >
      widget that lets you communicate with your website visitors using their
      preferred channels.
    </p>
  </div>
</template>

<script>
export default {
  name: 'BlogAbout',
}
</script>
