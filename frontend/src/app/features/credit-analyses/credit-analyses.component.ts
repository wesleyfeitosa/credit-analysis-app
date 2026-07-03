
import { Component, OnInit, inject, signal, ChangeDetectionStrategy } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { PageEvent } from '@angular/material/paginator';
import { MatSnackBar } from '@angular/material/snack-bar';
import { Sort } from '@angular/material/sort';
import { Router } from '@angular/router';

import {
  CreditAnalysis,
  ListParams,
  Page,
} from '../../core/models/credit-analysis.model';
import { SduiComponent, SduiScreen } from '../../core/models/sdui.model';
import { CreditAnalysisService } from '../../core/services/credit-analysis.service';
import { SduiService } from '../../core/services/sdui.service';
import { SduiFormComponent } from '../../sdui/sdui-form.component';
import { SduiTableComponent } from '../../sdui/sdui-table.component';
import { CreateAnalysisDialogComponent } from './create-analysis-dialog.component';

@Component({
    selector: 'app-credit-analyses',
    imports: [
    MatCardModule,
    MatButtonModule,
    MatIconModule,
    MatDialogModule,
    SduiFormComponent,
    SduiTableComponent
],
    templateUrl: './credit-analyses.component.html',
    changeDetection: ChangeDetectionStrategy.Eager,
    styleUrl: './credit-analyses.component.scss'
})
export class CreditAnalysesComponent implements OnInit {
  private readonly sdui = inject(SduiService);
  private readonly service = inject(CreditAnalysisService);
  private readonly router = inject(Router);
  private readonly dialog = inject(MatDialog);
  private readonly snackBar = inject(MatSnackBar);

  screen = signal<SduiScreen | null>(null);
  page = signal<Page<CreditAnalysis> | null>(null);
  loading = signal(false);
  error = signal<string | null>(null);

  // Filters chosen by the user, merged into every request.
  private filters: ListParams = {};
  private params: ListParams = { page: 1, pageSize: 20, sortBy: 'createdAt', sortDir: 'desc' };

  ngOnInit(): void {
    this.sdui.getScreen('credit-analyses').subscribe({
      next: (screen) => {
        this.screen.set(screen);
        this.load();
      },
      error: () => this.error.set('Não foi possível carregar a tela.'),
    });
  }

  get filterComponent(): SduiComponent | undefined {
    return this.screen()?.components.find((c) => c.type === 'filter');
  }

  get tableComponent(): SduiComponent | undefined {
    return this.screen()?.components.find((c) => c.type === 'table');
  }

  load(): void {
    this.loading.set(true);
    this.error.set(null);
    this.service.list({ ...this.params, ...this.filters }).subscribe({
      next: (page) => {
        this.page.set(page);
        this.loading.set(false);
      },
      error: () => {
        this.loading.set(false);
        this.error.set('Erro ao carregar análises.');
      },
    });
  }

  /**
   * Maps the generic keys emitted by the SDUI filter form to the API's param
   * names. The `createdAt` range field becomes dateFrom/dateTo; score already
   * matches scoreMin/scoreMax.
   */
  onFilter(values: Record<string, unknown>): void {
    const mapped: ListParams = {
      document: values['document'] as string,
      clientName: values['clientName'] as string,
      status: values['status'] as string,
      scoreMin: values['scoreMin'] as number,
      scoreMax: values['scoreMax'] as number,
      dateFrom: values['createdAtFrom'] as string,
      dateTo: values['createdAtTo'] as string,
    };
    this.filters = mapped;
    this.params.page = 1;
    this.load();
  }

  onSort(sort: Sort): void {
    this.params.sortBy = sort.direction ? sort.active : 'createdAt';
    this.params.sortDir = (sort.direction || 'desc') as 'asc' | 'desc';
    this.load();
  }

  onPage(event: PageEvent): void {
    this.params.page = event.pageIndex + 1;
    this.params.pageSize = event.pageSize;
    this.load();
  }

  onRowClick(row: Record<string, unknown>): void {
    this.router.navigate(['/credit-analyses', row['id']]);
  }

  openCreateDialog(): void {
    const ref = this.dialog.open(CreateAnalysisDialogComponent, { width: '420px' });
    ref.afterClosed().subscribe((created) => {
      if (created) {
        this.snackBar.open(
          `Análise criada: ${created.status} (score ${created.score})`,
          'OK',
          { duration: 4000 },
        );
        this.load();
      }
    });
  }

  get pageIndex(): number {
    return (this.params.page ?? 1) - 1;
  }
}
