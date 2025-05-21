<script lang="ts" setup>
import { useRouter } from 'vue-router';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { Button } from '@/components/ui/button';

const router = useRouter();

const goBack = () => {
  router.back();
};
</script>

<template>
  <div class="flex items-center justify-between px-2 py-2 bg-zinc-100 rounded-md">
    <Button variant="outline" size="sm" @click="goBack" class="text-sm">
      &leftarrow; Back
    </Button>
    <SidebarTrigger />
  </div>
</template>
