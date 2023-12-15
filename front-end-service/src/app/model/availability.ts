import { formatDate } from "@angular/common";

export class Availability{
    availabilityId: string;
    accommId: string;
    name: string;
    location: string;
    minCapacity: number;
    maxCapacity: number;
    startDate: Date;
    endDate: Date;
    constructor() {
        this.availabilityId = ""
        this.accommId = ""
        this.name = ""
        this.location = ""
        this.minCapacity = 0
        this.maxCapacity = 0
        this.startDate = new Date(formatDate(Date.now(),'yyyy-MM-dd','en_us'))
        this.endDate = new Date(formatDate(Date.now(),'yyyy-MM-dd','en_us'))
    }
}