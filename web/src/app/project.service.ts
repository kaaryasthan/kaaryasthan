import { Injectable } from '@angular/core';
import { HttpHeaders, HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';

import { Project } from './project';

export class ProjectModel {
    id = 0;
    name = "";
    description = "";
}

@Injectable()
export class ProjectService {
    private prjUrl = 'api/v1/projects';

    constructor(private http: HttpClient) { }


    getAll(): Observable<Project[]> {
        const httpOptions = {
            headers: new HttpHeaders({
                'Content-Type': 'application/vnd.api+json',
                'Authorization': 'Bearer ' + localStorage.getItem('currentUser'),
            })
        };

        console.log('Bearer ' + localStorage.getItem('currentUser'));

        return this.http
            .get(this.prjUrl, httpOptions)
            .map(data => {
                var prjList: ProjectModel[] = [];
                for (let i = 0; i < data['data'].length; i++) {
                    var o = data['data'][i];
                    var prj = new ProjectModel();
                    prj.id = o.id;
                    prj.name = o.attributes.name;
                    prj.description = o.attributes.description;
                    prjList.push(prj);
                }
                return prjList;
            });
    }

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
