import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AccountVerifPageComponent } from './account-verif-page.component';

describe('AccountVerifPageComponent', () => {
  let component: AccountVerifPageComponent;
  let fixture: ComponentFixture<AccountVerifPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ AccountVerifPageComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AccountVerifPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
