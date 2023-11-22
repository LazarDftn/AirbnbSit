import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SearchAccommodationPageComponent } from './search-accommodation-page.component';

describe('SearchAccommodationPageComponent', () => {
  let component: SearchAccommodationPageComponent;
  let fixture: ComponentFixture<SearchAccommodationPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SearchAccommodationPageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(SearchAccommodationPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
