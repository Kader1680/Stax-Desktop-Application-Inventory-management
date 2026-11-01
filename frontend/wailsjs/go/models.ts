export namespace domain {
	
	export class Quantity {
	
	
	    static createFrom(source: any = {}) {
	        return new Quantity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Product {
	    id: number;
	    name: string;
	    sku: string;
	    barcode?: string;
	    categoryId?: number;
	    categoryName?: string;
	    unit: string;
	    purchasePrice: number;
	    sellingPrice: number;
	    stockQty: Quantity;
	    minStockQty: Quantity;
	    isActive: boolean;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.sku = source["sku"];
	        this.barcode = source["barcode"];
	        this.categoryId = source["categoryId"];
	        this.categoryName = source["categoryName"];
	        this.unit = source["unit"];
	        this.purchasePrice = source["purchasePrice"];
	        this.sellingPrice = source["sellingPrice"];
	        this.stockQty = this.convertValues(source["stockQty"], Quantity);
	        this.minStockQty = this.convertValues(source["minStockQty"], Quantity);
	        this.isActive = source["isActive"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
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
	export class ProductInput {
	    name: string;
	    sku: string;
	    barcode: string;
	    categoryId?: number;
	    unit: string;
	    purchasePrice: number;
	    sellingPrice: number;
	    stockQty: string;
	    minStockQty: string;
	
	    static createFrom(source: any = {}) {
	        return new ProductInput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.sku = source["sku"];
	        this.barcode = source["barcode"];
	        this.categoryId = source["categoryId"];
	        this.unit = source["unit"];
	        this.purchasePrice = source["purchasePrice"];
	        this.sellingPrice = source["sellingPrice"];
	        this.stockQty = source["stockQty"];
	        this.minStockQty = source["minStockQty"];
	    }
	}
	
	export class StockMovement {
	    id: number;
	    productId: number;
	    type: string;
	    qty: Quantity;
	    reason?: string;
	    referenceId?: number;
	    referenceType?: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new StockMovement(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.productId = source["productId"];
	        this.type = source["type"];
	        this.qty = this.convertValues(source["qty"], Quantity);
	        this.reason = source["reason"];
	        this.referenceId = source["referenceId"];
	        this.referenceType = source["referenceType"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
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

