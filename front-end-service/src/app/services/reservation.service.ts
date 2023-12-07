import {Injectable} from '@angular/core';
import {map} from 'rxjs/operators';
import { HttpHeaders } from '@angular/common/http';
import { Accommodation } from '../model/accommodation';
import { ApiService } from './api.service';
import { ConfigService } from './config.service';
import { AuthService } from './auth.service';
import { formatDate } from '@angular/common';

@Injectable({
  providedIn: 'root'
})

export class ReservationService {

    constructor(
        private apiService: ApiService,
        private config: ConfigService
      ) {
      }

      getPriceByAccommId(id: string){
        return this.apiService.get(this.config.accommodation_price_url + id)
        .pipe(map((data) => {
          return data
        }));
      }

      getReservationsByAccommId(location: string, id: string){
        return this.apiService.get(this.config.reservations_by_accommodation_url + location + "/" + id)
        .pipe(map((data) => {
          return data
        }));
      }

      createPrice(price: number, payPer: string, accommId: string){

        var priceDTO = {
          accommId: accommId,
          price: price,
          payPer: payPer
        }

        const postHeaders = new HttpHeaders({
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'token': localStorage.getItem("airbnbToken") + ''
        });

        return this.apiService.post(this.config.accommodation_price_url, priceDTO, postHeaders)
        .pipe(map((data) => {
          return data
        }));
      }

      checkPrice(location: string, accommId: string, numOfPeople: number, startDate: Date, endDate: Date, price: number){

        var resDTO = {
          location: location,
          accommId: accommId,
          guestEmail: "",
          hostEmail: "",
          price: price,
          numOfPeople: numOfPeople,
          startDate: new Date(Date.UTC(startDate.getFullYear(), startDate.getMonth(), startDate.getDate(),
          startDate.getHours() - 2, startDate.getMinutes(), startDate.getMinutes())),
          endDate: new Date(Date.UTC(endDate.getFullYear(), endDate.getMonth(), endDate.getDate(),
          endDate.getHours() - 2, endDate.getMinutes(), endDate.getMinutes()))
        }

        const postHeaders = new HttpHeaders({
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'token': localStorage.getItem("airbnbToken") + ''
        });

        return this.apiService.post(this.config.check_reservation_price_url, resDTO, postHeaders)
        .pipe(map((data) => {
          return data
        }));
      }
    }