import { defineStore } from 'pinia';
import { ref, Ref } from 'vue';
import { 
    GetSetting,
    SetSetting,
    ValidateAndStoreSpotifyCredentials,
    HasValidSpotifyCredentials,
    CheckForUpdates
} from "../../wailsjs/go/main/App";

export const useSettingsStore = defineStore('settings', () => {
  const spotifyClientId = ref("");
  const spotifyClientSecret = ref("");
  const hasValidCredentials = ref(false);
  const showCredentialsModal = ref(false);
  const lastDownloadPath = ref("");
  const newReleasesDownloadPath = ref("");
  const appendArtistAlbumToPath = ref(false);
  const autoDownloadNewReleases = ref(false);

  // New state for update checking
  const updateAvailable = ref(false);
  const latestVersionTag = ref("");
  const updateReleaseURL = ref("");
  const updateCheckError = ref("");
  const showUpdateDialog = ref(false); // Controls visibility of the dialog

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

  async function fetchNewReleasesDownloadPath() {
    try {
      newReleasesDownloadPath.value = await GetSetting("newReleasesDownloadPath") || "";
    } catch (error) {
      console.error("Error fetching new releases download path:", error);
      newReleasesDownloadPath.value = "";
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

  // Generic loader for boolean settings stored as "true"/"false" strings
  async function loadBooleanSetting(key: string, targetRef: Ref<boolean>) {
    try {
      const stringValue = await GetSetting(key);
      targetRef.value = stringValue === "true";
    } catch (error) {
      console.error(`Error loading ${key} setting:`, error);
      targetRef.value = false;
    }
  }

  async function loadAppendArtistAlbumToPathSetting() {
    await loadBooleanSetting("appendArtistAlbumToPath", appendArtistAlbumToPath);
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

  // Load auto download setting
  async function loadAutoDownloadNewReleasesSetting() {
    await loadBooleanSetting("autoDownloadNewReleases", autoDownloadNewReleases);
  }

  // Toggle auto download setting
  async function toggleAutoDownloadNewReleases() {
    const oldValue = autoDownloadNewReleases.value;
    const newValue = !oldValue;
    autoDownloadNewReleases.value = newValue;
    try {
      await SetSetting("autoDownloadNewReleases", newValue ? "true" : "false");
    } catch (error) {
      console.error("Error saving autoDownloadNewReleases setting:", error);
      autoDownloadNewReleases.value = oldValue;
    }
  }

  // Action for checking updates
  async function performUpdateCheck() {
    try {
      const result = await CheckForUpdates();
      if (result.error) {
        updateCheckError.value = result.error;
        updateAvailable.value = false;
        showUpdateDialog.value = false;
        console.error("Update check failed:", result.error);
      } else {
        updateAvailable.value = result.updateAvailable;
        latestVersionTag.value = result.latestVersion;
        updateReleaseURL.value = result.releaseURL;
        updateCheckError.value = "";
        if (result.updateAvailable) {
          showUpdateDialog.value = true; // Automatically show dialog if update is available
        }
      }
    } catch (e) {
      console.error("Error calling CheckForUpdates:", e);
      updateCheckError.value = "Failed to check for updates.";
      updateAvailable.value = false;
      showUpdateDialog.value = false;
    }
  }

  async function initSettings() {
    await loadSpotifyCredentials();
    await fetchLastDownloadPath();
    await fetchNewReleasesDownloadPath();
    await loadAppendArtistAlbumToPathSetting();
    await loadAutoDownloadNewReleasesSetting();
    await performUpdateCheck();
  }

  return {
    initSettings,
    spotifyClientId,
    spotifyClientSecret,
    hasValidCredentials,
    showCredentialsModal,
    lastDownloadPath,
    newReleasesDownloadPath,
    fetchNewReleasesDownloadPath,
    appendArtistAlbumToPath,
    autoDownloadNewReleases,
    loadSpotifyCredentials,
    saveSpotifyCredentials,
    checkCredentialsValidity,
    fetchLastDownloadPath,
    updateLastDownloadPath,
    toggleAppendArtistAlbumToPath,
    loadAutoDownloadNewReleasesSetting,
    toggleAutoDownloadNewReleases,
    // Update checker state and actions
    updateAvailable,
    latestVersionTag,
    updateReleaseURL,
    updateCheckError,
    showUpdateDialog,
    performUpdateCheck
  };
}); 