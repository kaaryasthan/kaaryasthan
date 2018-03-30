import { Component, OnInit } from '@angular/core';

@Component({
    selector: 'app-project-create',
    templateUrl: './project-create.component.html',
    styleUrls: ['./project-create.component.css']
})
export class ProjectCreateComponent implements OnInit {

    constructor() { }

    ngOnInit() {
    }

    updateName(value: string) {
        //this.cred.email = value;
    }

    updateDescription(value: string) {
        //this.cred.email = value;
    }

    newProject() {
        //     console.log(this.cred);
        //     // store user details and jwt token in local storage to keep user logged in between page refreshes
        //     this.authService.register(this.cred)
        //         .subscribe(token => {
        //             localStorage.setItem('currentUser', token);
        //             this.router.navigate(['/']);
        //         });
    }
}
