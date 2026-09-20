import { useEffect, useMemo, useState } from "react"
import { AlertTriangle, BarChart3, CalendarDays, CheckCircle2, ChevronDown, Circle, Clock3, Download, LayoutDashboard, Menu, MessageSquare, MoreHorizontal, Plus, Rocket, Search, Settings, Trash2, Users, X } from "lucide-react"
import { api } from "@/lib/api"
import { cn, shortDate } from "@/lib/utils"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent } from "@/components/ui/card"
import { Input } from "@/components/ui/input"

const columns = [
  { id: "todo", name: "To do", icon: Circle, tone: "bg-slate-400" },
  { id: "in_progress", name: "In progress", icon: Clock3, tone: "bg-violet-500" },
  { id: "review", name: "Review", icon: MessageSquare, tone: "bg-amber-500" },
  { id: "done", name: "Done", icon: CheckCircle2, tone: "bg-emerald-500" },
]
const priorities = { low: "text-slate-500", medium: "text-blue-600", high: "text-orange-600", urgent: "text-red-600" }

function App() {
  const [d, setD] = useState({ tasks: [], projects: [], users: [], comments: [], stats: {} })
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")
  const [search, setSearch] = useState("")
  const [project, setProject] = useState(0)
  const [selected, setSelected] = useState(null)
  const [showNew, setShowNew] = useState(false)
  const [sidebar, setSidebar] = useState(false)
  const [toast, setToast] = useState("")
  const [form, setForm] = useState({ title: "", description: "", projectId: 1, assigneeId: 1, priority: "medium", estimate: 3, dueDate: "" })
  const [comment, setComment] = useState("")

  function load() {
    setLoading(true)
    api.dashboard().then((x) => { setD(x); setError("") }).catch((e) => setError(e.message)).finally(() => setLoading(false))
  }
  useEffect(() => { load() }, [])
  useEffect(() => {
    if (!toast) return
    const n = setTimeout(() => setToast(""), 2200)
    return () => clearTimeout(n)
  }, [toast])

  const tasks = useMemo(() => d.tasks.filter((t) => {
    const blob = `${t.title} ${t.description} ${(t.tags || []).join(" ")}`.toLowerCase()
    if (search && !blob.includes(search.toLowerCase())) return false
    if (project && t.projectId !== Number(project)) return false
    return true
  }), [d.tasks, search, project])

  async function move(id, status) {
    try {
      const updated = await api.updateTask(id, { status })
      setD({ ...d, tasks: d.tasks.map((x) => x.id === id ? updated : x) })
      if (selected?.id === id) setSelected(updated)
      setToast("Task moved")
    } catch (e) { setError(e.message) }
  }

  async function create(e) {
    e.preventDefault()
    try {
      const task = await api.createTask({ ...form, projectId: Number(form.projectId), assigneeId: Number(form.assigneeId), estimate: Number(form.estimate), status: "todo", tags: ["new"] })
      setD({ ...d, tasks: [...d.tasks, task] })
      setForm({ title: "", description: "", projectId: 1, assigneeId: 1, priority: "medium", estimate: 3, dueDate: "" })
      setShowNew(false)
      setToast("Task created")
    } catch (e2) { setError(e2.message) }
  }

  async function remove(id) {
    if (!window.confirm("Delete this task forever?")) return
    try {
      await api.deleteTask(id)
      setD({ ...d, tasks: d.tasks.filter((x) => x.id !== id) })
      setSelected(null)
      setToast("Task deleted")
    } catch (e) { setError(e.message) }
  }

  async function addComment(e) {
    e.preventDefault()
    try {
      const x = await api.addComment({ taskId: selected.id, authorId: 1, body: comment })
      setD({ ...d, comments: [...d.comments, x] })
      setComment("")
    } catch (e2) { setError(e2.message) }
  }

  function openResource(url) { window.location.href = url }
  function lookupUser(id) { return d.users.find((x) => x.id === id) || { name: "Unassigned", avatar: "?" } }
  function lookupProject(id) { return d.projects.find((x) => x.id === id) || { name: "Unknown", color: "slate" } }

  return <div className="min-h-screen bg-background text-foreground">
    <aside className={cn("fixed inset-y-0 left-0 z-40 w-64 border-r bg-[#17131f] text-white transition-transform lg:translate-x-0", sidebar ? "translate-x-0" : "-translate-x-full")}>
      <div className="flex h-20 items-center gap-3 border-b border-white/10 px-6">
        <div className="grid h-10 w-10 place-items-center rounded-xl bg-violet-500"><Rocket size={20} /></div>
        <div><div className="font-semibold tracking-tight">Orbit Board</div><div className="text-xs text-white/45">Moon launch team</div></div>
        <button className="ml-auto lg:hidden" onClick={() => setSidebar(false)}><X size={18} /></button>
      </div>
      <nav className="space-y-1 p-4 text-sm">
        <a className="nav-active" href="#board"><LayoutDashboard size={17} /> Board</a>
        <a className="nav-item" href="#insights"><BarChart3 size={17} /> Insights</a>
        <a className="nav-item" href="#calendar"><CalendarDays size={17} /> Calendar</a>
        <a className="nav-item" href="#team"><Users size={17} /> Team</a>
        <div className="px-3 pb-2 pt-7 text-[10px] font-bold uppercase tracking-[.18em] text-white/35">Projects</div>
        {d.projects.map((p) => <button key={p.id} onClick={() => { setProject(p.id); setSidebar(false) }} className="nav-item w-full"><span className={cn("h-2.5 w-2.5 rounded-full", p.color === "violet" ? "bg-violet-400" : "bg-cyan-400")} /> {p.name}</button>)}
      </nav>
      <div className="absolute bottom-5 left-4 right-4 rounded-xl border border-amber-300/20 bg-amber-300/10 p-3 text-xs text-amber-100">
        <div className="mb-1 flex items-center gap-2 font-semibold"><AlertTriangle size={14} /> Local preview</div>
        Do not deploy this application to production.
      </div>
    </aside>

    <main className="lg:pl-64">
      <header className="sticky top-0 z-30 flex h-20 items-center gap-3 border-b bg-background/90 px-4 backdrop-blur-xl md:px-8">
        <button className="lg:hidden" onClick={() => setSidebar(true)}><Menu /></button>
        <div className="relative max-w-md flex-1"><Search className="absolute left-3 top-3 text-muted-foreground" size={16} /><Input value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Search tasks, tags, people…" className="pl-9" /></div>
        <Button variant="outline" size="icon" onClick={() => openResource("/api/export")} title="Export data"><Download size={17} /></Button>
        <Button variant="ghost" size="icon"><Settings size={18} /></Button>
        <div className="grid h-9 w-9 place-items-center rounded-full bg-gradient-to-br from-violet-500 to-fuchsia-500 text-xs font-bold text-white">MC</div>
      </header>

      <div className="p-4 md:p-8">
        <section className="mb-7 flex flex-col justify-between gap-4 md:flex-row md:items-end">
          <div><p className="mb-2 text-sm font-medium text-violet-600">MONDAY, SEPTEMBER 20</p><h1 className="text-3xl font-bold tracking-tight md:text-4xl">Good morning, Maya.</h1><p className="mt-2 text-muted-foreground">You have {d.stats.inProgress || 0} tasks in progress and {d.stats.overdue || 0} overdue.</p></div>
          <Button onClick={() => setShowNew(true)}><Plus className="mr-2" size={17} /> Add task</Button>
        </section>

        {error && <div className="mb-5 flex items-center justify-between rounded-xl border border-red-200 bg-red-50 p-3 text-sm text-red-700"><span>{error}</span><button onClick={() => setError("")}><X size={16} /></button></div>}

        <section className="mb-8 grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
          <Stat label="Total tasks" value={d.stats.total || 0} foot={`${d.stats.points || 0} estimated points`} icon={LayoutDashboard} color="violet" />
          <Stat label="In progress" value={d.stats.inProgress || 0} foot="Across active projects" icon={Clock3} color="blue" />
          <Stat label="Completion" value={`${d.stats.completion || 0}%`} foot={`${d.stats.done || 0} tasks shipped`} icon={CheckCircle2} color="green" />
          <Stat label="Needs attention" value={(d.stats.urgent || 0) + (d.stats.overdue || 0)} foot="Urgent and overdue" icon={AlertTriangle} color="orange" />
        </section>

        <section id="board">
          <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
            <div><h2 className="text-xl font-semibold">Sprint board</h2><p className="text-sm text-muted-foreground">September launch · Week 3</p></div>
            <div className="flex items-center gap-2">
              <select className="select" value={project} onChange={(e) => setProject(Number(e.target.value))}><option value="0">All projects</option>{d.projects.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}</select>
              {(project !== 0 || search) && <Button variant="ghost" size="sm" onClick={() => { setProject(0); setSearch("") }}>Clear filters</Button>}
            </div>
          </div>

          {loading ? <Loading /> : <div className="grid gap-4 xl:grid-cols-4 md:grid-cols-2">
            {columns.map((col) => <div key={col.id} className="board-column">
              <div className="mb-3 flex items-center gap-2 px-1"><span className={cn("h-2 w-2 rounded-full", col.tone)} /><span className="text-sm font-semibold">{col.name}</span><span className="count">{tasks.filter((t) => t.status === col.id).length}</span><button className="ml-auto text-muted-foreground"><MoreHorizontal size={17} /></button></div>
              <div className="space-y-3">
                {tasks.filter((t) => t.status === col.id).map((t) => <TaskCard key={t.id} task={t} user={lookupUser(t.assigneeId)} project={lookupProject(t.projectId)} onOpen={() => setSelected(t)} onMove={move} />)}
                <button onClick={() => setShowNew(true)} className="flex w-full items-center gap-2 rounded-xl border border-dashed p-3 text-sm text-muted-foreground hover:border-violet-300 hover:text-violet-600"><Plus size={15} /> Add task</button>
              </div>
            </div>)}
          </div>}
        </section>
      </div>
    </main>

    {showNew && <Modal title="Create a task" onClose={() => setShowNew(false)}><form className="space-y-4" onSubmit={create}>
      <Field label="Title"><Input required value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} placeholder="What needs to happen?" /></Field>
      <Field label="Description"><textarea className="textarea" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} /></Field>
      <div className="grid grid-cols-2 gap-3">
        <Field label="Project"><select className="select w-full" value={form.projectId} onChange={(e) => setForm({ ...form, projectId: e.target.value })}>{d.projects.map((p) => <option key={p.id} value={p.id}>{p.name}</option>)}</select></Field>
        <Field label="Assignee"><select className="select w-full" value={form.assigneeId} onChange={(e) => setForm({ ...form, assigneeId: e.target.value })}>{d.users.map((u) => <option key={u.id} value={u.id}>{u.name}</option>)}</select></Field>
        <Field label="Priority"><select className="select w-full" value={form.priority} onChange={(e) => setForm({ ...form, priority: e.target.value })}><option>low</option><option>medium</option><option>high</option><option>urgent</option></select></Field>
        <Field label="Points"><Input type="number" min="0" value={form.estimate} onChange={(e) => setForm({ ...form, estimate: e.target.value })} /></Field>
      </div>
      <Field label="Due date"><Input type="date" value={form.dueDate} onChange={(e) => setForm({ ...form, dueDate: e.target.value })} /></Field>
      <div className="flex justify-end gap-2 pt-2"><Button type="button" variant="outline" onClick={() => setShowNew(false)}>Cancel</Button><Button>Create task</Button></div>
    </form></Modal>}

    {selected && <Modal title={`Task #${selected.id}`} wide onClose={() => setSelected(null)}><div className="space-y-5">
      <div><div className="mb-2 flex items-center gap-2"><Badge>{lookupProject(selected.projectId).name}</Badge><span className={cn("text-xs font-semibold uppercase", priorities[selected.priority])}>{selected.priority}</span></div><h3 className="text-2xl font-bold">{selected.title}</h3><p className="mt-2 text-sm leading-6 text-muted-foreground">{selected.description}</p></div>
      <div className="grid grid-cols-2 gap-3 rounded-xl bg-muted/50 p-4 text-sm"><Detail label="Assignee" value={lookupUser(selected.assigneeId).name} /><Detail label="Due date" value={shortDate(selected.dueDate)} /><Detail label="Estimate" value={`${selected.estimate} points`} /><Detail label="Status" value={selected.status.replace("_", " ")} /></div>
      <div><label className="label">Move to</label><div className="flex flex-wrap gap-2">{columns.map((x) => <Button key={x.id} size="sm" variant={selected.status === x.id ? "default" : "outline"} onClick={() => move(selected.id, x.id)}>{x.name}</Button>)}</div></div>
      <div><label className="label">Comments</label><div className="space-y-2">{d.comments.filter((x) => x.taskId === selected.id).map((x) => <div key={x.id} className="rounded-lg bg-muted p-3 text-sm"><b>{lookupUser(x.authorId).name}</b><div className="mt-1 text-muted-foreground" dangerouslySetInnerHTML={{ __html: x.body }} /></div>)}</div>
        <form onSubmit={addComment} className="mt-3 flex gap-2"><Input value={comment} onChange={(e) => setComment(e.target.value)} placeholder="Write a comment…" /><Button disabled={!comment}>Send</Button></form>
      </div>
      <div className="border-t pt-4"><Button variant="destructive" size="sm" onClick={() => remove(selected.id)}><Trash2 className="mr-2" size={15} /> Delete task</Button></div>
    </div></Modal>}

    {toast && <div className="fixed bottom-5 right-5 z-[70] rounded-xl bg-[#17131f] px-4 py-3 text-sm font-medium text-white shadow-2xl">{toast}</div>}
  </div>
}

