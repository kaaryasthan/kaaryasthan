import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';
import 'rxjs/add/operator/map'

import { ItemService } from '../item.service';
import { ProjectModel } from '../project.service';
import { ProjectService } from '../project.service';

class Item {
    id = 0;
    project_id = 0;
    title = '';
    description = '';
}

@Component({
    selector: 'app-item-create',
    templateUrl: './item-create.component.html',
    styleUrls: ['./item-create.component.css']
})
export class ItemCreateComponent implements OnInit {

    project = new ProjectModel();
    data = new Item();
    public projects: ProjectModel[] = [];

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public itemService: ItemService,
        public projectService: ProjectService) { }

    ngOnInit() {
        this.projectService.getAll().
            subscribe(data => {
                this.projects = data;
                console.log(this.projects);
            });
    }

    updateTitle(value: string) {
        this.data.title = value;
    }

    updateDescription(value: string) {
        this.data.description = value;
    }

    newItem() {
        this.data.project_id = this.project.id;
        this.itemService.create(this.data)
            .subscribe(token => {
                this.router.navigate(['/']);
            });
    }
}
