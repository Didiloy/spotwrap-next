
interface Album {
    id: string;
    name: string;
    release_date: string;
    total_tracks: number;
    artists: SimplifiedArtist[];
    images: Image[];
    external_urls: ExternalUrls;
    album_type: string;
  }
  
  interface Track {
    id: string;
    name: string;
    duration_ms: number;
    explicit: boolean;
    preview_url: string | null;
    artists: SimplifiedArtist[];
    album: SimplifiedAlbum;
    external_urls: ExternalUrls;
    track_number: number;
  }
  
  interface Artist {
    id: string;
    name: string;
    genres: string[];
    popularity: number;
    images: Image[];
    external_urls: ExternalUrls;
  }

  interface SimplifiedArtist {
    id: string;
    name: string;
    external_urls: ExternalUrls;
  }
  
  interface SimplifiedAlbum {
    id: string;
    name: string;
    images: Image[];
    external_urls: ExternalUrls;
  }
  
  interface Image {
    url: string;
    height: number;
    width: number;
  }
  
  interface ExternalUrls {
    spotify: string;
  }

  export interface SearchResults {
    albums?: { items: Album[] };
    tracks?: { items: Track[] };
    artists?: { items: Artist[] };
  }