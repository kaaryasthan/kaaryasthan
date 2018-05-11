import { Component, OnInit, Input } from '@angular/core';
import { CommentService } from '../comment.service';
import { CommentModel } from '../comment.service';

@Component({
    selector: 'app-comment-create',
    templateUrl: './comment-create.component.html',
    styleUrls: ['./comment-create.component.css']
})
export class CommentCreateComponent implements OnInit {

    @Input() item_id: number;

    body: string;
    data = new CommentModel();

    constructor(
        public commentService: CommentService) { }

    ngOnInit() {
    }

    updateCommentBody($event) {
        this.body = $event;
    }

    newComment() {
        this.data.body = this.body;
        this.data.item_id = this.item_id;
        console.log(this.data);
        this.commentService.create(this.data).subscribe();
    }

}
