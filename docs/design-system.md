# Design System

O design system do Feature Graveyard AI foi pensado como cockpit operacional para times de plataforma, arquitetura e produto.

## Principios

- Denso, mas legivel: informacao de decisao acima de ornamentacao.
- Executivo e tecnico ao mesmo tempo: metricas claras, racional visivel e sugestao de acao.
- Estados fortes: feature morta, em risco e ativa precisam ser diferenciaveis em segundos.
- Layout de ferramenta: primeira tela ja e o produto funcionando.

## Tokens

- `--bg`: base quente neutra para reduzir fadiga visual.
- `--surface`: paineis principais.
- `--brand`: verde institucional para acoes primarias.
- `--danger`: vermelho contido para features mortas.
- `--warning`: amarelo queimado para features em risco.
- `--positive`: verde/teal para features ativas.
- `--accent`: azul para metadados e contexto.

## Componentes

- Sidebar fixa com identidade e status da IA.
- Cards metricos para resumo executivo.
- Tabela operacional para comparacao rapida.
- Pílulas de status para leitura por escaneamento.
- Painel lateral de analise executiva.
- Formulario compacto para ingestao de logs.

## Regras de UI

- Raio maximo de 8px nos componentes.
- Tabelas para comparacao, cards apenas para metricas e paineis.
- Tipografia sem escala baseada em viewport.
- Sem hero marketing: a aplicacao abre diretamente no fluxo util.
- Estados de risco sempre combinam cor, texto e contexto.
