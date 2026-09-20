// WebMCP tools for this portfolio (https://github.com/webmachinelearning/webmcp).
// Requires a browser with document.modelContext — e.g. Chrome with
// chrome://flags/#enable-webmcp-testing enabled.

const modelContext = document.modelContext;
if (!modelContext) {
  console.debug("WebMCP: document.modelContext is not available in this browser.");
} else {
  const controller = new AbortController();
  const toolOptions = { signal: controller.signal };

  async function fetchJSON(path) {
    const response = await fetch(path);
    if (!response.ok) {
      throw new Error(`Request failed (${response.status}): ${path}`);
    }
    return response.json();
  }

  await modelContext.registerTool(
    {
      name: "listProjects",
      description:
        "List published portfolio projects. Returns title, slug, tags, and summary for each project.",
      inputSchema: {
        type: "object",
        properties: {},
      },
      async execute() {
        return fetchJSON("/data/projects.json");
      },
    },
    toolOptions,
  );

  await modelContext.registerTool(
    {
      name: "getProject",
      description:
        "Get full metadata for one published project by slug, including repo link and page URL.",
      inputSchema: {
        type: "object",
        properties: {
          slug: {
            type: "string",
            description: "Project slug, for example banking-platform or crossplane-study.",
          },
        },
        required: ["slug"],
      },
      async execute({ slug }) {
        const path = `/data/projects/${encodeURIComponent(slug)}.json`;
        const response = await fetch(path);
        if (response.status === 404) {
          return {
            content: [
              {
                type: "text",
                text: `No published project with slug "${slug}". Call listProjects to see available slugs.`,
              },
            ],
          };
        }
        if (!response.ok) {
          throw new Error(`Request failed (${response.status}): ${path}`);
        }
        return response.json();
      },
    },
    toolOptions,
  );

  await modelContext.registerTool(
    {
      name: "listResources",
      description:
        "List curated external resources grouped by category (Kubernetes, GitOps, observability, and more).",
      inputSchema: {
        type: "object",
        properties: {},
      },
      async execute() {
        return fetchJSON("/data/resources.json");
      },
    },
    toolOptions,
  );
}
