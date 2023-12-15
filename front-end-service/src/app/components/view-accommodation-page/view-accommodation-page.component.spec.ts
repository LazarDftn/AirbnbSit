import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewAccommodationPageComponent } from './view-accommodation-page.component';

describe('ViewAccommodationPageComponent', () => {
  let component: ViewAccommodationPageComponent;
  let fixture: ComponentFixture<ViewAccommodationPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ViewAccommodationPageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewAccommodationPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
