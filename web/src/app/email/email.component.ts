import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import 'rxjs/add/operator/filter';

import { AuthService } from '../auth.service';

class LoginCredentials {
    username = '';
    password = '';
    key = '';
}

@Component({
    selector: 'app-email',
    templateUrl: './email.component.html',
    styleUrls: ['./email.component.css']
})
export class EmailComponent implements OnInit {

    cred = new LoginCredentials();

    constructor(private route: ActivatedRoute,
        private router: Router,
        public authService: AuthService) { }

    ngOnInit() {
        this.route.queryParams
            .filter(params => params.key)
            .subscribe(params => {
                console.log(params); // {key: "popular"}
                this.cred.key = params.key;
                console.log(this.cred.key); // popular
            });
    }


    updateUsername(value: string) {
        this.cred.username = value
    }

    updatePassword(value: string) {
        this.cred.password = value
    }

    newLogin() {
        console.log(this.cred);
        // store user details and jwt token in local storage to keep user logged in between page refreshes
        this.authService.verifyemail(this.cred)
            .subscribe(token => {
                localStorage.setItem('currentUser', token);
                this.router.navigate(['/']);
            });
    }

}
