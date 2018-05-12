import { Component, OnInit, Input } from '@angular/core';
import { Observable } from 'rxjs/Observable';

import { Comment } from '../comment';
import { CommentService } from '../comment.service';

@Component({
    selector: 'app-comment-list',
    templateUrl: './comment-list.component.html',
    styleUrls: ['./comment-list.component.css']
})
export class CommentListComponent implements OnInit {

    @Input() item_id: number;

    public comments$: Observable<Comment[]>;

    constructor(
        public commentService: CommentService) { }

    ngOnInit() {
        this.comments$ = this.commentService.list(this.item_id);

    }



}
