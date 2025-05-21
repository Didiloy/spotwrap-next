<template>
    <div
        class="flex flex-col h-full w-full p-4 justify-start items-center overflow-auto"
    >
        <!-- Search Bar -->
        <div
            class="w-full max-w-2xl rounded-lg flex flex-row justify-center items-center gap-3 mb-6"
        >
            <Input
                type="text"
                v-model="search_query"
                :placeholder="i18n.t('Search.placeholder')"
                @keyup.enter="handleSearch"
                :disabled="is_loading"
                class="flex-1"
            />
            <Button
                @click="handleSearch"
                :disabled="is_loading || !search_query.trim()"
                class="shrink-0"
            >
                <template v-if="is_loading">
                    <LoaderCircle class="w-4 h-4 animate-spin mr-2" />
                    {{ i18n.t("Search.loading") }}
                </template>
                <template v-else>
                    <SearchIcon class="w-4 h-4 mr-2" />
                    {{ i18n.t("Search.button") }}
                </template>
            </Button>
        </div>

        <!-- Error Message -->
        <div
            v-if="error_message"
            class="w-full max-w-2xl bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-6 text-red-600 dark:text-red-300 flex items-start gap-3"
        >
            <AlertCircle class="w-5 h-5 mt-0.5 flex-shrink-0" />
            <div>
                <h3 class="font-medium">{{ i18n.t("Search.errorTitle") }}</h3>
                <p class="text-sm">{{ error_message }}</p>
                <Button
                    v-if="is_retry_available"
                    @click="handleSearch"
                    variant="ghost"
                    size="sm"
                    class="mt-2 text-red-600 dark:text-red-300 hover:bg-red-100 dark:hover:bg-red-900/30"
                >
                    {{ i18n.t("Search.retry") }}
                </Button>
            </div>
        </div>

        <!-- Empty State -->
        <div
            v-else-if="
                search_attempted &&
                !is_loading &&
                (!search_results || isEmptyResults(search_results))
            "
            class="w-full max-w-2xl text-center py-12"
        >
            <SearchIcon
                class="w-12 h-12 mx-auto text-gray-400 dark:text-gray-500 mb-4"
            />
            <h3
                class="text-lg font-medium text-gray-700 dark:text-gray-300 mb-2"
            >
                {{ i18n.t("Search.noResultsTitle") }}
            </h3>
            <p class="text-gray-500 dark:text-gray-400">
                {{ i18n.t("Search.noResultsMessage") }} "<span
                    class="font-medium"
                    >{{ search_query }}</span
                >"
            </p>
        </div>

        <!-- Search Results -->
        <div v-else-if="search_results" class="w-full space-y-8">
            <AlbumsRow
                v-if="search_results.albums?.items?.length"
                :albums="search_results.albums.items"
            />
            <TracksRow
                v-if="search_results.tracks?.items?.length"
                :tracks="search_results.tracks.items.slice(0, 10)"
            />
            <ArtistsRow
                v-if="search_results.artists?.items?.length"
                :artists="search_results.artists.items"
            />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, shallowRef, watch, onMounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import AlbumsRow from "@/components/search/AlbumsRow.vue";
import TracksRow from "@/components/search/TracksRow.vue";
import ArtistsRow from "@/components/search/ArtistsRow.vue";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import {
    Search as SearchIcon,
    LoaderCircle,
    AlertCircle,
} from "lucide-vue-next";
import { useI18n } from "vue-i18n";
import { Search } from "../../wailsjs/go/main/App";
import { SearchResults } from "@/interfaces/searchResult";

const i18n = useI18n();
const route = useRoute();
const router = useRouter();

const search_query = ref("");
const search_results = shallowRef<SearchResults | null>(null);

const is_loading = ref(false);
const error_message = ref("");
const search_attempted = ref(false);
const is_retry_available = ref(true);

const isEmptyResults = (results: any) => {
    return (
        !results.albums?.items?.length &&
        !results.tracks?.items?.length &&
        !results.artists?.items?.length
    );
};

const performSearch = async (term: string) => {
    if (term.trim() === "") {
        resetSearchState();
        return;
    }

    search_attempted.value = true;
    is_loading.value = true;
    error_message.value = "";
    search_query.value = term;

    try {
        search_results.value = await Search(term);
        if (isEmptyResults(search_results.value)) {
            error_message.value = i18n.t("Search.noResultsDetailed");
        }
    } catch (error) {
        console.error("Search error:", error);
        search_results.value = null;
        error_message.value = i18n.t("Search.errorMessage");
        is_retry_available.value = true;
    } finally {
        is_loading.value = false;
    }
};

const resetSearchState = () => {
    search_results.value = null;
    error_message.value = "";
    search_attempted.value = false;
};

const handleSearch = async () => {
    const term = search_query.value.trim();
    if (term === "") return;

    if (route.params.term !== term) {
        router.push({ name: 'search', params: { term } });
    } else {
        await performSearch(term);
    }
};

onMounted(async () => {
    if (route.params.term && typeof route.params.term === 'string') {
        search_query.value = route.params.term;
        await performSearch(route.params.term);
    }
});

watch(() => route.params.term, (newTerm) => {
    if (newTerm && typeof newTerm === 'string') {
        if (search_query.value !== newTerm) {
            search_query.value = newTerm;
            performSearch(newTerm);
        }
    } else if (!newTerm && search_query.value) {
        search_query.value = "";
        resetSearchState();
    }
});
</script>
