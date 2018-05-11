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
import { ProjectService } from './project.service';
import { ItemService } from './item.service';
import { CommentService } from './comment.service';

import { EmailComponent } from './email/email.component';
import { RegisterComponent } from './register/register.component';
import { ProjectCreateComponent } from './project-create/project-create.component';
import { ItemCreateComponent } from './item-create/item-create.component';
import { ItemListComponent } from './item-list/item-list.component';
import { ItemShowComponent } from './item-show/item-show.component';
import { ItemEditComponent } from './item-edit/item-edit.component';
import { CommentCreateComponent } from './comment-create/comment-create.component';
import { CommentListComponent } from './comment-list/comment-list.component';
import { CommentEditComponent } from './comment-edit/comment-edit.component';

@NgModule({
    declarations: [
        AppComponent,
        LoginComponent,
        HomeComponent,
        EmailComponent,
        RegisterComponent,
        ProjectCreateComponent,
        ItemCreateComponent,
        ItemListComponent,
        ItemShowComponent,
        ItemEditComponent,
        CommentCreateComponent,
        CommentListComponent,
        CommentEditComponent
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
        AuthService,
        ProjectService,
        ItemService,
        CommentService
    ],
    bootstrap: [AppComponent]
})
export class AppModule { }
