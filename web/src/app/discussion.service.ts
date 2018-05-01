import { Injectable } from '@angular/core';
import { HttpHeaders, HttpClient, HttpErrorResponse } from '@angular/common/http';

import { Observable } from 'rxjs/Observable';
import { ErrorObservable } from 'rxjs/observable/ErrorObservable';

import 'rxjs/add/operator/map';
import { catchError, retry } from 'rxjs/operators';

import { Discussion } from './discussion';

export class DiscussionModel {
    id = 0;
    body = "";
    item_id = 0;
}

@Injectable()
export class DiscussionService {
    private discussionUrl = 'api/v1/discussions';

    constructor(private http: HttpClient) { }

    private handleError(error: HttpErrorResponse) {
        if (error.error instanceof ErrorEvent) {
            // A client-side or network error occurred. Handle it accordingly.
            console.error('An error occurred:', error.error.message);
        } else {
            // The backend returned an unsuccessful response code.
            // The response body may contain clues as to what went wrong,
            console.error(
                `Backend returned code ${error.status}, ` +
                `body was: ${error.error}`);
        }
        // return an ErrorObservable with a user-facing error message
        return new ErrorObservable(
            'Something bad happened; please try again later.');
    };

    create(disc: Discussion): Observable<Discussion> {
        const entity = {
            data: {
                type: 'discussions',
                attributes: {
                    body: disc.body,
                    item_id: disc.item_id,
                },
            }
        };

        const httpOptions = {
            headers: new HttpHeaders({
                'Content-Type': 'application/vnd.api+json',
                'Authorization': 'Bearer ' + localStorage.getItem('currentUser'),
            })
        };

        return this.http
            .post<Discussion>(this.discussionUrl, entity, httpOptions)
            .pipe(catchError(this.handleError))
            .map(data => {
                return data['data'].attributes;
            },
            error => {
                console.log(error);
            });
    }

}
