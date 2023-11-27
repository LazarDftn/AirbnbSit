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
    // auth service checks if the user is logged in and has permission to view this page
    if (!this.authService.userIsLoggedIn()){
      this.router.navigate(['/login-page'])
    }
    if (!this.authService.userHasRole("GUEST")){
      this.router.navigate(['/welcome-page'])
    }
  }

}
