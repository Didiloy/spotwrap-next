<template>
    <div class="p-6 w-full h-full">
        <div
            v-if="loading"
            class="flex flex-col items-center justify-center py-12"
        >
            <div
                class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-500 mb-4"
            ></div>
            <p class="text-gray-400">{{ $t("Subscriptions.loading") }}</p>
        </div>

        <div
            v-else-if="artists.length === 0"
            class="flex flex-col items-center justify-center py-12"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-16 w-16 text-gray-400 mb-4"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
            >
                <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="1.5"
                    d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
            </svg>
            <h3 class="text-xl font-medium text-gray-300 mb-2">
                {{ $t("Subscriptions.no_subscriptions") }}
            </h3>
            <p class="text-gray-400 text-center max-w-md">
                {{ $t("Subscriptions.no_subscriptions_description") }}
            </p>
        </div>

        <div
            v-else
            class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6 overflow-y-auto"
            style="max-height: calc(100vh - 50px)"
        >
            <div
                v-for="artist in artists"
                :key="artist.id"
                class="bg-zinc-300 rounded-xl overflow-hidden shadow-lg hover:shadow-xl transition-shadow"
            >
                <div class="relative aspect-square">
                    <img
                        :src="artist.images[0]?.url"
                        :alt="artist.name"
                        class="w-full h-full object-cover"
                    />
                    <div
                        class="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent"
                    ></div>
                    <div class="absolute bottom-0 left-0 p-4">
                        <h3 class="text-xl font-bold text-white">
                            {{ artist.name }}
                        </h3>
                        <p class="text-gray-300 text-sm">
                            {{
                                artist.followers.total.toLocaleString() +
                                " " +
                                $t("Subscriptions.followers")
                            }}
                        </p>
                    </div>
                </div>
                <div class="p-4">
                    <Button
                        @click="unsubscribe(artist.id)"
                        variant="destructive"
                        class="w-full"
                    >
                        {{ $t("Subscriptions.unsubscribe") }}
                    </Button>
                </div>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { Button } from "@/components/ui/button";
import { GetArtist } from "../../wailsjs/go/main/App";
import {
    GetArtistsFromDB,
    RemoveArtist,
} from "../../wailsjs/go/database/Database";
import { useI18n } from "vue-i18n";
import { useToast } from "@/components/ui/toast/use-toast";

const { toast } = useToast();
const i18n = useI18n();

interface Artist {
    id: string;
    name: string;
    followers: {
        total: number;
    };
    images: Array<{
        url: string;
    }>;
}

const loading = ref(true);
const artists = ref<Artist[]>([]);

const loadArtists = async () => {
    try {
        loading.value = true;
        const dbArtists = await GetArtistsFromDB();
        const artistPromises = dbArtists.map((artist) =>
            GetArtist(artist.SpotifyID),
        );
        const artistData = await Promise.all(artistPromises);

        artists.value = artistData.map((data) => ({
            id: data.artist.id,
            name: data.artist.name,
            followers: data.artist.followers,
            images: data.artist.images,
        }));
    } catch (error) {
        console.error("Error loading artists:", error);
    } finally {
        loading.value = false;
    }
};

const unsubscribe = async (artistId: string) => {
    try {
        await RemoveArtist(artistId);
        artists.value = artists.value.filter((a) => a.id !== artistId);
    } catch (error) {
        console.error("Error unsubscribing:", error);
        toast({
            title: i18n.t("Subscriptions.error"),
            variant: "destructive",
        });
    }
};

onMounted(() => {
    loadArtists();
});
</script>
