import { formatDate } from "@angular/common";

export class Reservation{
    reservationId: string;
    accommodationId: string;
    location: string;
    guestEmail: string;
    hostEmail: string;
    price: number;
    numOfPeople: number;
    startDate: Date;
    endDate: Date;
    constructor() {
        this.reservationId = ""
        this.accommodationId = ""
        this.location = ""
        this.guestEmail = ""
        this.hostEmail = ""
        this.price = 0
        this.numOfPeople = 0
        this.startDate = new Date(formatDate(Date.now(),'yyyy-MM-dd','en_us'))
        this.endDate = new Date(formatDate(Date.now(),'yyyy-MM-dd','en_us'))
    }
}