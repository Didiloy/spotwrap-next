<template>
    <div
        class="flex flex-col h-full w-full p-2 justify-start items-center overflow-scroll"
    >
        <!-- Search Bar -->
        <div
            class="w-5/6 rounded-lg flex flex-row justify-center items-center gap-4"
        >
            <Input
                type="text"
                v-model="searchQuery"
                :placeholder="i18n.t('Search.placeholder')"
                @keyup.enter="handleSearch"
            />
            <Button @click="handleSearch">
                <SearchIcon class="w-4 h-4" />
                {{ i18n.t("Search.button") }}
            </Button>
        </div>

        <!-- Search Results -->
        <div v-if="searchResults" class="mt-6 w-full h-fit">
            <AlbumsRow
                v-if="searchResults.albums"
                :albums="searchResults.albums.items"
            />
            <TracksRow
                v-if="searchResults.tracks"
                :tracks="searchResults.tracks.items"
            />
            <!-- <ArtistsRow
                v-if="searchResults.artists"
                :artists="searchResults.artists.items"
            /> -->
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import AlbumsRow from "@/components/search/AlbumsRow.vue";
import TracksRow from "@/components/search/TracksRow.vue";
import ArtistsRow from "@/components/search/ArtistsRow.vue";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Search as SearchIcon } from "lucide-vue-next";
import { useI18n } from "vue-i18n";
import { Search } from "../../wailsjs/go/main/App";

const i18n = useI18n();
const searchQuery = ref("");
const searchResults = ref<any>(null);

const handleSearch = async () => {
    if (searchQuery.value.trim() === "") {
        console.log("Search query is empty");
        return;
    }
    console.log("Searching for:", searchQuery.value);
    searchResults.value = await Search(searchQuery.value);
    console.log("Search results:", searchResults.value);
};
</script>

<style scoped></style>
