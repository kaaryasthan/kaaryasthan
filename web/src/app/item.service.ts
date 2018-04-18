import { Injectable } from '@angular/core';
import { HttpHeaders, HttpClient } from '@angular/common/http';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';

import { Item } from './item';

export class ItemModel {
    id = 0;
    title = "";
    description = "";
    project_id = 0;
}

@Injectable()
export class ItemService {
    private itmUrl = 'api/v1/items';

    constructor(private http: HttpClient) { }

    create(itm: Item): Observable<Item> {
        const entity = {
            data: {
                type: 'items',
                attributes: {
                    title: itm.title,
                    description: itm.description,
                    project_id: itm.project_id,
                },
            }
        };

        console.log('itm.project_id', itm.project_id);
        const httpOptions = {
            headers: new HttpHeaders({
                'Content-Type': 'application/vnd.api+json',
                'Authorization': 'Bearer ' + localStorage.getItem('currentUser'),
            })
        };

        console.log('Bearer ' + localStorage.getItem('currentUser'));
        return this.http
            .post(this.itmUrl, entity, httpOptions)
            .map(data => data['data'].attributes);
    }

    search(text: string): Observable<Item[]> {
        const httpOptions = {
            headers: new HttpHeaders({
                'Content-Type': 'application/vnd.api+json',
                'Authorization': 'Bearer ' + localStorage.getItem('currentUser'),
            })
        };

        console.log('Bearer ' + localStorage.getItem('currentUser'));

        return this.http
            .get(this.itmUrl + "?Query=" + text, httpOptions)
            .map(data => {
                var itmList: ItemModel[] = [];
                for (let i = 0; i < data['data'].length; i++) {
                    var o = data['data'][i];
                    var itm = new ItemModel();
                    itm.id = o.id;
                    itm.title = o.attributes.title;
                    itm.description = o.attributes.description;
                    itm.project_id = o.attributes.project_id;
                    itmList.push(itm);
                }
                return itmList;
            });

    }

}
