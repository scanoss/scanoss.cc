export namespace entities {
	
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
	    ToggleSyncScrollPosition = "toggleSyncScrollPosition",
	    ShowKeyboardShortcutsModal = "showKeyboardShortcutsModal",
	    ScanWithOptions = "scanWithOptions",
	    OpenSettings = "openSettings",
	}
	export class ComponentFilter {
	    path?: string;
	    purl: string;
	    usage?: string;
	    comment?: string;
	    replace_with?: string;
	    license?: string;
	
	    static createFrom(source: any = {}) {
	        return new ComponentFilter(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.purl = source["purl"];
	        this.usage = source["usage"];
	        this.comment = source["comment"];
	        this.replace_with = source["replace_with"];
	        this.license = source["license"];
	    }
	}
	export class Bom {
	    include?: ComponentFilter[];
	    remove?: ComponentFilter[];
	    replace?: ComponentFilter[];
	
	    static createFrom(source: any = {}) {
	        return new Bom(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.include = this.convertValues(source["include"], ComponentFilter);
	        this.remove = this.convertValues(source["remove"], ComponentFilter);
	        this.replace = this.convertValues(source["replace"], ComponentFilter);
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
	    replace_with?: string;
	    license?: string;
	
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
	        this.replace_with = source["replace_with"];
	        this.license = source["license"];
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
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new FileDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.content = source["content"];
	        this.language = source["language"];
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
	export class InitialFilters {
	    Include: ComponentFilter[];
	    Remove: ComponentFilter[];
	    Replace: ComponentFilter[];
	
	    static createFrom(source: any = {}) {
	        return new InitialFilters(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Include = this.convertValues(source["Include"], ComponentFilter);
	        this.Remove = this.convertValues(source["Remove"], ComponentFilter);
	        this.Replace = this.convertValues(source["Replace"], ComponentFilter);
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
	export class License {
	    name: string;
	    licenseId: string;
	    reference: string;
	
	    static createFrom(source: any = {}) {
	        return new License(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.licenseId = source["licenseId"];
	        this.reference = source["reference"];
	    }
	}
	export class SortConfig {
	    option: string;
	    order: string;
	
	    static createFrom(source: any = {}) {
	        return new SortConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.option = source["option"];
	        this.order = source["order"];
	    }
	}
	export class RequestResultDTO {
	    match_type?: string;
	    query?: string;
	    sort?: SortConfig;
	
	    static createFrom(source: any = {}) {
	        return new RequestResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.match_type = source["match_type"];
	        this.query = source["query"];
	        this.sort = this.convertValues(source["sort"], SortConfig);
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
	export class ResultDTO {
	    path: string;
	    match_type: string;
	    workflow_state?: string;
	    filter_config?: FilterConfig;
	    comment?: string;
	    detected_purl?: string;
	    detected_purl_url?: string;
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
	        this.detected_purl_url = source["detected_purl_url"];
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
	export class ScanArgDef {
	    Name: string;
	    Shorthand: string;
	    Default: any;
	    Usage: string;
	    Tooltip: string;
	    Type: string;
	    IsCore: boolean;
	    IsFileSelector: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ScanArgDef(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Name = source["Name"];
	        this.Shorthand = source["Shorthand"];
	        this.Default = source["Default"];
	        this.Usage = source["Usage"];
	        this.Tooltip = source["Tooltip"];
	        this.Type = source["Type"];
	        this.IsCore = source["IsCore"];
	        this.IsFileSelector = source["IsFileSelector"];
	    }
	}
	export class SizesSkipSettings {
	    patterns?: string[];
	    min?: number;
	    max?: number;
	
	    static createFrom(source: any = {}) {
	        return new SizesSkipSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.patterns = source["patterns"];
	        this.min = source["min"];
	        this.max = source["max"];
	    }
	}
	export class Sizes {
	    scanning?: SizesSkipSettings[];
	    fingerprinting?: SizesSkipSettings[];
	
	    static createFrom(source: any = {}) {
	        return new Sizes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scanning = this.convertValues(source["scanning"], SizesSkipSettings);
	        this.fingerprinting = this.convertValues(source["fingerprinting"], SizesSkipSettings);
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
	export class SkipPatterns {
	    scanning?: string[];
	    fingerprinting?: string[];
	
	    static createFrom(source: any = {}) {
	        return new SkipPatterns(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scanning = source["scanning"];
	        this.fingerprinting = source["fingerprinting"];
	    }
	}
	export class SkipSettings {
	    patterns?: SkipPatterns;
	    sizes?: Sizes;
	
	    static createFrom(source: any = {}) {
	        return new SkipSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.patterns = this.convertValues(source["patterns"], SkipPatterns);
	        this.sizes = this.convertValues(source["sizes"], Sizes);
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
	export class ScanossSettingsSchema {
	    skip?: SkipSettings;
	
	    static createFrom(source: any = {}) {
	        return new ScanossSettingsSchema(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.skip = this.convertValues(source["skip"], SkipSettings);
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
	export class SettingsFile {
	    settings?: ScanossSettingsSchema;
	    bom?: Bom;
	
	    static createFrom(source: any = {}) {
	        return new SettingsFile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.settings = this.convertValues(source["settings"], ScanossSettingsSchema);
	        this.bom = this.convertValues(source["bom"], Bom);
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
	
	
	
	
	
	export class TreeNode {
	    id: string;
	    name: string;
	    path: string;
	    isFolder: boolean;
	    workflowState: string;
	    scanningSkipState: string;
	    children: TreeNode[];
	
	    static createFrom(source: any = {}) {
	        return new TreeNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isFolder = source["isFolder"];
	        this.workflowState = source["workflowState"];
	        this.scanningSkipState = source["scanningSkipState"];
	        this.children = this.convertValues(source["children"], TreeNode);
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

