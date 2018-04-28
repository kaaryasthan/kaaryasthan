import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';

import { ItemService } from '../item.service';
import { ItemModel } from '../item.service';

import { Item } from '../item';

@Component({
    selector: 'app-item-show',
    templateUrl: './item-show.component.html',
    styleUrls: ['./item-show.component.css']
})
export class ItemShowComponent implements OnInit {

    public item$: Observable<Item>;

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public itemService: ItemService) { }

    ngOnInit() {
        this.item$ = this.route.paramMap
            .switchMap((params: ParamMap) =>
                this.itemService.getItem(params.get('num')));
    }
}
