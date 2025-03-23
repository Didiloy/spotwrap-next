<template>
    <div class="flex flex-col items-center justify-center w-full bg-zinc-100">
        <img
            src="../assets/images/logo-universal.png"
            class="w-1/2 h-1/2 object-contain"
        />
        <Button @click="handleClick" class="mt-4">Click me!</Button>
    </div>
</template>

<script lang="ts" setup>
import { Button } from "@/components/ui/button";
function handleClick() {
    console.log("Button clicked!");
}
</script>

<style scoped></style>
