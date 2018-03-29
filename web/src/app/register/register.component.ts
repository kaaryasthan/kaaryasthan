import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { AuthService } from '../auth.service';

class Register {
    username = '';
    password = '';
    fullname = '';
    email = '';
}

@Component({
    selector: 'app-register',
    templateUrl: './register.component.html',
    styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {

    cred = new Register();

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public authService: AuthService) { }

    ngOnInit() {
    }

    updateUsername(value: string) {
        this.cred.username = value
    }

    updatePassword(value: string) {
        this.cred.password = value
    }

    updateFullname(value: string) {
        this.cred.fullname = value
    }

    updateEmail(value: string) {
        this.cred.email = value
    }

    newRegister() {
        console.log(this.cred);
        // store user details and jwt token in local storage to keep user logged in between page refreshes
        this.authService.register(this.cred)
            .subscribe(token => {
                localStorage.setItem('currentUser', token);
                this.router.navigate(['/']);
            });
    }
}
