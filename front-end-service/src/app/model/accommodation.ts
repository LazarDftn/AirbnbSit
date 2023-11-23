import {formatDate} from '@angular/common'
export class Accommodation{
    Owner: string; //change to type of User model later
    name: string;
    location: string;
    benefits: string;
    minCapacity: number;
    maxCapacity: number;
    price: number;
    discPrice: number;
    discDateStart: Date;
    discDateEnd: Date;
    discWeekend: number;
    payPer: number;
    constructor() {
        this.Owner = ""
        this.name = ""
        this.location = ""
        this.benefits = ""
        this.minCapacity = 0
        this.maxCapacity = 0
        this.price = 0
        this.discPrice = 0
        this.discDateStart = new Date(formatDate(Date.now(),'yyyy-MM-dd','en_us'))//da bi se format poklapao sa
        this.discDateEnd = new Date(formatDate(Date.now(),'yyyy-MM-dd','en_us'))//formatom iz date html elementa
        this.discWeekend= 0
        this.payPer = 0
      }
}
