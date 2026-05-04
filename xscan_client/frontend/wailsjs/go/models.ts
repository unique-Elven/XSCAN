export namespace main {
	
	export class AppConfigDTO {
	    modelStrategy: string;
	    soundEnabled: boolean;
	    language: string;
	    path2018: string;
	    path2024: string;
	    threshold2018: number;
	    threshold2024: number;
	
	    static createFrom(source: any = {}) {
	        return new AppConfigDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.modelStrategy = source["modelStrategy"];
	        this.soundEnabled = source["soundEnabled"];
	        this.language = source["language"];
	        this.path2018 = source["path2018"];
	        this.path2024 = source["path2024"];
	        this.threshold2018 = source["threshold2018"];
	        this.threshold2024 = source["threshold2024"];
	    }
	}
	export class FileScoreResult {
	    path: string;
	    score: number;
	    size: number;
	    modelUsed?: string;
	    errMsg?: string;
	
	    static createFrom(source: any = {}) {
	        return new FileScoreResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.score = source["score"];
	        this.size = source["size"];
	        this.modelUsed = source["modelUsed"];
	        this.errMsg = source["errMsg"];
	    }
	}
	export class HistorySaveEntry {
	    scannedAt: number;
	    filePath: string;
	    fileHash: string;
	    verdict: string;
	    fileSize: number;
	    engine: string;
	    score: number;
	
	    static createFrom(source: any = {}) {
	        return new HistorySaveEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.scannedAt = source["scannedAt"];
	        this.filePath = source["filePath"];
	        this.fileHash = source["fileHash"];
	        this.verdict = source["verdict"];
	        this.fileSize = source["fileSize"];
	        this.engine = source["engine"];
	        this.score = source["score"];
	    }
	}
	export class LightGBMModelStatus {
	    loaded: boolean;
	    modelPath: string;
	    nFeatures: number;
	
	    static createFrom(source: any = {}) {
	        return new LightGBMModelStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.loaded = source["loaded"];
	        this.modelPath = source["modelPath"];
	        this.nFeatures = source["nFeatures"];
	    }
	}
	export class ScanHistoryRow {
	    id: number;
	    scannedAt: number;
	    filePath: string;
	    fileHash: string;
	    verdict: string;
	    fileSize: number;
	    engine: string;
	    score: number;
	
	    static createFrom(source: any = {}) {
	        return new ScanHistoryRow(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.scannedAt = source["scannedAt"];
	        this.filePath = source["filePath"];
	        this.fileHash = source["fileHash"];
	        this.verdict = source["verdict"];
	        this.fileSize = source["fileSize"];
	        this.engine = source["engine"];
	        this.score = source["score"];
	    }
	}
	export class ScannerModelPaths {
	    path2018: string;
	    path2024: string;
	
	    static createFrom(source: any = {}) {
	        return new ScannerModelPaths(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path2018 = source["path2018"];
	        this.path2024 = source["path2024"];
	    }
	}

}

