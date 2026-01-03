# MVP Specification v2.0 - Funil Rápido

**Versão:** 2.0 (Atualizado com Feedback da Equipe)  
**Data:** 31 de Dezembro de 2025  
**Status:** Ready for Development  
**Prazo Estimado:** 6-8 semanas

---

## 📋 Visão Geral do MVP

### Objetivo do MVP

Validar a hipótese central de que **upsells fluidos com one-click checkout aumentam o AOV sem fricção para o usuário**, diferenciando-se de soluções existentes como order bumps do Hotmart.

### Hipótese a Validar

*"Infoprodutores brasileiros conseguem aumentar seu AOV em pelo menos 25% usando funis com páginas de upsell dedicadas e one-click checkout, comparado a não usar upsells ou usar order bumps tradicionais."*

### Escopo do MVP

**Fluxo Mínimo:**
```
Landing Page → Página de Vendas (checkout inline) → Upsell 1 → Confirmação
```

**Foco:** Simplicidade e validação rápida. Sequências múltiplas de upsells/downsells ficam para V1+.

---

## 🎯 Features do MVP

### Resumo Executivo

| # | Feature | Prioridade | Complexidade | Tempo Estimado |
|---|---------|------------|--------------|----------------|
| 1 | Editor No-Code Mobile-First | CRÍTICA | Alta | 2-3 semanas |
| 2 | Lógica de Funil com Proteção de URL | CRÍTICA | Média | 1 semana |
| 3 | Checkout Inline com PIX | CRÍTICA | Alta | 2 semanas |
| 4 | One-Click Upsell | CRÍTICA | Média | 1 semana |
| 5 | Analytics Básico do Funil | ALTA | Média | 1 semana |
| 6 | Custom Domain & SSL | ALTA | Baixa | 3 dias |
| 7 | Sistema de Templates | MÉDIA | Média | 1 semana |

**Total Estimado:** 6-8 semanas (considerando paralelização e imprevistos)

---

## 📝 Especificação Detalhada das Features

---

## Feature 1: Editor No-Code Mobile-First

### Descrição

Editor visual drag-and-drop para criação de páginas de funil, projetado com abordagem mobile-first. Permite criar Landing Pages, Páginas de Vendas, Páginas de Upsell e Páginas de Confirmação sem escrever código.

### Contexto e Motivação

**Feedback da Equipe:**
- Marcelo: "Também mobile first no design da landing page"
- Dante: Performance e mobile são baseline, não diferenciais
- Pesquisa de mercado: Brasil tem 3:1 mobile-to-desktop ratio

**Por que Mobile-First:**
- 75% do tráfego no Brasil vem de mobile
- Design responsivo tradicional (desktop → mobile) resulta em experiência mobile inferior
- Mobile-first garante experiência perfeita no dispositivo mais usado

### User Stories

**Como** infoprodutor iniciante (Juliana)  
**Quero** criar uma página de vendas profissional sem saber programar  
**Para que** eu possa lançar meu curso rapidamente

**Como** usuário mobile  
**Quero** que o editor funcione perfeitamente no meu celular  
**Para que** eu possa editar páginas de qualquer lugar

### Requisitos Funcionais

#### 1.1 Interface do Editor

**Descrição:** Interface WYSIWYG (What You See Is What You Get) com preview em tempo real.

**Componentes:**
- Canvas de edição (área principal)
- Barra lateral de componentes (biblioteca)
- Painel de propriedades (configurações do elemento selecionado)
- Barra superior (ações: salvar, preview, publicar)
- Seletor de dispositivo (mobile/desktop preview)

**Comportamento:**
- Drag-and-drop de componentes da biblioteca para o canvas
- Click para selecionar elemento e mostrar propriedades
- Double-click para editar texto inline
- Auto-save a cada 30 segundos
- Undo/Redo (até 50 ações)

#### 1.2 Biblioteca de Componentes

**Componentes Básicos (MVP):**

1. **Texto**
   - Título (H1, H2, H3)
   - Parágrafo
   - Lista (bullets, numerada)

2. **Mídia**
   - Imagem (upload ou URL)
   - Vídeo (embed YouTube/Vimeo)

3. **Botões**
   - Botão primário (CTA)
   - Botão secundário
   - Botão de link

4. **Formulários**
   - Campo de texto
   - Campo de email
   - Checkbox

5. **Layout**
   - Seção (container)
   - Coluna (grid 1-3 colunas)
   - Espaçador (vertical spacing)

6. **Social Proof**
   - Depoimento (texto + foto + nome)
   - Badge de garantia

**Componentes Especiais:**
- **Checkout Inline** (só para Página de Vendas)
- **Botão de Upsell** (só para Página de Upsell)

#### 1.3 Sistema de Propriedades

**Propriedades Globais (todos componentes):**
- Margem (top, bottom, left, right)
- Padding (top, bottom, left, right)
- Background (cor sólida ou gradiente)
- Visibilidade (mobile/desktop)

**Propriedades de Texto:**
- Fonte (3 opções: Sans-serif, Serif, Monospace)
- Tamanho (12-72px)
- Cor
- Alinhamento (esquerda, centro, direita)
- Peso (normal, bold)

**Propriedades de Botão:**
- Texto
- Cor de fundo
- Cor do texto
- Tamanho (pequeno, médio, grande)
- Largura (auto, 100%)
- Ação (link, submit form, checkout, accept upsell)

**Propriedades de Imagem:**
- URL ou upload
- Alt text
- Largura (%, px)
- Alinhamento
- Border radius

#### 1.4 Mobile-First Design System

**Breakpoints:**
- Mobile: 320-767px (design primário)
- Desktop: 768px+ (enhancement)

**Comportamento:**
- Editor mostra preview mobile por padrão
- Usuário pode alternar para preview desktop
- Mudanças feitas no mobile aplicam-se ao desktop (com ajustes automáticos)
- Mudanças no desktop NÃO afetam mobile (override)

**Grid System:**
- Mobile: 1 coluna (padrão)
- Desktop: até 3 colunas

**Touch Optimization:**
- Botões mínimo 44x44px (Apple HIG)
- Espaçamento mínimo 8px entre elementos clicáveis
- Forms com teclado otimizado (type="email", type="tel")

### Requisitos Não-Funcionais

#### Performance
- Editor carrega em < 3 segundos
- Drag-and-drop responde em < 100ms
- Auto-save não bloqueia UI

#### Usabilidade
- Usuário consegue criar página básica em < 10 minutos (primeira vez)
- Usuário consegue criar página em < 5 minutos (após familiarização)
- Taxa de erro < 5% (usuário não consegue completar ação desejada)

#### Compatibilidade
- Chrome 90+ (desktop e mobile)
- Safari 14+ (desktop e mobile)
- Firefox 88+
- Edge 90+

### Acceptance Criteria

#### AC1: Criar Página em Branco
```gherkin
DADO que estou logado no sistema
QUANDO clico em "Nova Página"
E seleciono tipo "Landing Page"
ENTÃO uma página em branco é criada
E o editor é aberto
E vejo o canvas vazio com mensagem "Arraste componentes aqui"
```

#### AC2: Adicionar Componente de Texto
```gherkin
DADO que estou no editor
QUANDO arrasto componente "Título" da biblioteca
E solto no canvas
ENTÃO o título aparece no canvas
E está selecionado (mostra borda azul)
E painel de propriedades abre à direita
E posso editar o texto clicando duas vezes
```

