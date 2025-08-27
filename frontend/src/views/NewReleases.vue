<template>
    <div class="h-full overflow-y-auto p-6">
        <div class="flex items-center justify-between mb-8">
            <div>
                <h1 class="text-3xl font-bold text-zinc-900 dark:text-white">
                    {{ $t("NewReleases.title") }}
                </h1>
                <p class="text-zinc-500 dark:text-zinc-400 mt-1">
                    {{ $t("NewReleases.subtitle") }}
                </p>
            </div>
            <div class="flex items-center gap-3">
                <Button
                    variant="outline"
                    size="sm"
                    class="rounded-md"
                    @click="reload"
                    :disabled="loading"
                >
                    <RefreshCw
                        class="w-4 h-4"
                        :class="{ 'animate-spin': loading }"
                    />
                    {{ $t("NewReleases.reload") }}
                </Button>
            </div>
        </div>

        <div v-if="loading" class="grid place-items-center py-16">
            <Loader />
            <p class="text-zinc-500 dark:text-zinc-400">
                {{ $t("NewReleases.loading") }}
            </p>
        </div>

        <div v-else>
            <div
                v-if="albums.length === 0"
                class="grid place-items-center py-16 text-center"
            >
                <Sparkles class="w-12 h-12 text-zinc-400 mb-4" />
                <h3
                    class="text-xl font-semibold text-zinc-800 dark:text-zinc-100 mb-2"
                >
                    {{ $t("NewReleases.emptyTitle") }}
                </h3>
                <p class="text-zinc-500 dark:text-zinc-400 max-w-md">
                    {{ $t("NewReleases.emptySubtitle") }}
                </p>
            </div>

            <div
                v-else
                class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6"
            >
                <Card
                    v-for="album in albums"
                    :key="album.id"
                    :album="album"
                    @view="viewAlbum"
                />
            </div>

            <div
                v-if="total > limit"
                class="flex items-center justify-center gap-3 mt-8"
            >
                <Button
                    variant="outline"
                    size="sm"
                    class="rounded-md"
                    :disabled="offset === 0 || loading"
                    @click="prevPage"
                    >{{ $t("NewReleases.prev") }}</Button
                >
                <span class="text-sm text-zinc-600 dark:text-zinc-400">{{
                    pageLabel
                }}</span>
                <Button
                    variant="outline"
                    size="sm"
                    class="rounded-md"
                    :disabled="offset + limit >= total || loading"
                    @click="nextPage"
                    >{{ $t("NewReleases.next") }}</Button
                >
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { Button } from "@/components/ui/button";
import Card from "@/components/new_releases/Card.vue";
import { RefreshCw, Sparkles } from "lucide-vue-next";
import { GetNewReleases } from "../../wailsjs/go/main/App";
import Loader from "@/components/common/Loader.vue";

const router = useRouter();
const loading = ref(false);
const limit = ref(24);
const offset = ref(0);
const total = ref(0);
const albums = ref<any[]>([]);

const pageLabel = computed(() => {
    const start = offset.value + 1;
    const end = Math.min(offset.value + limit.value, total.value);
    return `${start}-${end} / ${total.value}`;
});

async function fetchReleases() {
    loading.value = true;
    try {
        const res: any = await GetNewReleases(limit.value, offset.value);
        const payload = res?.albums || {};
        total.value = payload.total || 0;
        albums.value = payload.items || [];
    } finally {
        loading.value = false;
    }
}

function nextPage() {
    if (offset.value + limit.value < total.value) {
        offset.value += limit.value;
        fetchReleases();
    }
}

function prevPage() {
    if (offset.value >= limit.value) {
        offset.value -= limit.value;
    } else {
        offset.value = 0;
    }
    fetchReleases();
}

function reload() {
    fetchReleases();
}

onMounted(fetchReleases);

function viewAlbum(id: string) {
    router.push(`/album/${id}`);
}
</script>

<style scoped>
.line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
    line-clamp: 1;
}
.line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    line-clamp: 2;
}
</style>
