import { Injectable } from '@angular/core';
import { Response, Headers } from '@angular/http';
import { HttpClient } from '@angular/common/http';
import { HttpHeaders } from '@angular/common/http';
import * as jwt_decode from 'jwt-decode';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';
import { Login } from './login';
import { Register } from './register';

export const TOKEN_NAME: string = 'jwt_token';
const httpOptions = {
    headers: new HttpHeaders({
        'Content-Type': 'application/json',
        'Authorization': 'my-auth-token'
    })
};

@Injectable()
export class AuthService {
    private loginUrl: string = 'api/v1/login';
    private emailUrl: string = 'api/v1/emailverification';
    private registerUrl: string = 'api/v1/register';
    private headers = new Headers({ 'Content-Type': 'application/json' });
    private token: string;

    constructor(private http: HttpClient) { }

    getToken(): string {
        return localStorage.getItem(TOKEN_NAME);
    }

    setToken(token: string): void {
        localStorage.setItem(TOKEN_NAME, token);
    }

    getTokenExpirationDate(token: string): Date {
        const decoded = jwt_decode(token);

        if (decoded.exp === undefined) return null;

        const date = new Date(0);
        date.setUTCSeconds(decoded.exp);
        return date;
    }

    isTokenExpired(token?: string): boolean {
        if (!token) token = this.getToken();
        if (!token) return true;

        const date = this.getTokenExpirationDate(token);
        if (date === undefined) return false;
        return !(date.valueOf() > new Date().valueOf());
    }

    login(user: Login): Observable<string> {
        const entity = {
            data: {
                type: "logins",
                attributes: {
                    username: user.username,
                    password: user.password,
                },
            }
        }
        return this.http
            .post(this.loginUrl, entity, httpOptions)
            .map(data => data['data'].attributes.token);
    }

    verifyemail(user: Login): Observable<string> {
        const entity = {
            data: {
                type: "logins",
                attributes: {
                    username: user.username,
                    password: user.password,
                    email_verification_code: user.key,
                },
            }
        }
        return this.http
            .post(this.emailUrl, entity, httpOptions)
            .map(data => data['data'].attributes.token);
    }

    register(user: Register): Observable<string> {
        const entity = {
            data: {
                type: "logins",
                attributes: {
                    username: user.username,
                    password: user.password,
                    fullname: user.fullname,
                    email: user.email,
                },
            }
        }
        return this.http
            .post(this.registerUrl, entity, httpOptions)
            .map(data => data['data'].attributes.token);
    }

}
