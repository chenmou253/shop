import { useEffect, useState, type FormEvent } from "react";
import { Link, NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";
import { ChevronDown, ChevronUp, Mail, Menu, MessageCircle, Phone, Search, Send, X } from "lucide-react";
import { api } from "../api";
import { useInquiry } from "../context/InquiryContext";
import { useSite } from "../context/SiteContext";
import type { NavItem } from "../types";

export function Layout() {
  const { settings, categories, navigation } = useSite();
  const siteName = settings.site_name || "";
  const whatsappNumber = settings.whatsapp?.replace(/\D/g, "") || "";
  const hasFooterContact = Boolean(settings.phone || settings.whatsapp || settings.email || settings.address);
  const { items } = useInquiry();
  const [open, setOpen] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();
  const [query, setQuery] = useState("");
  useEffect(() => setOpen(false), [location.pathname, location.search]);

  function search(e: FormEvent) {
    e.preventDefault();
    if (query.trim()) navigate(`/search?q=${encodeURIComponent(query.trim())}`);
  }

  return <>
    <header className="site-header">
      <div className="header-inner shell">
        {settings.logo_url && <Link className="site-logo" to="/"><img src={settings.logo_url} alt={siteName} /></Link>}
        <button className="mobile-menu" title="Open navigation" onClick={() => setOpen(!open)}>{open ? <X /> : <Menu />}</button>
        <nav className={`main-nav ${open ? "open" : ""}`}>
          {navigation.filter((item) => item.location === "header").map((item) => item.children?.length
            ? <Dropdown key={item.id} title={item.label} to={item.path} links={item.children.map((child) => [child.label, child.path])} wide={item.path === "/products"} />
            : <NavLink key={item.id} to={item.path}>{item.label}</NavLink>)}
        </nav>
        <div className="language-badge"><span>EN</span> English</div>
      </div>
    </header>
    <main><Outlet /></main>
    <section className="pre-footer">
      <div className="shell search-contact">
        <form onSubmit={search}><Search size={20} /><input value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Please enter what you want to search" /><button>Search</button></form>
        <div><strong>Need technical support?</strong><span>Our engineering and sales team responds within 24 hours.</span></div>
        <Link className="button light" to="/feedback">Send Inquiry</Link>
      </div>
    </section>
    <footer className="site-footer">
      <div className="shell footer-grid">
        {hasFooterContact && <div className="footer-contact">
          <h3>Contact Us</h3>
          {settings.phone && <p><Phone size={16} /> Tel: {settings.phone}</p>}
          {settings.whatsapp && <p><MessageCircle size={16} /> WhatsApp: {settings.whatsapp}</p>}
          {settings.email && <p><Mail size={16} /> {settings.email}</p>}
          {settings.address && <p>{settings.address}</p>}
        </div>}
        <div className="footer-nav"><h3>Quick Navigation</h3><FooterLinks links={footerNavigation(navigation)} /></div>
        <div className="footer-products"><h3>Products</h3><FooterLinks links={categories.slice(0, 10).map((c) => [c.name, `/category/${c.slug}`])} /></div>
        {settings.qr_code_url && <div className="footer-qr"><h3>QR Code</h3><img className="qr" src={settings.qr_code_url} alt="QR Code" /></div>}
      </div>
      <div className="copyright shell">© {new Date().getFullYear()} {siteName} All Rights Reserved. <Link to="/page/privacy-policy">Privacy Policy</Link></div>
    </footer>
    <div className="mobile-actions">
      {settings.phone && <a href={`tel:${settings.phone}`} aria-label="Call us"><Phone /></a>}
      {whatsappNumber && <a href={`https://wa.me/${whatsappNumber}`} aria-label="Chat on WhatsApp"><MessageCircle /></a>}
      <button type="button" aria-label="Back to top" onClick={() => window.scrollTo({ top: 0, behavior: "smooth" })}><ChevronUp /></button>
      {settings.email && <a href={`mailto:${settings.email}`} aria-label="Email us"><Mail /></a>}
    </div>
    <Link className="floating-inquiry" to="/feedback"><Send size={20} /><span>Contact Supplier</span><b>({items.length}/10)</b></Link>
  </>;
}

function Dropdown({ title, to, links, wide }: { title: string; to: string; links: string[][]; wide?: boolean }) {
  const [expanded, setExpanded] = useState(false);
  return <div className={`nav-dropdown ${expanded ? "expanded" : ""}`}>
    <NavLink to={to}>{title}<ChevronDown className="desktop-dropdown-icon" size={14} /></NavLink>
    <button className="dropdown-toggle" type="button" aria-label={`Toggle ${title} submenu`} aria-expanded={expanded} onClick={() => setExpanded(!expanded)}><ChevronDown size={17} /></button>
    <div className={`dropdown-menu ${wide ? "wide" : ""}`}>{links.map(([label, path]) => <Link key={`${label}-${path}`} to={path}>{label}</Link>)}</div>
  </div>;
}

function FooterLinks({ links }: { links: string[][] }) {
  return <ul className="footer-links">{links.map(([label, path]) => <li key={`${label}-${path}`}><Link to={path}>{label}</Link></li>)}</ul>;
}

function footerNavigation(items: NavItem[]): string[][] {
  const footer = items.filter((item) => item.location === "footer");
  const selected = footer.length ? footer : items.filter((item) => item.location === "header");
  return selected.map((item) => [item.label, item.path]);
}

export function SubscribeBox() {
  const [email, setEmail] = useState("");
  const [message, setMessage] = useState("");
  function submit(e: FormEvent) {
    e.preventDefault();
    api<{ message: string }>("/subscriptions", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ email }) })
      .then((data) => { setMessage(data.message); setEmail(""); })
      .catch((error) => setMessage(error.message));
  }
  return <form className="subscribe" onSubmit={submit}><h3>Email Subscribe</h3><p>Get product updates and practical fitting knowledge.</p><input type="email" required value={email} onChange={(e) => setEmail(e.target.value)} placeholder="Your email address" /><button className="button">Subscribe</button>{message && <small>{message}</small>}</form>;
}
