import { CommonModule } from '@angular/common';
import { Component, OnInit, inject, signal, ChangeDetectionStrategy } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MatCardModule } from '@angular/material/card';
import { MatIconModule } from '@angular/material/icon';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { ActivatedRoute, Router } from '@angular/router';

import { CreditAnalysisDetail } from '../../core/models/credit-analysis.model';
import { CreditAnalysisService } from '../../core/services/credit-analysis.service';

@Component({
    selector: 'app-analysis-detail',
    imports: [
        CommonModule,
        MatCardModule,
        MatButtonModule,
        MatIconModule,
        MatProgressBarModule,
    ],
    templateUrl: './analysis-detail.component.html',
    changeDetection: ChangeDetectionStrategy.Eager,
    styleUrl: './analysis-detail.component.scss'
})
export class AnalysisDetailComponent implements OnInit {
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);
  private readonly service = inject(CreditAnalysisService);

  detail = signal<CreditAnalysisDetail | null>(null);
  loading = signal(true);
  error = signal<string | null>(null);

  ngOnInit(): void {
    const id = Number(this.route.snapshot.paramMap.get('id'));
    if (!id) {
      this.error.set('Análise inválida.');
      this.loading.set(false);
      return;
    }
    this.service.get(id).subscribe({
      next: (detail) => {
        this.detail.set(detail);
        this.loading.set(false);
      },
      error: (err) => {
        this.loading.set(false);
        this.error.set(
          err.status === 404 ? 'Análise não encontrada.' : 'Erro ao carregar a análise.',
        );
      },
    });
  }

  back(): void {
    this.router.navigate(['/credit-analyses']);
  }
}
