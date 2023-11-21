import { Attribute, Component } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, Validator, Validators } from '@angular/forms';
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
    private fb: FormBuilder){}

  ngOnInit(): void {
    this.buildForm();
  }

  public onSubmit(){
    this.submitted = true
    if (this.userForm.valid && this.validatePassword(this.user.password)
    && this.validateEmail(this.user.email)) {
  console.log(this.user)
      this.user.type = "HOST" //create a question for user what he wants to be
      this.authService.signup(this.user).subscribe(data => {
        
      }
      , error =>{
        this.signupError = error.error.error;
      });
    } else {
      console.log(this.userForm.valid)
    }
  }

  private buildForm() {
    this.userForm = this.fb.group({
      fname: ["", Validators.pattern('^[a-zA-Z ]*$')],
      lname: ["", Validators.pattern('^[a-zA-Z ]*$')],
      email: ["", Validators.required],
      address: ["", Validators.required],
      username: ["", Validators.required],
      password: ["", Validators.required]
    });
  }

  validateEmail(email: string): boolean {
    if (!email.includes("@") || !email.includes(".")){
      return false
    }
    return true
  }

  validatePassword(password: string): boolean {

    if (password.length < 12){
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

  
