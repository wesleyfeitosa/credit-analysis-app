import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from '../../../environments/environment';
import { SduiScreen } from '../models/sdui.model';

/** Fetches SDUI screen contracts that describe how a screen should render. */
@Injectable({ providedIn: 'root' })
export class SduiService {
  private readonly http = inject(HttpClient);

  getScreen(name: string): Observable<SduiScreen> {
    return this.http.get<SduiScreen>(`${environment.apiBaseUrl}/sdui/screens/${name}`);
  }
}
