import { Router, ActivatedRoute } from '@angular/router';
import { Component } from '@angular/core';
import { AuthService } from './auth.service';

@Component({
    selector: 'app',
    templateUrl: './app.component.html',
    styleUrls: ['./app.component.css']
})
export class AppComponent {
    title = 'app';

    constructor(
        private route: ActivatedRoute,
        private router: Router,
        public authService: AuthService) { }

    logout() {
        localStorage.removeItem('currentUser');
        this.router.navigate(['/login']);
    }
}
