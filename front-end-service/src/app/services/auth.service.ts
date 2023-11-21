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

      var userDTO = {
        first_name: user.fName,
        last_name: user.lName,
        username: user.username,
        password: user.password,
        email: user.email,
        address: user.address,
        token: "",
        user_type: user.type,
        refresh_token:""
      }
      
      console.log(user)
      const postHeaders = new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      });
      return this.apiService.post(this.config.signup_url, JSON.stringify(userDTO), postHeaders)
      .pipe(map(() => {
      }));
    } 
  }
