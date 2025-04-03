<template>
    <div class="flex items-center space-x-2 mb-6">
        <Label for="sort-method">{{ $t("Subscriptions.Sort.sort_by") }}</Label>
        <Select v-model="sortMethod" @update:modelValue="emitSortChange">
            <SelectTrigger class="w-[180px]">
                <SelectValue
                    :placeholder="$t('Subscriptions.Sort.select_sort')"
                />
            </SelectTrigger>
            <SelectContent>
                <SelectItem value="name-asc">
                    {{ $t("Subscriptions.Sort.name_asc") }}
                </SelectItem>
                <SelectItem value="name-desc">
                    {{ $t("Subscriptions.Sort.name_desc") }}
                </SelectItem>
                <SelectItem value="date-asc">
                    {{ $t("Subscriptions.Sort.date_asc") }}
                </SelectItem>
                <SelectItem value="date-desc">
                    {{ $t("Subscriptions.Sort.date_desc") }}
                </SelectItem>
            </SelectContent>
        </Select>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from "@/components/ui/select";
import { Label } from "@/components/ui/label";

const emit = defineEmits(["sort-change"]);

const sortMethod = ref("date-desc");

const emitSortChange = () => {
    emit("sort-change", sortMethod.value);
};
</script>
