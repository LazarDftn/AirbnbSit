import { Component, OnInit } from '@angular/core';
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
    price = 0

  accomm: Accommodation = new Accommodation()
  reservations: Reservation[] = []
  submitted = false
  
  ngOnInit(): void {

    var accommId = null
    
    this.route.paramMap.subscribe((params: ParamMap) => {
      
      accommId = params.get("id")
    })

    console.log(accommId)

    this.accommService.getAccommById(accommId + "").subscribe(data => {
      this.accomm = data
      this.accomm.Owner = data.owner
      this.accomm.Id = data.id
    })

    this.reservationService.getPriceByAccommId(accommId + "").subscribe(data => {
      this.accomm.price = data.price
      this.accomm.payPer = data.payPer
      this.price = data.price
    })

    this.reservationService.getReservationsByAccommId(accommId + "").subscribe(data => {
      console.log(data)
      this.reservations = data
    })
  }

  addPrice(){
    console.log(this.accomm.price, this.accomm.payPer, this.accomm.Id)
    this.reservationService.createPrice(this.accomm.price, this.accomm.payPer, this.accomm.Id).subscribe
    (data => {window.location.reload()})
  }

}
