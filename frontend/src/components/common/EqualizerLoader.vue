<template>
    <div class="flex items-end gap-1 h-8">
        <span
            class="w-1.5 rounded-sm animate-bar1 bg-gradient-to-t from-red-500 via-yellow-500 to-green-500"
        ></span>
        <span
            class="w-1.5 rounded-sm animate-bar2 bg-gradient-to-t from-purple-500 via-pink-500 to-red-500"
        ></span>
        <span
            class="w-1.5 rounded-sm animate-bar3 bg-gradient-to-t from-green-500 via-blue-500 to-purple-500"
        ></span>
        <span
            class="w-1.5 rounded-sm animate-bar4 bg-gradient-to-t from-yellow-500 via-orange-500 to-red-500"
        ></span>
        <span
            class="w-1.5 rounded-sm animate-bar5 bg-gradient-to-t from-pink-500 via-purple-500 to-blue-500"
        ></span>
    </div>
</template>

<style scoped>
/* bounce bars up and down */
@keyframes bounce {
    0%,
    100% {
        height: 20%;
    }
    50% {
        height: 100%;
    }
}

/* shift gradient colors */
@keyframes rainbowShift {
    0% {
        background-position: 0% 50%;
    }
    50% {
        background-position: 100% 50%;
    }
    100% {
        background-position: 0% 50%;
    }
}

/* apply animations with different durations for natural effect */
.animate-bar1 {
    animation:
        bounce 1s ease-in-out infinite,
        rainbowShift 3s linear infinite;
    background-size: 200% 200%;
}
.animate-bar2 {
    animation:
        bounce 1.2s ease-in-out infinite,
        rainbowShift 3.5s linear infinite;
    background-size: 200% 200%;
}
.animate-bar3 {
    animation:
        bounce 0.9s ease-in-out infinite,
        rainbowShift 4s linear infinite;
    background-size: 200% 200%;
}
.animate-bar4 {
    animation:
        bounce 1.3s ease-in-out infinite,
        rainbowShift 3.8s linear infinite;
    background-size: 200% 200%;
}
.animate-bar5 {
    animation:
        bounce 1.1s ease-in-out infinite,
        rainbowShift 4.2s linear infinite;
    background-size: 200% 200%;
}
</style>
