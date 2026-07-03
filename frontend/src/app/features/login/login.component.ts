
import { Component, OnInit, inject, signal, ChangeDetectionStrategy } from '@angular/core';
import { MatCardModule } from '@angular/material/card';
import { MatProgressBarModule } from '@angular/material/progress-bar';
import { Router } from '@angular/router';

import { AuthService } from '../../core/services/auth.service';
import { SduiService } from '../../core/services/sdui.service';
import { SduiComponent, SduiScreen } from '../../core/models/sdui.model';
import { SduiFormComponent } from '../../sdui/sdui-form.component';

@Component({
    selector: 'app-login',
    imports: [MatCardModule, MatProgressBarModule, SduiFormComponent],
    templateUrl: './login.component.html',
    changeDetection: ChangeDetectionStrategy.Eager,
    styleUrl: './login.component.scss'
})
export class LoginComponent implements OnInit {
  private readonly sdui = inject(SduiService);
  private readonly auth = inject(AuthService);
  private readonly router = inject(Router);

  screen = signal<SduiScreen | null>(null);
  loading = signal(false);
  error = signal<string | null>(null);

  ngOnInit(): void {
    this.sdui.getScreen('login').subscribe({
      next: (screen) => this.screen.set(screen),
      error: () => this.error.set('Não foi possível carregar a tela de login.'),
    });
  }

  get formComponent(): SduiComponent | undefined {
    return this.screen()?.components.find((c) => c.type === 'form');
  }

  onSubmit(values: Record<string, unknown>): void {
    this.error.set(null);
    this.loading.set(true);
    this.auth.login(String(values['email']), String(values['password'])).subscribe({
      next: () => this.router.navigate(['/credit-analyses']),
      error: (err) => {
        this.loading.set(false);
        this.error.set(
          err.status === 401
            ? 'Credenciais inválidas.'
            : 'Erro ao autenticar. Tente novamente.',
        );
      },
    });
  }
}
