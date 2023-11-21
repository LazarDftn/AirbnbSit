import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validator, Validators } from '@angular/forms';

@Component({
  selector: 'login-page',
  templateUrl: './login-page.component.html',
  styleUrls: ['./login-page.component.css']
})
export class LoginPageComponent implements OnInit {
  protected aFormGroup!: FormGroup;

  constructor(private formBuilder: FormBuilder){}

  ngOnInit() {
    this.aFormGroup = this.formBuilder.group({
      recaptcha: ['', Validators.required]
    });
  }



    siteKey: string = "6LeMdRcpAAAAACPUVy8p8_rT9z1dkpHqDcVu07AV";

}
