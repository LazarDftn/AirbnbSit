import {formatDate} from '@angular/common'
export class Accommodation{
    Owner: string; //change to type of User model later
    name: string;
    location: string;
    benefits: string;
    minCapacity: number;
    maxCapacity: number;
    price: number;
    payPer: number;
    constructor() {
        this.Owner = ""
        this.name = ""
        this.location = ""
        this.benefits = ""
        this.minCapacity = 0
        this.maxCapacity = 0
        this.price = 0
        this.payPer = 0
      }
}
