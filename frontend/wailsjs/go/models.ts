export namespace device {
	
	export class BackupInfo {
	    id: string;
	    device_udid: string;
	    device_name: string;
	    backup_path: string;
	    // Go type: time
	    created_at: any;
	    size: number;
	    is_encrypted: boolean;
	    ios_version: string;
	
	    static createFrom(source: any = {}) {
	        return new BackupInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.device_udid = source["device_udid"];
	        this.device_name = source["device_name"];
	        this.backup_path = source["backup_path"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.size = source["size"];
	        this.is_encrypted = source["is_encrypted"];
	        this.ios_version = source["ios_version"];
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
	export class BackupProgress {
	    status: string;
	    progress: number;
	    current_file: string;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new BackupProgress(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.status = source["status"];
	        this.progress = source["progress"];
	        this.current_file = source["current_file"];
	        this.error = source["error"];
	    }
	}
	export class Device {
	    udid: string;
	    name: string;
	    model: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Device(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.udid = source["udid"];
	        this.name = source["name"];
	        this.model = source["model"];
	        this.status = source["status"];
	    }
	}

}

