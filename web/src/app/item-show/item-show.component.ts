import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { Component, OnInit, AfterViewInit, ViewChild } from '@angular/core';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';

import { ItemService } from '../item.service';
import { ItemModel } from '../item.service';

import { Item } from '../item';

import { CommentListComponent } from '../comment-list/comment-list.component';

@Component({
    selector: 'app-item-show',
    templateUrl: './item-show.component.html',
    styleUrls: ['./item-show.component.css']
})
export class ItemShowComponent implements AfterViewInit {


    @ViewChild(CommentListComponent)
    private commentListComponent: CommentListComponent;

    public item$: Observable<Item>;

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public itemService: ItemService) { }

    ngAfterViewInit() {
        this.item$ = this.route.paramMap
            .switchMap((params: ParamMap) =>
                this.itemService.getItem(params.get('num')));
    }

    onCommentAdded(itmID: string) {
        console.log('Item updated:' + itmID);
        this.commentListComponent.refreshComments();
    }

    refreshComments() {
        this.commentListComponent.refreshComments();
    }
}
