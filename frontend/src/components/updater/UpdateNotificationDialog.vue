<script lang="ts" setup>
import { computed } from 'vue';
import { useSettingsStore } from '@/store/settings';
import { Button } from '@/components/ui/button';
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from '@/components/ui/dialog';
import { BrowserOpenURL } from '../../../wailsjs/runtime/runtime';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();
const settingsStore = useSettingsStore();

const isOpen = computed({
    get: () => settingsStore.showUpdateDialog,
    set: (value) => {
        settingsStore.showUpdateDialog = value;
    },
});

function openReleasePage() {
    if (settingsStore.updateReleaseURL) {
        BrowserOpenURL(settingsStore.updateReleaseURL);
    }
    settingsStore.showUpdateDialog = false;
}

function dismissDialog() {
    settingsStore.showUpdateDialog = false;
}
</script>

<template>
    <Dialog v-model:open="isOpen">
        <DialogContent class="sm:max-w-[425px]">
            <DialogHeader>
                <DialogTitle>{{ t('UpdateDialog.title') }}</DialogTitle>
                <DialogDescription>
                    {{ t('UpdateDialog.description', { version: settingsStore.latestVersionTag }) }}
                </DialogDescription>
            </DialogHeader>
            <DialogFooter class="pt-4">
                <Button variant="outline" @click="dismissDialog">
                    {{ t('UpdateDialog.laterButton') }}
                </Button>
                <Button @click="openReleasePage">
                    {{ t('UpdateDialog.viewButton') }}
                </Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template> 