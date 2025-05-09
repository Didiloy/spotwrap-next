<script setup lang="ts">
import { Home, Settings, Search, Bell } from "lucide-vue-next";
import ProgressCard from "@/components/sidebar/ProgressCard.vue";
import {
    Sidebar,
    SidebarHeader,
    SidebarContent,
    SidebarGroup,
    SidebarGroupContent,
    SidebarGroupLabel,
    SidebarMenu,
    SidebarMenuButton,
    SidebarMenuItem,
    SidebarFooter,
    SidebarTrigger,
    useSidebar
} from "@/components/ui/sidebar";
import logo from "../../assets/images/appicon.png";
import { useDownloadStore } from "@/store/download";
import { storeToRefs } from "pinia";
import infos from "../../../package.json";
import { useI18n } from "vue-i18n";
import { ref, watch, onMounted, computed } from "vue";
import { useToast } from "@/components/ui/toast/use-toast";

const { toast } = useToast();

const i18n = useI18n();
const downloadStore = useDownloadStore();
const { downloadMessages, isDownloading } = storeToRefs(downloadStore);

// Get the last message using a computed property
const lastMessage = computed(() => {
    return downloadMessages.value.length > 0
        ? downloadMessages.value[downloadMessages.value.length - 1]
        : "";
});

// Reactive navigation items
const items = ref([
    {
        title: i18n.t("AppSidebar.home"),
        url: "/",
        icon: Home,
    },
    {
        title: i18n.t("AppSidebar.search"),
        url: "/search",
        icon: Search,
    },
    {
        title: i18n.t("AppSidebar.subscriptions"),
        url: "/subscriptions",
        icon: Bell,
    },
    {
        title: i18n.t("AppSidebar.settings"),
        url: "/settings",
        icon: Settings,
    },
]);

// Watch for locale changes
watch(
    () => i18n.locale.value,
    () => {
        items.value = [
            {
                title: i18n.t("AppSidebar.home"),
                url: "/",
                icon: Home,
            },
            {
                title: i18n.t("AppSidebar.search"),
                url: "/search",
                icon: Search,
            },
            {
                title: i18n.t("AppSidebar.subscriptions"),
                url: "/subscriptions",
                icon: Bell,
            },
            {
                title: i18n.t("AppSidebar.settings"),
                url: "/settings",
                icon: Settings,
            },
        ];
    },
);

onMounted(() => {
    downloadStore.setupEventListener();
});

watch(downloadStore.downloadMessages, (messages) => {
    if (messages.includes("fatal_error")) {
        toast({
            description: i18n.t("AppSidebar.download_error"),
            variant: "destructive",
        });
        downloadStore.isDownloading = false;
    }
});
</script>

<template>
    <Sidebar>
        <SidebarHeader class="p-4">
            <div class="flex flex-row items-center justify-between w-full">
                <SidebarTrigger variant="sidebar" class="mr-2" />
                <div class="flex items-center">
                    <img :src="logo" alt="Logo" class="w-8 h-8 mr-2" />
                    <span class="text-xl font-bold">spotwrap-next</span>
                </div>
            </div>
        </SidebarHeader>
        <SidebarContent class="px-2">
            <SidebarGroup>
                <SidebarGroupContent>
                    <SidebarMenu>
                        <SidebarMenuItem
                            v-for="item in items"
                            :key="item.title"
                            class="my-1"
                        >
                            <SidebarMenuButton asChild>
                                <router-link :to="item.url">
                                    <component :is="item.icon" />
                                    <span>{{ item.title }}</span>
                                </router-link>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>
        </SidebarContent>
        <SidebarFooter class="flex flex-col items-center justify-center p-4 space-y-3">
            <ProgressCard
                :progress="50"
                :showProgress="isDownloading"
                :message="lastMessage"
                :maxLength="30"
                class="w-full mb-2"
            />
            <span class="text-xs text-sidebar-foreground/70">{{ $t("AppSidebar.version") }} {{ infos.version }}</span>
        </SidebarFooter>
    </Sidebar>
</template>
