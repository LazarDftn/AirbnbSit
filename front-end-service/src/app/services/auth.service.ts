import {Injectable} from '@angular/core';
import {map} from 'rxjs/operators';
import { HttpHeaders } from '@angular/common/http';
import { User } from '../model/user';
import { ApiService } from './api.service';
import { ConfigService } from './config.service';
import { Router } from '@angular/router';

@Injectable({
  providedIn: 'root'
})

export class AuthService {

    constructor(
        private apiService: ApiService,
        private config: ConfigService,
        private router: Router
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
        refresh_token: "",
        is_verified: false
      }
      
      const postHeaders = new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      });
      return this.apiService.post(this.config.signup_url, JSON.stringify(userDTO), postHeaders)
      .pipe(map(() => {
      }));
    } 

    verifyEmail(username: any, code: any){

      var userDTO = {
        verifUsername: username,
        code: code
      }
      
      const postHeaders = new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      });
      return this.apiService.post(this.config.verify_email_url, JSON.stringify(userDTO), postHeaders)
      .pipe(map(() => {
      }));
    }

    verifyEmailForPassword(email: string){

      var emailDTO = {
        email: email
      }

      const postHeaders = new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      });
      return this.apiService.post(this.config.password_code_url, JSON.stringify(emailDTO), postHeaders)
      .pipe(map(() => {
      }));
    }

    forgotPassword(email: string, code: string, password: string){

      var forgotPasswordDTO = {
        email: email,
        code: code,
        password: password
      }

      const postHeaders = new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      });
      return this.apiService.post(this.config.forgot_password_url, JSON.stringify(forgotPasswordDTO), postHeaders)
      .pipe(map(() => {
      }));
    }

    login(data: any){

      var loginUserDTO = {
        recaptcha: data.recaptcha,
        email: data.email,
        password: data.password
      }

      const postHeaders = new HttpHeaders({
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      });

      return this.apiService.post(this.config.login_url, JSON.stringify(loginUserDTO), postHeaders)
      .pipe(map((data) => {
        return data
      }));
    }
    
    logout(){
      localStorage.removeItem("airbnbToken")
      localStorage.removeItem("airbnbUsername")
      localStorage.removeItem("airbnbRole")

      this.router.navigate(['/login-page'])
    }

    /* details about the Logged user are kept in the local storage.
       This method should be implemented on pages that require any Role */
    userIsLoggedIn(): boolean{
      var role = localStorage.getItem("airbnbRole")
      if (role == "" || role == null || role == undefined){
        return false
      }
      return true
    }

    // This method should be implemented on pages that require a certain Role   
    userHasRole(requiredRole: string): boolean{
      var role = localStorage.getItem("airbnbRole")
      if (role == requiredRole){
        return true
      }
      return false
    }
  }
