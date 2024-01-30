import { Component } from '@angular/core';
import { User } from 'src/app/model/user';
import { UserProfile } from 'src/app/model/userProfile';
import { AuthService } from 'src/app/services/auth.service';


@Component({
  selector: 'app-profile-page',
  templateUrl: './profile-page.component.html',
  styleUrls: ['./profile-page.component.css']
})
export class ProfilePageComponent {
  
  constructor(private authService: AuthService){}

  // booleans for the navbar to check the users role and restrict access to pages
  isUserLoggedIn = this.authService.userIsLoggedIn()
  isUserHost = this.authService.userHasRole("HOST")
  isUserGuest = this.authService.userHasRole("GUEST")

  username = localStorage.getItem("airbnbUsername")
  email = localStorage.getItem("airbnbEmail");


  
  deleteAccount(){

    this.authService.deleteAccount().subscribe(data => {
      console.log(data)
    }, err => {
      alert(err.error.error)
    })
  }

}
