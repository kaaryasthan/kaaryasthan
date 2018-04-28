import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DiscussionEditComponent } from './discussion-edit.component';

describe('DiscussionEditComponent', () => {
  let component: DiscussionEditComponent;
  let fixture: ComponentFixture<DiscussionEditComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DiscussionEditComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DiscussionEditComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
