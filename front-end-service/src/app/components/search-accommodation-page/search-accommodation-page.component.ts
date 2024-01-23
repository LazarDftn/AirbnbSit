import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Accommodation } from 'src/app/model/accommodation';
import { Availability } from 'src/app/model/availability';
import { UserProfile } from 'src/app/model/userProfile';
import { AccommodationService } from 'src/app/services/accommodation.service';
import { AuthService } from 'src/app/services/auth.service';
import { ReservationService } from 'src/app/services/reservation.service';

@Component({
  selector: 'search-accommodation-page',
  templateUrl: './search-accommodation-page.component.html',
  styleUrls: ['./search-accommodation-page.component.css']
})
export class SearchAccommodationPageComponent implements OnInit {

  constructor(private authService: AuthService,
    private router: Router,
    private accommService: AccommodationService,
    private reservationService: ReservationService){}

  accommodations: Accommodation[] = []
  availabilities: Availability[] = []
  availability: Availability = new Availability()
  startDate: Date = new Date()
  endDate: Date = new Date()

  ngOnInit(): void {

    this.accommService.getAll().subscribe(data => {
      this.accommodations = data
      console.log(this.accommodations)

    })
  }

  goTo(accomm: Accommodation){
    this.router.navigate(['/accommodation/' + accomm.id])
  }

  goToAvailable(av: Availability){
    window.open('http://localhost:4200/accommodation/' + av.accommId, '_blank')
  }

  search(){
    this.availabilities = []
    this.accommodations = []
    this.availability.availabilityId = "00000000-0000-0000-0000-000000000000"
    this.availability.startDate = new Date(this.startDate)
    this.availability.endDate = new Date(this.endDate)

    this.reservationService.searchAccommodations(this.availability)
    .subscribe(data => {
      if (data.body != null){
        this.availabilities = data.body
      }
    })
  }

}
