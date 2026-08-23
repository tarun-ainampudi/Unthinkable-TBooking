import React, { useState } from "react";
import { Mail, Phone, Edit2, LogOut, Ticket } from "lucide-react";

export default function ProfilePage({ user, onLogin, onLogout }) {
  const [editing, setEditing] = useState(false);
  const [loginForm, setLoginForm] = useState({ email: "aarav.shah@example.com", password: "ticket123" });
  const [loginError, setLoginError] = useState("");
  const [busy, setBusy] = useState(false);
  const [form, setForm] = useState(user || {
    name: "",
    email: "",
    phone: "",
    memberSince: "",
  });

  React.useEffect(() => {
    if (user) setForm(user);
  }, [user]);

  const initials = (form.name || "U")
    .split(" ")
    .map((p) => p[0])
    .join("")
    .slice(0, 2)
    .toUpperCase();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoginError("");
    setBusy(true);

    try {
      await onLogin(loginForm);
    } catch (error) {
      setLoginError(error.message || "Login failed. Please try again.");
    } finally {
      setBusy(false);
    }
  };

  if (!user) {
    return (
      <div className="page page-profile login-view">
        <h1 className="page-title">Login</h1>

        <div className="profile-layout login-layout">
          <div className="profile-card login-card">
            <div className="profile-avatar"><Ticket size={20} strokeWidth={2.25} /></div>
            <h2>Welcome back</h2>
            <p className="profile-member-since">Sign in to book tickets, hold seats, and access your bookings.</p>

            <form className="field-group" onSubmit={handleSubmit}>
              <label className="field">
                <span>Email</span>
                <input
                  type="email"
                  value={loginForm.email}
                  onChange={(e) => setLoginForm({ ...loginForm, email: e.target.value })}
                  placeholder="you@example.com"
                />
              </label>
              <label className="field">
                <span>Password</span>
                <input
                  type="password"
                  value={loginForm.password}
                  onChange={(e) => setLoginForm({ ...loginForm, password: e.target.value })}
                  placeholder="Enter password"
                />
              </label>

              {loginError && <p className="error-text">{loginError}</p>}

              <button type="submit" className="btn btn-primary btn-block" disabled={busy}>
                {busy ? "Logging in..." : "Login"}
              </button>
            </form>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="page page-profile profile-view">
      <h1 className="page-title">Profile</h1>

      <div className="profile-layout profile-layout-wide">
        <div className="profile-card">
          <div className="profile-avatar">{initials}</div>
          {!editing ? (
            <>
              <h2>{form.name || "Guest User"}</h2>
              <p className="profile-member-since">Member since {form.memberSince || "recently"}</p>
              <div className="profile-detail-list">
                <div className="profile-detail">
                  <Mail size={15} strokeWidth={2.25} />
                  <span>{form.email || "—"}</span>
                </div>
                <div className="profile-detail">
                  <Phone size={15} strokeWidth={2.25} />
                  <span>{form.phone || "—"}</span>
                </div>
              </div>
              <button className="btn btn-secondary btn-block" onClick={() => setEditing(true)}>
                <Edit2 size={14} strokeWidth={2.25} /> Edit profile
              </button>
            </>
          ) : (
            <form
              className="field-group"
              onSubmit={(e) => {
                e.preventDefault();
                setEditing(false);
              }}
            >
              <label className="field">
                <span>Full name</span>
                <input value={form.name || ""} onChange={(e) => setForm({ ...form, name: e.target.value })} />
              </label>
              <label className="field">
                <span>Email</span>
                <input value={form.email || ""} onChange={(e) => setForm({ ...form, email: e.target.value })} />
              </label>
              <label className="field">
                <span>Phone</span>
                <input value={form.phone || ""} onChange={(e) => setForm({ ...form, phone: e.target.value })} />
              </label>
              <div className="field-row">
                <button type="button" className="btn btn-secondary btn-block" onClick={() => setEditing(false)}>
                  Cancel
                </button>
                <button type="submit" className="btn btn-primary btn-block">
                  Save changes
                </button>
              </div>
            </form>
          )}
        </div>

        <div className="profile-stats">
          <div className="stat-card">
            <span className="stat-number">{user ? 1 : 0}</span>
            <span className="stat-label">Total bookings</span>
          </div>
          <div className="stat-card">
            <span className="stat-number">0</span>
            <span className="stat-label">Shows attended</span>
          </div>
          <div className="stat-card">
            <span className="stat-number">{(form.memberSince || "").split(" ")[1] || "Now"}</span>
            <span className="stat-label">Customer since</span>
          </div>
          <button className="btn btn-secondary btn-block logout-btn" onClick={onLogout}>
            <LogOut size={14} strokeWidth={2.25} /> Sign out
          </button>
        </div>
      </div>
    </div>
  );
}
