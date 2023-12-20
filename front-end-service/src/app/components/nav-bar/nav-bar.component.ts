import { Component } from '@angular/core';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'nav-bar',
  templateUrl: './nav-bar.component.html',
  styleUrls: ['./nav-bar.component.css']
})
export class NavBarComponent {

  constructor(private authService: AuthService){}

  // booleans for the navbar to check the users role and restrict access to pages
  isUserLoggedIn = this.authService.userIsLoggedIn()
  isUserHost = this.authService.userHasRole("HOST")
  isUserGuest = this.authService.userHasRole("GUEST")

  logout(){
    this.authService.logout()
  }

  email = localStorage.getItem("airbnbEmail")

}
