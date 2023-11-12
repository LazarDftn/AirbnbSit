export class Accommodation{
    Owner: string; //change to type of User model later
    name: string;
    location: string;
    benefits: string;
    minCapacity: number;
    maxCapacity: number;
    constructor() {
        this.Owner = ""
        this.name = ""
        this.location = ""
        this.benefits = ""
        this.minCapacity = 0
        this.maxCapacity = 0
      }
}