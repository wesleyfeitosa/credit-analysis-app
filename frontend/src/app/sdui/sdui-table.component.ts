import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output, ChangeDetectionStrategy } from '@angular/core';
import { MatChipsModule } from '@angular/material/chips';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { MatSortModule, Sort } from '@angular/material/sort';
import { MatTableModule } from '@angular/material/table';

import { SduiColumn } from '../core/models/sdui.model';

/**
 * Renders an SDUI `table` component: columns are described by the contract and
 * rows come from the API. The `status` and `createdAt` keys get special
 * rendering (status chip, formatted date); everything else prints its raw
 * value. Sorting and pagination are driven by the server via the outputs.
 */
@Component({
    selector: 'app-sdui-table',
    imports: [
        CommonModule,
        MatTableModule,
        MatSortModule,
        MatPaginatorModule,
        MatChipsModule,
        MatProgressBarModule,
    ],
    templateUrl: './sdui-table.component.html',
    changeDetection: ChangeDetectionStrategy.Eager,
    styleUrl: './sdui-table.component.scss'
})
export class SduiTableComponent {
  @Input({ required: true }) columns: SduiColumn[] = [];
  @Input() data: Record<string, unknown>[] = [];
  @Input() total = 0;
  @Input() pageIndex = 0;
  @Input() pageSize = 20;
  @Input() loading = false;

  @Output() sortChange = new EventEmitter<Sort>();
  @Output() pageChange = new EventEmitter<PageEvent>();
  @Output() rowClick = new EventEmitter<Record<string, unknown>>();

  get displayedColumns(): string[] {
    return this.columns.map((c) => c.key);
  }
}
