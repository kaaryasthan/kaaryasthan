import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';

import { LoginService } from './login.service';
import { AuthGuard } from './auth.guard';
import { routing } from './app-routing.module';
import { HomeComponent } from './home/home.component';

import { AuthService } from './auth.service';
import { EmailComponent } from './email/email.component';
import { RegisterComponent } from './register/register.component';
import { ProjectCreateComponent } from './project-create/project-create.component';

@NgModule({
    declarations: [
        AppComponent,
        LoginComponent,
        HomeComponent,
        EmailComponent,
        RegisterComponent,
        ProjectCreateComponent
    ],
    imports: [
        BrowserModule,
        HttpClientModule,
        FormsModule,
        routing
    ],
    providers: [
        LoginService,
        AuthGuard,
        AuthService
    ],
    bootstrap: [AppComponent]
})
export class AppModule { }
