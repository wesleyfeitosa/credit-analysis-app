
import { Component, inject, signal, ChangeDetectionStrategy } from '@angular/core';
import { FormBuilder, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import {
  MatDialogModule,
  MatDialogRef,
} from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatProgressBarModule } from '@angular/material/progress-bar';

import { CreditAnalysisService } from '../../core/services/credit-analysis.service';

/**
 * Dialog to register a new analysis. On submit it calls POST /credit-analyses,
 * which emulates the full lifecycle server-side, and closes with the created
 * detail so the list can refresh.
 */
@Component({
    selector: 'app-create-analysis-dialog',
    imports: [
    ReactiveFormsModule,
    MatDialogModule,
    MatFormFieldModule,
    MatInputModule,
    MatButtonModule,
    MatProgressBarModule
],
    changeDetection: ChangeDetectionStrategy.Eager,
    templateUrl: './create-analysis-dialog.component.html'
})
export class CreateAnalysisDialogComponent {
  private readonly fb = inject(FormBuilder);
  private readonly service = inject(CreditAnalysisService);
  private readonly dialogRef = inject(MatDialogRef<CreateAnalysisDialogComponent>);

  loading = signal(false);
  error = signal<string | null>(null);

  form = this.fb.group({
    document: ['', [Validators.required, Validators.minLength(11)]],
    clientName: ['', [Validators.required, Validators.minLength(2)]],
  });

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }
    this.loading.set(true);
    this.error.set(null);
    const { document, clientName } = this.form.getRawValue();
    this.service.create(document!, clientName!).subscribe({
      next: (detail) => this.dialogRef.close(detail),
      error: () => {
        this.loading.set(false);
        this.error.set('Não foi possível criar a análise.');
      },
    });
  }
}
