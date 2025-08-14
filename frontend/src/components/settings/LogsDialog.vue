<template>
    <Dialog :open="open" @update:open="onUpdateOpen">
        <DialogContent class="max-h-[80vh] overflow-y-auto">
            <DialogHeader>
                <DialogTitle>{{ t("Settings.logs_title") }}</DialogTitle>
                <DialogDescription>
                    {{ t("Settings.logs_description") }}
                </DialogDescription>
            </DialogHeader>

            <div class="font-mono text-sm space-y-1">
                <div
                    v-for="(log, index) in downloadStore.downloadMessages"
                    :key="index"
                    class="py-1 border-b"
                >
                    {{ log }}
                </div>
                <div v-if="downloadStore.downloadMessages.length === 0" class="text-muted-foreground">
                    {{ t("Settings.no_logs") }}
                </div>
            </div>

            <DialogFooter>
                <Button @click="onUpdateOpen(false)">
                    {{ t("Settings.close") }}
                </Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
</template>

<script setup lang="ts">
import { useI18n } from "vue-i18n";
import { useDownloadStore } from "@/store/download";
import { Button } from "@/components/ui/button";
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog";

const props = defineProps<{ open: boolean }>();
const emit = defineEmits<{ (e: "update:open", value: boolean): void }>();

const { t } = useI18n();
const downloadStore = useDownloadStore();

function onUpdateOpen(value: boolean) {
    emit("update:open", value);
}
</script>