#### AC3: Editar Propriedades de Componente
```gherkin
DADO que tenho um título no canvas
QUANDO clico no título
E altero cor no painel de propriedades para vermelho (#FF0000)
ENTÃO o título no canvas muda para vermelho imediatamente
E a mudança é salva automaticamente em 30 segundos
```

#### AC4: Preview Mobile/Desktop
```gherkin
DADO que estou editando uma página
QUANDO clico no botão "Mobile" na barra superior
ENTÃO o canvas mostra preview mobile (375px width)
QUANDO clico no botão "Desktop"
ENTÃO o canvas mostra preview desktop (1200px width)
E todos os componentes se adaptam responsivamente
```

#### AC5: Salvar e Publicar
```gherkin
DADO que criei uma página com pelo menos 1 componente
QUANDO clico em "Salvar"
ENTÃO vejo mensagem "Página salva com sucesso"
QUANDO clico em "Publicar"
ENTÃO página fica disponível na URL pública
E vejo mensagem "Página publicada em [URL]"
```

#### AC6: Undo/Redo
```gherkin
DADO que adicionei um título
E adicionei um parágrafo
QUANDO clico em "Undo" (Ctrl+Z)
ENTÃO o parágrafo é removido
QUANDO clico em "Redo" (Ctrl+Y)
ENTÃO o parágrafo volta
```

#### AC7: Mobile-First Behavior
```gherkin
DADO que estou no preview mobile
QUANDO adiciono um componente de 2 colunas
ENTÃO no mobile aparece como 1 coluna (stacked)
QUANDO mudo para preview desktop
ENTÃO aparece como 2 colunas lado a lado
```

### Casos de Teste

#### Teste 1: Fluxo Completo de Criação
1. Login no sistema
2. Criar nova Landing Page
3. Adicionar título "Curso de Nutrição"
4. Adicionar vídeo (embed YouTube)
5. Adicionar botão "Comprar Agora"
6. Configurar cor do botão (verde)
7. Preview mobile
8. Preview desktop
9. Salvar
10. Publicar

**Resultado Esperado:** Página publicada em < 10 minutos, sem erros

#### Teste 2: Performance com Muitos Componentes
1. Criar página
2. Adicionar 50 componentes (mix de texto, imagem, botão)
3. Medir tempo de resposta do drag-and-drop
4. Medir tempo de auto-save

**Resultado Esperado:** 
- Drag-and-drop < 100ms
- Auto-save < 2 segundos

#### Teste 3: Responsividade
1. Criar página com grid 3 colunas
2. Adicionar conteúdo em cada coluna
3. Preview mobile

**Resultado Esperado:** 
- Mobile mostra 1 coluna (stacked)
- Conteúdo legível e botões clicáveis

### Dependências Técnicas

**Frontend:**
- React 18+ (UI framework)
- React DnD ou dnd-kit (drag-and-drop)
- TailwindCSS (styling, mobile-first)
- Monaco Editor ou similar (para edição de texto inline)

**Backend:**
- API para salvar/carregar páginas
- Storage para uploads de imagem (S3 ou similar)

**Infraestrutura:**
- CDN para assets do editor (rápido carregamento)

### Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Editor lento com muitos componentes | Média | Alto | Virtualização de lista, lazy loading |
| Drag-and-drop não funciona em mobile | Baixa | Alto | Testar em dispositivos reais, fallback touch |
| Usuários não entendem como usar | Alta | Médio | Tutorial interativo, tooltips, vídeos |

### Métricas de Sucesso

**Quantitativas:**
- 80% dos usuários publicam primeira página em < 15 minutos
- Taxa de abandono do editor < 30%
- Tempo médio de criação de página < 10 minutos

**Qualitativas:**
- NPS do editor > 50
- Feedback positivo sobre facilidade de uso
- Usuários conseguem criar páginas "profissionais" sem ajuda

---

## Feature 2: Lógica de Funil com Proteção de URL

### Descrição

Sistema de gerenciamento de fluxo de funil que garante que usuários só possam acessar páginas na sequência correta, impedindo acesso direto via URL e retorno a páginas anteriores após avançar.

### Contexto e Motivação

**Feedback da Equipe (do Google Slides):**
- "Não ser possível acessar se não estiver no funil"
- "Quem estiver no funil não consegue retornar se não aceitar a proposta"

**Por que isso importa:**
- Protege a integridade da oferta (evita "gaming" do sistema)
- Mantém escassez artificial (one-time offer de verdade)
- Melhora conversão (usuário não pode "pensar melhor" voltando)

### User Stories

**Como** criador de funil  
**Quero** que minha oferta de upsell seja realmente "one-time"  
**Para que** o cliente não possa voltar e aceitar depois

**Como** sistema  
**Quero** impedir acesso direto a páginas de upsell via URL  
**Para que** apenas quem comprou o produto principal veja a oferta

### Requisitos Funcionais

#### 2.1 Sistema de Sessão de Funil

**Descrição:** Cada visitante que entra no funil recebe uma sessão única que rastreia seu progresso.

**Comportamento:**
1. Usuário acessa Landing Page → sessão criada (cookie + server-side)
2. Usuário avança para Página de Vendas → sessão atualizada (estado: "viewing_sales_page")
3. Usuário completa compra → sessão atualizada (estado: "purchased", produto_id)
4. Usuário vê Upsell → sessão atualizada (estado: "viewing_upsell")
5. Usuário aceita/recusa → sessão atualizada (estado: "upsell_accepted" ou "upsell_declined")

**Dados da Sessão:**
```json
{
  "session_id": "uuid",
  "funnel_id": "uuid",
  "user_ip": "192.168.1.1",
  "current_page": "upsell_1",
  "state": "viewing_upsell",
  "purchase_data": {
    "product_id": "uuid",
    "amount": 197.00,
    "payment_method": "pix",
    "payment_token": "encrypted_token"
  },
  "history": [
    {"page": "landing", "timestamp": "2025-12-31T10:00:00Z"},
    {"page": "sales", "timestamp": "2025-12-31T10:05:00Z"},
    {"page": "upsell_1", "timestamp": "2025-12-31T10:10:00Z"}
  ],
  "created_at": "2025-12-31T10:00:00Z",
  "expires_at": "2025-12-31T22:00:00Z"
}
```

#### 2.2 Middleware de Proteção de URL

**Descrição:** Middleware que valida se usuário pode acessar página solicitada baseado no estado da sessão.

**Regras de Acesso:**

| Página | Requer Estado | Ação se Inválido |
|--------|---------------|------------------|
| Landing Page | Nenhum | Permitir sempre |
| Página de Vendas | Visitou Landing OU acesso direto permitido | Redirecionar para Landing |
| Página de Upsell | Comprou produto principal | Redirecionar para Landing |
| Confirmação | Completou funil | Redirecionar para Landing |

**Comportamento de Retorno (Back Button):**
- Usuário em Página de Vendas tenta voltar → permitido (volta para Landing)
- Usuário em Upsell tenta voltar → **bloqueado** (mostra mensagem ou redireciona para frente)
- Usuário em Confirmação tenta voltar → bloqueado

