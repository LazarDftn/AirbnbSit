import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EditAccommodationPageComponent } from './edit-accommodation-page.component';

describe('EditAccommodationPageComponent', () => {
  let component: EditAccommodationPageComponent;
  let fixture: ComponentFixture<EditAccommodationPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ EditAccommodationPageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EditAccommodationPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
