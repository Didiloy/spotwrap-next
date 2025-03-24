<template>
    <div class="w-full">
        <h2 class="text-lg font-semibold mb-3">
            {{ i18n.t("AlbumsRow.title") }}
        </h2>
        <div class="flex w-full justify-center items-center">
            <Carousel
                class="w-5/6"
                :opts="{
                    align: 'start',
                }"
            >
                <CarouselContent>
                    <CarouselItem
                        v-for="album in props.albums"
                        :key="album.id"
                        class="max-w-[200px] flex flex-col items-center justify-center hover:cursor-pointer"
                    >
                        <img
                            :src="album.images[0]?.url"
                            alt="Album Cover"
                            class="w-40 min-w-40 h-40 object-cover rounded-md shadow-2xl"
                        />
                        <p
                            class="w-full text-[var(--accent-color)] mt-2 text-center text-bold truncate"
                            :title="album.name"
                        >
                            {{ album.name }}
                        </p>
                        <p class="text-zinc-500 mt-2 text-center text-sm">
                            {{ formatReleaseDate(album.release_date) }}
                        </p>
                    </CarouselItem>
                </CarouselContent>
                <CarouselPrevious />
                <CarouselNext />
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
const props = defineProps<{ albums: any[] }>();

const formatReleaseDate = (releaseDate: string): string => {
    const date = new Date(releaseDate);
    const day = String(date.getDate()).padStart(2, "0");
    const month = String(date.getMonth() + 1).padStart(2, "0"); // Months are 0-indexed
    const year = date.getFullYear();
    return `${day}/${month}/${year}`;
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
