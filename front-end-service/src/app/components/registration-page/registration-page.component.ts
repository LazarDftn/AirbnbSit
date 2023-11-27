import { Attribute, Component } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validator, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { User } from 'src/app/model/user';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'registration-page',
  templateUrl: './registration-page.component.html',
  styleUrls: ['./registration-page.component.css']
})
export class RegistrationPageComponent {

  public user: User = new User();
  signupError: string = ""
  userForm!: FormGroup;
  submitted: boolean = false;

  constructor(private authService: AuthService,
    private fb: FormBuilder,
    private router: Router){}

  ngOnInit(): void {
    this.buildForm();
  }

  public onSubmit() {

    this.submitted = true

    if (this.userForm.valid && this.validatePassword(this.user.password)
    && this.validateEmail(this.user.email)) {

      this.authService.signup(this.user).subscribe(data => {
        alert("Success! Please go confirm your Email")
        this.router.navigate(['/login-page'])
      }
        , error => {
          this.signupError = error.error.error;
        });
    }
  }

  private buildForm() {
    this.userForm = this.fb.group({
      fname: ["", Validators.pattern('^[a-zA-Z ]*$')],
      lname: ["", Validators.pattern('^[a-zA-Z ]*$')],
      email: ["", Validators.required],
      address: ["", Validators.required],
      username: ["", Validators.required],
      password: ["", Validators.required],
      role : ["", Validators.required]
    });
  }

  validateEmail(email: string): boolean {
    if (!email.includes("@") || !email.includes(".")) {
      return false
    }
    return true
  }

  validatePassword(password: string): boolean {

    if (password.length < 12) {
      return false;
    }

    const lowercaseRegex = new RegExp("(?=.*[a-z])");// has at least one lower case letter
    if (!lowercaseRegex.test(password)) {
      return false;
    }

    const uppercaseRegex = new RegExp("(?=.*[A-Z])"); //has at least one upper case letter
    if (!uppercaseRegex.test(password)) {
      return false;
    }

    const numRegex = new RegExp("(?=.*\\d)"); // has at least one number
    if (!numRegex.test(password)) {
      return false;
    }

    const specialcharRegex = new RegExp("[!@#$%^&*(),.?\":{}|<>]");
    if (!specialcharRegex.test(password)) {
      return false;
    }
    return true
  }
}