**Implementação Técnica:**
```javascript
// Middleware Express.js (exemplo)
function funnelProtectionMiddleware(req, res, next) {
  const session = req.session.funnel;
  const requestedPage = req.params.page;
  
  // Validar se pode acessar
  if (!canAccessPage(session, requestedPage)) {
    return res.redirect('/landing');
  }
  
  // Prevenir retorno
  if (isGoingBackward(session, requestedPage)) {
    return res.redirect(session.current_page);
  }
  
  next();
}
```

#### 2.3 Prevenção de Retorno (Back Button)

**Técnicas:**

1. **History Manipulation (Frontend):**
```javascript
// Adicionar entrada fake no history
window.history.pushState(null, '', window.location.href);

// Detectar back button
window.addEventListener('popstate', function(event) {
  // Impedir retorno
  window.history.pushState(null, '', window.location.href);
  
  // Mostrar modal
  showModal('Você não pode voltar. Aceite ou recuse a oferta.');
});
```

2. **Server-Side Validation:**
- Cada página valida sessão no backend
- Se estado inválido → redireciona

3. **Session Expiration:**
- Sessão expira em 12 horas
- Após expiração, usuário deve começar do zero

#### 2.4 Acesso Direto via URL

**Cenário:** Usuário tenta acessar `https://funil.com/upsell-1` diretamente

**Comportamento:**
1. Sistema verifica sessão
2. Se sessão não existe ou estado inválido → redireciona para Landing Page
3. Se sessão válida mas não comprou ainda → redireciona para Página de Vendas
4. Se sessão válida e comprou → permite acesso

### Requisitos Não-Funcionais

#### Segurança
- Sessão armazenada server-side (não apenas cookie)
- Cookie com HttpOnly, Secure, SameSite=Strict
- Sessão expira após 12 horas de inatividade
- Proteção contra session hijacking (validar IP + User-Agent)

#### Performance
- Validação de sessão < 50ms
- Cache de regras de acesso

### Acceptance Criteria

#### AC1: Criar Sessão ao Entrar no Funil
```gherkin
DADO que sou um visitante novo
QUANDO acesso a Landing Page do funil
ENTÃO uma sessão é criada
E recebo um cookie "funnel_session_id"
E sessão é salva no servidor com estado "viewing_landing"
```

#### AC2: Bloquear Acesso Direto a Upsell
```gherkin
DADO que NÃO tenho sessão ativa
QUANDO tento acessar URL do Upsell diretamente
ENTÃO sou redirecionado para Landing Page
E vejo mensagem "Acesse o funil desde o início"
```

#### AC3: Permitir Acesso a Upsell Após Compra
```gherkin
DADO que completei compra na Página de Vendas
E minha sessão tem estado "purchased"
QUANDO sou redirecionado para Upsell
ENTÃO a página carrega normalmente
E vejo a oferta de upsell
```

#### AC4: Bloquear Retorno de Upsell para Página de Vendas
```gherkin
DADO que estou na página de Upsell
QUANDO clico no botão "Voltar" do navegador
ENTÃO permaneço na página de Upsell
E vejo modal "Você não pode voltar. Aceite ou recuse a oferta."
```

#### AC5: Expirar Sessão Após 12 Horas
```gherkin
DADO que criei uma sessão há 12 horas
QUANDO tento acessar qualquer página do funil
ENTÃO sou redirecionado para Landing Page
E uma nova sessão é criada
```

#### AC6: Prevenir Session Hijacking
```gherkin
DADO que tenho uma sessão ativa no IP 192.168.1.1
QUANDO alguém tenta usar meu session_id de outro IP
ENTÃO o acesso é negado
E sessão é invalidada
E ambos são redirecionados para Landing Page
```

### Casos de Teste

#### Teste 1: Fluxo Normal
1. Acessar Landing Page → sessão criada
2. Clicar em "Comprar" → vai para Página de Vendas
3. Completar compra → vai para Upsell
4. Aceitar upsell → vai para Confirmação
5. Tentar voltar → bloqueado

**Resultado Esperado:** Fluxo completo sem possibilidade de retorno

#### Teste 2: Acesso Direto Malicioso
1. Copiar URL do Upsell
2. Abrir em aba anônima (sem sessão)
3. Tentar acessar

**Resultado Esperado:** Redirecionado para Landing Page

#### Teste 3: Session Hijacking
1. Usuário A cria sessão
2. Usuário B rouba cookie
3. Usuário B tenta acessar de outro IP

**Resultado Esperado:** Acesso negado, sessão invalidada

### Dependências Técnicas

**Backend:**
- Redis ou similar (armazenamento de sessão rápido)
- Express-session ou similar (gerenciamento de sessão)

**Frontend:**
- JavaScript para prevenir back button
- Modal/Toast para mensagens de bloqueio

### Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Usuário limpa cookies e perde sessão | Média | Médio | Permitir recuperação via email/telefone |
| Back button prevention não funciona em todos browsers | Baixa | Médio | Validação server-side como backup |
| Falso positivo de session hijacking (IP dinâmico) | Média | Baixo | Usar fingerprinting mais sofisticado |

### Métricas de Sucesso

**Quantitativas:**
- 0% de acessos diretos bem-sucedidos a páginas protegidas
- < 1% de sessões perdidas por expiração prematura
- 0 casos de session hijacking

**Qualitativas:**
- Usuários entendem que oferta é "one-time"
- Criadores de funil confiam na proteção

---

## Feature 3: Checkout Inline com PIX

### Descrição

Sistema de checkout integrado diretamente na Página de Vendas, sem redirecionar para página externa, com suporte a pagamento via PIX através do Mercado Pago.

### Contexto e Motivação

**Feedback da Equipe (do Google Slides):**
- "Checkout dentro da Pg de Vendas, para não enviar o usuário para uma nova página"
- Dúvida: "pq os sites jogam as compras para uma página que contém o checkout de pagamento? É uma questão de segurança?"

**Nossa Decisão:** Checkout inline reduz fricção e melhora conversão, mantendo contexto da oferta.

**Por que PIX:**
- Método de pagamento mais rápido no Brasil
- Confirmação instantânea
- Menor taxa (0.5-1% vs. 3-5% cartão)
- Adoção crescente (40-50% das transações)

### User Stories

**Como** comprador  
**Quero** pagar sem sair da página de vendas  
**Para que** eu não perca o contexto da oferta

**Como** comprador  
**Quero** pagar via PIX  
**Para que** a compra seja confirmada instantaneamente

**Como** vendedor  
**Quero** que o checkout seja rápido e sem fricção  
**Para que** eu maximize minha taxa de conversão

### Requisitos Funcionais

#### 3.1 Interface de Checkout Inline

**Localização:** Seção na Página de Vendas (geralmente após descrição do produto)

**Layout (Mobile-First):**
```
┌─────────────────────────────────┐
│  [Imagem do Produto]            │
│                                 │
│  GARANTIA 7 DIAS                │
│                                 │
│  ┌─────────────────────────┐   │
│  │ CHECKOUT                │   │
│  │                         │   │
│  │ Nome: ________________  │   │
│  │                         │   │
│  │ Email: _______________  │   │
│  │                         │   │
│  │ Telefone: ____________  │   │
│  │                         │   │
│  │ ○ PIX (Instantâneo)     │   │
│  │                         │   │
│  │ Total: R$ 197,00        │   │
│  │                         │   │
│  │ [ FINALIZAR COMPRA ]    │   │
│  └─────────────────────────┘   │
│                                 │
│  🔒 Compra 100% Segura          │
└─────────────────────────────────┘
```

