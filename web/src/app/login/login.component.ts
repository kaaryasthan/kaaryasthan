import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { Login } from '../login';
import { AuthService } from '../auth.service';

class LoginCredentials {
    username = '';
    password = '';
    key = '';
}

@Component({
    templateUrl: './login.component.html',
    styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

    cred = new LoginCredentials();

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public authService: AuthService) { }

    ngOnInit() {
    }

    updateUsername(value: string) {
        this.cred.username = value;
    }

    updatePassword(value: string) {
        this.cred.password = value;
    }

    newLogin() {
        console.log(this.cred);
        // store user details and jwt token in local storage to keep user logged in between page refreshes
        this.authService.login(this.cred)
            .subscribe(token => {
                localStorage.setItem('currentUser', token);
                this.router.navigate(['/']);
            });
    }
}
