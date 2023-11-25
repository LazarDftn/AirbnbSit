import { Component } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { Accommodation } from 'src/app/model/accommodation';
import { AccommodationService } from 'src/app/services/accommodation.service';

@Component({
  selector: 'app-create-accommodation-page',
  templateUrl: './create-accommodation-page.component.html',
  styleUrls: ['./create-accommodation-page.component.css']
})
export class CreateAccommodationPageComponent {

  public accommodation: Accommodation = new Accommodation();
  accommForm!: FormGroup;
  submitted: boolean = false;

  constructor(private accommodationService: AccommodationService,
    private fb: FormBuilder){}

  ngOnInit(): void {
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
      name: ["Name field is requried!", Validators.required],
      location: ["Location field is requried!", Validators.required],
      benefits: [],
      minCapacity: ["Minimum capacity can't be below 1!", Validators.min(1)],
      maxCapacity: []
    });
  }
  
}
