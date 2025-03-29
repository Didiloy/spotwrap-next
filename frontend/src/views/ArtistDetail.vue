<template>
    <div class="artist-detail h-full overflow-y-auto">
        <!-- Hero Section -->
        <div
            class="relative w-full aspect-square max-h-[500px] bg-gray-800 overflow-hidden"
        >
            <img
                :src="artist?.images[0]?.url"
                alt="Artist image"
                class="w-full h-full object-cover opacity-80"
            />
            <div
                class="absolute inset-0 bg-gradient-to-t from-black/80 to-transparent flex items-end p-6"
            >
                <div class="w-full">
                    <h1 class="text-4xl md:text-6xl font-bold text-white mb-2">
                        {{ artist?.name }}
                    </h1>
                    <p class="text-gray-300 mb-6">
                        {{ formatNumber(artist?.followers?.total) }} followers •
                        {{ artist?.genres?.slice(0, 3).join(", ") }}
                    </p>
                    <Button
                        @click="toggleFollow"
                        class="px-8 py-3 text-lg font-bold rounded-full"
                        :variant="isFollowing ? 'default' : 'outline'"
                    >
                        {{
                            isFollowing
                                ? i18n.t("ArtistDetails.following")
                                : i18n.t("ArtistDetails.follow")
                        }}
                    </Button>
                </div>
            </div>
        </div>

        <!-- Content Section -->
        <div class="p-6 space-y-8">
            <!-- Albums Section -->
            <div v-if="albums.length">
                <AlbumsRow :albums="albums" />
            </div>

            <!-- Popular Tracks -->
            <TracksRow
                v-if="topTracks.length"
                :tracks="topTracks"
                :title="'Popular Tracks'"
            />

            <!-- Singles Section -->
            <div v-if="singles.length">
                <AlbumsRow
                    :albums="singles"
                    :title="i18n.t('ArtistDetails.singles')"
                />
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { useRoute } from "vue-router";
import { Button } from "@/components/ui/button";
import { GetArtist } from "../../wailsjs/go/main/App";
import AlbumsRow from "@/components/search/AlbumsRow.vue";
import TracksRow from "@/components/search/TracksRow.vue";
import { useI18n } from "vue-i18n";

const route = useRoute();
const artist = ref<any>(null);
const artistData = ref<any>({});
const isFollowing = ref(false);
const i18n = useI18n();

const topTracks = computed(() => {
    return (
        artistData.value.top_tracks?.map((track: any) => ({
            ...track,
            album: {
                images: track.album?.images || artist.value?.images,
            },
        })) || []
    );
});

const albums = computed(() => {
    return (
        artistData.value.albums?.filter(
            (a: any) => a.album_group === "album",
        ) || []
    );
});

const singles = computed(() => {
    return (
        artistData.value.albums?.filter(
            (a: any) => a.album_group === "single",
        ) || []
    );
});

// Helper functions
const formatNumber = (num: number) => {
    return new Intl.NumberFormat().format(num);
};

const toggleFollow = () => {
    isFollowing.value = !isFollowing.value;
};

onMounted(async () => {
    const artistId = route.params.id as string;
    const data = await getArtistDetails(artistId);
    artistData.value = data;
    artist.value = data.artist;
});

const getArtistDetails = async (id: string) => {
    try {
        const artistData = await GetArtist(id);
        return artistData;
    } catch (error) {
        console.error("Error fetching artist details:", error);
        return {};
    }
};
</script>

<style scoped>
.artist-detail {
    scrollbar-width: none; /* Firefox */
}
.artist-detail::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}
</style>
