import { StrictMode, useEffect, useState } from "react";
import { createRoot } from "react-dom/client";
import "./styles.css";

const types = ["Membership", "Society Day", "Event"];

function Layout({ children }) {
  return <><nav><a href="/">Enquiry form</a><a href="/staff">Staff</a></nav><main>{children}</main></>;
}

function Form() {
  const [message, setMessage] = useState("");
  const [error, setError] = useState(false);

  async function submit(event) {
    event.preventDefault();
    const response = await fetch("/api/enquiries", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(Object.fromEntries(new FormData(event.target))),
    });
    setError(!response.ok);
    setMessage(response.ok ? "Thanks — we have received your enquiry." : await response.text());
    if (response.ok) event.target.reset();
  }

  return <Layout><h1>Club Enquiry</h1><form onSubmit={submit}>
    <label>Name<input required name="name" /></label>
    <label>Email<input required type="email" name="email" /></label>
    <label>Enquiry type<select name="type">{types.map((type) => <option key={type}>{type}</option>)}</select></label>
    <label>Message<textarea required name="message" rows="5" /></label>
    <button>Submit enquiry</button><p className={error ? "error" : "success"}>{message}</p>
  </form></Layout>;
}

function Staff() {
  const [enquiries, setEnquiries] = useState([]);
  const [selected, setSelected] = useState("");
  const [message, setMessage] = useState("");
  const load = () => fetch("/api/enquiries").then((r) => r.ok ? r.json() : Promise.reject()).then(setEnquiries).catch(() => setMessage("Could not load enquiries"));
  useEffect(() => { load(); }, []);

  async function backdate(event) {
    event.preventDefault();
    if (!selected) return setMessage("Select an enquiry to backdate");
    await fetch(`/api/enquiries/${selected}/backdate`, { method: "POST" });
    setMessage("");
    load();
  }

  async function contact(id) {
    await fetch(`/api/enquiries/${id}`, { method: "PATCH", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ status: "Contacted" }) });
    load();
  }

  return <Layout><h1>Staff — enquiries</h1>
    <form className="dev-tools" onSubmit={backdate}><label>Dev feature: backdate by 3 days<select value={selected} onChange={(e) => setSelected(e.target.value)}><option value="">Select enquiry</option>{enquiries.map((e) => <option key={e.id} value={e.id}>#{e.id} {e.name}</option>)}</select></label><button>Backdate selected enquiry</button></form>
    <p className="error">{message}</p><table><thead><tr><th>Name</th><th>Email</th><th>Type</th><th>Message</th><th>Status</th><th>Created</th><th>Follow-Up Due</th><th /></tr></thead><tbody>
      {enquiries.map((e) => <tr className={e.follow_up_due ? "due" : ""} key={e.id}><td>{e.name}</td><td>{e.email}</td><td>{e.type}</td><td>{e.message}</td><td>{e.status}</td><td>{new Date(e.created_at).toLocaleString()}</td><td>{e.follow_up_due ? "Yes" : "No"}</td><td>{e.status === "New" && <button onClick={() => contact(e.id)}>Mark as Contacted</button>}</td></tr>)}
      {!enquiries.length && <tr><td colSpan="8">No enquiries yet.</td></tr>}
    </tbody></table>
  </Layout>;
}

function App() { return window.location.pathname === "/staff" ? <Staff /> : <Form />; }
createRoot(document.getElementById("root")).render(<App />);
