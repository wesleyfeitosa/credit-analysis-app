// Mirrors the backend SDUI contract (internal/sdui). The frontend renders these
// descriptions dynamically instead of hard-coding screen layouts.

export type SduiFieldType = 'text' | 'password' | 'select' | 'dateRange' | 'numberRange';

export interface SduiField {
  key: string;
  label: string;
  type: SduiFieldType;
  options?: string[];
}

export interface SduiColumn {
  key: string;
  label: string;
}

export type SduiComponentType = 'form' | 'filter' | 'table';

export interface SduiComponent {
  type: SduiComponentType;
  fields?: SduiField[];
  columns?: SduiColumn[];
  // Previously saved state for this component (keyed by field/control key),
  // sent by the backend so the frontend can render it pre-filled. Only present
  // on the filter component when the user has saved preferences.
  values?: Record<string, unknown>;
}

export interface SduiScreen {
  screen: string;
  title: string;
  components: SduiComponent[];
}
