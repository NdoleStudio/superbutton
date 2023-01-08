<template>
  <v-container>
    <v-row class="mt-8">
      <v-col cols="12" md="9">
        <v-hover v-slot="{ hover }">
          <v-card
            v-for="post in blogPosts"
            :key="post.route"
            :to="post.route"
            :elevation="hover ? 4 : 1"
          >
            <v-card-title class="text-h4 text-break">{{
              post.title
            }}</v-card-title>
            <v-card-subtitle>
              <span class="text-uppercase primary--text-color">{{
                post.date
              }}</span>
              • <span class="text-uppercase">{{ post.readTime }}</span>
            </v-card-subtitle>
            <v-card-text class="mt-n2">
              <p class="text--primary subtitle-1">{{ post.description }}</p>
              <div class="d-flex mt-n2">
                <v-avatar class="mb-n2">
                  <v-img :src="post.authorImage"></v-img>
                </v-avatar>
                <div class="ml-2">
                  <p class="subtitle-1">{{ post.authorName }}</p>
                  <p class="mt-n5 mb-n4">
                    {{ post.authorTwitter }}
                    <v-icon color="#1DA1F2" small>{{ mdiTwitter }}</v-icon>
                  </p>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-hover>
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
  name: 'BlogIndex',
  data() {
    return {
      mdiTwitter,
      blogPosts: [
        {
          route: '/blog/configurable-action-button',
          title: 'Configurable floating action button on your website',
          date: 'Jan 08, 2023',
          readTime: '2 min read',
          authorImage: require('@/static/arnold.png'),
          description:
            'Learn how to add a super button so you can communicate with your website visitors on whatsApp, email, contact forms, and also direct them to external links like your documentation pages.',
          authorName: 'Acho Arnold',
          authorTwitter: 'acho_arnold',
        },
      ],
    }
  },
  head() {
    return {
      title: 'Blog Posts - SuperButton',
    }
  },
}
</script>
