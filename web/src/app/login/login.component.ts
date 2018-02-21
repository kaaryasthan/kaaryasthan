import { Component, OnInit } from '@angular/core';

import { Login } from '../login';
import { LoginService } from '../login.service';

class LoginCredentials {
    username = '';
    password = '';
}

@Component({
    selector: 'app-login',
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

    cred = new LoginCredentials()

    constructor(public loginService: LoginService) { }

    ngOnInit() {
    }

    updateUsername(value: string) {
        this.cred.username = value
    }

    updatePassword(value: string) {
        this.cred.password = value
    }

    newLogin() {
        console.log(this.cred)
        this.loginService.loginUser(this.cred).subscribe();
    }
}
