<template>
    <div class="settings-page p-6 max-w-3xl mx-auto space-y-8">
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

        <!-- Autostart -->
        <AutoStartToggle />

        <!-- Append Artist/Album Path Toggle -->
        <div class="space-y-2">
            <Label class="text-base">{{ $t("Settings.downloadPathOptionsTitle") }}</Label>
            <div class="flex items-center justify-between">
                <div>
                    <p class="text-sm font-medium">{{ $t("Settings.appendPathTitle") }}</p>
                    <p class="text-xs text-muted-foreground">
                        {{ $t("Settings.appendPathDescription") }}
                    </p>
                </div>
                <Switch 
                    :model-value="settingsStore.appendArtistAlbumToPath"
                    @update:model-value="settingsStore.toggleAppendArtistAlbumToPath"
                />
            </div>
        </div>

        <!-- Spotify API Credentials -->
        <div class="space-y-2">
            <Label>{{ $t("Settings.spotify_api") }}</Label>
            <div class="flex flex-col gap-2">
                <div class="flex items-center justify-between">
                    <div>
                        <p class="text-sm font-medium">{{ $t("Settings.spotify_credentials") }}</p>
                        <p class="text-xs text-muted-foreground">
                            {{ settingsStore.hasValidCredentials ? $t("Settings.spotify_credentials_set") : $t("Settings.spotify_credentials_not_set") }}
                        </p>
                    </div>
                    <Button @click="settingsStore.showCredentialsModal = true">
                        {{ settingsStore.hasValidCredentials ? $t("Settings.update_credentials") : $t("Settings.set_credentials") }}
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
        <Dialog v-model:open="showLogsDialog">
            <DialogContent class="max-h-[80vh] overflow-y-auto">
                <DialogHeader>
                    <DialogTitle>{{ $t("Settings.logs_title") }}</DialogTitle>
                    <DialogDescription>
                        {{ $t("Settings.logs_description") }}
                    </DialogDescription>
                </DialogHeader>

                <div class="font-mono text-sm space-y-1">
                    <div
                        v-for="(log, index) in logs"
                        :key="index"
                        class="py-1 border-b"
                    >
                        {{ log }}
                    </div>
                    <div v-if="logs.length === 0" class="text-muted-foreground">
                        {{ $t("Settings.no_logs") }}
                    </div>
                </div>

                <DialogFooter>
                    <Button @click="showLogsDialog = false">
                        {{ $t("Settings.close") }}
                    </Button>
                </DialogFooter>
            </DialogContent>
        </Dialog>

        <!-- Spotify Credentials Modal -->
        <SpotifyCredentialsModal v-model:open="settingsStore.showCredentialsModal" @saved="checkCredentials" />
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
    SelectLabel,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import AutoStartToggle from "@/components/settings/AutoStartToggle.vue";
import SpotifyCredentialsModal from "@/components/settings/SpotifyCredentialsModal.vue";
import { useDownloadStore } from "@/store/download";
import { useSettingsStore } from "@/store/settings";

const { locale, t } = useI18n();
const downloadStore = useDownloadStore();
const settingsStore = useSettingsStore();
// Available languages
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

// Spotify credentials state
function checkCredentials() {
    settingsStore.checkCredentialsValidity();
}

// Set initial language
onMounted(async () => {
    if (currentLanguage.value) {
        locale.value = currentLanguage.value;
    }
    logs.value = [...downloadStore.downloadMessages];
});

// Logs dialog state
const showLogsDialog = ref(false);
const logs = ref<string[]>([]);

watch(
    () => downloadStore.downloadMessages,
    (newLogs) => {
        logs.value = [...newLogs];
    },
    { deep: true },
);

function clearLogs() {
    downloadStore.clearMessages();
    logs.value = [];
}
</script>

<style scoped>
.settings-page {
    min-height: calc(100vh - 128px);
}
</style>
