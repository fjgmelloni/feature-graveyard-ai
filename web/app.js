const state = {
  report: null,
  selectedFeature: null,
};

const elements = {
  windowDays: document.querySelector("#windowDays"),
  refreshButton: document.querySelector("#refreshButton"),
  totalFeatures: document.querySelector("#totalFeatures"),
  deadFeatures: document.querySelector("#deadFeatures"),
  atRiskFeatures: document.querySelector("#atRiskFeatures"),
  removalCandidates: document.querySelector("#removalCandidates"),
  generatedAt: document.querySelector("#generatedAt"),
  featureRows: document.querySelector("#featureRows"),
  executiveInsight: document.querySelector("#executiveInsight"),
  generatedBy: document.querySelector("#generatedBy"),
  modelStatus: document.querySelector("#modelStatus"),
  usageForm: document.querySelector("#usageForm"),
  formFeedback: document.querySelector("#formFeedback"),
};

const today = new Date().toISOString().slice(0, 10);
elements.usageForm.elements.lastAccess.value = today;

elements.refreshButton.addEventListener("click", loadReport);
elements.windowDays.addEventListener("change", loadReport);
elements.usageForm.addEventListener("submit", async (event) => {
  event.preventDefault();
  const form = new FormData(elements.usageForm);
  const payload = {
    logs: [
      {
        feature: form.get("feature"),
        userId: form.get("userId"),
        lastAccess: form.get("lastAccess"),
        totalAccess: Number(form.get("totalAccess")),
      },
    ],
  };

  elements.formFeedback.textContent = "Enviando log...";
  const response = await fetch("/api/usage-logs", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    const error = await response.json();
    elements.formFeedback.textContent = error.error || "Nao foi possivel enviar o log.";
    return;
  }

  elements.formFeedback.textContent = "Log inserido. Relatorio recalculado.";
  await loadReport();
});

async function loadReport() {
  const windowDays = elements.windowDays.value;
  const response = await fetch(`/api/graveyard/report?windowDays=${windowDays}`);
  if (!response.ok) {
    elements.modelStatus.textContent = "Erro ao carregar relatorio";
    return;
  }

  state.report = await response.json();
  const sorted = [...state.report.analyses].sort(compareFeatures);
  state.selectedFeature =
    sorted.find((item) => item.feature === state.selectedFeature?.feature) || sorted[0] || null;

  renderMetrics();
  renderRows(sorted);
  renderInsight(state.selectedFeature);
}

function compareFeatures(a, b) {
  const rank = { DEAD_FEATURE: 0, AT_RISK: 1, ACTIVE: 2 };
  return rank[a.status] - rank[b.status] || b.daysSinceAccess - a.daysSinceAccess;
}

function renderMetrics() {
  elements.totalFeatures.textContent = state.report.totalFeatures;
  elements.deadFeatures.textContent = state.report.deadFeatures;
  elements.atRiskFeatures.textContent = state.report.atRiskFeatures;
  elements.removalCandidates.textContent = state.report.removalCandidates;
  elements.generatedAt.textContent = formatDateTime(state.report.generatedAt);
}

function renderRows(analyses) {
  elements.featureRows.innerHTML = "";
  for (const analysis of analyses) {
    const row = document.createElement("tr");
    row.className = analysis.feature === state.selectedFeature?.feature ? "selected" : "";
    row.innerHTML = `
      <td><strong>${escapeHTML(analysis.feature)}</strong></td>
      <td><span class="status ${analysis.status}">${labelStatus(analysis.status)}</span></td>
      <td>${labelRisk(analysis.risk)}</td>
      <td>${analysis.lastAccess} <span class="muted">(${analysis.daysSinceAccess}d)</span></td>
      <td>${analysis.uniqueUsers}</td>
      <td>${analysis.frequencyPerMonth.toFixed(2)}</td>
    `;
    row.addEventListener("click", () => {
      state.selectedFeature = analysis;
      renderRows(analyses);
      renderInsight(analysis);
    });
    elements.featureRows.appendChild(row);
  }
}

function renderInsight(analysis) {
  if (!analysis) {
    elements.executiveInsight.innerHTML = "<strong>Nenhuma feature</strong><p>Envie logs para gerar a primeira analise.</p>";
    return;
  }

  elements.generatedBy.textContent = analysis.generatedBy;
  elements.modelStatus.textContent =
    analysis.generatedBy === "gemini" ? "Gemini ativo" : "Analise local ativa";
  elements.executiveInsight.innerHTML = `
    <strong>${escapeHTML(analysis.feature)}</strong>
    <p>${escapeHTML(analysis.summary)}</p>
    <dl>
      <div>
        <dt>Impacto</dt>
        <dd>${escapeHTML(analysis.businessImpact)}</dd>
      </div>
      <div>
        <dt>Acao sugerida</dt>
        <dd>${escapeHTML(analysis.suggestedAction)}</dd>
      </div>
      <div>
        <dt>Racional</dt>
        <dd>${escapeHTML(analysis.executiveRationale)}</dd>
      </div>
      <div>
        <dt>Confianca</dt>
        <dd>${Math.round(analysis.confidence * 100)}%</dd>
      </div>
    </dl>
  `;
}

function labelStatus(status) {
  return {
    DEAD_FEATURE: "Morta",
    AT_RISK: "Em risco",
    ACTIVE: "Ativa",
  }[status] || status;
}

function labelRisk(risk) {
  return {
    LOW: "Baixo",
    MEDIUM: "Medio",
    HIGH: "Alto",
    CRITICAL: "Critico",
  }[risk] || risk;
}

function formatDateTime(value) {
  return new Intl.DateTimeFormat("pt-BR", {
    dateStyle: "short",
    timeStyle: "short",
  }).format(new Date(value));
}

function escapeHTML(value) {
  return String(value)
    .replaceAll("&", "&amp;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;")
    .replaceAll('"', "&quot;")
    .replaceAll("'", "&#039;");
}

loadReport();
