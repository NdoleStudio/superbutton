<template>
  <dashboard-loading></dashboard-loading>
</template>
<script>
export default {
  name: 'IndexPage',
  layout: 'default',
  async mounted() {
    if (!this.$store.getters.authUser) {
      return this.$router.push('/login')
    }

    await Promise.all([
      this.$store.dispatch('loadProjects'),
      this.$store.dispatch('loadUser'),
    ])

    if (!this.$store.getters.hasProjects) {
      await this.$router.push('/projects/create')
      return
    }

    await this.$router.push({
      name: 'projects-id-settings',
      params: {
        id: this.$store.getters.activeProjectId,
      },
    })
  },
}
</script>
