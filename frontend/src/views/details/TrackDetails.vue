<template>
    <div class="track-detail h-full overflow-y-auto">
        <!-- Hero Section -->
        <div class="relative w-full" :style="heroSectionStyle">
            <div
                class="container mx-auto px-6 py-12 flex flex-col md:flex-row items-start gap-8"
            >
                <div
                    class="relative w-48 h-48 md:w-64 md:h-64 lg:w-80 lg:h-80 flex-shrink-0"
                >
                    <img
                        :src="trackDetails?.track?.album?.images[0]?.url || ''"
                        alt="Track cover"
                        class="w-full h-full object-cover rounded-lg shadow-2xl"
                    />
                    <div
                        class="absolute inset-0 rounded-lg ring-1 ring-inset ring-white/10"
                    ></div>
                </div>

                <!-- Track Info -->
                <div class="flex-1 space-y-4 text-white pt-4">
                    <div>
                        <h1 class="text-4xl font-bold mt-1 mb-2">
                            {{ trackDetails?.track?.name }}
                        </h1>
                        <div
                            class="flex flex-wrap items-center gap-x-4 gap-y-2 text-gray-300"
                        >
                            <span class="flex items-center">
                                <span class="text-purple-400 mr-2">{{
                                    trackDetails?.track?.artists[0]?.name
                                }}</span>
                            </span>
                            <span class="text-sm">
                                {{
                                    formatDuration(
                                        trackDetails?.track?.duration_ms,
                                    )
                                }}
                            </span>
                        </div>
                    </div>

                    <!-- Album Info -->
                    <div class="pt-2 border-t border-white/10">
                        <p class="text-gray-300 mb-1">
                            {{ i18n.t("TrackDetails.from_album") }}
                        </p>
                        <div class="flex items-center gap-3">
                            <img
                                :src="
                                    trackDetails?.track?.album?.images[2]?.url
                                "
                                alt="Album cover"
                                class="w-12 h-12 rounded"
                            />
                            <div>
                                <router-link
                                    :to="
                                        '/album/' +
                                        trackDetails?.track?.album?.id
                                    "
                                    class="font-medium underline"
                                >
                                    {{ trackDetails?.track?.album?.name }}
                                </router-link>
                                <p class="text-sm text-gray-300">
                                    {{
                                        new Date(
                                            trackDetails?.track?.album?.release_date,
                                        ).getFullYear()
                                    }}
                                    •
                                    {{
                                        trackDetails?.track?.album?.total_tracks
                                    }}
                                    tracks
                                </p>
                            </div>
                        </div>
                    </div>

                    <!-- Action Buttons -->
                    <div class="flex flex-col items-start gap-4 pt-6">
                        <div class="flex flex-wrap items-center gap-4">
                            <Button
                                @click="downloadTrack"
                                class="px-8 py-3 text-lg font-bold rounded-full bg-purple-600 hover:bg-purple-700 transition-colors"
                            >
                                {{ i18n.t("TrackDetails.download") }}
                            </Button>
                            <Button
                                @click="selectDownloadPath"
                                variant="outline"
                                class="px-8 py-3 text-lg font-bold rounded-full bg-white/10 transition-colors"
                            >
                                {{ i18n.t("TrackDetails.select_path") }}
                            </Button>
                        </div>
                        <div class="w-full">
                            <span
                                v-if="downloadOptions.path"
                                class="inline-block px-4 py-2 text-sm font-medium rounded-full bg-green-600/20 text-green-400 transition-colors"
                            >
                                {{ downloadOptions.path }}
                            </span>
                            <span
                                v-else
                                class="inline-block px-4 py-2 text-sm font-medium rounded-full bg-red-600/20 text-red-400 transition-colors"
                            >
                                {{ i18n.t("TrackDetails.no_path_selected") }}
                            </span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- Download Options -->
        <div class="p-6">
            <div class="flex flex-wrap items-center gap-4">
                <div>
                    <Label for="bitrate">{{
                        i18n.t("TrackDetails.bitrate")
                    }}</Label>
                    <Select id="bitrate" v-model="downloadOptions.bitrate">
                        <SelectTrigger class="w-[180px]">
                            <SelectValue placeholder="Select bitrate" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem
                                v-for="bitrate in BITRATE_OPTIONS"
                                :value="bitrate"
                                >{{ bitrate }} kbps</SelectItem
                            >
                        </SelectContent>
                    </Select>
                </div>

                <div>
                    <Label for="format">{{
                        i18n.t("TrackDetails.format")
                    }}</Label>
                    <Select id="format" v-model="downloadOptions.format">
                        <SelectTrigger class="w-[180px]">
                            <SelectValue placeholder="Select format" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem
                                v-for="format in FORMAT_OPTIONS"
                                :value="format"
                                >{{ format }}</SelectItem
                            >
                        </SelectContent>
                    </Select>
                </div>
            </div>
        </div>

        <!-- Track Info -->
        <div class="p-6 border-t border-gray-800">
            <h2 class="text-xl font-bold mb-4">
                {{ i18n.t("TrackDetails.about") }}
            </h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                    <h3 class="font-semibold mb-2">
                        {{ i18n.t("TrackDetails.artists") }}
                    </h3>
                    <div class="flex flex-wrap gap-2">
                        <span
                            v-for="artist in trackDetails?.track?.artists"
                            :key="artist.id"
                            class="px-3 py-1 rounded-full bg-white/10 text-sm"
                        >
                            {{ artist.name }}
                        </span>
                    </div>
                </div>
                <div>
                    <h3 class="font-semibold mb-2">
                        {{ i18n.t("TrackDetails.release_date") }}
                    </h3>
                    <p>
                        {{
                            formatDate(trackDetails?.track?.album?.release_date)
                        }}
                    </p>
                </div>
                <div>
                    <h3 class="font-semibold mb-2">
                        {{ i18n.t("TrackDetails.album") }}
                    </h3>
                    <p>{{ trackDetails?.track?.album?.name }}</p>
                </div>
                <div>
                    <h3 class="font-semibold mb-2">
                        {{ i18n.t("TrackDetails.track_number") }}
                    </h3>
                    <p>
                        Track {{ trackDetails?.track?.track_number }} of
                        {{ trackDetails?.track?.album?.total_tracks }}
                    </p>
                </div>
                <div>
                    <h3 class="font-semibold mb-2">
                        {{ i18n.t("TrackDetails.duration") }}
                    </h3>
                    <p>
                        {{ formatDuration(trackDetails?.track?.duration_ms) }}
                    </p>
                </div>
                <div>
                    <h3 class="font-semibold mb-2">
                        {{ i18n.t("TrackDetails.popularity") }}
                    </h3>
                    <div class="flex items-center">
                        <span class="mr-2"
                            >{{ trackDetails?.track?.popularity }}%</span
                        >
                        <Progress
                            :modelValue="trackDetails?.track?.popularity"
                            class="w-[60%]"
                        />
                    </div>
                </div>
            </div>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { GetDominantColor } from "../../../wailsjs/go/utils/Utils";
