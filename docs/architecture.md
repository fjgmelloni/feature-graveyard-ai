# Arquitetura

Feature Graveyard AI usa uma arquitetura DDD enxuta, pensada para evoluir de demo para produto interno.

## Fluxo principal

1. A API recebe logs de uso em `POST /api/usage-logs`.
2. A camada de aplicacao valida e persiste os logs via `UsageRepository`.
3. O relatorio agrega eventos por feature.
4. O dominio classifica cada feature como `ACTIVE`, `AT_RISK` ou `DEAD_FEATURE`.
5. O adaptador de IA gera uma analise executiva com Gemini ou fallback local.
6. O frontend consome `GET /api/graveyard/report`.

## Decisoes

- Persistencia em memoria: suficiente para demo, facil de trocar por PostgreSQL.
- Gemini atras de interface: mantem o dominio testavel e permite fallback offline.
- Classificacao deterministica: IA nao decide o risco; ela explica o resultado calculado.
- Frontend estatico: reduz atrito de setup e deixa o projeto simples para portfólio.

## Evolucoes naturais

- Repository PostgreSQL com historico por workspace.
- Importacao de logs via CSV, Segment, Datadog ou OpenTelemetry.
- Pesos configuraveis por modulo critico.
- Workflow de aprovacao de sunset.
- Exportacao de relatorio executivo em PDF.
