export namespace device {
	
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

