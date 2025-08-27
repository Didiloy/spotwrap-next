<template>
    <div
        class="settings-page p-6 max-w-3xl mx-auto space-y-8 overflow-y-auto h-full"
    >
        <h1 class="text-2xl font-bold">{{ $t("Settings.title") }}</h1>

        <!-- Language Selector -->
        <div class="space-y-2 flex flex-row items-center justify-between">
            <Label for="language-select">{{ $t("Settings.language") }}</Label>
            <Select v-model="currentLanguage">
                <SelectTrigger class="w-[180px]">
                    <SelectValue />
                </SelectTrigger>
                <SelectContent>
                    <SelectGroup>
                        <SelectItem
                            v-for="lang in availableLanguages"
                            :key="lang.code"
                            :value="lang.code"
                        >
                            {{ lang.name }}
                        </SelectItem>
                    </SelectGroup>
                </SelectContent>
            </Select>
        </div>

        <!-- Background Mode Section -->
        <div class="space-y-2">
            <Label class="text-purple-800">{{
                $t("Settings.background_mode")
            }}</Label>
            <div class="space-y-4 ml-2">
                <AutoStartToggle />
                <AutoDownloadToggle />

                <!-- Preferred Download Path Selector -->
                <div class="flex flex-wrap items-center gap-4">
                    <Button @click="selectPreferredPath" variant="outline">
                        {{ t("AlbumDetails.select_path") }}
                    </Button>
                    <span
                        v-if="settingsStore.newReleasesDownloadPath"
                        class="text-sm text-green-600"
                    >
                        {{ settingsStore.newReleasesDownloadPath }}
                    </span>
                    <span v-else class="text-sm text-muted-foreground">
                        {{ t("AlbumDetails.no_path_selected") }}
                    </span>
                </div>
            </div>
        </div>

        <!-- Append Artist/Album Path Toggle -->
        <div class="space-y-2">
            <Label class="text-purple-800">{{
                $t("Settings.downloadPathOptionsTitle")
            }}</Label>
            <div class="flex items-center justify-between ml-2">
                <div>
                    <h3 class="font-medium">
                        {{ $t("Settings.appendPathTitle") }}
                    </h3>
                    <p class="text-sm text-muted-foreground">
                        {{ $t("Settings.appendPathDescription") }}
                    </p>
                </div>
                <Switch
                    :model-value="settingsStore.appendArtistAlbumToPath"
                    @update:model-value="
                        settingsStore.toggleAppendArtistAlbumToPath
                    "
                />
            </div>
        </div>

        <!-- Spotify API Credentials -->
        <div class="space-y-2">
            <Label class="text-purple-800">{{
                $t("Settings.spotify_api")
            }}</Label>
            <div class="flex flex-col gap-2 ml-2">
                <div class="flex items-center justify-between">
                    <div>
                        <h3 class="font-medium">
                            {{ $t("Settings.spotify_credentials") }}
                        </h3>
                        <p class="text-sm text-muted-foreground">
                            {{
                                settingsStore.hasValidCredentials
                                    ? $t("Settings.spotify_credentials_set")
                                    : $t("Settings.spotify_credentials_not_set")
                            }}
                        </p>
                    </div>
                    <Button @click="settingsStore.showCredentialsModal = true">
                        {{
                            settingsStore.hasValidCredentials
                                ? $t("Settings.update_credentials")
                                : $t("Settings.set_credentials")
                        }}
                    </Button>
                </div>
            </div>
        </div>

        <!-- Logs Section -->
        <div class="space-y-2">
            <Label>{{ $t("Settings.logs") }}</Label>
            <div class="flex gap-2">
                <Button @click="showLogsDialog = true">
                    {{ $t("Settings.view_logs") }}
                </Button>
                <Button variant="destructive" @click="clearLogs">
                    {{ $t("Settings.clear_logs") }}
                </Button>
            </div>
        </div>

        <!-- Logs Dialog -->
        <LogsDialog v-model:open="showLogsDialog" />

        <!-- Spotify Credentials Modal -->
        <SpotifyCredentialsModal
            v-model:open="settingsStore.showCredentialsModal"
            @saved="checkCredentials"
        />
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import { useI18n } from "vue-i18n";
import {
    Select,
    SelectContent,
    SelectGroup,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import LogsDialog from "@/components/settings/LogsDialog.vue";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import AutoStartToggle from "@/components/settings/AutoStartToggle.vue";
import AutoDownloadToggle from "@/components/settings/AutoDownloadToggle.vue";
import SpotifyCredentialsModal from "@/components/settings/SpotifyCredentialsModal.vue";
import { useDownloadStore } from "@/store/download";
import { useSettingsStore } from "@/store/settings";
import { ChooseDirectory, SetSetting } from "../../wailsjs/go/main/App";

const { locale, t } = useI18n();
const downloadStore = useDownloadStore();
const settingsStore = useSettingsStore();
const availableLanguages = [
    { code: "en", name: "English" },
    { code: "fr", name: "Français" },
];

// Current language (default to browser language or English)
const currentLanguage = ref(
    localStorage.getItem("lang") || navigator.language.split("-")[0] || "en",
);

watch(currentLanguage, (newLang) => {
    locale.value = newLang;
    localStorage.setItem("lang", newLang);
});

function checkCredentials() {
    settingsStore.checkCredentialsValidity();
}

onMounted(async () => {
    if (currentLanguage.value) {
        locale.value = currentLanguage.value;
    }
});

const showLogsDialog = ref(false);

function clearLogs() {
    downloadStore.clearMessages();
}

const selectPreferredPath = async () => {
    const selectedPath = await ChooseDirectory();
    if (selectedPath) {
        await SetSetting("newReleasesDownloadPath", selectedPath);
        await settingsStore.fetchNewReleasesDownloadPath();
    }
};
</script>

<style scoped>
.settings-page {
    min-height: calc(100vh - 128px);
}
</style>
