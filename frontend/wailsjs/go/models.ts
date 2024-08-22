export namespace adapter {
	
	export class FileDTO {
	    name: string;
	    path: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new FileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.content = source["content"];
	    }
	}
	export class GitFileDTO {
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new GitFileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	    }
	}

}

export namespace common {
	
	export class RequestResultDTO {
	    matchType?: string;
	
	    static createFrom(source: any = {}) {
	        return new RequestResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.matchType = source["matchType"];
	    }
	}
	export class ResultDTO {
	    file: string;
	    matchType: string;
	
	    static createFrom(source: any = {}) {
	        return new ResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.matchType = source["matchType"];
	    }
	}

}

export namespace entities {
	
	export class ComponentDTO {
	    id: string;
	    lines?: string;
	    oss_lines?: string;
	    matched?: string;
	    file_hash?: string;
	    source_hash?: string;
	    file_url?: string;
	    purl?: string[];
	    vendor?: string;
	    component?: string;
	    version?: string;
	    latest?: string;
	    url?: string;
	    status?: string;
	    release_date?: string;
	    file?: string;
	    url_hash?: string;
	    // Go type: struct {}
	    url_stats?: any;
	    provenance?: string;
	    licenses?: any[];
	    // Go type: struct { Version string "json:\"version,omitempty\""; KbVersion struct { Monthly string "json:\"monthly,omitempty\""; Daily string "json:\"daily,omitempty\"" } "json:\"kb_version\""; Hostname string "json:\"hostname,omitempty\""; Flags string "json:\"flags,omitempty\""; Elapsed string "json:\"elapsed,omitempty\"" }
	    server: any;
	
	    static createFrom(source: any = {}) {
	        return new ComponentDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.lines = source["lines"];
	        this.oss_lines = source["oss_lines"];
	        this.matched = source["matched"];
	        this.file_hash = source["file_hash"];
	        this.source_hash = source["source_hash"];
	        this.file_url = source["file_url"];
	        this.purl = source["purl"];
	        this.vendor = source["vendor"];
	        this.component = source["component"];
	        this.version = source["version"];
	        this.latest = source["latest"];
	        this.url = source["url"];
	        this.status = source["status"];
	        this.release_date = source["release_date"];
	        this.file = source["file"];
	        this.url_hash = source["url_hash"];
	        this.url_stats = this.convertValues(source["url_stats"], Object);
	        this.provenance = source["provenance"];
	        this.licenses = source["licenses"];
	        this.server = this.convertValues(source["server"], Object);
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
	export class ComponentFilterDTO {
	    path: string;
	    purl: string;
	    usage?: string;
	    version?: string;
	    action: string;
	
	    static createFrom(source: any = {}) {
	        return new ComponentFilterDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.purl = source["purl"];
	        this.usage = source["usage"];
	        this.version = source["version"];
	        this.action = source["action"];
	    }
	}

}

