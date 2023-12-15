import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { Accommodation } from 'src/app/model/accommodation';
import { AccommodationService } from 'src/app/services/accommodation.service';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'search-accommodation-page',
  templateUrl: './search-accommodation-page.component.html',
  styleUrls: ['./search-accommodation-page.component.css']
})
export class SearchAccommodationPageComponent implements OnInit {

  constructor(private authService: AuthService,
    private router: Router,
    private accommService: AccommodationService){}

  accommodations: Accommodation[] = []

  ngOnInit(): void {

    this.accommService.getAll().subscribe(data => {
      this.accommodations = data
      console.log(this.accommodations)

    })
  }

  goTo(accomm: Accommodation){
    this.router.navigate(['/accommodation/' + accomm.id])
  }

}
