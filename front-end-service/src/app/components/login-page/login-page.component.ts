import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validator, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { ToastrService } from 'ngx-toastr';
import { User } from 'src/app/model/user';
import { UserProfile } from 'src/app/model/userProfile';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.css']
})
export class LoginPageComponent implements OnInit {
  protected loginForm!: FormGroup;

  constructor(private formBuilder: FormBuilder,
    private toastr: ToastrService,
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
      this.authService.login(loginData).subscribe(data => {
        var user: UserProfile = data.body
        localStorage.setItem("airbnbToken", user.token) //set the token and user data in the localStorage when he logs in
        localStorage.setItem("airbnbEmail", user.email)
        localStorage.setItem("airbnbRole", data.body.user_type)
        localStorage.setItem("airbnbUsername", data.body.username)
        this.router.navigate(['welcome-page'])
      }, err => {
        this.toastr.error(err.error.error, "Error");
      })
    } else {
      this.toastr.warning("Please fill out all fields and check that you're not a robot!", "Warning");
    }
  }

    siteKey: string = "6LeMdRcpAAAAACPUVy8p8_rT9z1dkpHqDcVu07AV";

}
