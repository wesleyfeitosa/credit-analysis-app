// Package sdui builds the Server Driven UI screen contracts returned to the
// frontend. Changing layout here changes the rendered screens without a
// frontend deploy.
package sdui

import "encoding/json"

// Field describes a single filter/form input.
type Field struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Type    string   `json:"type"`
	Options []string `json:"options,omitempty"`
}

// Column describes a single table column.
type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

// Component is a renderable unit of a screen (form, filter, table, ...).
type Component struct {
	Type    string   `json:"type"`
	Fields  []Field  `json:"fields,omitempty"`
	Columns []Column `json:"columns,omitempty"`
	// Values carries the user's previously saved state for this component
	// (keyed by field key), so the frontend can render it pre-filled. Only set
	// on the filter component when the user has saved preferences; omitted
	// otherwise.
	Values json.RawMessage `json:"values,omitempty"`
}

// Screen is the full SDUI contract for one screen.
type Screen struct {
	Screen     string      `json:"screen"`
	Title      string      `json:"title"`
	Components []Component `json:"components"`
}

// LoginScreen returns the SDUI contract for the login screen.
func LoginScreen() Screen {
	return Screen{
		Screen: "login",
		Title:  "Acesso ao Sistema",
		Components: []Component{
			{
				Type: "form",
				Fields: []Field{
					{Key: "email", Label: "E-mail", Type: "text"},
					{Key: "password", Label: "Senha", Type: "password"},
				},
			},
		},
	}
}

// CreditAnalysesScreen returns the SDUI contract for the listing screen,
// describing both the filter and the table. When savedFilters is non-nil it is
// embedded as the filter component's values so the screen renders pre-filled
// with the user's previously saved preferences.
func CreditAnalysesScreen(savedFilters json.RawMessage) Screen {
	return Screen{
		Screen: "credit-analyses",
		Title:  "Análises de Crédito App",
		Components: []Component{
			{
				Type: "filter",
				Fields: []Field{
					{Key: "document", Label: "CPF/CNPJ", Type: "text"},
					{Key: "clientName", Label: "Cliente", Type: "text"},
					{Key: "status", Label: "Status", Type: "select", Options: []string{
						"APROVADO", "REPROVADO", "EM_ANALISE", "PENDENTE",
					}},
					{Key: "createdAt", Label: "Período", Type: "dateRange"},
					{Key: "score", Label: "Score", Type: "numberRange"},
				},
				Values: savedFilters,
			},
			{
				Type: "table",
				Columns: []Column{
					{Key: "clientName", Label: "Cliente"},
					{Key: "document", Label: "CPF/CNPJ"},
					{Key: "status", Label: "Status da análise"},
					{Key: "score", Label: "Score"},
					{Key: "createdAt", Label: "Data da Análise"},
				},
			},
		},
	}
}
