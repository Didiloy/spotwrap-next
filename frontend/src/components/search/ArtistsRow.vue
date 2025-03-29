<template>
    <div class="mb-8">
        <div class="flex items-center justify-between mb-5">
            <h2
                class="relative text-xl font-bold uppercase tracking-widest text-gray-800 dark:text-gray-200 font-montserrat pb-2 after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-12 after:h-[2px] after:bg-[var(--accent-color)]"
            >
                {{ i18n.t("ArtistsRow.title") }}
            </h2>
        </div>

        <div class="relative w-full select-none">
            <Carousel
                class="w-full"
                :opts="{
                    align: 'start',
                    dragFree: true,
                    slidesToScroll: 'auto',
                }"
            >
                <CarouselContent class="-ml-1">
                    <CarouselItem
                        v-for="artist in artists"
                        :key="artist.id"
                        class="pl-1 basis-[160px] max-w-[160px]"
                        @click="navigateToArtist(artist.id)"
                    >
                        <div
                            class="flex flex-col items-center p-2 transition-all duration-200 group hover:scale-105 active:scale-95"
                        >
                            <div class="relative">
                                <img
                                    :src="
                                        artist.images?.[0]?.url ||
                                        default_artist
                                    "
                                    alt="Artist Image"
                                    class="w-32 h-32 object-cover rounded-full shadow-lg ring-2 ring-transparent group-hover:ring-[var(--accent-color)] transition-all duration-300"
                                />
                                <div
                                    class="absolute inset-0 rounded-full bg-black/10 opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                ></div>
                            </div>
                            <p
                                class="w-full mt-3 text-center font-medium truncate px-1 text-gray-900 dark:text-gray-100 group-hover:text-[var(--accent-color)] transition-colors duration-200"
                                :title="artist.name"
                            >
                                {{ artist.name }}
                            </p>
                            <p
                                v-if="artist.genres"
                                class="text-xs text-gray-500 dark:text-gray-400 truncate max-w-full"
                            >
                                {{ artist.genres[0] }}
                            </p>
                        </div>
                    </CarouselItem>
                </CarouselContent>

                <CarouselPrevious class="-left-4 top-1/3 size-10" />
                <CarouselNext class="-right-4 top-1/3 size-10" />
            </Carousel>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { useI18n } from "vue-i18n";
import { useRouter } from "vue-router";
import {
    Carousel,
    CarouselContent,
    CarouselItem,
    CarouselNext,
    CarouselPrevious,
} from "@/components/ui/carousel";
import default_artist from "@/assets/images/default_artist.png";

const i18n = useI18n();
const router = useRouter();

defineProps<{
    artists: Array<{
        id: string;
        name: string;
        images?: Array<{ url: string }>;
        genres?: string[];
    }>;
}>();

const navigateToArtist = (artistId: string) => {
    router.push(`/artist/${artistId}`);
};
</script>
