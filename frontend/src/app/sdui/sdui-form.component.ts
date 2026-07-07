
import {
  Component,
  EventEmitter,
  Input,
  OnInit,
  Output,
  inject,
  ChangeDetectionStrategy
} from '@angular/core';
import { FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MatNativeDateModule } from '@angular/material/core';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';

import { SduiComponent, SduiField } from '../core/models/sdui.model';

/**
 * Renders an SDUI `form` or `filter` component: it builds a reactive form from
 * the field descriptions and emits the collected values on submit. Field types
 * map to controls as follows:
 *   text | password -> input       select -> mat-select
 *   dateRange       -> two dates (<key>From / <key>To, emitted as ISO strings)
 *   numberRange     -> two numbers (<key>Min / <key>Max)
 *
 * For `form` components every field is required; for `filter` components fields
 * are optional. The consumer maps the emitted keys to API params.
 */
@Component({
    selector: 'app-sdui-form',
    imports: [
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    MatDatepickerModule,
    MatNativeDateModule,
    MatButtonModule,
    MatIconModule
],
    templateUrl: './sdui-form.component.html',
    changeDetection: ChangeDetectionStrategy.Eager,
    styleUrl: './sdui-form.component.scss'
})
export class SduiFormComponent implements OnInit {
  @Input({ required: true }) component!: SduiComponent;
  @Input() submitLabel = 'Aplicar';
  @Input() loading = false;
  @Output() submitForm = new EventEmitter<Record<string, unknown>>();

  private readonly fb = inject(FormBuilder);
  form!: FormGroup;

  get fields(): SduiField[] {
    return this.component.fields ?? [];
  }

  get isFilter(): boolean {
    return this.component.type === 'filter';
  }

  ngOnInit(): void {
    const controls: Record<string, unknown> = {};
    const validators = this.isFilter ? [] : [Validators.required];
    // Saved preferences echoed back by the backend, keyed by control key.
    const saved = this.component.values ?? {};

    for (const field of this.fields) {
      switch (field.type) {
        case 'dateRange':
          // Dates are stored as ISO strings; the datepicker needs Date objects.
          controls[`${field.key}From`] = [toDate(saved[`${field.key}From`])];
          controls[`${field.key}To`] = [toDate(saved[`${field.key}To`])];
          break;
        case 'numberRange':
          controls[`${field.key}Min`] = [saved[`${field.key}Min`] ?? null];
          controls[`${field.key}Max`] = [saved[`${field.key}Max`] ?? null];
          break;
        default:
          controls[field.key] = [saved[field.key] ?? '', validators];
      }
    }
    this.form = this.fb.group(controls);
  }

  submit(): void {
    if (this.form.invalid) {
      this.form.markAllAsTouched();
      return;
    }

    const raw = this.form.value as Record<string, unknown>;
    const out: Record<string, unknown> = {};
    for (const [key, value] of Object.entries(raw)) {
      if (value === null || value === '' || value === undefined) {
        continue;
      }
      // Material datepicker yields Date objects; emit RFC3339 for the API.
      out[key] = value instanceof Date ? value.toISOString() : value;
    }
    this.submitForm.emit(out);
  }

  reset(): void {
    this.form.reset();
    this.submitForm.emit({});
  }
}

/** Parses a saved ISO date string into a Date, or null when absent/invalid. */
function toDate(value: unknown): Date | null {
  if (!value) {
    return null;
  }
  const d = new Date(value as string);
  return isNaN(d.getTime()) ? null : d;
}