function TaskCard({ task, user, project, onOpen, onMove }) {
  return <Card className="group cursor-pointer border-white bg-white/90 transition hover:-translate-y-0.5 hover:border-violet-200" onClick={onOpen}>
    <CardContent className="p-4"><div className="mb-3 flex items-start justify-between gap-2"><Badge variant="secondary">{project.name}</Badge><button onClick={(e) => { e.stopPropagation(); const next = columns[(columns.findIndex((x) => x.id === task.status) + 1) % columns.length]; onMove(task.id, next.id) }} className="opacity-0 transition group-hover:opacity-100"><ChevronDown size={16} /></button></div>
      <h3 className="font-semibold leading-5">{task.title}</h3><p className="mt-1 line-clamp-2 text-xs leading-5 text-muted-foreground">{task.description}</p>
      <div className="mt-4 flex flex-wrap gap-1">{(task.tags || []).map((x) => <span key={x} className="rounded bg-violet-50 px-1.5 py-0.5 text-[10px] font-medium text-violet-700">{x}</span>)}</div>
      <div className="mt-4 flex items-center border-t pt-3"><span className={cn("mr-2 text-[10px] font-bold uppercase", priorities[task.priority])}>{task.priority}</span><CalendarDays size={13} className="ml-auto mr-1 text-muted-foreground" /><span className="text-xs text-muted-foreground">{shortDate(task.dueDate)}</span><span className="ml-3 grid h-7 w-7 place-items-center rounded-full bg-[#272032] text-[9px] font-bold text-white" title={user.name}>{user.avatar}</span></div>
    </CardContent>
  </Card>
}

