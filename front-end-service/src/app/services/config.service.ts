import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private _api_url = 'http://localhost:8000';
  private _signup_url = this._api_url + '/users/signup';
  private _accommodations_url = this._api_url + '/accommodations/create';
  private _verify_email_url = this._api_url + '/users/verify-account';
  private _password_code_url = this._api_url + '/users/password-code';
  private _forgot_password_url = this._api_url + '/users/forgot-password';

  get accommodations_url(): string {
    return this._accommodations_url;
  }

  get signup_url(): string {
    return this._signup_url;
  }

  get verify_email_url(): string {
    return this._verify_email_url
  }

  get password_code_url(): string {
    return this._password_code_url
  }

  get forgot_password_url(): string {
    return this._forgot_password_url
  }
}