import { useEffect, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { api } from "../api";
import { PageHero, ProductCard } from "../components/Common";
import type { Article, Product } from "../types";

export function SearchPage() {
  const [params] = useSearchParams(); const query = params.get("q") || "";
  const [data, setData] = useState<{ products: Product[]; articles: Article[] }>({ products: [], articles: [] });
  useEffect(() => { api<{ products: Product[]; articles: Article[] }>(`/search?q=${encodeURIComponent(query)}`).then(setData); }, [query]);
  return <><PageHero title="Search Results" /><section className="shell page-space search-results"><p className="lead">Results for: <b>{query}</b></p><h2>Products</h2>{data.products.length ? <div className="product-grid catalog-grid">{data.products.map((p) => <ProductCard key={p.id} product={p} inquiry />)}</div> : <p>No matching products.</p>}<h2>Articles</h2>{data.articles.length ? <ul>{data.articles.map((a) => <li key={a.id}><Link to={`/article/${a.slug}`}>{a.title}</Link><span>{a.summary}</span></li>)}</ul> : <p>No matching articles.</p>}</section></>;
}