**Campos do Formulário:**
1. **Nome Completo** (obrigatório)
   - Validação: mínimo 3 caracteres
   - Placeholder: "Seu nome completo"

2. **Email** (obrigatório)
   - Validação: formato email válido
   - Placeholder: "seu@email.com"
   - Type: "email" (teclado otimizado mobile)

3. **Telefone** (obrigatório)
   - Validação: formato brasileiro (11) 99999-9999
   - Máscara automática
   - Placeholder: "(11) 99999-9999"
   - Type: "tel" (teclado numérico mobile)

4. **Método de Pagamento** (MVP: apenas PIX)
   - Radio button: "PIX (Pagamento Instantâneo)"
   - Ícone do PIX
   - Badge: "Confirmação em segundos"

**Botão de Finalizar:**
- Texto: "FINALIZAR COMPRA"
- Cor: Verde (call-to-action)
- Tamanho: Grande (mínimo 44px altura)
- Largura: 100% (mobile)
- Estado loading: Spinner + "Processando..."

#### 3.2 Fluxo de Pagamento PIX

**Passo 1: Usuário Preenche Formulário**
- Validação em tempo real (ao sair do campo)
- Mensagens de erro abaixo do campo
- Botão desabilitado se form inválido

**Passo 2: Usuário Clica "Finalizar Compra"**
- Botão entra em estado loading
- Requisição para backend: criar pedido

**Passo 3: Backend Cria Pedido no Mercado Pago**
```javascript
// Exemplo de chamada API Mercado Pago
const payment = await mercadopago.payment.create({
  transaction_amount: 197.00,
  description: "Curso de Nutrição",
  payment_method_id: "pix",
  payer: {
    email: "cliente@email.com",
    first_name: "João",
    last_name: "Silva"
  }
});

// Retorna QR Code e código PIX
return {
  qr_code_base64: payment.point_of_interaction.transaction_data.qr_code_base64,
  qr_code: payment.point_of_interaction.transaction_data.qr_code,
  payment_id: payment.id
};
```

**Passo 4: Exibir QR Code PIX**

**Layout:**
```
┌─────────────────────────────────┐
│  PAGAMENTO VIA PIX              │
│                                 │
│  ┌─────────────────────────┐   │
│  │                         │   │
│  │    [QR CODE IMAGE]      │   │
│  │                         │   │
│  └─────────────────────────┘   │
│                                 │
│  Escaneie o QR Code com o       │
│  app do seu banco               │
│                                 │
│  OU                             │
│                                 │
│  PIX Copia e Cola:              │
│  ┌─────────────────────────┐   │
│  │ 00020126580014br.gov... │   │
│  │ [ COPIAR CÓDIGO ]       │   │
│  └─────────────────────────┘   │
│                                 │
│  ⏱️ Aguardando pagamento...     │
│  (atualiza automaticamente)     │
│                                 │
│  Valor: R$ 197,00               │
│  Válido por: 15 minutos         │
└─────────────────────────────────┘
```

**Funcionalidades:**
- QR Code exibido como imagem (base64)
- Código PIX em campo de texto (readonly)
- Botão "Copiar Código" (clipboard API)
- Timer de expiração (15 minutos countdown)
- Polling a cada 3 segundos para verificar pagamento

**Passo 5: Polling de Status**

**Frontend:**
```javascript
// Verificar status a cada 3 segundos
const pollPaymentStatus = async (paymentId) => {
  const interval = setInterval(async () => {
    const status = await fetch(`/api/payment/${paymentId}/status`);
    const data = await status.json();
    
    if (data.status === 'approved') {
      clearInterval(interval);
      // Redirecionar para Upsell
      window.location.href = '/upsell-1';
    }
    
    if (data.status === 'rejected' || data.status === 'expired') {
      clearInterval(interval);
      showError('Pagamento não confirmado. Tente novamente.');
    }
  }, 3000);
};
```

**Backend:**
- Webhook do Mercado Pago notifica quando pagamento aprovado
- Atualizar status do pedido no banco de dados
- API `/payment/:id/status` retorna status atual

**Passo 6: Pagamento Confirmado**
- Mostrar mensagem "Pagamento Confirmado! ✅"
- Salvar dados de pagamento na sessão (para one-click upsell)
- Redirecionar para Página de Upsell em 2 segundos

#### 3.3 Integração com Mercado Pago

**Credenciais:**
- Access Token (variável de ambiente)
- Public Key (para frontend, se necessário)

**Endpoints Usados:**
1. **POST /v1/payments** - Criar pagamento PIX
2. **GET /v1/payments/:id** - Consultar status
3. **Webhook** - Receber notificações de mudança de status

**Webhook Configuration:**
```javascript
// Rota para receber webhook
app.post('/webhooks/mercadopago', async (req, res) => {
  const { type, data } = req.body;
  
  if (type === 'payment') {
    const paymentId = data.id;
    
    // Buscar detalhes do pagamento
    const payment = await mercadopago.payment.get(paymentId);
    
    if (payment.status === 'approved') {
      // Atualizar pedido no banco
      await updateOrder(payment.external_reference, {
        status: 'paid',
        payment_id: paymentId
      });
      
      // Enviar email de confirmação (opcional para MVP)
    }
  }
  
  res.sendStatus(200);
});
```

#### 3.4 Tratamento de Erros

**Erros Possíveis:**

1. **Erro de Validação de Formulário**
   - Mensagem: "Por favor, preencha todos os campos corretamente"
   - Destacar campos com erro em vermelho

2. **Erro ao Criar Pagamento (API Mercado Pago)**
   - Mensagem: "Erro ao processar pagamento. Tente novamente."
   - Botão: "Tentar Novamente"
   - Log do erro no backend

3. **Timeout de Pagamento (15 minutos)**
   - Mensagem: "Tempo expirado. Gere um novo código PIX."
   - Botão: "Gerar Novo PIX"

4. **Pagamento Rejeitado**
   - Mensagem: "Pagamento não confirmado. Verifique com seu banco."
   - Botão: "Tentar Novamente"

### Requisitos Não-Funcionais

#### Segurança
- Comunicação HTTPS obrigatória
- Access Token do Mercado Pago em variável de ambiente (nunca no código)
- Validação de webhook (verificar assinatura do Mercado Pago)
- Dados do cliente criptografados em trânsito e em repouso

#### Performance
- Criação de pagamento PIX < 2 segundos
- QR Code exibido < 1 segundo após criação
- Polling não bloqueia UI

#### Disponibilidade
- Fallback se Mercado Pago estiver fora: mostrar mensagem de manutenção
- Retry automático (3 tentativas) se API falhar

### Acceptance Criteria

#### AC1: Exibir Formulário de Checkout
```gherkin
DADO que estou na Página de Vendas
QUANDO rolo até a seção de checkout
ENTÃO vejo formulário com campos Nome, Email, Telefone
E vejo opção de pagamento PIX selecionada
E vejo botão "FINALIZAR COMPRA"
```

#### AC2: Validar Formulário
```gherkin
DADO que estou no formulário de checkout
QUANDO deixo campo Email vazio
E clico em "FINALIZAR COMPRA"
ENTÃO vejo mensagem "Email é obrigatório"
E botão permanece desabilitado
```

