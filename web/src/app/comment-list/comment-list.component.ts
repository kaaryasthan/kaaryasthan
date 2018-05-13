import { Component, OnInit, OnDestroy, Input } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { Comment } from '../comment';
import { CommentService } from '../comment.service';

@Component({
    selector: 'app-comment-list',
    templateUrl: './comment-list.component.html',
    styleUrls: ['./comment-list.component.css']
})
export class CommentListComponent implements OnInit, OnDestroy {

    @Input() item_id: number;

    public comments$: Observable<Comment[]>;

    ngOnDestroy() { }

    constructor(
        public commentService: CommentService) { }

    ngOnInit() {
        this.comments$ = this.commentService.list(this.item_id);
    }

    refreshComments() {
        this.comments$ = this.commentService.list(this.item_id);
    }

}
