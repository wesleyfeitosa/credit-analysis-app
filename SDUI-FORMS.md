# Por que SDUI é difícil em formulários (ex.: tela de login)

> Resumo: **SDUI é ótimo para descrever _exibição_; formulários são difíceis
> porque tratam de _comportamento, validação e efeitos colaterais_ — e JSON
> descreve isso muito mal.**

Uma tabela é uma função pura dos dados: o servidor envia as linhas, o cliente
renderiza. Uma direção só. Um formulário é bidirecional e imperativo — ele
_coleta_ entrada, _valida_ e _faz algo_ com ela. Esse "faz algo" é justamente a
parte que um contrato de layout não consegue expressar de forma limpa.

Exemplo concreto: neste projeto o formulário de login já é renderizado via SDUI
(`frontend/src/app/sdui/sdui-form.component.ts`), então dá para apontar as
dificuldades reais.

## 1. Validação não tem lugar no contrato

O campo SDUI é apenas `{key, label, type}` — **não há onde descrever regras**
("e-mail válido", "senha ≥ 4 caracteres"). Então elas ficaram hard-coded: o
`SduiFormComponent` marca todo campo de um componente `form` como `required`, e o
backend valida de novo. Para tornar a validação realmente server-driven seria
preciso um **DSL de validação** no JSON (required, minLength, regex, cross-field)
_mais_ um motor genérico para interpretá-lo — e ainda assim seria necessário
validar no servidor por segurança. No fim, mantém-se a regra em dois dialetos que
podem divergir.

## 2. A _ação_ do formulário não é declarativa

O login não apenas envia dados — ele guarda um JWT e redireciona para
`/credit-analyses`. Esse comportamento está hard-coded no `onSubmit()` e
basicamente **não pode** viver no contrato. SDUI descreve estrutura bem e
_efeitos colaterais_ mal: para onde fazer POST, o que fazer em caso de sucesso vs.
erro, o que persistir. Todo formulário tem lógica pós-submit específica.

## 3. Tipos de campo são um vocabulário fixo

O servidor só pode pedir widgets que o cliente já sabe renderizar (`text`,
`select`, `dateRange`…). Quer uma máscara de CPF, um autocomplete ou validação
assíncrona ("e-mail já cadastrado")? Isso é **código novo no frontend** — o que
anula a promessa central do SDUI de "mudar a UI sem deploy". Células read-only não
têm esse problema; inputs interativos têm.

## 4. Descasamento entre chave e payload

As chaves dos campos do formulário precisam mapear para o corpo esperado pela API.
Isso apareceu já no formulário de filtro: a chave do campo é `createdAt`, mas o
endpoint espera `dateFrom`/`dateTo`, então há código de tradução no componente.
Formulários criam um contrato apertado entre a chave do campo e o nome do
parâmetro da API que o JSON não captura — e acaba-se escrevendo código de mapeamento
mesmo assim.

## 5. O login especificamente tem o _pior_ custo/benefício

- **Benefício baixo:** a tela de login quase nunca muda de layout, e a lógica real
  dela (autenticação, rate limiting) já é totalmente server-side e precisa ser —
  SDUI não move nada disso.
- **Custo real:** ainda é preciso construir o motor genérico de formulário
  reativo, o interpretador de validação e o tratamento hard-coded de sucesso/erro.

Ou seja, paga-se pela infraestrutura na única tela que menos se beneficia.

---

## Conclusão pragmática

SDUI brilha em superfícies **declarativas e read-only** — a listagem (tabela) e a
tela de detalhe são encaixes naturais (UI = função dos dados). Para formulários, o
padrão sensato é o **híbrido** que usamos: deixar o contrato descrever a
_estrutura_ (campos, labels, tipos) para o layout seguir flexível, mas manter as
_regras de validação e o comportamento de submit_ no cliente e no endpoint. A
"dificuldade" é exatamente a parte que o SDUI não cobre bem — regras e ações — não
a renderização.

Se o time realmente quiser formulários server-driven, a resposta honesta de escopo
é: o contrato precisa ganhar um **schema de validação** e um **descritor de ação de
submit**, e aceita-se **validação duplicada (cliente + servidor)**. É um
investimento real — vale para muitos formulários dinâmicos espalhados por várias
telas, raramente vale para uma tela de login.