#### AC3: Gerar QR Code PIX
```gherkin
DADO que preenchi formulário corretamente
QUANDO clico em "FINALIZAR COMPRA"
ENTÃO vejo spinner "Processando..."
E em até 2 segundos vejo QR Code PIX
E vejo código PIX copia-e-cola
E vejo timer "Válido por 15:00"
```

#### AC4: Copiar Código PIX
```gherkin
DADO que vejo código PIX
QUANDO clico em "COPIAR CÓDIGO"
ENTÃO código é copiado para clipboard
E vejo mensagem "Código copiado!"
```

#### AC5: Detectar Pagamento Aprovado
```gherkin
DADO que gerei QR Code PIX
QUANDO pago via app do banco
ENTÃO em até 10 segundos vejo "Pagamento Confirmado! ✅"
E sou redirecionado para Página de Upsell
```

#### AC6: Expirar Pagamento Após 15 Minutos
```gherkin
DADO que gerei QR Code PIX
QUANDO passam 15 minutos sem pagar
ENTÃO vejo mensagem "Tempo expirado"
E vejo botão "Gerar Novo PIX"
```

### Casos de Teste

#### Teste 1: Fluxo Completo de Compra
1. Acessar Página de Vendas
2. Preencher formulário (Nome, Email, Telefone)
3. Clicar "FINALIZAR COMPRA"
4. Ver QR Code
5. Pagar via PIX (ambiente de testes Mercado Pago)
6. Ver confirmação
7. Ser redirecionado para Upsell

**Resultado Esperado:** Fluxo completo em < 30 segundos

#### Teste 2: Erro de API
1. Desligar API do Mercado Pago (simular)
2. Tentar finalizar compra
3. Ver mensagem de erro
4. Clicar "Tentar Novamente"

**Resultado Esperado:** Erro tratado graciosamente, retry funciona

#### Teste 3: Timeout de Pagamento
1. Gerar QR Code
2. Aguardar 15 minutos sem pagar
3. Ver expiração

**Resultado Esperado:** Timer chega a 00:00, mostra botão "Gerar Novo PIX"

### Dependências Técnicas

**Backend:**
- Mercado Pago SDK (Node.js ou Python)
- Banco de dados para armazenar pedidos
- Webhook endpoint público (ngrok para dev, domínio real para prod)

**Frontend:**
- Biblioteca de validação de formulário (Formik + Yup ou similar)
- QR Code como imagem (base64 do Mercado Pago)
- Clipboard API (para copiar código PIX)
- Polling (setInterval ou React Query)

**Infraestrutura:**
- SSL certificate (obrigatório para PIX)
- Webhook deve ser acessível publicamente

### Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| Mercado Pago API fora do ar | Baixa | Alto | Retry automático, mensagem clara ao usuário |
| Webhook não chega (firewall, etc.) | Média | Alto | Polling como backup, logs robustos |
| Usuário não sabe usar PIX | Baixa | Médio | Instruções claras, vídeo tutorial |
| Pagamento demora mais que 10s | Média | Baixo | Aumentar timeout de polling, mostrar mensagem "Pode demorar até 1 minuto" |

### Métricas de Sucesso

**Quantitativas:**
- Taxa de conversão do checkout > 65%
- Tempo médio de pagamento < 30 segundos
- Taxa de erro < 2%

**Qualitativas:**
- Usuários acham checkout "fácil" e "rápido"
- NPS do checkout > 70

---

## Feature 4: One-Click Upsell

### Descrição

Página de oferta adicional (upsell) que aparece imediatamente após a compra do produto principal, permitindo que o cliente adicione o produto ao pedido com um único clique, sem reinserir dados de pagamento.

### Contexto e Motivação

**Feedback da Equipe:**
- "One-click purchase sem inserir dados novamente" (Killer Feature do Google Slides)
- Dante: "várias páginas depois da pagina de vendas, por exemplo, como upsells, downsells etc...td na mesma lógica de one-time-offer"

**Por que isso é diferencial:**
- Hotmart tem order bump (checkbox no checkout)
- Nossa solução: página dedicada, experiência fluida, sem fricção
- Aumenta AOV em 20-40% (benchmark de mercado)

### User Stories

**Como** comprador que acabou de comprar  
**Quero** adicionar produto complementar com um clique  
**Para que** eu não precise preencher dados novamente

**Como** vendedor  
**Quero** oferecer upsell sem fricção  
**Para que** eu aumente meu ticket médio

### Requisitos Funcionais

#### 4.1 Interface da Página de Upsell

**Quando Aparece:** Imediatamente após pagamento confirmado do produto principal

**Layout (Mobile-First):**
```
┌─────────────────────────────────┐
│  🎉 PARABÉNS!                   │
│  Sua compra foi confirmada!     │
│                                 │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│                                 │
│  OFERTA ESPECIAL SÓ PARA VOCÊ   │
│                                 │
│  [Imagem do Produto Upsell]     │
│                                 │
│  E-book de Receitas Saudáveis   │
│                                 │
│  ✓ 50 receitas exclusivas       │
│  ✓ Plano alimentar 30 dias      │
│  ✓ Lista de compras pronta      │
│                                 │
│  De R$ 97 por apenas            │
│  R$ 47                          │
│  (51% OFF - só hoje!)           │
│                                 │
│  [ SIM, ADICIONAR AO PEDIDO ]   │
│  (pagamento instantâneo)        │
│                                 │
│  [ Não, obrigado ]              │
│                                 │
│  ⏱️ Oferta expira em 10:00      │
└─────────────────────────────────┘
```

**Elementos Obrigatórios:**
1. **Mensagem de Parabéns** - Confirmar compra anterior
2. **Imagem do Produto** - Visual atrativo
3. **Título e Descrição** - Benefícios claros (bullets)
4. **Preço** - Mostrar desconto (De X por Y)
5. **Botão de Aceitar** - CTA forte, verde
6. **Botão de Recusar** - Discreto, link ou botão secundário
7. **Timer de Escassez** - Countdown (opcional mas recomendado)

#### 4.2 Fluxo de Aceitação (One-Click)

**Passo 1: Usuário Clica "SIM, ADICIONAR AO PEDIDO"**
- Botão entra em estado loading
- Texto muda para "Processando..."

**Passo 2: Backend Processa Upsell**

**Dados Necessários (já salvos na sessão):**
```json
{
  "customer": {
    "name": "João Silva",
    "email": "joao@email.com",
    "phone": "(11) 99999-9999"
  },
  "payment": {
    "method": "pix",
    "token": "encrypted_payment_token",
    "payer_id": "mercado_pago_payer_id"
  },
  "original_purchase": {
    "order_id": "uuid",
    "amount": 197.00
  }
}
```

**Lógica de Processamento:**
```javascript
async function processUpsell(sessionData, upsellProduct) {
  // 1. Criar pagamento no Mercado Pago usando dados salvos
  const payment = await mercadopago.payment.create({
    transaction_amount: upsellProduct.price,
    description: upsellProduct.name,
    payment_method_id: "pix",
    payer: {
      id: sessionData.payment.payer_id,
      email: sessionData.customer.email
    },
    // Usar token salvo para one-click
    token: sessionData.payment.token
  });
  
  // 2. Adicionar ao pedido original
  await addItemToOrder(sessionData.original_purchase.order_id, {
    product_id: upsellProduct.id,
    amount: upsellProduct.price,
    payment_id: payment.id
  });
  
  // 3. Atualizar sessão
  sessionData.upsell_accepted = true;
  sessionData.total_amount += upsellProduct.price;
  
  return { success: true, order_id: sessionData.original_purchase.order_id };
}
```

