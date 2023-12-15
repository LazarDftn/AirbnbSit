import { formatDate } from "@angular/common";

export class PriceVariation{
    variationId: string;
    location: string;
    accommId: string;
    startDate: Date;
    endDate: Date;
    percentage: number;
    constructor() {
        this.variationId = ""
        this.location = ""
        this.accommId = ""
        this.percentage = 0
        this.startDate = new Date(formatDate(Date.now(),'yyyy-MM-dd','en_us'))
        this.endDate = new Date(formatDate(Date.now(),'yyyy-MM-dd','en_us'))
    }
}