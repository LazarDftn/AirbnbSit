import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CommonModule } from '@angular/common';
import { LoginPageComponent } from './components/login-page/login-page.component';
import { AppComponent } from './app.component';
import { WelcomePageComponent } from './components/welcome-page/welcome-page.component';
import { RegistrationPageComponent } from './components/registration-page/registration-page.component';
import { CreateAccommodationPageComponent } from './components/create-accommodation-page/create-accommodation-page.component';
import { AccountVerifPageComponent } from './components/account-verif-page/account-verif-page.component';
import { ForgotPasswordPageComponent } from './components/forgot-password-page/forgot-password-page.component';
import { SearchAccommodationPageComponent } from './components/search-accommodation-page/search-accommodation-page.component';
import { ViewAccommodationPageComponent } from './components/view-accommodation-page/view-accommodation-page.component';
import { EditAccommodationPageComponent } from './components/edit-accommodation-page/edit-accommodation-page.component';
import { ProfilePageComponent } from './components/profile-page/profile-page.component';

const routes: Routes = [
  {path: '', pathMatch:'full', component: WelcomePageComponent},
  {path: 'login-page', component: LoginPageComponent},
  {path: 'registration-page', component: RegistrationPageComponent},
  {path: 'create-accommodation', component: CreateAccommodationPageComponent},
  {path: 'account/:id/:id2', component: AccountVerifPageComponent},
  {path: 'forgot-password-page', component: ForgotPasswordPageComponent},
  {path: 'welcome-page', component: WelcomePageComponent},
  {path: 'profile-page', component: ProfilePageComponent},
  {path: 'accommodation/:id', component: ViewAccommodationPageComponent},
  {path: 'accommodation/edit/:id', component: EditAccommodationPageComponent}]
;

@NgModule({
  imports: [CommonModule, RouterModule.forRoot(routes,{scrollPositionRestoration: 'enabled'})],
  exports: [RouterModule]
})
export class AppRoutingModule { }