function Stat({ label, value, foot, icon: Icon, color }) {
  const colors = { violet: "bg-violet-100 text-violet-700", blue: "bg-blue-100 text-blue-700", green: "bg-emerald-100 text-emerald-700", orange: "bg-orange-100 text-orange-700" }
  return <Card className="border-white"><CardContent className="flex items-center gap-4 p-5"><div className={cn("grid h-11 w-11 place-items-center rounded-xl", colors[color])}><Icon size={20} /></div><div><div className="text-2xl font-bold">{value}</div><div className="text-sm font-medium">{label}</div><div className="mt-0.5 text-[11px] text-muted-foreground">{foot}</div></div></CardContent></Card>
}
function Modal({ title, onClose, children, wide }) {
  return <div className="fixed inset-0 z-50 grid place-items-center bg-[#120d19]/60 p-4 backdrop-blur-sm" onMouseDown={onClose}><div className={cn("max-h-[90vh] w-full overflow-auto rounded-2xl bg-white p-6 shadow-2xl", wide ? "max-w-2xl" : "max-w-lg")} onMouseDown={(e) => e.stopPropagation()}><div className="mb-5 flex items-center justify-between"><h2 className="text-xl font-bold">{title}</h2><Button variant="ghost" size="icon" onClick={onClose}><X size={18} /></Button></div>{children}</div></div>
}
function Field({ label, children }) { return <label className="block"><span className="label">{label}</span>{children}</label> }
function Detail({ label, value }) { return <div><div className="text-xs text-muted-foreground">{label}</div><div className="mt-1 font-medium capitalize">{value}</div></div> }
function Loading() { return <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">{[1, 2, 3, 4].map((x) => <div key={x} className="space-y-3">{[1, 2].map((y) => <div key={y} className="h-44 animate-pulse rounded-xl bg-muted" />)}</div>)}</div> }

export default App
