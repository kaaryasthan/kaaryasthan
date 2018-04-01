import { TestBed, inject } from '@angular/core/testing';

import { DiscussionService } from './discussion.service';

describe('DiscussionService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [DiscussionService]
    });
  });

  it('should be created', inject([DiscussionService], (service: DiscussionService) => {
    expect(service).toBeTruthy();
  }));
});