**Passo 3: Confirmação Instantânea**
- Mostrar mensagem "Adicionado com sucesso! ✅"
- Redirecionar para Página de Confirmação em 2 segundos

**Nota Técnica - Limitação do PIX:**
PIX não suporta "tokenização" tradicional como cartão de crédito. Para MVP, temos duas opções:

**Opção A (Recomendada para MVP):** Gerar novo QR Code PIX
- Usuário clica "SIM"
- Sistema gera novo QR Code para valor do upsell
- Usuário paga novamente via PIX (rápido, já está com app aberto)
- Ainda é "one-click" no sentido de não preencher formulário

**Opção B (V1+):** Integrar cartão de crédito
- Cartão permite tokenização real
- Upsell é cobrado automaticamente no cartão salvo
- Verdadeiro "one-click" sem interação adicional

**Para MVP, usar Opção A** (PIX com novo QR Code, mas sem formulário)

#### 4.3 Fluxo de Recusa

**Passo 1: Usuário Clica "Não, obrigado"**
- Registrar recusa no analytics
- Atualizar sessão (upsell_declined = true)

**Passo 2: (Opcional) Downsell**
- Se configurado, mostrar oferta menor
- Ex: "Que tal só o plano alimentar por R$ 27?"
- Mesma lógica de one-click

**Passo 3: Redirecionar para Confirmação**
- Se não há downsell ou usuário recusou downsell
- Ir para Página de Confirmação

#### 4.4 Timer de Escassez

**Funcionalidade:** Countdown de 10 minutos para criar urgência

**Implementação:**
```javascript
// Countdown timer
let timeLeft = 600; // 10 minutos em segundos

const timer = setInterval(() => {
  timeLeft--;
  
  const minutes = Math.floor(timeLeft / 60);
  const seconds = timeLeft % 60;
  
  document.getElementById('timer').textContent = 
    `${minutes}:${seconds.toString().padStart(2, '0')}`;
  
  if (timeLeft <= 0) {
    clearInterval(timer);
    // Redirecionar para Confirmação (oferta expirou)
    window.location.href = '/confirmacao';
  }
}, 1000);
```

**Comportamento ao Expirar:**
- Timer chega a 00:00
- Página redireciona automaticamente para Confirmação
- Sessão marcada como "upsell_expired"

#### 4.5 Proteção contra Acesso Direto

**Validação:** Usuário só pode acessar Upsell se:
1. Completou compra do produto principal
2. Sessão tem payment_token salvo
3. Não aceitou nem recusou upsell ainda

**Se validação falhar:** Redirecionar para Landing Page

### Requisitos Não-Funcionais

#### Performance
- Página carrega em < 1 segundo (usuário já está "quente")
- Processamento de upsell < 3 segundos
- Timer não trava UI

#### Conversão
- Taxa de aceitação de upsell > 15% (benchmark)
- Aumento de AOV > 20%

### Acceptance Criteria

#### AC1: Exibir Página de Upsell Após Compra
```gherkin
DADO que completei compra do produto principal
QUANDO pagamento é confirmado
ENTÃO sou redirecionado para Página de Upsell
E vejo mensagem "PARABÉNS! Sua compra foi confirmada"
E vejo oferta de upsell com imagem e preço
```

#### AC2: Aceitar Upsell com One-Click
```gherkin
DADO que estou na Página de Upsell
QUANDO clico em "SIM, ADICIONAR AO PEDIDO"
ENTÃO vejo spinner "Processando..."
E em até 3 segundos vejo "Adicionado com sucesso! ✅"
E sou redirecionado para Confirmação
E NÃO preciso preencher dados novamente
```

#### AC3: Recusar Upsell
```gherkin
DADO que estou na Página de Upsell
QUANDO clico em "Não, obrigado"
ENTÃO sou redirecionado para Confirmação imediatamente
E meu pedido contém apenas o produto principal
```

#### AC4: Timer de Escassez
```gherkin
DADO que estou na Página de Upsell
QUANDO a página carrega
ENTÃO vejo timer "Oferta expira em 10:00"
E timer decrementa a cada segundo
QUANDO timer chega a 00:00
ENTÃO sou redirecionado para Confirmação automaticamente
```

#### AC5: Bloquear Acesso Direto
```gherkin
DADO que NÃO completei compra do produto principal
QUANDO tento acessar URL do Upsell diretamente
ENTÃO sou redirecionado para Landing Page
E vejo mensagem "Acesse o funil desde o início"
```

### Casos de Teste

#### Teste 1: Fluxo Completo de Upsell Aceito
1. Completar compra do produto principal (R$ 197)
2. Ver página de Upsell
3. Clicar "SIM, ADICIONAR"
4. Ver confirmação
5. Ir para Confirmação
6. Verificar pedido total = R$ 197 + R$ 47 = R$ 244

**Resultado Esperado:** Upsell adicionado, total correto

#### Teste 2: Fluxo de Upsell Recusado
1. Completar compra do produto principal
2. Ver página de Upsell
3. Clicar "Não, obrigado"
4. Ir para Confirmação
5. Verificar pedido total = R$ 197 (apenas produto principal)

**Resultado Esperado:** Upsell não adicionado

#### Teste 3: Expiração do Timer
1. Ver página de Upsell
2. Aguardar 10 minutos (ou simular com timer acelerado)
3. Ver redirecionamento automático

**Resultado Esperado:** Timer expira, redireciona para Confirmação

### Dependências Técnicas

**Backend:**
- Sistema de sessão (salvar payment_token)
- API para processar upsell
- Integração com Mercado Pago

**Frontend:**
- Timer (JavaScript)
- Loading states
- Redirecionamento automático

### Riscos e Mitigações

| Risco | Probabilidade | Impacto | Mitigação |
|-------|---------------|---------|-----------|
| PIX não permite tokenização real | Alta | Médio | Usar novo QR Code (ainda sem formulário) |
| Usuário não entende que não precisa preencher dados | Média | Baixo | Texto claro "pagamento instantâneo" |
| Erro ao processar upsell | Baixa | Alto | Retry automático, fallback para formulário |

### Métricas de Sucesso

**Quantitativas:**
- Taxa de aceitação de upsell > 15%
- Aumento de AOV > 20%
- Tempo médio de decisão < 30 segundos

**Qualitativas:**
- Usuários acham processo "rápido" e "sem fricção"
- NPS do upsell > 60

---

## Feature 5: Analytics Básico do Funil

### Descrição

Dashboard simples mostrando métricas essenciais de performance do funil: visitantes, conversões, receita, e taxa de aceitação de upsell.

### Contexto e Motivação

**Do Google Slides:** "Analytics completo do funil" é feature indispensável

**Para MVP:** Focar em métricas essenciais, não analytics avançado

### User Stories

**Como** criador de funil  
**Quero** ver quantas pessoas estão convertendo  
**Para que** eu saiba se meu funil está funcionando

