import {formatDate} from '@angular/common'
export class Accommodation{
    Id: string;
    Owner: string; //change to type of User model later
    name: string;
    location: string;
    benefits: string;
    minCapacity: number;
    maxCapacity: number;
    price: number;
    payPer: string;
    constructor() {
        this.Id = ""
        this.Owner = ""
        this.name = ""
        this.location = ""
        this.benefits = ""
        this.minCapacity = 0
        this.maxCapacity = 0
        this.price = 0
        this.payPer = ""
      }
}
