import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'search-accommodation-page',
  templateUrl: './search-accommodation-page.component.html',
  styleUrls: ['./search-accommodation-page.component.css']
})
export class SearchAccommodationPageComponent implements OnInit {

  constructor(private authService: AuthService,
    private router: Router){}

  ngOnInit(): void {
  }

}
