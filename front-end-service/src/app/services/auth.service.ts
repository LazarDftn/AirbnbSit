import {Injectable} from '@angular/core';
import {map} from 'rxjs/operators';
import { HttpHeaders } from '@angular/common/http';
import { User } from '../model/user';
import { ApiService } from './api.service';
import { ConfigService } from './config.service';

@Injectable({
  providedIn: 'root'
})

export class AuthService {

    constructor(
        private apiService: ApiService,
        private config: ConfigService
      ) {
      }
    
    signup(user: User){
      
      const postHeaders = new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      });
      return this.apiService.post(this.config.signup_url, JSON.stringify(user), postHeaders)
      .pipe(map(() => {
      }));
    } 
  }
