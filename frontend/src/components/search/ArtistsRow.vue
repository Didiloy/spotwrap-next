<template>
    <div>
        <h2 class="text-lg font-semibold mb-3">Artists</h2>
        <div class="grid grid-cols-4 gap-4">
            <div
                v-for="artist in artists"
                :key="artist.id"
                class="p-2 bg-zinc-800 rounded-lg shadow-md"
            >
                <img
                    :src="artist.images[0]?.url"
                    alt="Artist Image"
                    class="w-full rounded-md"
                />
                <p class="text-white mt-2 text-center">{{ artist.name }}</p>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
defineProps<{ artists: any[] }>();
</script>
