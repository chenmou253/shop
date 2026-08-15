import { useEffect, useState } from "react";
import { CalendarDays } from "lucide-react";
import { Link, useParams } from "react-router-dom";
import { api } from "../api";
import { PageHero } from "../components/Common";
import { SubscribeBox } from "../components/Layout";
import type { Article } from "../types";

export function ArticleDetail() {
  const { slug } = useParams();
  const [data, setData] = useState<{ article: Article; latest: Article[] } | null>(null);
  useEffect(() => { api<{ article: Article; latest: Article[] }>(`/articles/${slug}`).then(setData); }, [slug]);
  if (!data) return <div className="state-page"><span className="loader" />Loading...</div>;
  const { article } = data;
  return <><PageHero title={article.title} crumbs={[label(article.content_type)]} /><div className="shell two-column page-space"><article className="main-column article-detail"><span className="article-category">{article.category}</span><h1>{article.title}</h1><time><CalendarDays size={16} />{new Date(article.published_at).toLocaleDateString("en-US", { month: "long", day: "numeric", year: "numeric" })}</time>{article.cover_image && <img className="article-hero-image" src={article.cover_image} alt={article.title} />}<div className="rich-content" dangerouslySetInnerHTML={{ __html: article.content }} /><div className="article-inquiry"><h3>Have a project in mind?</h3><p>Send drawings, sample pictures, a BOM or specifications. We reply within 24 hours.</p><Link className="button" to="/feedback">Send Inquiry</Link></div></article><aside className="sidebar-content"><div className="side-card"><h3>Latest {label(article.content_type)}</h3><ul className="plain-links">{data.latest.map((item) => <li key={item.id}><Link to={`/article/${item.slug}`}>{item.title}</Link></li>)}</ul></div><SubscribeBox /></aside></div></>;
}

function label(type: string) { return type === "news" ? "News" : type === "knowledge" ? "Knowledge" : "Blog"; }
