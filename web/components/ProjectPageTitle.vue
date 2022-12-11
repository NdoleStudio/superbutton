<template>
  <v-row>
    <v-col class="d-flex" :class="{ 'mt-n4': !$vuetify.breakpoint.lgAndUp }">
      <back-button
        :icon="true"
        :large="true"
        :class="{ 'mt-n1': !$vuetify.breakpoint.lgAndUp }"
      ></back-button>
      <h3
        class="ml-2"
        :class="{
          'text-h4 mb-4': $vuetify.breakpoint.lgAndUp,
          'text-h5': !$vuetify.breakpoint.lgAndUp,
        }"
      >
        <slot></slot>
      </h3>
    </v-col>
  </v-row>
</template>

<script lang="ts">
import { Vue, Component } from 'vue-property-decorator'
@Component
export default class ProjectPageTitle extends Vue {}
</script>