### Requisitos Funcionais

#### 5.1 Métricas Rastreadas

**Eventos:**
1. `landing_page_view` - Visitou Landing Page
2. `sales_page_view` - Visitou Página de Vendas
3. `checkout_started` - Iniciou checkout (preencheu formulário)
4. `purchase_completed` - Completou compra
5. `upsell_view` - Viu página de Upsell
6. `upsell_accepted` - Aceitou upsell
7. `upsell_declined` - Recusou upsell

**Propriedades de Cada Evento:**
- `session_id` - ID da sessão
- `funnel_id` - ID do funil
- `timestamp` - Data/hora
- `device` - Mobile ou Desktop
- `source` - UTM source (se disponível)

#### 5.2 Dashboard

**Métricas Exibidas:**

1. **Visitantes Únicos**
   - Total de sessões únicas
   - Filtro: hoje, últimos 7 dias, últimos 30 dias

2. **Taxa de Conversão (Landing → Compra)**
   - Fórmula: (purchase_completed / landing_page_view) * 100
   - Exibir em %

3. **Receita Total**
   - Soma de todos os pagamentos aprovados
   - Exibir em R$

4. **Ticket Médio**
   - Fórmula: receita_total / número_de_compras
   - Exibir em R$

5. **Taxa de Aceitação de Upsell**
   - Fórmula: (upsell_accepted / upsell_view) * 100
   - Exibir em %

6. **Funil de Conversão (Visual)**
```
Landing Page:     1000 visitantes (100%)
                    ↓
Página de Vendas:  500 visitantes (50%)
                    ↓
Compras:           100 compras (10% do total)
                    ↓
Upsell Visto:      100 visualizações (100% dos compradores)
                    ↓
Upsell Aceito:     20 aceitações (20% dos que viram)
```

**Layout (Simplificado):**
```
┌─────────────────────────────────────────┐
│  ANALYTICS - Últimos 7 dias             │
│                                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐  │
│  │ 1,234   │ │ 10.5%   │ │ R$ 24K  │  │
│  │ Visitas │ │ Conversão│ │ Receita │  │
│  └─────────┘ └─────────┘ └─────────┘  │
│                                         │
│  FUNIL DE CONVERSÃO                     │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━  │
│  Landing:        1,234 (100%)           │
│  Vendas:           618 (50%)            │
│  Compras:          130 (10.5%)          │
│  Upsell Aceito:     26 (20%)            │
│                                         │
│  TICKET MÉDIO: R$ 184                   │
│  (R$ 197 produto + R$ 47 upsell × 20%)  │
└─────────────────────────────────────────┘
```

#### 5.3 Rastreamento de Eventos

**Implementação (Frontend):**
```javascript
// Função para rastrear evento
function trackEvent(eventName, properties = {}) {
  fetch('/api/analytics/track', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      event: eventName,
      session_id: getSessionId(),
      funnel_id: getFunnelId(),
      timestamp: new Date().toISOString(),
      device: isMobile() ? 'mobile' : 'desktop',
      ...properties
    })
  });
}

// Exemplos de uso
trackEvent('landing_page_view');
trackEvent('purchase_completed', { amount: 197.00, product_id: 'abc123' });
trackEvent('upsell_accepted', { amount: 47.00, product_id: 'xyz789' });
```

**Backend:**
- Salvar eventos em banco de dados (tabela `analytics_events`)
- Agregar dados para dashboard (pode ser calculado em tempo real ou pré-agregado)

### Requisitos Não-Funcionais

#### Performance
- Dashboard carrega em < 2 segundos
- Rastreamento de evento não bloqueia UI (assíncrono)

#### Precisão
- 100% dos eventos críticos rastreados (compra, upsell)
- Tolerância de 5% de perda em eventos não-críticos (views)

### Acceptance Criteria

#### AC1: Rastrear Visualização de Página
```gherkin
DADO que acesso a Landing Page
QUANDO a página carrega
ENTÃO evento "landing_page_view" é registrado
E aparece no dashboard em até 1 minuto
```

#### AC2: Exibir Métricas no Dashboard
```gherkin
DADO que tenho 100 visitantes e 10 compras
QUANDO acesso o dashboard
ENTÃO vejo "100 Visitas"
E vejo "10% Taxa de Conversão"
E vejo "R$ 1.970 Receita" (10 × R$ 197)
```

#### AC3: Calcular Taxa de Upsell
```gherkin
DADO que 10 pessoas viram upsell
E 2 aceitaram
QUANDO acesso o dashboard
ENTÃO vejo "20% Taxa de Aceitação de Upsell"
```

### Casos de Teste

#### Teste 1: Fluxo Completo de Rastreamento
1. Visitar Landing Page → verificar evento
2. Ir para Página de Vendas → verificar evento
3. Completar compra → verificar evento
4. Ver Upsell → verificar evento
5. Aceitar Upsell → verificar evento
6. Abrir dashboard → verificar todas as métricas

**Resultado Esperado:** Todos os eventos rastreados, métricas corretas

### Dependências Técnicas

**Backend:**
- Banco de dados para eventos (PostgreSQL ou similar)
- API para rastreamento e consulta

**Frontend:**
- Biblioteca de charts (Chart.js ou similar) para visualização do funil

### Métricas de Sucesso

**Quantitativas:**
- 100% de eventos críticos rastreados
- Dashboard acessado por 80%+ dos criadores de funil

---

## Feature 6: Custom Domain & SSL

### Descrição

Permitir que usuários conectem seus próprios domínios personalizados (ex: `oferta.seusite.com.br`) ao funil, com SSL automático.

### Contexto e Motivação

**Importância:** Domínio personalizado aumenta credibilidade e confiança do comprador.

### User Stories

**Como** criador de funil  
**Quero** usar meu próprio domínio  
**Para que** meu funil pareça mais profissional

### Requisitos Funcionais

#### 6.1 Configuração de Domínio

**Fluxo:**
1. Usuário acessa "Configurações" do funil
2. Clica em "Domínio Personalizado"
3. Insere domínio (ex: `oferta.meusite.com.br`)
4. Sistema mostra instruções de DNS:
   - Adicionar registro CNAME apontando para `funnels.funilrapido.com.br`
