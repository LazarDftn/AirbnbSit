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
    discWeekend: boolean;
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
        this.discDateStart = new Date()
        this.discDateEnd = new Date()
        this.discWeekend= false
        this.payPer = 0
      }
}
