package sdui

import "encoding/json"

type Field struct {
	Key     string   `json:"key"`
	Label   string   `json:"label"`
	Type    string   `json:"type"`
	Options []string `json:"options,omitempty"`
}

type Column struct {
	Key   string `json:"key"`
	Label string `json:"label"`
}

type Component struct {
	Type    string          `json:"type"`
	Fields  []Field         `json:"fields,omitempty"`
	Columns []Column        `json:"columns,omitempty"`
	Values  json.RawMessage `json:"values,omitempty"`
}

type Screen struct {
	Screen     string      `json:"screen"`
	Title      string      `json:"title"`
	Components []Component `json:"components"`
}

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
