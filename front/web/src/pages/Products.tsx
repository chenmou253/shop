import { useEffect, useState } from "react";
import { Search } from "lucide-react";
import { useParams, useSearchParams } from "react-router-dom";
import { api } from "../api";
import { CategorySidebar, PageHero, Pagination, ProductCard } from "../components/Common";
import type { Product } from "../types";

type Response = { items: Product[]; total: number; page: number; page_size: number };

export function Products() {
  const { slug } = useParams();
  const [params, setParams] = useSearchParams();
  const [data, setData] = useState<Response>({ items: [], total: 0, page: 1, page_size: 12 });
  const [loading, setLoading] = useState(true);
  const [query, setQuery] = useState(params.get("q") || "");
  const page = Number(params.get("page") || 1);
  useEffect(() => {
    setLoading(true);
    const search = new URLSearchParams(params);
    search.set("page", String(page));
    if (slug) search.set("category", slug);
    api<Response>(`/products?${search}`).then(setData).finally(() => setLoading(false));
  }, [slug, params.toString(), page]);
  const title = slug ? data.items[0]?.category?.name || slug.split("-").map(capitalize).join(" ") : "Products";
  function submit(e: React.FormEvent) { e.preventDefault(); const next = new URLSearchParams(params); query ? next.set("q", query) : next.delete("q"); next.set("page", "1"); setParams(next); }

  return <><PageHero title={title} /><div className="shell two-column page-space"><section className="main-column"><form className="catalog-filter" onSubmit={submit}><div><Search size={18} /><input value={query} onChange={(e) => setQuery(e.target.value)} placeholder="Search by product, material or standard" /></div><input placeholder="Material" value={params.get("material") || ""} onChange={(e) => changeFilter("material", e.target.value, params, setParams)} /><input placeholder="Size" value={params.get("size") || ""} onChange={(e) => changeFilter("size", e.target.value, params, setParams)} /><button className="button">Search</button></form>{loading ? <div className="state-page compact"><span className="loader" />Loading...</div> : data.items.length ? <><div className="product-grid catalog-grid">{data.items.map((p) => <ProductCard key={p.id} product={p} inquiry />)}</div><Pagination page={data.page} pageSize={data.page_size} total={data.total} onChange={(n) => { const next = new URLSearchParams(params); next.set("page", String(n)); setParams(next); window.scrollTo({ top: 240, behavior: "smooth" }); }} /></> : <div className="empty-state">No matching products were found.</div>}</section><CategorySidebar latest={data.items.slice(0, 6)} /></div></>;
}

function capitalize(value: string) { return value.charAt(0).toUpperCase() + value.slice(1); }
function changeFilter(key: string, value: string, params: URLSearchParams, setParams: (next: URLSearchParams) => void) { const next = new URLSearchParams(params); value ? next.set(key, value) : next.delete(key); next.set("page", "1"); setParams(next); }
