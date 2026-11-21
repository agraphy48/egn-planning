<template>
  <div class="h-full w-full">
    <!-- Public Layout: Login, Forgot Password, Reset Password -->
    <PublicLayout v-if="isPublicLayout">
      <RouterView />
    </PublicLayout>

    <!-- Protected Layout: Dashboard, Projects, Tasks, Team, Client Portal -->
    <ProtectedLayout v-else>
      <RouterView />
    </ProtectedLayout>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import PublicLayout from './layouts/PublicLayout.vue'
import ProtectedLayout from './layouts/ProtectedLayout.vue'

const route = useRoute()

// Strict layout separation: public routes must NEVER render protected layout or navigation
const isPublicLayout = computed(() => {
  return (
    route.meta.layout === 'public' ||
    route.meta.requiresAuth === false ||
    route.path === '/login'
  )
})
</script>
