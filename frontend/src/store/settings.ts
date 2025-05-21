import { defineStore } from 'pinia';
import { ref } from 'vue';
import { 
    GetSpotifyCredentials, 
    SetSpotifyCredentials, 
    HasValidSpotifyCredentials,
    SaveLastDownloadPath,
    GetLastDownloadPath
} from "../../wailsjs/go/main/App";

export const useSettingsStore = defineStore('settings', () => {
  const spotifyClientId = ref("");
  const spotifyClientSecret = ref("");
  const hasValidCredentials = ref(false);
  const showCredentialsModal = ref(false);
  const lastDownloadPath = ref("");

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

  async function fetchLastDownloadPath() {
    try {
      const path = await GetLastDownloadPath();
      lastDownloadPath.value = path || "";
    } catch (error) {
      console.error("Error fetching last download path:", error);
      lastDownloadPath.value = "";
    }
  }

  async function updateLastDownloadPath(newPath: string) {
    try {
      await SaveLastDownloadPath(newPath);
      lastDownloadPath.value = newPath;
    } catch (error) {
      console.error("Error updating last download path:", error);
    }
  }

  return {
    spotifyClientId,
    spotifyClientSecret,
    hasValidCredentials,
    showCredentialsModal,
    lastDownloadPath,
    loadSpotifyCredentials, 
    saveSpotifyCredentials,
    checkCredentialsValidity,
    fetchLastDownloadPath,
    updateLastDownloadPath
  };
}); 