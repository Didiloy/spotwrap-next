import { defineStore } from 'pinia';
import { ref } from 'vue';
import { 
    GetSpotifyCredentials, 
    SetSpotifyCredentials, 
    HasValidSpotifyCredentials,
    SaveLastDownloadPath,
    GetLastDownloadPath,
    SaveAppendArtistAlbumToPath,
    GetAppendArtistAlbumToPath
} from "../../wailsjs/go/main/App";

export const useSettingsStore = defineStore('settings', () => {
  const spotifyClientId = ref("");
  const spotifyClientSecret = ref("");
  const hasValidCredentials = ref(false);
  const showCredentialsModal = ref(false);
  const lastDownloadPath = ref("");
  const appendArtistAlbumToPath = ref(false);

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

  async function loadAppendArtistAlbumToPathSetting() {
    try {
      appendArtistAlbumToPath.value = await GetAppendArtistAlbumToPath();
    } catch (error) {
      console.error("Error loading appendArtistAlbumToPath setting:", error);
      appendArtistAlbumToPath.value = false;
    }
  }

  async function toggleAppendArtistAlbumToPath() {
    const oldValue = appendArtistAlbumToPath.value;
    appendArtistAlbumToPath.value = !oldValue;
    try {
      await SaveAppendArtistAlbumToPath(appendArtistAlbumToPath.value);
    } catch (error) {
      console.error("Error saving appendArtistAlbumToPath setting:", error);
      appendArtistAlbumToPath.value = oldValue;
    }
  }

  async function initSettings() {
    await loadSpotifyCredentials();
    await fetchLastDownloadPath();
    await loadAppendArtistAlbumToPathSetting();
  }


  return {
    initSettings,
    spotifyClientId,
    spotifyClientSecret,
    hasValidCredentials,
    showCredentialsModal,
    lastDownloadPath,
    appendArtistAlbumToPath,
    loadSpotifyCredentials, 
    saveSpotifyCredentials,
    checkCredentialsValidity,
    fetchLastDownloadPath,
    updateLastDownloadPath,
    toggleAppendArtistAlbumToPath
  };
}); 