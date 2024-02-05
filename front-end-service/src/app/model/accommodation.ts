import {formatDate} from '@angular/common'
export class Accommodation{
    id: string;
    Owner: string; //change to type of User model later
    ownerId: string;
    name: string;
    location: string;
    benefits: string;
    minCapacity: number;
    maxCapacity: number;
    price: number;
    payPer: string;
    constructor() {
        this.id = ""
        this.Owner = ""
        this.ownerId = ""
        this.name = ""
        this.location = ""
        this.benefits = ""
        this.minCapacity = 0
        this.maxCapacity = 0
        this.price = 0
        this.payPer = ""
      }
}
