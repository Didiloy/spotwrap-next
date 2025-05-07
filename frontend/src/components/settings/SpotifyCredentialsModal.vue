<script lang="ts" setup>
import { ref, watch } from 'vue';
import { useSettingsStore } from '@/store/settings';
import { useI18n } from 'vue-i18n';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { useToast } from '@/components/ui/toast/use-toast';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';

const { t } = useI18n();
const { toast } = useToast();
const settingsStore = useSettingsStore();

const props = defineProps<{
  open: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void;
  (e: 'saved'): void;
}>();

const isSaving = ref(false);
const errorMessage = ref("");

watch(() => props.open, async (newValue) => {
  if (newValue) {
    await settingsStore.loadSpotifyCredentials();
    errorMessage.value = "";
  }
});

function onClose() {
  emit('update:open', false);
}

async function saveCredentials() {
  if (!settingsStore.spotifyClientId || !settingsStore.spotifyClientSecret) {
    errorMessage.value = t('Settings.spotify_credentials_empty');
    return;
  }

  isSaving.value = true;
  errorMessage.value = "";
  
  try {
    const success = await settingsStore.saveSpotifyCredentials();
    
    if (success) {
      toast({
        title: t('Settings.spotify_credentials_saved_title'),
        description: t('Settings.spotify_credentials_saved_message'),
      });
      emit('saved');
      onClose();
    } else {
      errorMessage.value = t('Settings.spotify_credentials_invalid');
    }
  } catch (error) {
    errorMessage.value = t('Settings.spotify_credentials_error');
    console.error("Error saving Spotify credentials:", error);
  } finally {
    isSaving.value = false;
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>{{ t('Settings.spotify_credentials_title') }}</DialogTitle>
        <DialogDescription>
          {{ t('Settings.spotify_credentials_description') }}
        </DialogDescription>
      </DialogHeader>
      
      <div class="space-y-4 py-4">
        <div class="space-y-2">
          <Label for="client-id">{{ t('Settings.spotify_client_id') }}</Label>
          <Input 
            id="client-id" 
            v-model="settingsStore.spotifyClientId" 
            :placeholder="t('Settings.enter_client_id')"
          />
        </div>
        
        <div class="space-y-2">
          <Label for="client-secret">{{ t('Settings.spotify_client_secret') }}</Label>
          <Input 
            id="client-secret" 
            v-model="settingsStore.spotifyClientSecret" 
            :placeholder="t('Settings.enter_client_secret')"
            type="password"
          />
        </div>
        
        <div v-if="errorMessage" class="text-red-500 text-sm">
          {{ errorMessage }}
        </div>
        
        <div class="text-sm text-muted-foreground">
          {{ t('Settings.spotify_credentials_help') }}
        </div>
      </div>
      
      <DialogFooter>
        <Button 
          type="button" 
          variant="outline" 
          @click="onClose"
        >
          {{ t('Settings.cancel') }}
        </Button>
        <Button 
          type="button" 
          @click="saveCredentials" 
          :disabled="isSaving"
        >
          <span v-if="isSaving">{{ t('Settings.saving') }}</span>
          <span v-else>{{ t('Settings.save') }}</span>
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template> 