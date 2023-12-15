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

export class AccommodationService {

    constructor(
        private apiService: ApiService,
        private config: ConfigService,
        private authService: AuthService
      ) {
      }

      create(accommodation: Accommodation){
        
        var accommDTO = {
            Owner: localStorage.getItem("airbnbEmail"),
            Name: accommodation.name,
            Location: accommodation.location,
            Benefits: accommodation.benefits,
            MinCapacity: accommodation.minCapacity,
            MaxCapacity: accommodation.maxCapacity
        }

        const postHeaders = new HttpHeaders({
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'token': localStorage.getItem("airbnbToken") + '' // for now send auth tokens like this, intercept all requests later
          });
          return this.apiService.post(this.config.accommodations_url, JSON.stringify(accommDTO), postHeaders)
            .pipe(map((data) => {
              return data
            }));
      }

      getAccommById(id: string){
        
        return this.apiService.get(this.config.accommodation_url + id)
        .pipe(map((data) => {
          return data
        }));
      }

      getAll(){

        return this.apiService.get(this.config.accommodation_url)
        .pipe(map((data) => {
          return data
        }));
      }
    }
