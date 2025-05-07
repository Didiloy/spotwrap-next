import { defineStore } from 'pinia';
import { ref } from 'vue';
import { 
    GetSpotifyCredentials, 
    SetSpotifyCredentials, 
    HasValidSpotifyCredentials 
} from "../../wailsjs/go/main/App";

export const useSettingsStore = defineStore('settings', () => {
  const spotifyClientId = ref("");
  const spotifyClientSecret = ref("");
  const hasValidCredentials = ref(false);
  const showCredentialsModal = ref(false);

  async function loadSpotifyCredentials() {
    const credentials = await GetSpotifyCredentials();
    spotifyClientId.value = credentials.clientId || "";
    spotifyClientSecret.value = credentials.clientSecret || "";
    await checkCredentialsValidity();
  }

  async function saveSpotifyCredentials() {
    const success = await SetSpotifyCredentials(
      spotifyClientId.value,
      spotifyClientSecret.value
    );
    
    hasValidCredentials.value = success;
    return success;
  }

  async function checkCredentialsValidity() {
    hasValidCredentials.value = await HasValidSpotifyCredentials();
    return hasValidCredentials.value;
  }

  return {
    spotifyClientId,
    spotifyClientSecret,
    hasValidCredentials,
    showCredentialsModal,
    loadSpotifyCredentials, 
    saveSpotifyCredentials,
    checkCredentialsValidity
  };
}); 