<template>
    <div
        class="group relative rounded-2xl overflow-hidden border border-white/10 dark:border-white/10 bg-gradient-to-br from-white/60 to-white/30 dark:from-white/5 dark:to-white/0 backdrop-blur-md transition-transform cursor-pointer hover:bg-zinc-300 dark:hover:bg-white/5"
    >
        <div class="relative" @click="onViewClick">
            <img
                :src="album.images?.[0]?.url || fallbackImage"
                :alt="album.name"
                class="w-full aspect-square object-cover"
            />
            <div
                class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"
            ></div>
            <button
                class="absolute bottom-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity"
                @click.stop="onViewClick"
            >
                <Button
                    size="sm"
                    class="rounded-full bg-purple-600 hover:bg-purple-700"
                >
                    {{ $t("NewReleases.view") }}
                </Button>
            </button>
        </div>
        <div class="p-4">
            <h3
                class="font-semibold text-zinc-900 dark:text-white leading-tight line-clamp-2"
            >
                {{ album.name }}
            </h3>
            <p
                class="text-sm text-purple-600 dark:text-purple-400 mt-1 line-clamp-1"
            >
                <template v-for="(a, idx) in album.artists" :key="a.id">
                    <router-link
                        :to="`/artist/${a.id}`"
                        @click.stop
                        class="hover:underline"
                    >
                        {{ a.name }}
                    </router-link>
                    <span v-if="idx < album.artists.length - 1">, </span>
                </template>
            </p>
            <div
                class="mt-3 flex items-center justify-between text-xs text-zinc-600 dark:text-zinc-400"
            >
                <span>{{
                    formattedDate(
                        album.release_date,
                        album.release_date_precision,
                    )
                }}</span>
                <span
                    class="px-2 py-0.5 rounded-full bg-white/70 dark:bg-white/10"
                >
                    {{ (album.album_type || "").toUpperCase() }}
                </span>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { Button } from "@/components/ui/button";
import fallbackImage from "@/assets/images/default_artist.png";

interface AlbumArtist {
    id: string;
    name: string;
}

interface Album {
    id: string;
    name: string;
    images?: Array<{ url?: string }>;
    artists: AlbumArtist[];
    album_type?: string;
    release_date: string;
    release_date_precision?: string;
}

const props = defineProps<{ album: Album }>();
const emit = defineEmits<{ (e: "view", id: string): void }>();

function onViewClick() {
    emit("view", props.album.id);
}

function formattedDate(date: string, precision?: string) {
    if (!date) return "";
    const normalized =
        precision === "year"
            ? `${date}-01-01`
            : precision === "month" && date.length === 7
              ? `${date}-01`
              : date;
    try {
        return new Date(normalized).toLocaleDateString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric",
        });
    } catch {
        return date;
    }
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
