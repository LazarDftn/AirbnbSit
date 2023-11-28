import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-forgot-password-page',
  templateUrl: './forgot-password-page.component.html',
  styleUrls: ['./forgot-password-page.component.css']
})
export class ForgotPasswordPageComponent {
  isEmailDisabled = false
  isPassDisabled = true
  submitted = false
  email = ""
  code = ""
  password = ""
  repPassword = ""
  error = ""

  constructor(private authService: AuthService,
    private router: Router){}

  sendEmail(){
    this.authService.verifyEmailForPassword(this.email).subscribe(data => {
      this.isEmailDisabled = true
      this.isPassDisabled = false
      alert("Check email for code! You only have 60 seconds")
    }, err => {
      this.error = err.error.error
    })
  }

  sendPassword(){
    this.submitted = true

    if (this.validatePassword(this.password) && this.comparePasswords(this.password, this.repPassword)){
    this.authService.forgotPassword(this.email, this.code, this.password).subscribe(data => {
    alert("Password successfully changed!")
    this.router.navigate(['/login-page'])
    }, err => {
      this.isEmailDisabled = false
      this.isPassDisabled = true
      this.code = ""
      this.error = err.error.error + " Send another code to your email"
    })
  }
  }

  validatePassword(password: string): boolean {

    if (password.length < 12){
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

  comparePasswords(password: string, repPassword: string): boolean {
    if (password != repPassword){
      return false
    }
    return true
  }

}
