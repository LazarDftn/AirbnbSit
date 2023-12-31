import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validator, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { User } from 'src/app/model/user';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.css']
})
export class LoginPageComponent implements OnInit {
  protected loginForm!: FormGroup;
  error = ""

  constructor(private formBuilder: FormBuilder,
    private authService: AuthService,
    private router: Router){}

  ngOnInit() {
    this.loginForm = this.formBuilder.group({
      recaptcha: ['', Validators.required],
      email: ['', Validators.required],
      password: ['', Validators.required]
    });
  }

  onSubmit(loginData: any){

    if (this.loginForm.valid){
      this.error = ""
      this.authService.login(loginData).subscribe(data => {
        var user: User = data.body
        localStorage.setItem("airbnbToken", user.token) //set the token and user data in the localStorage when he logs in
        localStorage.setItem("airbnbEmail", user.email)
        localStorage.setItem("airbnbRole", data.body.user_type)
        this.router.navigate(['welcome-page'])
      }, err => {
        this.error = err.error.error
      })
    } else {
      this.error = "Please fill out all fields and check that you're not a robot!"
    }
  }

    siteKey: string = "6LeMdRcpAAAAACPUVy8p8_rT9z1dkpHqDcVu07AV";

}
