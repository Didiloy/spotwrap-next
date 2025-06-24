import { ref } from "vue";
import { defineStore } from "pinia";
import {
    GetArtist,
    GetArtistsFromDB,
    IsANewRelease,
    AddArtist,
} from "../../wailsjs/go/main/App";
import { GetDominantColor } from "../../wailsjs/go/utils/Utils";

export interface TimelineItem {
    artist: {
        id: string;
        name: string;
        images?: Array<{ url: string }>;
    };
    album: {
        id: string;
        name: string;
        album_type: string;
        release_date: string;
        total_tracks: number;
        images: Array<{ url: string }>;
    };
    type: "album" | "single";
    date: Date;
    dominantColors?: string[];
    isNewRelease?: boolean;
}

export const useTimelineStore = defineStore("timeline", () => {
    const timelineItems = ref<(TimelineItem & { dominantColors?: string[] })[]>([]);
    const loading = ref(false);
    const initialLoaded = ref(false);
    const currentCheckingArtist = ref<string | null>("");
    const totalArtistsToCheck = ref(0);
    const currentArtistProgressIndex = ref(0);
    const markingAsSeenIndex = ref<number | null>(null);

    async function fetchTimelineItems(force = false) {
        if (loading.value) {
            return;
        }

        if (timelineItems.value.length > 0 && !force) {
            return;
        }

        loading.value = true;
        currentCheckingArtist.value = "";
        totalArtistsToCheck.value = 0;
        currentArtistProgressIndex.value = 0;

        try {
            const artists = await GetArtistsFromDB();

            if (artists.length === 0) {
                loading.value = false;
                initialLoaded.value = true;
                timelineItems.value = [];
                return;
            }

            totalArtistsToCheck.value = artists.length;
            const allAlbums: TimelineItem[] = [];

            for (let i = 0; i < artists.length; i++) {
                const artist = artists[i];
                currentArtistProgressIndex.value = i + 1;

                const artistData = await GetArtist(artist.SpotifyID);
                currentCheckingArtist.value = artistData.artist.name;

                if (artistData.albums) {
                    for (const album of artistData.albums) {
                        const isNewRelease = await IsANewRelease(
                            artist.SpotifyID,
                            album,
                        );

                        allAlbums.push({
                            artist: {
                                id: artistData.artist.id,
                                name: artistData.artist.name,
                                images: artistData.artist.images,
                            },
                            album: {
                                id: album.id,
                                name: album.name,
                                album_type: album.album_type,
                                release_date: album.release_date,
                                total_tracks: album.total_tracks,
                                images: album.images,
                            },
                            type:
                                album.album_type === "album"
                                    ? "album"
                                    : "single",
                            date: new Date(album.release_date),
                            isNewRelease: isNewRelease,
                        });
                    }
                }
            }

            timelineItems.value = allAlbums
                .sort((a, b) => {
                    if (a.isNewRelease && !b.isNewRelease) return -1;
                    if (!a.isNewRelease && b.isNewRelease) return 1;
                    return b.date.getTime() - a.date.getTime();
                })
                .slice(0, 20);

            currentCheckingArtist.value = "";
            loading.value = false;
            initialLoaded.value = true;

            timelineItems.value.forEach((item, index) => {
                loadDominantColorsForItem(item, index);
            });
        } catch (error) {
            console.error("Error fetching artist data:", error);
            currentCheckingArtist.value = "";
            totalArtistsToCheck.value = 0;
            currentArtistProgressIndex.value = 0;
            loading.value = false;
        }
    }

    async function loadDominantColorsForItem(item: TimelineItem, index: number) {
        try {
            if (item.album.images?.[0]?.url) {
                const colors = await GetDominantColor(item.album.images[0].url);
                if (timelineItems.value[index]) {
                    timelineItems.value[index].dominantColors = colors;
                }
            }
        } catch (error) {
            console.error("Error loading dominant colors:", error);
        }
    }

    async function markAsSeen(artistId: string, index: number) {
        markingAsSeenIndex.value = index;
        try {
            await AddArtist(artistId);
            if (timelineItems.value[index]) {
                timelineItems.value[index].isNewRelease = false;
            }
        } catch (error) {
            console.error("Error marking release as seen:", error);
        } finally {
            markingAsSeenIndex.value = null;
        }
    }

    return {
        timelineItems,
        loading,
        initialLoaded,
        currentCheckingArtist,
        totalArtistsToCheck,
        currentArtistProgressIndex,
        markingAsSeenIndex,
        fetchTimelineItems,
        loadDominantColorsForItem,
        markAsSeen,
    };
}); 