import { HttpClient } from '@angular/common/http';
import { Injectable, inject, signal } from '@angular/core';
import { Observable, map, tap } from 'rxjs';

import { environment } from '../../../environments/environment';

interface LoginResponse {
  token: string;
}

/**
 * AuthService handles login, JWT persistence and auth state. The token is kept
 * in localStorage so the session survives reloads, and exposed as a signal so
 * the UI (e.g. the toolbar) reacts to login/logout.
 */
@Injectable({ providedIn: 'root' })
export class AuthService {
  private readonly http = inject(HttpClient);
  private readonly tokenKey = 'ca_token';

  readonly isAuthenticated = signal<boolean>(!!this.token);

  login(email: string, password: string): Observable<void> {
    return this.http
      .post<LoginResponse>(`${environment.apiBaseUrl}/auth/login`, { email, password })
      .pipe(
        tap((res) => {
          localStorage.setItem(this.tokenKey, res.token);
          this.isAuthenticated.set(true);
        }),
        map(() => void 0),
      );
  }

  logout(): void {
    localStorage.removeItem(this.tokenKey);
    this.isAuthenticated.set(false);
  }

  get token(): string | null {
    return localStorage.getItem(this.tokenKey);
  }
}
