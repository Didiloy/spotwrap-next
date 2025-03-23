<script setup lang="ts">
import { Home, Settings } from "lucide-vue-next";
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
} from "@/components/ui/sidebar";
import logo from "../assets/images/appicon.png";
import { useI18n } from "vue-i18n";

const i18n = useI18n();

const items = [
    {
        title: i18n.t("AppSidebar.home"),
        url: "#",
        icon: Home,
    },
    {
        title: i18n.t("AppSidebar.settings"),
        url: "#",
        icon: Settings,
    },
];
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
                <!-- <SidebarGroupLabel>Application</SidebarGroupLabel> -->
                <SidebarGroupContent>
                    <SidebarMenu>
                        <SidebarMenuItem
                            v-for="item in items"
                            :key="item.title"
                        >
                            <SidebarMenuButton asChild>
                                <a :href="item.url">
                                    <component :is="item.icon" />
                                    <span>{{ item.title }}</span>
                                </a>
                            </SidebarMenuButton>
                        </SidebarMenuItem>
                    </SidebarMenu>
                </SidebarGroupContent>
            </SidebarGroup>
        </SidebarContent>
    </Sidebar>
</template>
