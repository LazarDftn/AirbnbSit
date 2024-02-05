import { AfterViewInit, Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap, Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { Accommodation } from 'src/app/model/accommodation';
import { Availability } from 'src/app/model/availability';
import { PriceVariation } from 'src/app/model/priceVariation';
import { Reservation } from 'src/app/model/reservation';
import { AccommodationService } from 'src/app/services/accommodation.service';
import { AuthService } from 'src/app/services/auth.service';
import { ReservationService } from 'src/app/services/reservation.service';

@Component({
  selector: 'view-accommodation-page',
  templateUrl: './view-accommodation-page.component.html',
  styleUrls: ['./view-accommodation-page.component.css']
})
export class ViewAccommodationPageComponent implements OnInit {

  constructor(public authService: AuthService,
    private accommService: AccommodationService,
    private reservationService: ReservationService,
    private router: Router,
    private route: ActivatedRoute,
    private toastr: ToastrService){}
    
  accommId = ""
  price = 0 //accommodation price per day
  finalPrice = 0 //after calculating the days (and guests if that's the payment method)
  accomm: Accommodation = new Accommodation()
  reservation: Reservation = new Reservation()
  reservations: Reservation[] = []
  availabilities: Availability[] = []
  priceVariations: PriceVariation[] = []
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
      this.accomm.ownerId = data.ownerId
      this.accomm.id = data.id
      this.accomm.minCapacity = data.minCapacity
      this.accomm.maxCapacity = data.maxCapacity
      this.accomm.price = this.price

      this.reservationService.getReservationsByAccommId(this.accomm.location, this.accommId + "").subscribe(data => {
        this.reservations = data
        console.log(this.reservations)
      })

      this.reservationService.getAvailabilities(this.accomm.location, this.accommId + "").subscribe(data => {
        this.availabilities = data
      })

      this.reservationService.getPriceVariations(this.accomm.location, this.accommId + "").subscribe(data => {
        this.priceVariations = data
      })
    })

    this.reservationService.getPriceByAccommId(this.accommId + "").subscribe(data => {
      this.accomm.price = data.price
      this.accomm.payPer = data.payPer
      this.price = data.price
      this.loadedPrice = true
    })

  }

  addPrice(){
    this.reservationService.createPrice(this.accomm.price, this.accomm.payPer, this.accomm.id).subscribe
    (data => {reloadTimeOut();
      this.toastr.success("Successfully added price!", "Success!")
    })
  }

  // before Guest makes the reservation, check service for any price variations for given period and calculate 
  checkPrice(){

    this.priceMessage = ""
    this.additionalPriceMessage = ""

    if (!(this.reservation.numOfPeople >= this.accomm.minCapacity) ||
    !(this.reservation.numOfPeople <= this.accomm.maxCapacity)){

      this.toastr.warning("Number of people is above or below capacity!", "Warning!")
      return
    }

    this.submitted = true

    this.reservationService.checkPrice(this.accomm.location, this.accomm.id, this.reservation.numOfPeople, new Date(this.reservation.startDate),
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

  deleteAvailability(av: Availability){

    this.reservationService.deleteAvailability(av).subscribe(data => {
        reloadTimeOut();
        this.toastr.success("Successfully deleted this Availability!", "Success");
    }, err => {
      this.toastr.warning("Can't delete availability because there are reservations during this period", "Warning!")
    })

  }

  deletePriceVariation(pv: PriceVariation){

    this.reservationService.deletePriceVariation(pv).subscribe(data => {
        reloadTimeOut();
        this.toastr.success("Successfully deleted this Price Variation!", "Success");
    }, err => {
      this.toastr.warning("Can't delete price variation because there are reservations during this period", "Warning!")
    })
  }

  makeReservation(){
    
    this.reservation.guestId = localStorage.getItem("airbnbId")!
    this.reservation.hostId = this.accomm.ownerId
    this.reservation.guestEmail = localStorage.getItem("airbnbEmail")!
    this.reservation.hostEmail = this.accomm.Owner
    this.reservation.accommodationId = this.accomm.id
    this.reservation.location = this.accomm.location
    this.reservation.reservationId = "00000000-0000-0000-0000-000000000000"
    this.reservation.startDate = new Date(this.reservation.startDate)
    this.reservation.endDate = new Date(this.reservation.endDate)

    this.reservationService.createReservation(this.reservation).subscribe(data => {
      if (data.body != null){
        this.toastr.warning("Home is not available during this time", "Warning!")
      } else {
        reloadTimeOut();
        this.toastr.success("Successfuly Booked the Accommadation!", "Success!")
      }
    })
  }

  edit(id: string){
    this.router.navigate(['accommodation/edit/' + id])
  }

  cancelReservation(res: Reservation){

    res.accommodationId = this.accomm.id
    
    this.reservationService.cancelReservation(res).subscribe(data => {
      reloadTimeOut();
      this.toastr.success("Successfully canceled your reservation!", "Success");
      
    }, err => {
      this.toastr.warning("This reservation has already started!", "Warning");
    })
  }
  
  
  
}

//reload page after 3 seconds
function reloadTimeOut(){
  setTimeout(() => window.location.reload(), 3000)
}