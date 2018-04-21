import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { ItemService } from '../item.service';
import { ItemModel } from '../item.service';

@Component({
    selector: 'app-item-list',
    templateUrl: './item-list.component.html',
    styleUrls: ['./item-list.component.css']
})
export class ItemListComponent implements OnInit {

    public newQuery: string;;
    public query: string;
    public items: ItemModel[] = [];

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public itemService: ItemService) {
    }

    ngOnInit() {
        this.query = "";
        this.route.queryParams.subscribe(params => {
            if ("q" in params) {
                this.query = params["q"];
            }
            console.log("query: ", this.query);
            this.itemService.search(this.query)
                .subscribe(data => {
                    this.items = data;
                    console.log(data);
                });
        });
    }

    updateSearch(value: string) {
        this.newQuery = value;
    }

    newSearch() {
        console.log("new query", this.newQuery);
        this.router.navigate(["/items"], { queryParams: { q: this.newQuery } });
        // this.itemService.search(this.newQuery)
        //     .subscribe(token => {
        //         this.router.navigate(["/items"], { queryParams: { q: this.newQuery } });
        //     });
    }
}
