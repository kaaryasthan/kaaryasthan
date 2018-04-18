import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { ItemService } from '../item.service';

@Component({
    selector: 'app-item-list',
    templateUrl: './item-list.component.html',
    styleUrls: ['./item-list.component.css']
})
export class ItemListComponent implements OnInit {

    data;

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public itemService: ItemService) { }

    ngOnInit() {
    }

    updateSearch(value: string) {
        this.data = value;
    }

    newSearch() {
        console.log("aaaaaaaaaaaaaaa", this.data)
        this.itemService.search(this.data)
            .subscribe(token => {
                this.router.navigate(['/']);
            });

    }
}
