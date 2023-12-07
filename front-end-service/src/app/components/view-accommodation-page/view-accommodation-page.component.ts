import { AfterViewInit, Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { Accommodation } from 'src/app/model/accommodation';
import { Reservation } from 'src/app/model/reservation';
import { AccommodationService } from 'src/app/services/accommodation.service';
import { AuthService } from 'src/app/services/auth.service';
import { ReservationService } from 'src/app/services/reservation.service';

@Component({
  selector: 'app-view-accommodation-page',
  templateUrl: './view-accommodation-page.component.html',
  styleUrls: ['./view-accommodation-page.component.css']
})
export class ViewAccommodationPageComponent implements OnInit {

  constructor(public authService: AuthService,
    private accommService: AccommodationService,
    private reservationService: ReservationService,
    private router: Router,
    private route: ActivatedRoute){}
    
  accommId = ""
  price = 0 //accommodation price per day
  finalPrice = 0 //after calculating the days (and guests if that's the payment method)
  accomm: Accommodation = new Accommodation()
  reservation: Reservation = new Reservation()
  reservations: Reservation[] = []
  submitted = false
  loadedPrice = false
  priceMessage = ""
  additionalPriceMessage = ""

  ngOnInit(): void {
    
    this.route.paramMap.subscribe((params: ParamMap) => {
      
      this.accommId = params.get("id")!

    })
    
    this.accommService.getAccommById(this.accommId + "").subscribe(data => {
      this.accomm.name = data.name
      this.accomm.location = data.location
      this.accomm.benefits = data.benefits
      this.accomm.Owner = data.owner
      this.accomm.Id = data.id
      this.accomm.minCapacity = data.minCapacity
      this.accomm.maxCapacity = data.maxCapacity
      this.accomm.price = this.price

      this.reservationService.getReservationsByAccommId(this.accomm.location, this.accommId + "").subscribe(data => {
        this.reservations = data
      })
    })

    this.reservationService.getPriceByAccommId(this.accommId + "").subscribe(data => {
      console.log(data)
      this.accomm.price = data.price
      this.accomm.payPer = data.payPer
      this.price = data.price
      this.loadedPrice = true
    })

  }

  addPrice(){
    this.reservationService.createPrice(this.accomm.price, this.accomm.payPer, this.accomm.Id).subscribe
    (data => {window.location.reload()})
  }

  // before Guest makes the reservation, check service for any price variations for given period and calculate 
  checkPrice(){
    this.priceMessage = ""
    this.additionalPriceMessage = ""

    this.submitted = true

    this.reservationService.checkPrice(this.accomm.location, this.accomm.Id, this.reservation.numOfPeople, new Date(this.reservation.startDate),
      new Date(this.reservation.endDate), this.accomm.price).subscribe(data => {
        
      if (data.body.Percentages != null){

        this.reservation.price = data.body.Price

        if (this.accomm.payPer == "per guest"){

          this.reservation.price = this.reservation.price * this.reservation.numOfPeople
        }

        this.finalPrice = this.reservation.price

        this.priceMessage = "The total price will be $" + this.reservation.price.toString()

        for (var i = 0; i < data.body.Percentages.length; i++){
          this.additionalPriceMessage = this.additionalPriceMessage + data.body.Percentages[i].percentage + 
          "% increase for the " + data.body.OverlapDays[i] + " days between " + 
          data.body.Percentages[i].startDate.substring(0, 10) + " and " + data.body.Percentages[i].endDate.substring(0, 10)
        }
      } else {

        this.reservation.price = this.accomm.price

        if (this.accomm.payPer == "per guest"){

          this.reservation.price = this.reservation.price * this.reservation.numOfPeople
        }

        this.reservation.price = this.reservation.price * data.body.Days

        this.priceMessage = "The total price will be $" + this.reservation.price
      }
    })
  }

}
