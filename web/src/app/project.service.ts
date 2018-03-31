import { Injectable } from '@angular/core';
import { HttpHeaders, HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';

import { Project } from './project';

@Injectable()
export class ProjectService {
    private prjUrl = 'api/v1/projects';

    constructor(private http: HttpClient) { }

    create(prj: Project): Observable<Project> {
        const entity = {
            data: {
                type: 'projects',
                attributes: {
                    name: prj.name,
                    description: prj.description,
                },
            }
        };

        const httpOptions = {
            headers: new HttpHeaders({
                'Content-Type': 'application/vnd.api+json',
                'Authorization': 'Bearer ' + localStorage.getItem('currentUser'),
            })
        };

        console.log('Bearer ' + localStorage.getItem('currentUser'));
        return this.http
            .post(this.prjUrl, entity, httpOptions)
            .map(data => data['data'].attributes);
    }

}