import { GetTrack, ChooseDirectory } from "../../../wailsjs/go/main/App";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";
import { useI18n } from "vue-i18n";
import { useToast } from "@/components/ui/toast/use-toast";
import { useDownloadStore } from "@/store/download";
import { Download } from "../../../wailsjs/go/spotdl/Downloader";
const downloadStore = useDownloadStore();
const i18n = useI18n();
const route = useRoute();
const router = useRouter();
const { toast } = useToast();
const trackDetails = ref<any>({});

const downloadOptions = ref({
    bitrate: "320",
    format: "mp3",
    path: "",
});

const BITRATE_OPTIONS = ["320", "256", "192", "128", "96"];
const FORMAT_OPTIONS = ["mp3", "flac", "m4a", "ogg", "opus"];

onMounted(async () => {
    const trackId = route.params.id as string;
    trackDetails.value = await getTrackDetails(trackId);
    if (trackDetails.value?.track?.album?.images?.[0]?.url) {
        dominantColors.value = await GetDominantColor(
            trackDetails.value.track.album.images[0].url,
        );
    }
});

const formatDuration = (ms?: number) => {
    if (!ms) return "0:00";
    const minutes = Math.floor(ms / 60000);
    const seconds = ((ms % 60000) / 1000).toFixed(0);
    return `${minutes}:${seconds.padStart(2, "0")}`;
};

const formatDate = (dateString?: string) => {
    if (!dateString) return "";
    const options: Intl.DateTimeFormatOptions = {
        year: "numeric",
        month: "long",
        day: "numeric",
    };
    return new Date(dateString).toLocaleDateString(undefined, options);
};

const error = ref<string | null>(null);
const selectDownloadPath = async () => {
    const path = await ChooseDirectory();
    if (path) {
        downloadOptions.value.path = path;
    }
};

const downloadTrack = async () => {
    if (!downloadOptions.value.path) {
        toast({
            title: i18n.t("TrackDetails.error_title"),
            description: i18n.t("TrackDetails.no_path_selected"),
            variant: "destructive",
        });
        return;
    }
    
    try {
        // Setup event listener for download updates
        downloadStore.clearMessages();
        downloadStore.setupEventListener();
        
        // Download track - now returns a boolean
        const success = await Download(
            trackDetails.value.track.external_urls.spotify,
            downloadOptions.value.path,
            downloadOptions.value.format,
            downloadOptions.value.bitrate + "k",
            [],
        );
        
        // Show success toast immediately if Download function returned true
        if (success) {
            toast({
                title: i18n.t("TrackDetails.download_complete"),
                description: i18n.t("TrackDetails.download_complete_message", {
                    name: trackDetails.value.track.name
                }),
                variant: "default",
            });
        } else {
            // Show error toast immediately if Download function returned false
            toast({
                title: i18n.t("TrackDetails.download_error"),
                description: i18n.t("TrackDetails.download_error_message"),
                variant: "destructive",
            });
        }
        
        console.log("Download result:", success);
    } catch (error) {
        console.error("Download error:", error);
        toast({
            title: i18n.t("TrackDetails.download_error"),
            description: i18n.t("TrackDetails.download_error_message"),
            variant: "destructive",
        });
    }
};

const dominantColors = ref<string[]>([]);
const heroSectionStyle = computed(() => {
    if (dominantColors.value.length >= 2) {
        return {
            background: `linear-gradient(135deg, ${dominantColors.value[0]} 0%, ${dominantColors.value[1]} 100%)`,
        };
    }
    return {
        background: "linear-gradient(135deg, #4f46e5 0%, #1e40af 100%)",
    };
});

const getTrackDetails = async (id: string) => {
    try {
        const trackData = await GetTrack(id);
        return trackData;
    } catch (error) {
        console.error("Error fetching track details:", error);
        toast({
            title: i18n.t("TrackDetails.error_title"),
            description: i18n.t("TrackDetails.error_getting"),
            variant: "destructive",
        });
        return {};
    }
};
</script>

<style scoped>
.track-detail {
    scrollbar-width: none; /* Firefox */
}
.track-detail::-webkit-scrollbar {
    display: none; /* Chrome/Safari */
}

.relative::after {
    content: "";
    position: absolute;
    inset: 0;
    border-radius: 0;
    pointer-events: none;
}
</style>
