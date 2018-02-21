import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { HttpHeaders } from '@angular/common/http';

import { Observable } from 'rxjs/Observable';

import { Login } from './login';

const httpOptions = {
    headers: new HttpHeaders({
        'Content-Type': 'application/json',
        'Authorization': 'my-auth-token'
    })
};

@Injectable()
export class LoginService {
    loginUrl = 'api/v1/login';

    constructor(
        private http: HttpClient) { }

    loginUser(loginCredentials: Login): Observable<Login> {
        const entity = {
            data: {
                type: "logins",
                attributes: {
                    username: loginCredentials.username,
                    password: loginCredentials.password,
                },
            }
        }
        return this.http.post<Login>(this.loginUrl, entity, httpOptions);
    }

}
