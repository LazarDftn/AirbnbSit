import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private _api_url = 'http://localhost:8000';
  private _signup_url = this._api_url + '/users/signup';
  private _accommodations_url = this._api_url + '/accommodations/create';
  private _accommodation_url = this._api_url + '/accommodations/'
  private _verify_email_url = this._api_url + '/users/verify-account';
  private _password_code_url = this._api_url + '/users/password-code';
  private _forgot_password_url = this._api_url + '/users/forgot-password';
  private _login_url = this._api_url + '/users/login';
  private _accommodation_price_url = this._api_url + '/reservation/accommodation-price/'
  private _reservations_by_accommodation_url = this._api_url + '/reservation/accommodation/'
  private _check_reservation_price_url = this._api_url + "/reservation/check-price/"
  private _create_accommodation_price_url = this._api_url + "/reservation/host/accommodation-price/"
  private _price_variation_url = this._api_url + "/reservation/accommodation/price-variation/"
  private _availability_url = this._api_url + "/reservation/availability/"
  private _create_price_variation_url = this._api_url + "/reservation/host/accommodation/price-variation/"
  private _create_availability_url = this._api_url + "/reservation/host/availability/create"
  private _delete_availability_url = this._api_url + "/reservation/host/availability/delete"
  private _delete_price_variation_url = this._api_url + "/reservation/host/accommodation/price-variation/delete"
  private _reservation_url = this._api_url + "/reservation/guest/"
  private _search_accomm_url = this._api_url + "/reservation/search/"
  private _delete_account_url = this._api_url + "/users/delete-account"
  private _cancel_reservation_url = this._api_url + "/reservation/guest/delete/"

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

  get login_url(): string {
    return this._login_url
  }

  get accommodation_url(): string {
    return this._accommodation_url
  }

  get accommodation_price_url(): string {
    return this._accommodation_price_url
  }

  get reservations_by_accommodation_url(): string {
    return this._reservations_by_accommodation_url
  }

  get check_reservation_price_url(): string {
    return this._check_reservation_price_url
  }

  get price_variation_url(): string {
    return this._price_variation_url
  }

  get availabillity_url(): string {
    return this._availability_url
  }

  get delete_availability_url(): string {
    return this._delete_availability_url
  }

  get delete_price_variation_url(): string {
    return this._delete_price_variation_url
  }

  get reservation_url(): string {
    return this._reservation_url
  }

  get create_accommodation_price(): string {
    return this._create_accommodation_price_url
  }

  get create_price_variation_url(): string {
    return this._create_price_variation_url
  }

  get create_availability_url(): string {
    return this._create_availability_url
  }

  get search_accomm_url(): string {
    return this._search_accomm_url
  }

  get delete_account_url(): string {
    return this._delete_account_url
  }

  get cancel_reservation_url(): string {
    return this._cancel_reservation_url
  }
}