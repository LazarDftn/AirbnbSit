import { Component } from '@angular/core';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-forgot-password-page',
  templateUrl: './forgot-password-page.component.html',
  styleUrls: ['./forgot-password-page.component.css']
})
export class ForgotPasswordPageComponent {
  isEmailDisabled = false
  isPassDisabled = true
  email = ""
  code = ""
  password = ""
  repPassword = ""
  error = ""

  constructor(private authService: AuthService){}

  sendEmail(){
    this.authService.verifyEmailForPassword(this.email).subscribe(data => {
      this.isEmailDisabled = true
      this.isPassDisabled = false
    }, err => {
      this.error = err.error.error
    })
  }

  sendPassword(){
    this.authService.forgotPassword(this.email, this.code, this.password).subscribe(data => {
    }, err => {
      this.isEmailDisabled = false
      this.isPassDisabled = true
      this.code = ""
      this.error = err.error.error + " Send another code to your email"
    })
  }
}
