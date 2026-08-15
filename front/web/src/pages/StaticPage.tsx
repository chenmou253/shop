import { useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { api } from "../api";
import { CategorySidebar, PageHero } from "../components/Common";
import type { PageData } from "../types";

export function StaticPage({ fixedSlug }: { fixedSlug?: string }) {
  const params = useParams();
  const slug = fixedSlug || params.slug || "about-us";
  const [page, setPage] = useState<PageData | null>(null);
  useEffect(() => { setPage(null); api<{ page: PageData }>(`/pages/${slug}`).then((data) => setPage(data.page)); }, [slug]);
  if (!page) return <div className="state-page"><span className="loader" />Loading...</div>;
  return <><PageHero title={page.title} /><div className="shell two-column page-space"><article className="main-column static-page">{page.cover_image && <img className="page-cover" src={page.cover_image} alt={page.title} />}<h1>{page.title}</h1>{page.subtitle && <p className="lead">{page.subtitle}</p>}<div className="rich-content" dangerouslySetInnerHTML={{ __html: page.content }} /><div className="article-inquiry"><h3>Talk to a Depo engineer</h3><p>Engineering support is available for OEM / ODM projects.</p><Link className="button" to="/feedback">Contact Now</Link></div></article><CategorySidebar /></div></>;
}
