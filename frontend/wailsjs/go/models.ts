export namespace entities {
	
	export class ComponentDTO {
	    id: string;
	    lines?: string;
	    oss_lines?: string;
	    matched?: string;
	    file_hash?: string;
	    source_hash?: string;
	    file_url?: string;
	    purl: string[];
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
	    path?: string;
	    purl: string;
	    usage?: string;
	    action: string;
	    comment?: string;
	    replace_with_purl?: string;
	    replace_with_name?: string;
	
	    static createFrom(source: any = {}) {
	        return new ComponentFilterDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.purl = source["purl"];
	        this.usage = source["usage"];
	        this.action = source["action"];
	        this.comment = source["comment"];
	        this.replace_with_purl = source["replace_with_purl"];
	        this.replace_with_name = source["replace_with_name"];
	    }
	}
	export class DeclaredComponent {
	    name: string;
	    purl: string;
	
	    static createFrom(source: any = {}) {
	        return new DeclaredComponent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.purl = source["purl"];
	    }
	}
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
	export class FilterConfig {
	    action?: string;
	    type?: string;
	
	    static createFrom(source: any = {}) {
	        return new FilterConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.action = source["action"];
	        this.type = source["type"];
	    }
	}
	export class RequestResultDTO {
	    match_type?: string;
	    query?: string;
	
	    static createFrom(source: any = {}) {
	        return new RequestResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.match_type = source["match_type"];
	        this.query = source["query"];
	    }
	}
	export class ResultDTO {
	    path: string;
	    match_type: string;
	    workflow_state?: string;
	    filter_config?: FilterConfig;
	    comment?: string;
	    detected_purl?: string;
	    concluded_purl?: string;
	    concluded_purl_url?: string;
	    concluded_name?: string;
	
	    static createFrom(source: any = {}) {
	        return new ResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.match_type = source["match_type"];
	        this.workflow_state = source["workflow_state"];
	        this.filter_config = this.convertValues(source["filter_config"], FilterConfig);
	        this.comment = source["comment"];
	        this.detected_purl = source["detected_purl"];
	        this.concluded_purl = source["concluded_purl"];
	        this.concluded_purl_url = source["concluded_purl_url"];
	        this.concluded_name = source["concluded_name"];
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

export namespace keyboard {
	
	export enum Action {
	    Undo = "undo",
	    Redo = "redo",
	    Save = "save",
	    Confirm = "confirm",
	    FocusSearch = "focusSearch",
	    SelectAll = "selectAll",
	    MoveUp = "moveUp",
	    MoveDown = "moveDown",
	    IncludeFileWithoutComments = "includeFileWithoutComments",
	    IncludeFileWithComments = "includeFileWithComments",
	    IncludeComponentWithoutComments = "includeComponentWithoutComments",
	    IncludeComponentWithComments = "includeComponentWithComments",
	    DismissFileWithoutComments = "dismissFileWithoutComments",
	    DismissFileWithComments = "dismissFileWithComments",
	    DismissComponentWithoutComments = "dismissComponentWithoutComments",
	    DismissComponentWithComments = "dismissComponentWithComments",
	    ReplaceFileWithoutComments = "replaceFileWithoutComments",
	    ReplaceFileWithComments = "replaceFileWithComments",
	    ReplaceComponentWithoutComments = "replaceComponentWithoutComments",
	    ReplaceComponentWithComments = "replaceComponentWithComments",
	}
	export class Shortcut {
	    name: string;
	    description: string;
	    accelerator?: keys.Accelerator;
	    alternativeAccelerator?: keys.Accelerator;
	    keys: string;
	    group?: string;
	    action?: Action;
	
	    static createFrom(source: any = {}) {
	        return new Shortcut(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.description = source["description"];
	        this.accelerator = this.convertValues(source["accelerator"], keys.Accelerator);
	        this.alternativeAccelerator = this.convertValues(source["alternativeAccelerator"], keys.Accelerator);
	        this.keys = source["keys"];
	        this.group = source["group"];
	        this.action = source["action"];
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

export namespace keys {
	
	export class Accelerator {
	    Key: string;
	    Modifiers: string[];
	
	    static createFrom(source: any = {}) {
	        return new Accelerator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Key = source["Key"];
	        this.Modifiers = source["Modifiers"];
	    }
	}

}

export namespace scanoss_settings {
	
	export class Module {
	    Controller: any;
	    Service: any;
	
	    static createFrom(source: any = {}) {
	        return new Module(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Controller = source["Controller"];
	        this.Service = source["Service"];
	    }
	}

}

