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
} from "@/components/ui/sidebar";
import logo from "../../assets/images/appicon.png";
import infos from "../../../package.json";
import { useI18n } from "vue-i18n";
import { ref, watch } from "vue";
const i18n = useI18n();

watch(i18n.locale, () => {
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
});

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
</script>

<template>
    <Sidebar>
        <SidebarHeader>
            <div class="flex flex-row items-center justify-center w-full">
                <img :src="logo" alt="Logo" class="w-8 h-8 mr-2" />
                <span class="text-xl font-bold">spotwrap-next</span>
            </div>
        </SidebarHeader>
        <SidebarContent>
            <SidebarGroup>
                <SidebarGroupContent>
                    <SidebarMenu>
                        <SidebarMenuItem
                            v-for="item in items"
                            :key="item.title"
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
        <SidebarFooter class="flex items-center justify-center">
            <ProgressCard :progress="75" message="Loading..." />
            <span>{{ $t("AppSidebar.version") }} {{ infos.version }}</span>
        </SidebarFooter>
    </Sidebar>
</template>
