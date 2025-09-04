export namespace main {
	
	export class FileInfo {
	    name: string;
	    path: string;
	    isDir: boolean;
	    size: number;
	    modified: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.isDir = source["isDir"];
	        this.size = source["size"];
	        this.modified = source["modified"];
	    }
	}
	export class FileExplorerState {
	    currentDir: FileInfo;
	    selectedFile: FileInfo;
	
	    static createFrom(source: any = {}) {
	        return new FileExplorerState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentDir = this.convertValues(source["currentDir"], FileInfo);
	        this.selectedFile = this.convertValues(source["selectedFile"], FileInfo);
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
	
	export class HurlTimings {
	    app_connect: number;
	    begin_call: string;
	    connect: number;
	    end_call: string;
	    name_lookup: number;
	    pre_transfer: number;
	    start_transfer: number;
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new HurlTimings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.app_connect = source["app_connect"];
	        this.begin_call = source["begin_call"];
	        this.connect = source["connect"];
	        this.end_call = source["end_call"];
	        this.name_lookup = source["name_lookup"];
	        this.pre_transfer = source["pre_transfer"];
	        this.start_transfer = source["start_transfer"];
	        this.total = source["total"];
	    }
	}
	export class HurlCertificate {
	    expire_date: string;
	    issuer: string;
	    serial_number: string;
	    start_date: string;
	    subject: string;
	
	    static createFrom(source: any = {}) {
	        return new HurlCertificate(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.expire_date = source["expire_date"];
	        this.issuer = source["issuer"];
	        this.serial_number = source["serial_number"];
	        this.start_date = source["start_date"];
	        this.subject = source["subject"];
	    }
	}
	export class HurlResponse {
	    body: string;
	    bodyContent: string;
	    certificate?: HurlCertificate;
	    cookies: HurlCookie[];
	    headers: HurlHeader[];
	    http_version: string;
	    status: number;
	
	    static createFrom(source: any = {}) {
	        return new HurlResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.body = source["body"];
	        this.bodyContent = source["bodyContent"];
	        this.certificate = this.convertValues(source["certificate"], HurlCertificate);
	        this.cookies = this.convertValues(source["cookies"], HurlCookie);
	        this.headers = this.convertValues(source["headers"], HurlHeader);
	        this.http_version = source["http_version"];
	        this.status = source["status"];
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
	export class HurlQueryParam {
	    name: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new HurlQueryParam(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class HurlHeader {
	    name: string;
	    value: string;
	
	    static createFrom(source: any = {}) {
	        return new HurlHeader(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	    }
	}
	export class HurlCookie {
	    name: string;
	    value: string;
	    domain?: string;
	    path?: string;
	    expires?: string;
	    max_age?: number;
	    http_only?: boolean;
	    secure?: boolean;
	    same_site?: string;
	
	    static createFrom(source: any = {}) {
	        return new HurlCookie(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.value = source["value"];
	        this.domain = source["domain"];
	        this.path = source["path"];
	        this.expires = source["expires"];
	        this.max_age = source["max_age"];
	        this.http_only = source["http_only"];
	        this.secure = source["secure"];
	        this.same_site = source["same_site"];
	    }
	}
	export class HurlRequest {
	    cookies: HurlCookie[];
	    headers: HurlHeader[];
	    method: string;
	    query_string: HurlQueryParam[];
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new HurlRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cookies = this.convertValues(source["cookies"], HurlCookie);
	        this.headers = this.convertValues(source["headers"], HurlHeader);
	        this.method = source["method"];
	        this.query_string = this.convertValues(source["query_string"], HurlQueryParam);
	        this.url = source["url"];
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
	export class HurlCall {
	    request: HurlRequest;
	    response: HurlResponse;
	    timings: HurlTimings;
	
	    static createFrom(source: any = {}) {
	        return new HurlCall(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.request = this.convertValues(source["request"], HurlRequest);
	        this.response = this.convertValues(source["response"], HurlResponse);
	        this.timings = this.convertValues(source["timings"], HurlTimings);
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
	
	
	export class HurlEntry {
	    asserts: any[];
	    calls: HurlCall[];
	    captures: any[];
	    curl_cmd: string;
	    index: number;
	    line: number;
	    time: number;
	
	    static createFrom(source: any = {}) {
	        return new HurlEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.asserts = source["asserts"];
	        this.calls = this.convertValues(source["calls"], HurlCall);
	        this.captures = source["captures"];
	        this.curl_cmd = source["curl_cmd"];
	        this.index = source["index"];
	        this.line = source["line"];
	        this.time = source["time"];
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
	
	
	
	
	export class HurlSession {
	    cookies: HurlCookie[];
	    entries: HurlEntry[];
	    filename: string;
	    success: boolean;
	    time: number;
	
	    static createFrom(source: any = {}) {
	        return new HurlSession(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cookies = this.convertValues(source["cookies"], HurlCookie);
	        this.entries = this.convertValues(source["entries"], HurlEntry);
	        this.filename = source["filename"];
	        this.success = source["success"];
	        this.time = source["time"];
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
	
	export class ReturnValue {
	    fileContent?: string;
	    fileExplorer: FileExplorerState;
	    files: FileInfo[];
	    error?: string;
	    hurlReport?: HurlSession[];
	    envs?: string[];
	    envFilePath?: string;
	
	    static createFrom(source: any = {}) {
	        return new ReturnValue(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fileContent = source["fileContent"];
	        this.fileExplorer = this.convertValues(source["fileExplorer"], FileExplorerState);
	        this.files = this.convertValues(source["files"], FileInfo);
	        this.error = source["error"];
	        this.hurlReport = this.convertValues(source["hurlReport"], HurlSession);
	        this.envs = source["envs"];
	        this.envFilePath = source["envFilePath"];
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

