<template>
  <v-simple-table v-if="integrations.length">
    <template #default>
      <thead class="text-uppercase">
        <tr>
          <th
            v-if="$vuetify.breakpoint.mdAndDown"
            class="text-left"
            style="width: 80%"
          >
            Name
          </th>
          <th v-else class="text-left" style="width: 30%">Name</th>
          <th
            v-if="$vuetify.breakpoint.lgAndUp"
            class="text-left"
            style="width: 50%"
          >
            Identifier
          </th>
          <th class="">Action</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in integrations" :key="item.id">
          <td>{{ item.name }}</td>
          <td v-if="$vuetify.breakpoint.lgAndUp">{{ item.integration_id }}</td>
          <td>
            <v-btn
              small
              class="secondary"
              :to="{
                name: routeName,
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
</template>

<script lang="ts">
import { Vue, Component, Prop } from 'vue-property-decorator'
import { mdiSquareEditOutline } from '@mdi/js'

interface Integration {
  name: string
  integration_id: string
}

@Component
export default class ProjectIntegrationTable extends Vue {
  @Prop({ required: true, type: String }) routeName!: string
  @Prop({ required: true, type: Array }) integrations!: Array<Integration>
  mdiSquareEditOutline: string = mdiSquareEditOutline
}
</script>
