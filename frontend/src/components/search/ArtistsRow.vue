<template>
    <div>
        <h2
            class="relative text-lg font-semibold uppercase tracking-wider text-gray-700 dark:text-gray-300 mb-4 font-montserrat pb-1 after:content-[''] after:absolute after:bottom-0 after:left-0 after:w-10 after:h-0.5 after:bg-[var(--accent-color)]"
        >
            {{ i18n.t("ArtistsRow.title") }}
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
                        v-for="artist in artists"
                        :key="artist.id"
                        class="max-w-[180px] flex flex-col items-center justify-center hover:cursor-pointer group"
                    >
                        <div class="p-1">
                            <img
                                :src="artist.images[0]?.url"
                                alt="Artist Image"
                                class="w-32 h-32 object-cover rounded-full shadow-2xl transition-all duration-300 group-hover:shadow-[var(--accent-color)] group-hover:scale-105"
                            />
                        </div>
                        <p
                            class="w-full text-[var(--accent-color)] mt-3 text-center font-medium truncate px-2"
                            :title="artist.name"
                        >
                            {{ artist.name }}
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

defineProps<{ artists: any[] }>();
</script>
