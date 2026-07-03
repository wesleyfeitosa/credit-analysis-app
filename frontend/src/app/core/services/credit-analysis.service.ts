import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from '../../../environments/environment';
import {
  CreditAnalysis,
  CreditAnalysisDetail,
  ListParams,
  Page,
} from '../models/credit-analysis.model';

@Injectable({ providedIn: 'root' })
export class CreditAnalysisService {
  private readonly http = inject(HttpClient);
  private readonly base = `${environment.apiBaseUrl}/credit-analyses`;

  list(params: ListParams): Observable<Page<CreditAnalysis>> {
    let httpParams = new HttpParams();
    for (const [key, value] of Object.entries(params)) {
      if (value !== undefined && value !== null && value !== '') {
        httpParams = httpParams.set(key, String(value));
      }
    }
    return this.http.get<Page<CreditAnalysis>>(this.base, { params: httpParams });
  }

  get(id: number): Observable<CreditAnalysisDetail> {
    return this.http.get<CreditAnalysisDetail>(`${this.base}/${id}`);
  }

  create(document: string, clientName: string): Observable<CreditAnalysisDetail> {
    return this.http.post<CreditAnalysisDetail>(this.base, { document, clientName });
  }
}
