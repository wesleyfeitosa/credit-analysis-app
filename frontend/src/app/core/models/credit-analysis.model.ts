export type AnalysisStatus = 'APROVADO' | 'REPROVADO' | 'EM_ANALISE' | 'PENDENTE';

export interface CreditAnalysis {
  id: number;
  document: string;
  clientName: string;
  status: AnalysisStatus;
  score: number;
  createdAt: string;
}

export interface CreditAnalysisEvent {
  id: number;
  analysisId: number;
  status: AnalysisStatus;
  note: string;
  createdAt: string;
}

export interface CreditAnalysisDetail extends CreditAnalysis {
  events: CreditAnalysisEvent[];
}

export interface Page<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

// Query parameters accepted by GET /credit-analyses.
export interface ListParams {
  document?: string;
  clientName?: string;
  status?: string;
  scoreMin?: number;
  scoreMax?: number;
  dateFrom?: string;
  dateTo?: string;
  page?: number;
  pageSize?: number;
  sortBy?: string;
  sortDir?: 'asc' | 'desc';
}
