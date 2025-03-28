<template>
    <div class="w-full space-y-4">
        <h2
            class="relative text-lg font-semibold uppercase tracking-widest text-gray-800 dark:text-gray-200 mb-6 font-montserrat pb-2 after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-12 after:h-[3px] after:bg-[var(--accent-color)] after:rounded-full"
        >
            {{ i18n.t("AlbumsRow.title") }}
        </h2>

        <div class="flex w-full justify-center">
            <Carousel
                class="w-full max-w-[95%]"
                :opts="{ align: 'start', dragFree: true }"
            >
                <CarouselContent class="-ml-1">
                    <CarouselItem
                        v-for="album in props.albums"
                        :key="album.id"
                        class="pl-1 basis-[180px] max-w-[180px]"
                    >
                        <div
                            class="group relative p-2 transition-all duration-300 hover:scale-[1.03] active:scale-95"
                        >
                            <div
                                class="relative overflow-hidden rounded-lg aspect-square"
                            >
                                <img
                                    :src="album.images[0]?.url"
                                    alt="Album Cover"
                                    class="w-full h-full object-cover shadow-lg transition-all duration-500 group-hover:shadow-[0_8px_24px_rgba(0,0,0,0.15)] group-hover:brightness-110"
                                />
                                <div
                                    class="absolute inset-0 bg-gradient-to-t from-black/40 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"
                                ></div>
                            </div>

                            <div class="mt-3 space-y-1 px-1">
                                <p
                                    class="font-medium text-gray-900 dark:text-gray-100 truncate text-sm"
                                >
                                    {{ album.name }}
                                </p>
                                <p
                                    class="text-zinc-500 dark:text-zinc-400 text-xs"
                                >
                                    {{ formatReleaseDate(album.release_date) }}
                                </p>
                                <p
                                    v-if="album.artists"
                                    class="text-zinc-400 dark:text-zinc-500 text-xs truncate"
                                >
                                    {{
                                        album.artists
                                            .map((a) => a.name)
                                            .join(", ")
                                    }}
                                </p>
                            </div>
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
const i18n = useI18n();
import {
    Carousel,
    CarouselContent,
    CarouselItem,
    CarouselNext,
    CarouselPrevious,
} from "@/components/ui/carousel";

const props = defineProps<{
    albums: Array<{
        id: string;
        name: string;
        release_date: string;
        images: Array<{ url: string }>;
        artists?: Array<{ name: string }>;
    }>;
}>();

const formatReleaseDate = (releaseDate: string): string => {
    try {
        const options: Intl.DateTimeFormatOptions = {
            year: "numeric",
            month: "short",
            day: "numeric",
        };
        return new Date(releaseDate).toLocaleDateString(undefined, options);
    } catch {
        return releaseDate; // Fallback to raw string if date parsing fails
    }
};
</script>

<style scoped>
/* Hide scrollbar but keep scroll functionality */
.scrollbar-hide::-webkit-scrollbar {
    display: none;
}
.scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
}
</style>
