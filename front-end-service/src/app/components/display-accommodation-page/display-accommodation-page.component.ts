import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'display-accommodation-page',
  templateUrl: './display-accommodation-page.component.html',
  styleUrls: ['./display-accommodation-page.component.css']
})

export class DisplayAccommodationPageComponent{
  name = localStorage.getItem("name")//nznm da li ovo treba da koristi localStorage ili nesto drugo
  location = localStorage.getItem("location")
  owner=localStorage.getItem("owner")
  benefits=localStorage.getItem("benefits")
  minCapacity=localStorage.getItem("minCapacity")
  maxCapacity=localStorage.getItem("maxCapacity")
  price=localStorage.getItem("price")
  discPrice=localStorage.getItem("discPrice")
  discDateStart=localStorage.getItem("discDateStart")
  discDateEnd=localStorage.getItem("discDateEnd")
  discWeekend=localStorage.getItem("discWeekend")
  payPer=localStorage.getItem("payPer")
}