5. Usuário configura DNS no provedor (Registro.br, Cloudflare, etc.)
6. Sistema verifica DNS (pode demorar até 48h)
7. Quando DNS propaga, SSL é gerado automaticamente (Let's Encrypt)
8. Domínio fica ativo

#### 6.2 SSL Automático

**Implementação:**
- Usar Let's Encrypt (gratuito)
- Renovação automática a cada 90 dias
- Certificado gerado via Certbot ou similar

### Acceptance Criteria

#### AC1: Adicionar Domínio Personalizado
```gherkin
DADO que tenho um funil criado
QUANDO acesso "Configurações" > "Domínio"
E insiro "oferta.meusite.com.br"
ENTÃO vejo instruções de DNS
E vejo status "Aguardando configuração DNS"
```

#### AC2: Ativar Domínio Após DNS Configurado
```gherkin
DADO que configurei DNS corretamente
QUANDO DNS propaga (até 48h)
ENTÃO sistema detecta automaticamente
E gera certificado SSL
E status muda para "Ativo ✅"
E funil fica acessível em https://oferta.meusite.com.br
```

### Casos de Teste

#### Teste 1: Configuração Completa
1. Adicionar domínio `teste.funilrapido.com.br`
2. Configurar CNAME no DNS
3. Aguardar propagação
4. Verificar SSL ativo
5. Acessar funil via domínio personalizado

**Resultado Esperado:** Funil acessível via HTTPS no domínio personalizado

### Dependências Técnicas

**Infraestrutura:**
- Servidor web (Nginx ou similar) com suporte a múltiplos domínios
- Certbot para Let's Encrypt
- Verificação de DNS (dig ou similar)

### Métricas de Sucesso

- 60%+ dos usuários conectam domínio personalizado
- 100% dos domínios têm SSL ativo

---

## Feature 7: Sistema de Templates

### Descrição

Biblioteca de 2-3 templates pré-construídos para acelerar criação de funis.

### Contexto e Motivação

**Feedback:** Usuários (especialmente Juliana) precisam de ponto de partida, não tela em branco.

### User Stories

**Como** usuário iniciante  
**Quero** começar com um template pronto  
**Para que** eu não precise criar tudo do zero

### Requisitos Funcionais

#### 7.1 Templates Disponíveis (MVP)

**Template 1: Infoproduto (Curso Online)**
- Landing Page: Vídeo de vendas + bullets + CTA
- Página de Vendas: Descrição longa + depoimentos + checkout
- Upsell: E-book complementar
- Confirmação: Instruções de acesso

**Template 2: E-commerce (Produto Físico)**
- Landing Page: Imagem do produto + benefícios + CTA
- Página de Vendas: Galeria de fotos + especificações + checkout
- Upsell: Produto complementar
- Confirmação: Rastreamento de envio

**Template 3: Serviço Local (Consultoria)**
- Landing Page: Formulário de contato + benefícios
- Página de Vendas: Sobre o profissional + casos de sucesso + checkout
- Upsell: Sessão extra com desconto
- Confirmação: Agendamento

#### 7.2 Usar Template

**Fluxo:**
1. Usuário clica "Novo Funil"
2. Escolhe template da galeria
3. Template é duplicado para conta do usuário
4. Usuário edita textos, imagens, cores
5. Publica

### Acceptance Criteria

#### AC1: Selecionar Template
```gherkin
DADO que clico em "Novo Funil"
QUANDO vejo galeria de templates
E clico em "Infoproduto"
ENTÃO template é criado na minha conta
E posso editá-lo no editor
```

#### AC2: Personalizar Template
```gherkin
DADO que selecionei template
QUANDO edito título para "Meu Curso"
E troco imagem
E mudo cor do botão
ENTÃO mudanças são salvas
E template personalizado está pronto para publicar
```

### Casos de Teste

#### Teste 1: Usar Template Completo
1. Criar funil a partir de template "Infoproduto"
2. Editar apenas textos (não layout)
3. Publicar
4. Testar funil completo (LP → Vendas → Upsell → Confirmação)

**Resultado Esperado:** Funil funcional em < 15 minutos

### Dependências Técnicas

**Backend:**
- Templates armazenados como JSON (estrutura de página)
- Função de duplicação de template

### Métricas de Sucesso

- 80%+ dos usuários começam com template (vs. tela em branco)
- Tempo de criação com template < 15 minutos

---

## 🎯 Priorização e Roadmap

### Sprint 1 (Semanas 1-2): Fundação
- Feature 2: Lógica de Funil com Proteção de URL
- Feature 6: Custom Domain & SSL (setup básico)

### Sprint 2 (Semanas 3-4): Editor
- Feature 1: Editor No-Code Mobile-First (versão básica)
- Feature 7: Sistema de Templates (2 templates)

### Sprint 3 (Semanas 5-6): Pagamentos
- Feature 3: Checkout Inline com PIX
- Feature 4: One-Click Upsell

### Sprint 4 (Semanas 7-8): Analytics e Polish
- Feature 5: Analytics Básico
- Testes end-to-end
- Bug fixes
- Documentação

---

## 📊 Critérios de Sucesso do MVP

### Métricas de Validação

**Objetivo:** Validar hipótese de que upsells fluidos aumentam AOV

**Critérios de Sucesso:**
1. ✅ 20 usuários ativos (criaram e publicaram funil)
2. ✅ Taxa de aceitação de upsell > 15%
3. ✅ Aumento de AOV > 20% (com upsell vs. sem upsell)
4. ✅ NPS > 50
5. ✅ 80% dos usuários publicam primeiro funil em < 30 minutos

**Se atingir critérios:** Escalar para V1+ (múltiplos upsells, parcelamento, etc.)

**Se não atingir:** Pivotar ou iterar baseado em feedback

---

## 🚨 Fora do Escopo do MVP

**Features que NÃO estão no MVP:**
- ❌ Múltiplos upsells/downsells em sequência (V1+)
- ❌ Parcelamento (cartão de crédito) (V1+)
- ❌ Boleto bancário (V1+)
- ❌ Integração com Hotmart (V1+)
- ❌ Email automation (V1+)
- ❌ A/B testing (V1+)
- ❌ Heatmaps (V2)
- ❌ White-label (Enterprise)
- ❌ Multi-client agency dashboard (V1+)

---

## 📝 Notas Técnicas

### Stack Recomendado

**Frontend:**
- React 18+ (UI)
- TailwindCSS (styling mobile-first)
- React DnD (drag-and-drop editor)
- React Query (data fetching)
- Zustand ou Context API (state management)

**Backend:**
- Node.js + Express (API)
- PostgreSQL (banco de dados)
- Redis (sessões)
- Mercado Pago SDK (pagamentos)

**Infraestrutura:**
- Vercel ou Railway (hosting)
- Cloudflare (CDN + DNS)
- AWS S3 (uploads de imagem)
- Let's Encrypt (SSL)

### Estimativa de Custos (MVP)

**Desenvolvimento:**
- 1 desenvolvedor full-stack × 8 semanas = ~R$ 40.000

**Infraestrutura (mensal):**
- Hosting: R$ 200
- Banco de dados: R$ 100
- CDN: R$ 50
- S3: R$ 50
- **Total:** ~R$ 400/mês

**Serviços:**
- Mercado Pago: 0.5-1% por transação (variável)

---

## ✅ Checklist de Entrega do MVP

### Funcional
- [ ] Editor no-code funciona em mobile e desktop
- [ ] Usuário consegue criar funil completo (LP → Vendas → Upsell → Confirmação)
- [ ] Checkout inline com PIX funciona
- [ ] One-click upsell funciona
- [ ] Proteção de URL impede acesso direto
- [ ] Analytics mostra métricas básicas
- [ ] Custom domain funciona com SSL

### Qualidade
- [ ] Performance: páginas carregam em < 2s
- [ ] Mobile-first: experiência perfeita em mobile
- [ ] Sem bugs críticos
- [ ] Testes end-to-end passando

### Documentação
- [ ] Guia de uso para usuários
- [ ] Documentação técnica para desenvolvedores
- [ ] Vídeos tutoriais (opcional para MVP)

### Validação
- [ ] 10 beta testers usaram e deram feedback
- [ ] Ajustes baseados em feedback implementados

---

**Versão:** 2.0  
**Última Atualização:** 31 de Dezembro de 2025  
**Status:** ✅ Ready for Development
