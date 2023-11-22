import { Component, OnInit } from '@angular/core';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { AuthService } from 'src/app/services/auth.service';

@Component({
  selector: 'app-account-verif-page',
  templateUrl: './account-verif-page.component.html',
  styleUrls: ['./account-verif-page.component.css']
})
export class AccountVerifPageComponent implements OnInit {
  
  message: string = ""
  constructor(private authService: AuthService,
    private route: ActivatedRoute){}
  ngOnInit(): void {
    var username = null
    var code = null

    this.route.paramMap.subscribe((params: ParamMap) => {
      username = params.get("id")
      code = params.get("id2")
    })

    this.authService.verifyEmail(username, code).subscribe(data => {
      this.message = "You have successfully verified your account!"
    }, error => {this.message = error.error.error})
  }

}
