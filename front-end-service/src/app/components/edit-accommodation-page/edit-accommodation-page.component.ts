import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup } from '@angular/forms';
import { ActivatedRoute, ParamMap, Route, Router } from '@angular/router';
import { Accommodation } from 'src/app/model/accommodation';
import { Availability } from 'src/app/model/availability';
import { PriceVariation } from 'src/app/model/priceVariation';
import { AccommodationService } from 'src/app/services/accommodation.service';
import { AuthService } from 'src/app/services/auth.service';
import { ReservationService } from 'src/app/services/reservation.service';

@Component({
  selector: 'app-edit-accommodation-page',
  templateUrl: './edit-accommodation-page.component.html',
  styleUrls: ['./edit-accommodation-page.component.css']
})
export class EditAccommodationPageComponent implements OnInit{

  accomm: Accommodation = new Accommodation();
  pv: PriceVariation = new PriceVariation();
  av: Availability = new Availability();
  submitError = ""
  accommForm!: FormGroup;

  constructor(private accommService: AccommodationService,
    private router: Router,
    private route: ActivatedRoute,
    private reservationService: ReservationService,
    private authService: AuthService,
    private fb: FormBuilder){}

  ngOnInit(): void {

    this.buildForm()
    
    this.route.paramMap.subscribe((params: ParamMap) => {
      
      this.accomm.id = params.get("id")!

    })

    this.accommService.getAccommById(this.accomm.id + "").subscribe(data => {

      this.accomm = data
      this.accomm.Owner = data.owner
      console.log(this.accomm.Owner)

      if (!this.authService.userHasEmail(this.accomm.Owner)){
        this.router.navigate(['welcome-page'])
      }

    })
  }

  addAvailability(){

    this.submitError = ""

    this.av.accommId = this.accomm.id
    this.av.availabilityId = "00000000-0000-0000-0000-000000000000"
    this.av.location = this.accomm.location
    this.av.name = this.accomm.name
    this.av.minCapacity = this.av.minCapacity
    this.av.maxCapacity = this.av.maxCapacity
    this.av.startDate = new Date(this.av.startDate)
    this.av.endDate = new Date(this.av.endDate)

    this.reservationService.createAvailability(this.av).subscribe(data => {

      if (data.body != "Changed"){
        this.submitError = data.body
      } else {
        this.router.navigate(['accommodation/' + this.accomm.id])
      }
    })
  }

  addPriceIncrease(){

    this.submitError = ""

    this.pv.variationId = "00000000-0000-0000-0000-000000000000"
    this.pv.accommId = this.accomm.id
    this.pv.location = this.accomm.location
    this.pv.startDate = new Date(this.pv.startDate)
    this.pv.endDate = new Date(this.pv.endDate)

    this.reservationService.createPriceVariation(this.pv).subscribe(data => {
      if (data.body != "Created"){
        this.submitError = data.body
      } else {
        this.router.navigate(['accommodation/' + this.accomm.id])
      }
    })
  }

  private buildForm() {
    this.accommForm = this.fb.group({
      avStartDate: [""],
      avEndDate: [""],
      pvStartDate: [""],
      pvEndDate: [""],
      percentage: [""]
    });
  }


}
