<template>
    <div class="artist-detail h-full overflow-y-auto">
        <!-- Hero Section -->
        <div class="w-full flex justify-center">
            <div class="w-11/12 max-w-[1800px] mt-2 rounded-xl">
                <div
                    class="relative w-full aspect-square max-h-[300px] overflow-hidden rounded-xl"
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
                            <h1
                                class="text-4xl md:text-6xl font-bold text-white mb-2"
                            >
                                {{ artist?.name }}
                            </h1>
                            <p class="text-gray-300 mb-6">
                                {{ formatNumber(artist?.followers?.total) }}
                                followers •
                                {{ artist?.genres?.slice(0, 5).join(", ") }}
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
import { useToast } from "@/components/ui/toast/use-toast";
import { GetArtist } from "../../../wailsjs/go/main/App";
import {
    AddArtist,
    RemoveArtist,
    GetArtistsFromDB,
} from "../../../wailsjs/go/database/Database";
import AlbumsRow from "@/components/search/AlbumsRow.vue";
import TracksRow from "@/components/search/TracksRow.vue";
import { useI18n } from "vue-i18n";

const route = useRoute();
const artist = ref<any>(null);
const artistData = ref<any>({});
const isFollowing = ref(false);
const i18n = useI18n();
const { toast } = useToast();

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

const toggleFollow = async () => {
    if (isFollowing.value) {
        const success = await RemoveArtist(artist.value.id);
        if (!success) {
            console.error("Failed to remove artist");
            toast({
                title: i18n.t("ArtistDetails.error_title"),
                description: i18n.t("ArtistDetails.error_unfollowing"),
                variant: "destructive",
            });
            return;
        }
    } else {
        const success = await AddArtist(artist.value.id);
        if (!success) {
            console.error("Failed to add artist");
            toast({
                title: i18n.t("ArtistDetails.error_title"),
                description: i18n.t("ArtistDetails.error_following"),
                variant: "destructive",
            });
            return;
        }
    }
    isFollowing.value = !isFollowing.value;
};

onMounted(async () => {
    const artistId = route.params.id as string;
    const data = await getArtistDetails(artistId);
    const subbed_artists = await GetArtistsFromDB();
    let isSubbed = false;
    isSubbed = subbed_artists.some(
        (artist: { SpotifyID: string; LastChecked: Date }) =>
            artist.SpotifyID === artistId,
    );
    artistData.value = data;
    artist.value = data.artist;
    isFollowing.value = isSubbed;
});

const getArtistDetails = async (id: string) => {
    try {
        const artistData = await GetArtist(id);
        return artistData;
    } catch (error) {
        console.error("Error fetching artist details:", error);
        toast({
            title: i18n.t("ArtistDetails.error_title"),
            description: i18n.t("ArtistDetails.error_getting"),
            variant: "destructive",
        });
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
