import { defineStore } from 'pinia';
import { ref } from 'vue';
import { 
    GetSetting,
    SetSetting,
    ValidateAndStoreSpotifyCredentials,
    HasValidSpotifyCredentials
} from "../../wailsjs/go/main/App";

export const useSettingsStore = defineStore('settings', () => {
  const spotifyClientId = ref("");
  const spotifyClientSecret = ref("");
  const hasValidCredentials = ref(false);
  const showCredentialsModal = ref(false);
  const lastDownloadPath = ref("");
  const appendArtistAlbumToPath = ref(false);

  async function loadSpotifyCredentials() {
    try {
      spotifyClientId.value = await GetSetting("spotify_client_id") || "";
      spotifyClientSecret.value = await GetSetting("spotify_client_secret") || "";
      await checkCredentialsValidity();
    } catch (error) {
      console.error("Error loading Spotify credentials:", error);
      spotifyClientId.value = "";
      spotifyClientSecret.value = "";
      hasValidCredentials.value = false;
    }
  }


  async function saveSpotifyCredentials() {
    const success = await ValidateAndStoreSpotifyCredentials(
      spotifyClientId.value,
      spotifyClientSecret.value
    );
    hasValidCredentials.value = success;
    if (success) {
      await checkCredentialsValidity();
    }
    return success;
  }

  async function checkCredentialsValidity() {
    hasValidCredentials.value = await HasValidSpotifyCredentials();
    return hasValidCredentials.value;
  }

  async function fetchLastDownloadPath() {
    try {
      lastDownloadPath.value = await GetSetting("lastDownloadPath") || "";
    } catch (error) {
      console.error("Error fetching last download path:", error);
      lastDownloadPath.value = "";
    }
  }

  async function updateLastDownloadPath(newPath: string) {
    try {
      await SetSetting("lastDownloadPath", newPath);
      lastDownloadPath.value = newPath;
    } catch (error) {
      console.error("Error updating last download path:", error);
    }
  }

  async function loadAppendArtistAlbumToPathSetting() {
    try {
      const stringValue = await GetSetting("appendArtistAlbumToPath");
      appendArtistAlbumToPath.value = stringValue === "true";
    } catch (error) {
      console.error("Error loading appendArtistAlbumToPath setting:", error);
      appendArtistAlbumToPath.value = false;
    }
  }

  async function toggleAppendArtistAlbumToPath() {
    const oldValue = appendArtistAlbumToPath.value;
    const newValue = !oldValue;
    appendArtistAlbumToPath.value = newValue;
    try {
      await SetSetting("appendArtistAlbumToPath", newValue ? "true" : "false");
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