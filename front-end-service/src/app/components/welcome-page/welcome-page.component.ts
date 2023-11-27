import { Component, Input } from '@angular/core';



@Component({
  selector: 'welcome-page',
  templateUrl: './welcome-page.component.html',
  styleUrls: ['./welcome-page.component.css']
})
export class WelcomePageComponent {
  user = localStorage.getItem("airbnbUsername")
}
let fruits: string[] = ['Apple', 'Orange', 'Banana']; 
if (fruits.includes('Apples')) {
  console.log('✅ The value is contained in array');
} else {
  console.log('⛔️ The value is NOT contained in array');
}
// console.log(fruits);

