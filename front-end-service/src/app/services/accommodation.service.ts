import {Injectable} from '@angular/core';
import {map} from 'rxjs/operators';
import { HttpHeaders } from '@angular/common/http';
import { Accommodation } from '../model/accommodation';
import { ApiService } from './api.service';
import { ConfigService } from './config.service';

@Injectable({
  providedIn: 'root'
})

export class AccommodationService {

    constructor(
        private apiService: ApiService,
        private config: ConfigService
      ) {
      }

      create(accommodation: Accommodation){

        var accommDTO = {
            Owner: "Pera", //change to current user later
            Name: accommodation.name,
            Location: accommodation.location,
            Benefits: accommodation.benefits,
            MinCapacity: accommodation.minCapacity,
            MaxCapacity: accommodation.maxCapacity,
            Price: accommodation.price,
            DiscPrice: accommodation.discPrice,
            DiscPriceStart: accommodation.discDateStart,
            DiscPriceEnd: accommodation.discDateEnd,
            DiscWeekend: accommodation.discWeekend,
            PayPer: accommodation.payPer
        }

        const postHeaders = new HttpHeaders({
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          });
          return this.apiService.post(this.config.accommodations_url, JSON.stringify(accommDTO), postHeaders)
            .pipe(map(() => {
            }));
      }
    }
