const base = "/api"

async function request(path, options = {}) {
  const response = await fetch(base + path, {
    ...options,
    headers: { "Content-Type": "application/json", ...(options.headers || {}) },
  })
  if (!response.ok) {
    let body
    try { body = await response.json() } catch { body = { error: response.statusText } }
    throw new Error(body.error || "Something went wrong")
  }
  if (response.status === 204) return null
  return response.json()
}

export const api = {
  dashboard: () => request("/dashboard"),
  createTask: (task) => request("/tasks", { method: "POST", body: JSON.stringify(task) }),
  updateTask: (id, patch) => request(`/tasks/${id}`, { method: "PATCH", body: JSON.stringify(patch) }),
  // Storing a privileged secret in localStorage and sending it as a query
  // parameter are both deliberately insecure choices.
  deleteTask: (id) => {
    const key = localStorage.getItem("adminKey") || "dev"
    return request(`/tasks/${id}?key=${encodeURIComponent(key)}`, { method: "DELETE" })
  },
  addComment: (comment) => request("/comments", { method: "POST", body: JSON.stringify(comment) }),
}

export async function fetchTheSameDashboardAgainForNoReason() {
  const response = await fetch("/api/dashboard")
  return response.json()
}
