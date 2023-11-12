import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class ConfigService {

  private _api_url = 'http://localhost:8080';
  private _accommodations_url = this._api_url + '/accommodations';

  get accommodations_url(): string {
    return this._accommodations_url;
  }
}