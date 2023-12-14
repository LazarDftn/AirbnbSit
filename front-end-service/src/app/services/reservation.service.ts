import {Injectable} from '@angular/core';
import {map} from 'rxjs/operators';
import { HttpHeaders } from '@angular/common/http';
import { Accommodation } from '../model/accommodation';
import { ApiService } from './api.service';
import { ConfigService } from './config.service';
import { AuthService } from './auth.service';
import { formatDate } from '@angular/common';
import { Availability } from '../model/availability';
import { PriceVariation } from '../model/priceVariation';
import { Reservation } from '../model/reservation';

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

      getPriceVariations(location: string, id: string){
        return this.apiService.get(this.config.price_variation_url + location + "/" + id)
        .pipe(map((data) => {
          return data
        }));
      }

      getAvailabilities(location: string, id: string){
        return this.apiService.get(this.config.availabillity_url + location + "/" + id)
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

      deleteAvailability(av: Availability){


        const postHeaders = new HttpHeaders({
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'token': localStorage.getItem("airbnbToken") + ''
        });

        return this.apiService.post(this.config.delete_availability_url, av, postHeaders)
        .pipe(map((data) => {
          return data
        }));
      }

      deletePriceVariation(pv: PriceVariation){


        const postHeaders = new HttpHeaders({
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'token': localStorage.getItem("airbnbToken") + ''
        });

        return this.apiService.post(this.config.delete_price_variation_url, pv, postHeaders)
        .pipe(map((data) => {
          return data
        }));
      }

      createReservation(res: Reservation){

        var resDTO = {
          accommId: res.accommodationId,
          location: res.location,
          guestEmail: res.guestEmail,
          hostEmail: res.hostEmail,
          price: res.price,
          numOfPeople: res.numOfPeople,
          startDate: new Date(Date.UTC(res.startDate.getFullYear(), res.startDate.getMonth(), res.startDate.getDate(),
          res.startDate.getHours() - 2, res.startDate.getMinutes(), res.startDate.getMinutes())),
          endDate: new Date(Date.UTC(res.endDate.getFullYear(), res.endDate.getMonth(), res.endDate.getDate(),
          res.endDate.getHours() - 2, res.endDate.getMinutes(), res.endDate.getMinutes()))
        }
        
        const postHeaders = new HttpHeaders({
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'token': localStorage.getItem("airbnbToken") + ''
        });

        return this.apiService.post(this.config.reservation_url, resDTO, postHeaders)
        .pipe(map((data) => {
          return data
        }));

      }
    }