import { ref } from "vue";
import { defineStore } from "pinia";
import {
    GetArtistsFromDB,
    IsANewRelease,
    AddArtist,
    GetSeveralArtists,
    GetArtistAlbums,
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

            // Process in chunks of 50 artists
            const chunkSize = 50;
            for (let start = 0; start < artists.length; start += chunkSize) {
                const end = Math.min(start + chunkSize, artists.length);
                const batch = artists.slice(start, end);

                const ids = batch.map((a) => a.SpotifyID);
                // Fetch basic artist info for names/images (batched up to 50)
                const artistsResp: any = await GetSeveralArtists(ids);
                const artistObjs: any[] = (artistsResp?.artists as any[]) || [];

                for (let bi = 0; bi < artistObjs.length; bi++) {
                    const artistObj = artistObjs[bi];
                    if (!artistObj || !artistObj.id) continue;

                    // Update progress immediately per artist
                    currentArtistProgressIndex.value = start + bi + 1;
                    currentCheckingArtist.value = artistObj.name || "";

                    // Fetch recent albums/singles for this artist
                    const albumsResp: any = await GetArtistAlbums(artistObj.id);
                    const albums: any[] = (albumsResp?.items as any[]) || [];

                    for (const album of albums) {
                        const isNewRelease = await IsANewRelease(
                            artistObj.id,
                            album as any,
                        );

                        allAlbums.push({
                            artist: {
                                id: artistObj.id,
                                name: artistObj.name,
                                images: artistObj.images,
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