import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Accommodation } from 'src/app/model/accommodation';
import { AccommodationService } from 'src/app/services/accommodation.service';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-create-accommodation-page',
  templateUrl: './create-accommodation-page.component.html',
  styleUrls: ['./create-accommodation-page.component.css']
})
export class CreateAccommodationPageComponent implements OnInit {

  public accommodation: Accommodation = new Accommodation();
  accommForm!: FormGroup;
  submitted: boolean = false;

  constructor(private accommodationService: AccommodationService,
    private authService: AuthService,
    private fb: FormBuilder,
    private router: Router){}

  ngOnInit(): void {

    // auth service checks if the user is logged in and has permission to view this page
    if (!this.authService.userIsLoggedIn()){
      this.router.navigate(['/login-page'])
    }
    if (!this.authService.userHasRole("HOST")){
      this.router.navigate(['/welcome-page'])
    }
    this.buildForm();
  }


  public onSubmit(){
    this.submitted = true
    if (this.accommForm.valid && this.accommodation.maxCapacity >= this.accommodation.minCapacity) {
    this.accommodationService.create(this.accommodation).subscribe(data => {window.location.reload()});
    } else {
      console.error("Form is invalid!");
    }
  }

  private buildForm() {
    this.accommForm = this.fb.group({
      name: ["Name field is required!", Validators.required],
      location: ["Location field is required!", Validators.required],
      benefits: [],
      minCapacity: ["Minimum capacity can't be below 1!", Validators.min(1)],
      maxCapacity: [],
      price: ["Price can't be negative!", Validators.min(0)],
      discPrice: ["Price can't be negative!", Validators.min(0)],
      discDateStart: [],
      discDateEnd: ["Start date can't be after end date!", Validators.min(this.accommodation.discDateEnd.getTime())],
      discWeekend:[],
      payPer: []
    });
  }
  
}
