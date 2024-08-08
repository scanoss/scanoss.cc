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

