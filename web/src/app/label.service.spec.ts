import { TestBed, inject } from '@angular/core/testing';

import { LabelService } from './label.service';

describe('LabelService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [LabelService]
    });
  });

  it('should be created', inject([LabelService], (service: LabelService) => {
    expect(service).toBeTruthy();
  }));
});
