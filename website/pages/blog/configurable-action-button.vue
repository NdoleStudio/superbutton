<template>
  <v-container>
    <v-row class="mt-8">
      <v-col cols="12" md="9">
        <h1 class="text-h3">
          Configurable floating action button on your website
        </h1>
        <p class="subtitle-2 mt-2">
          <span class="text-uppercase primary--text-color">{{ postDate }}</span>
          • <span class="text-uppercase">{{ readTime }}</span>
        </p>
        <p class="text--secondary mt-2">
          Add a SuperButton so you can communicate with your website visitors on
          whatsApp, email, contact forms, or direct them to external links like
          your documentation pages.
        </p>
        <h3 class="text-h4">What is a floating action button (FAB)?</h3>
        <p>
          A floating action button is a configurable button on a website or
          mobile app. It is usually found on the bottom right portion of a
          website or app and when you click on this button you get a list of
          control options or actions depending on the configuration. You can see
          an example floating button below 👇
        </p>
        <vue-glow
          color="#8338ec"
          mode="hex"
          elevation="14"
          :intensity="1.07"
          intense
        >
          <v-img :src="require('@/static/blog/fab.svg')"></v-img>
        </vue-glow>
        <h3 class="text-h4 mt-12">Create a project.</h3>
        <p>
          To integrate SuperButton to your website, you need to
          <a class="text-decoration-none" :href="$config.dashboardURL"
            >create a project</a
          >. You can configure the color, icon, position, and text of the
          floating action button so that it matches the brand colors on your
          website.
        </p>
        <vue-glow
          color="#8338ec"
          mode="hex"
          elevation="14"
          :intensity="1.07"
          intense
        >
          <v-img
            :src="require('@/static/blog/dashboard-project-page.png')"
          ></v-img>
        </vue-glow>
        <h3 class="text-h4 mt-12">Add integrations.</h3>
        <p>
          From the integrations menu, you can add multiple integrations so your
          website visitors can call your phone, send you a whatsapp message,
          send you an email or even visit an external link. You can also add
          content to the widget for example your top FAQ items can be embedded
          directly in the widget.
        </p>
        <vue-glow
          color="#8338ec"
          mode="hex"
          elevation="14"
          :intensity="1.07"
          intense
        >
          <v-img
            :src="require('@/static/blog/dashboard-integrations.png')"
          ></v-img>
        </vue-glow>
        <h3 class="text-h4 mt-12">Re-order integrations.</h3>
        <p>
          Integrations buttons in the widget are sorted based on when they are
          created by default. You can alter the sort order in which integrations
          appear in the widget according to your preference. Some people will
          like the option to contact them via email first on the widget while
          you may prefer to add a link to your documentation page as the first
          item. You can change the order of integrations by clicking
          <b>Settings</b> → <b>Integrations</b>
        </p>
        <vue-glow
          color="#8338ec"
          mode="hex"
          elevation="14"
          :intensity="1.07"
          intense
        >
          <v-img
            :src="require('@/static/blog/dashboard-integrations-re-order.png')"
          ></v-img>
        </vue-glow>
        <h3 class="text-h4 mt-12">Future Plans</h3>
        <p>
          <a class="text-decoration-none" :href="$config.dashboardURL"
            >SuperButton</a
          >
          is an open source software and I intend to keep it that way for ever.
          Right now, I'm focussed on adding more and more features and
          integrations to the project. I'm also building superbutton in public.
          If you'll love to share any feedback with me, don't hesitate to
          contact the official SuperButton account on twitter
          <a
            href="https://twitter.com/superbuttonHQ"
            class="text-decoration-none"
            >@SuperButtonHQ.</a
          >
        </p>
        <p>Until the next time✌️</p>
        <div class="d-flex mb-6 mt-8">
          <v-avatar class="mb-n2">
            <v-img :src="authorImage"></v-img>
          </v-avatar>
          <div class="ml-2">
            <p class="subtitle-1 mb-n1">{{ authorName }}</p>
            <a
              class="mb-n4 text-decoration-none text--primary"
              href="https://twitter.com/acho_arnold"
            >
              {{ authorTwitter }}
              <v-icon color="#1DA1F2" small>{{ mdiTwitter }}</v-icon>
            </a>
          </div>
        </div>
        <v-divider class="mx-16"></v-divider>
        <div class="text-center mt-8 mb-4">
          <back-button></back-button>
        </div>
      </v-col>
      <v-col v-if="$vuetify.breakpoint.mdAndUp" md="3">
        <blog-about></blog-about>
      </v-col>
    </v-row>
  </v-container>
</template>

<script>
import { mdiTwitter } from '@mdi/js'
export default {
  name: 'BlogConfigurableActionButton',
  data() {
    return {
      mdiTwitter,
      authorImage: require('@/static/arnold.png'),
      authorName: 'Acho Arnold',
      postDate: 'January 08, 2023',
      readTime: '3 min read',
      authorTwitter: 'acho_arnold',
    }
  },
  head() {
    return {
      title:
        'Add a configurable floating action button to your website - SuperButton',
      meta: [
        {
          hid: 'og-title',
          property: 'og:title',
          content: 'Add a configurable floating action button to your website',
        },
        {
          hid: 'og-desc',
          property: 'og:description',
          content:
            'Add a SuperButton so you can communicate with your website visitors on whatsApp, email, contact forms, or direct them to external links like your documentation.',
        },
        {
          hid: 'og-image',
          property: 'og:image',
          content: 'https://superbutton.app/header.png',
        },
        {
          hid: 'twitter-card',
          name: 'twitter:card',
          content: 'summary_large_image',
        },
        {
          hid: 'og-url',
          property: 'og:url',
          content: 'https://superbutton.app/blog/configurable-action-button',
        },
      ],
    }
  },
}
</script>
