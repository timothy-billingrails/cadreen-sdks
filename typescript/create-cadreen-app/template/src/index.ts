import "dotenv/config";
import { Cadreen } from "@cadreen/sdk";

async function main() {
  const cadreen = new Cadreen({ apiKey: process.env.CADREEN_API_KEY });

  // ── 1. Setup your workspace ──
  console.log("Setting up workspace...");
  const workspace = await cadreen.setup({
    purpose: process.env.CADREEN_PURPOSE || "A general-purpose assistant that can answer questions and follow instructions",
  });

  if (workspace.proposals && workspace.proposals.length > 0) {
    console.log("Proposed resources:", workspace.proposals.map((p) => p.description).join(", "));
  }

  // ── 2. Make your first intent call ──
  console.log("Calling intent...");
  const result = await cadreen.intent.invoke({
    messages: [{ role: "user", content: "Hello! What can you help me with?" }],
  });

  console.log(`Status: ${result.status}`);
  console.log(`Type: ${result.type}`);

  if (result.type === "direct") {
    console.log(`Response: ${result.message.content}`);
  }

  // ── 3. Read the trace ──
  console.log(`Trace: ${result.intelligence.summary}`);
  console.log(`Confidence: ${result.intelligence.governance.confidence}`);
}

main().catch(console.error);
