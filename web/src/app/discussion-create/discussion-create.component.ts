import { Component, OnInit, Input } from '@angular/core';
import { DiscussionService } from '../discussion.service';
import { DiscussionModel } from '../discussion.service';

@Component({
    selector: 'app-discussion-create',
    templateUrl: './discussion-create.component.html',
    styleUrls: ['./discussion-create.component.css']
})
export class DiscussionCreateComponent implements OnInit {

    @Input() item_id: number;

    body: string;
    data = new DiscussionModel();

    constructor(
        public discussionService: DiscussionService) { }

    ngOnInit() {
    }

    updateDiscussionBody($event) {
        this.body = $event;
    }

    newDiscussion() {
        this.data.body = this.body;
        this.data.item_id = this.item_id;
        console.log(this.data);
        this.discussionService.create(this.data).subscribe();
    }

}
