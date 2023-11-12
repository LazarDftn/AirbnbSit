import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateAccommodationPageComponent } from './create-accommodation-page.component';

describe('CreateAccommodationPageComponent', () => {
  let component: CreateAccommodationPageComponent;
  let fixture: ComponentFixture<CreateAccommodationPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ CreateAccommodationPageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateAccommodationPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
