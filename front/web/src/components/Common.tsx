import { ChevronRight, Plus } from "lucide-react";
import { Link } from "react-router-dom";
import { useInquiry } from "../context/InquiryContext";
import { useSite } from "../context/SiteContext";
import type { Category, Product } from "../types";

export function PageHero({ title, crumbs = [] }: { title: string; crumbs?: string[] }) {
  const { settings } = useSite();
  const heroStyle = settings.page_hero_url ? { backgroundImage: `linear-gradient(rgba(4, 34, 58, .72), rgba(4, 34, 58, .72)), url(${settings.page_hero_url})` } : undefined;
  return <section className="page-hero" style={heroStyle}><div className="shell"><div className="breadcrumbs"><Link to="/">Home</Link><ChevronRight size={14} />{crumbs.map((item) => <span key={item}>{item}<ChevronRight size={14} /></span>)}<b>{title}</b></div><h1>{title}</h1></div></section>;
}

export function SectionTitle({ children }: { children: React.ReactNode }) { return <h2 className="section-title">{children}</h2>; }

export function ProductCard({ product, inquiry = false }: { product: Product; inquiry?: boolean }) {
  const { add } = useInquiry();
  return <article className="product-card">
    <Link className="product-image" to={`/product/${product.slug}`}><img src={product.main_image} alt={product.name} loading="lazy" /></Link>
    <div className="product-card-body"><Link className="product-name" to={`/product/${product.slug}`}>{product.name}</Link><p>{product.summary}</p>
      <div className="card-actions"><Link to={`/product/${product.slug}`}>Learn More <ChevronRight size={15} /></Link>{inquiry && <button onClick={() => add(product)}><Plus size={15} /> Add to Inquiry</button>}</div>
    </div>
  </article>;
}

export function CategorySidebar({ latest = [] }: { latest?: Product[] }) {
  const { categories, settings } = useSite();
  return <aside className="sidebar-content">
    <div className="side-card"><h3>Categories</h3><CategoryList categories={categories} /></div>
    {latest.length > 0 && <div className="side-card"><h3>Latest Products</h3><ul className="latest-list">{latest.slice(0, 6).map((p) => <li key={p.id}><Link to={`/product/${p.slug}`}><img src={p.main_image} alt="" /><span>{p.name}</span></Link></li>)}</ul></div>}
    <div className="side-card contact-card"><h3>Contact Us</h3>{settings.phone && <p>Tel: {settings.phone}</p>}{settings.whatsapp && <p>Phone: {settings.whatsapp}</p>}{settings.email && <p>Email: {settings.email}</p>}<Link className="button" to="/feedback">Send Inquiry</Link></div>
  </aside>;
}

function CategoryList({ categories }: { categories: Category[] }) {
  return <ul className="category-list">{categories.map((category) => <li key={category.id}><Link to={`/category/${category.slug}`}>{category.name}<ChevronRight size={14} /></Link>{category.children?.length ? <CategoryList categories={category.children} /> : null}</li>)}</ul>;
}

export function Pagination({ page, pageSize, total, onChange }: { page: number; pageSize: number; total: number; onChange: (page: number) => void }) {
  const pages = Math.max(1, Math.ceil(total / pageSize));
  if (pages <= 1) return null;
  const nums = Array.from({ length: Math.min(7, pages) }, (_, i) => Math.max(1, Math.min(pages - 6, page - 3)) + i).filter((n) => n <= pages);
  return <div className="pagination"><button disabled={page === 1} onClick={() => onChange(1)}>First</button>{nums.map((n) => <button className={n === page ? "active" : ""} key={n} onClick={() => onChange(n)}>{n}</button>)}<button disabled={page === pages} onClick={() => onChange(pages)}>Last</button><span>{page}/{pages}</span></div>;
}
