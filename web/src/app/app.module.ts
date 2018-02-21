import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';

import { AppComponent } from './app.component';
import { LoginComponent } from './login/login.component';

import { LoginService } from './login.service';


@NgModule({
    declarations: [
        AppComponent,
        LoginComponent
    ],
    imports: [
        BrowserModule,
        HttpClientModule,
        FormsModule
    ],
    providers: [LoginService],
    bootstrap: [AppComponent]
})
export class AppModule { }
