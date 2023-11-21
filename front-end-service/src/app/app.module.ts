import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { LoginPageComponent } from './components/login-page/login-page.component';
import { RegistrationPageComponent } from './components/registration-page/registration-page.component';
import { WelcomePageComponent } from './components/welcome-page/welcome-page.component';
import { NavBarComponent } from './components/nav-bar/nav-bar.component';
import { CreateAccommodationPageComponent } from './components/create-accommodation-page/create-accommodation-page.component';
import { ApiService } from './services/api.service';
import { AccommodationService } from './services/accommodation.service';
import { ConfigService } from './services/config.service';
import { HttpClient, HttpClientModule, HttpHandler } from '@angular/common/http';
import { AuthService } from './services/auth.service';


@NgModule({
  declarations: [
    AppComponent,
    LoginPageComponent,
    RegistrationPageComponent,
    WelcomePageComponent,
    NavBarComponent,
    CreateAccommodationPageComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule
  ],
  providers: [ApiService,
    AccommodationService,
    ConfigService,
    AuthService],
  bootstrap: [AppComponent]
})
export class AppModule { }
