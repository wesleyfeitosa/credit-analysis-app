import { Routes } from '@angular/router';

import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  {
    path: 'login',
    loadComponent: () =>
      import('./features/login/login.component').then((m) => m.LoginComponent),
  },
  {
    path: 'credit-analyses',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/credit-analyses/credit-analyses.component').then(
        (m) => m.CreditAnalysesComponent,
      ),
  },
  {
    path: 'credit-analyses/:id',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/credit-analyses/analysis-detail.component').then(
        (m) => m.AnalysisDetailComponent,
      ),
  },
  { path: '', pathMatch: 'full', redirectTo: 'credit-analyses' },
  { path: '**', redirectTo: 'credit-analyses' },
];
