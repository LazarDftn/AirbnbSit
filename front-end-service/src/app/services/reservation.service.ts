import {Injectable} from '@angular/core';
import {map} from 'rxjs/operators';
import { HttpHeaders } from '@angular/common/http';
import { Accommodation } from '../model/accommodation';
import { ApiService } from './api.service';
import { ConfigService } from './config.service';
import { AuthService } from './auth.service';

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

      getReservationsByAccommId(id: string){
        return this.apiService.get(this.config.reservations_by_accommodation_url + id)
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
    }