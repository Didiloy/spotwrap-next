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
import { Label } from "@/components/ui/label";

const { locale, t } = useI18n();

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
    console.log("Language changed to:", newLang);
});

// Set initial language
onMounted(() => {
    if (currentLanguage.value) {
        locale.value = currentLanguage.value;
    }
});
</script>

<style scoped>
.settings-page {
    min-height: calc(100vh - 128px);
}
</style>
