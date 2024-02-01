import { Component, OnInit } from '@angular/core';
import { ToastrService } from 'ngx-toastr';
import { User } from 'src/app/model/user';
import { UserProfile } from 'src/app/model/userProfile';
import { AuthService } from 'src/app/services/auth.service';


@Component({
  selector: 'app-profile-page',
  templateUrl: './profile-page.component.html',
  styleUrls: ['./profile-page.component.css']
})
export class ProfilePageComponent implements OnInit{
  
  constructor(private authService: AuthService,
    private toastr: ToastrService,){}

  // booleans for the navbar to check the users role and restrict access to pages
  isUserLoggedIn = this.authService.userIsLoggedIn()
  isUserHost = this.authService.userHasRole("HOST")
  isUserGuest = this.authService.userHasRole("GUEST")

  username = localStorage.getItem("airbnbUsername")
  email = localStorage.getItem("airbnbEmail");

  password = ""
  repeatPassword = ""
  oldPassword = ""

  user: User = new User()
  userToEdit: User = new User()

  ngOnInit(): void {
    
    this.authService.getProfile(localStorage.getItem("airbnbId")!).subscribe(data => {
      this.user = data
      this.user.firstName = data.first_name
      this.user.lastName = data.last_name
    }, err => {
      console.log(err)
    })
  }

  deleteAccount(){

    this.authService.deleteAccount().subscribe(data => {
      console.log(data)
    }, err => {
      alert(err.error.error)
    })
  }

  editProfile(){

    this.userToEdit.ID = this.user.ID

    if (this.user.address == "" || this.user.firstName == "" || this.user.lastName == ""){
      this.toastr.warning("Fields can't be empty!", "Warning")
      return
    }

    if (this.user.email != this.email){
      if (this.user.email == ""){
        this.toastr.warning("Please enter a valid mail!", "Warning")
        return
      }
      this.userToEdit.email = this.user.email
      if (this.oldPassword == ""){
        this.toastr.warning("Enter your current password!", "Warning")
        return
      }
      this.userToEdit.password = this.oldPassword
    } else {
      this.userToEdit.email = ""
    }

    if (this.user.username != this.username){
      if (this.user.username == ""){
        this.toastr.warning("Please enter a valid username!", "Warning")
        return
      }
      this.userToEdit.username = this.user.username
    } else {
      this.userToEdit.username = ""
    }

    this.userToEdit.firstName = this.user.firstName
    this.userToEdit.lastName = this.user.lastName
    this.userToEdit.address = this.user.address

    if (this.password != ""){
      
      if (this.password != this.repeatPassword){
        this.toastr.warning("Passwords don't match!", "Warning")
        return
      }

      if (this.oldPassword == ""){
        this.toastr.warning("Enter old password!", "Warning")
        return
      }

      if (!this.validatePassword(this.password)){
        this.toastr.warning("Please enter a valid new password!", "Warning")
        return
      }
      this.userToEdit.password = this.oldPassword
    } else {
      if (this.repeatPassword != ""){
        this.toastr.warning("Passwords don't match!", "Warning")
        return
      }
    }

    this.authService.editProfile(this.userToEdit, this.password).subscribe(data => {
      if (this.userToEdit.username != ""){
        localStorage.setItem("airbnbUsername", this.userToEdit.username)
      }
      if (this.userToEdit.email != ""){
        localStorage.setItem("airbnbEmail", this.userToEdit.email)
      }
      this.userToEdit = new User()
      window.location.reload()
      this.toastr.success("profile edited", "Success")
    }, err => {
      this.toastr.error(err.error.error, "Error")
      this.userToEdit = new User()
    })

  }

  validatePassword(password: string): boolean {

    if (password.length < 12) {
      return false;
    }

    const lowercaseRegex = new RegExp("(?=.*[a-z])");// has at least one lower case letter
    if (!lowercaseRegex.test(password)) {
      return false;
    }

    const uppercaseRegex = new RegExp("(?=.*[A-Z])"); //has at least one upper case letter
    if (!uppercaseRegex.test(password)) {
      return false;
    }

    const numRegex = new RegExp("(?=.*\\d)"); // has at least one number
    if (!numRegex.test(password)) {
      return false;
    }

    const specialcharRegex = new RegExp("[!@#$%^&*(),.?\":{}|<>]");
    if (!specialcharRegex.test(password)) {
      return false;
    }
    return true
  }

}
