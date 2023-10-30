import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { DashboardBodyComponent } from './dashboard-body/dashboard-body.component';
import { DashboardRoutingModule } from './dashboard-routing.module';
import { PrimengModule } from 'src/app/shared/primeng.module';
import { DashboardHeaderComponent } from './dashboard-header/dashboard-header.component';
import { DashboardDatePickerComponent } from './dashboard-date-picker/dashboard-date-picker.component';
import { DashboardResumeComponent } from './dashboard-resume/dashboard-resume.component';



@NgModule({
  declarations: [
    DashboardBodyComponent,
    DashboardHeaderComponent,
    DashboardDatePickerComponent,
    DashboardResumeComponent
  ],
  imports: [
    CommonModule,
    DashboardRoutingModule,
    PrimengModule
  ]
})
export class DashboardModule { }
