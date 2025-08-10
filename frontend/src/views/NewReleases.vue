<template>
  <div class="h-full overflow-y-auto p-6">
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-3xl font-bold text-zinc-900 dark:text-white">{{ $t('NewReleases.title') }}</h1>
        <p class="text-zinc-500 dark:text-zinc-400 mt-1">{{ $t('NewReleases.subtitle') }}</p>
      </div>
      <div class="flex items-center gap-3">
        <Button variant="outline" size="sm" class="rounded-md" @click="reload" :disabled="loading">
          <RefreshCw class="w-4 h-4" :class="{ 'animate-spin': loading }" />
          {{ $t('NewReleases.reload') }}
        </Button>
      </div>
    </div>

    <div v-if="loading" class="grid place-items-center py-16">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-purple-500 mb-4"></div>
      <p class="text-zinc-500 dark:text-zinc-400">{{ $t('NewReleases.loading') }}</p>
    </div>

    <div v-else>
      <div v-if="albums.length === 0" class="grid place-items-center py-16 text-center">
        <Sparkles class="w-12 h-12 text-zinc-400 mb-4" />
        <h3 class="text-xl font-semibold text-zinc-800 dark:text-zinc-100 mb-2">{{ $t('NewReleases.emptyTitle') }}</h3>
        <p class="text-zinc-500 dark:text-zinc-400 max-w-md">{{ $t('NewReleases.emptySubtitle') }}</p>
      </div>

      <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        <div v-for="album in albums" :key="album.id" class="group relative rounded-2xl overflow-hidden border border-white/10 dark:border-white/10 bg-gradient-to-br from-white/60 to-white/30 dark:from-white/5 dark:to-white/0 backdrop-blur-md transition-transform cursor-pointer hover:bg-zinc-300 dark:hover:bg-white/5">
          <div class="relative" @click="viewAlbum(album.id)">
            <img :src="album.images?.[0]?.url || fallbackImage" :alt="album.name" class="w-full aspect-square object-cover" />
            <div class="absolute inset-0 bg-gradient-to-t from-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity"></div>
            <button class="absolute bottom-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity" @click.stop="viewAlbum(album.id)">
              <Button size="sm" class="rounded-full bg-purple-600 hover:bg-purple-700">
                {{ $t('NewReleases.view') }}
              </Button>
            </button>
          </div>
          <div class="p-4">
            <h3 class="font-semibold text-zinc-900 dark:text-white leading-tight line-clamp-2">{{ album.name }}</h3>
            <p class="text-sm text-purple-600 dark:text-purple-400 mt-1 line-clamp-1">
              <template v-for="(a, idx) in album.artists" :key="a.id">
                <router-link :to="`/artist/${a.id}`" @click.stop class="hover:underline">
                  {{ a.name }}
                </router-link>
                <span v-if="idx < album.artists.length - 1">, </span>
              </template>
            </p>
            <div class="mt-3 flex items-center justify-between text-xs text-zinc-600 dark:text-zinc-400">
              <span>{{ formattedDate(album.release_date, album.release_date_precision) }}</span>
              <span class="px-2 py-0.5 rounded-full bg-white/70 dark:bg-white/10">
                {{ (album.album_type || '').toUpperCase() }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div v-if="total > limit" class="flex items-center justify-center gap-3 mt-8">
        <Button variant="outline" size="sm" class="rounded-md" :disabled="offset === 0 || loading" @click="prevPage">{{ $t('NewReleases.prev') }}</Button>
        <span class="text-sm text-zinc-600 dark:text-zinc-400">{{ pageLabel }}</span>
        <Button variant="outline" size="sm" class="rounded-md" :disabled="offset + limit >= total || loading" @click="nextPage">{{ $t('NewReleases.next') }}</Button>
      </div>
    </div>
  </div>
  
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Button } from '@/components/ui/button'
import { RefreshCw, Sparkles } from 'lucide-vue-next'
import fallbackImage from '@/assets/images/default_artist.png'
import { GetNewReleases } from '../../wailsjs/go/main/App'

const router = useRouter()
const loading = ref(false)
const limit = ref(24)
const offset = ref(0)
const total = ref(0)
const albums = ref<any[]>([])

const pageLabel = computed(() => {
  const start = offset.value + 1
  const end = Math.min(offset.value + limit.value, total.value)
  return `${start}-${end} / ${total.value}`
})

function formattedDate(date: string, precision?: string) {
  if (!date) return ''
  // Spotify may return YYYY or YYYY-MM
  const normalized = precision === 'year' ? `${date}-01-01` : precision === 'month' && date.length === 7 ? `${date}-01` : date
  try {
    return new Date(normalized).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return date
  }
}

async function fetchReleases() {
  loading.value = true
  try {
    const res: any = await GetNewReleases(limit.value, offset.value)
    const payload = res?.albums || {}
    total.value = payload.total || 0
    albums.value = payload.items || []
  } finally {
    loading.value = false
  }
}

function nextPage() {
  if (offset.value + limit.value < total.value) {
    offset.value += limit.value
    fetchReleases()
  }
}

function prevPage() {
  if (offset.value >= limit.value) {
    offset.value -= limit.value
  } else {
    offset.value = 0
  }
  fetchReleases()
}

function reload() {
  fetchReleases()
}

onMounted(fetchReleases)

function viewAlbum(id: string) {
  router.push(`/album/${id}`)
}
</script>

<style scoped>
.line-clamp-1 { display: -webkit-box; -webkit-line-clamp: 1; -webkit-box-orient: vertical; overflow: hidden; line-clamp: 1; }
.line-clamp-2 { display: -webkit-box; -webkit-line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; line-clamp: 2; }
</style>


