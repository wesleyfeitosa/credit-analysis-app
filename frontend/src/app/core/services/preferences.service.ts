import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from '../../../environments/environment';

/**
 * Persists the authenticated user's filter preferences so they are restored
 * between sessions. The body is the raw SDUI filter form state (keyed by field
 * control key), which the backend echoes back as the filter component's
 * `values` on the next screen load.
 */
@Injectable({ providedIn: 'root' })
export class PreferencesService {
  private readonly http = inject(HttpClient);
  private readonly url = `${environment.apiBaseUrl}/users/preferences/filters`;

  saveFilters(filters: Record<string, unknown>): Observable<void> {
    return this.http.post<void>(this.url, filters);
  }
}
