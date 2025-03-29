export namespace database {
	
	export class Artist {
	    SpotifyID: string;
	    // Go type: time
	    LastChecked: any;
	
	    static createFrom(source: any = {}) {
	        return new Artist(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.SpotifyID = source["SpotifyID"];
	        this.LastChecked = this.convertValues(source["LastChecked"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

