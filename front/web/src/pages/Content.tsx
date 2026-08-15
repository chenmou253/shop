import { useEffect, useState } from "react";
import { ArrowRight, CalendarDays } from "lucide-react";
import { Link, useSearchParams } from "react-router-dom";
import { api } from "../api";
import { PageHero, Pagination } from "../components/Common";
import { SubscribeBox } from "../components/Layout";
import type { Article } from "../types";

type ContentType = "news" | "knowledge" | "blog";
type Response = { items: Article[]; total: number; page: number; page_size: number };

export function ContentList({ type }: { type: ContentType }) {
  const [params, setParams] = useSearchParams();
  const [data, setData] = useState<Response>({ items: [], total: 0, page: 1, page_size: 10 });
  const page = Number(params.get("page") || 1);
  const title = type === "news" ? "Customer Cases & News" : type === "knowledge" ? "Knowledge" : "Blog";
  useEffect(() => { const search = new URLSearchParams(params); search.set("type", type); search.set("page", String(page)); search.set("page_size", "10"); api<Response>(`/articles?${search}`).then(setData); }, [type, params.toString(), page]);
  return <><PageHero title={title} /><div className="shell two-column page-space"><section className="main-column"><div className="article-list">{data.items.map((article) => <article className="article-row" key={article.id}><Link className="article-cover" to={`/article/${article.slug}`}><img src={article.cover_image} alt={article.title} /></Link><div><span className="article-category">{article.category}</span><h2><Link to={`/article/${article.slug}`}>{article.title}</Link></h2><time><CalendarDays size={15} />{new Date(article.published_at).toLocaleDateString("en-US", { month: "long", day: "numeric", year: "numeric" })}</time><p>{article.summary}</p><Link className="more-link" to={`/article/${article.slug}`}>Read More <ArrowRight size={16} /></Link></div></article>)}</div><Pagination page={data.page} pageSize={data.page_size} total={data.total} onChange={(n) => { const next = new URLSearchParams(params); next.set("page", String(n)); setParams(next); }} /></section><aside className="sidebar-content"><div className="side-card"><h3>Content Categories</h3><ul className="plain-links">{categories[type].map((item) => <li key={item}><Link to={`/${type}?category=${encodeURIComponent(item)}`}>{item}</Link></li>)}</ul></div><SubscribeBox /><div className="side-card"><h3>Popular Posts</h3><ul className="plain-links">{data.items.slice(0, 5).map((item) => <li key={item.id}><Link to={`/article/${item.slug}`}>{item.title}</Link></li>)}</ul></div></aside></div></>;
}

const categories: Record<ContentType, string[]> = {
  news: ["Customer orders case", "shipment news", "industrial knowledge", "Company culture", "Customer feedback"],
  knowledge: ["brass valve", "Valve", "fittings", "material", "hydraulic fittings", "Stainless steel valve", "Syphon tube"],
  blog: ["Product Guide", "Engineering", "Materials", "Maintenance"]
};
