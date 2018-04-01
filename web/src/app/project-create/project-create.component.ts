import { Router, ActivatedRoute } from '@angular/router';
import { Component, OnInit } from '@angular/core';

import { ProjectService } from '../project.service';
import { ProjectModel } from '../project.service';

// class Project {
//     name = '';
//     description = '';
// }

@Component({
    selector: 'app-project-create',
    templateUrl: './project-create.component.html',
    styleUrls: ['./project-create.component.css']
})
export class ProjectCreateComponent implements OnInit {

    data = new ProjectModel();

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public projectService: ProjectService) { }

    ngOnInit() {
    }

    updateName(value: string) {
        this.data.name = value;
    }

    updateDescription(value: string) {
        this.data.description = value;
    }

    newProject() {
        console.log(this.data);
        this.projectService.create(this.data)
            .subscribe(token => {
                this.router.navigate(['/']);
            });
    }
}
