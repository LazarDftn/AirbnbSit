import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

import { ToastrModule } from 'ngx-toastr';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

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
import { AccountVerifPageComponent } from './components/account-verif-page/account-verif-page.component';

import { NgxCaptchaModule } from 'ngx-captcha';
import { ForgotPasswordPageComponent } from './components/forgot-password-page/forgot-password-page.component';
import { SearchAccommodationPageComponent } from './components/search-accommodation-page/search-accommodation-page.component';
import { ViewAccommodationPageComponent } from './components/view-accommodation-page/view-accommodation-page.component';
import { EditAccommodationPageComponent } from './components/edit-accommodation-page/edit-accommodation-page.component';
import { ProfilePageComponent } from './components/profile-page/profile-page.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginPageComponent,
    RegistrationPageComponent,
    WelcomePageComponent,
    NavBarComponent,
    CreateAccommodationPageComponent,
    AccountVerifPageComponent,
    ForgotPasswordPageComponent,
    SearchAccommodationPageComponent,
    ViewAccommodationPageComponent,
    EditAccommodationPageComponent,
    ProfilePageComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    FormsModule,
    ReactiveFormsModule,
    HttpClientModule,
    NgxCaptchaModule,
    ToastrModule.forRoot({
      positionClass:"toast-top-center",
      preventDuplicates: true,
      closeButton: true
    }),
    BrowserAnimationsModule
  ],
  providers: [ApiService,
    AccommodationService,
    ConfigService,
    AuthService],
  bootstrap: [AppComponent]
})
export class AppModule { }
